package main

import (
	"context"
	"database/sql"
	"flag"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/database"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/app"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/config"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/handlers"
	internalgrpc "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/server/grpc"
	internalhttp "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/server/http"
	mware "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/server/middleware"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/storage"
	memorystorage "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/storage/memory"
	sqlstorage "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/storage/sql"
	"github.com/pressly/goose/v3"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/pkg/logger"
	"github.com/gorilla/mux"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		PrintVersion()
		return
	}
	config := config.MustLoad(configFile)

	// logger
	filename := filepath.Join(config.Logger.LogFolder, time.Now().Format("2006-01-02_15-04-05")+".txt")
	file, err := logger.OpenLogFile(filename)
	if err != nil {
		log.Println("cant open log file", err)
	}
	defer file.Close()
	multOutput := io.MultiWriter(os.Stdout, file)
	multLogger := logger.New(config.Logger.Level, multOutput)
	log := logger.New(config.Logger.Level, os.Stdout)

	// app
	ctx := context.Background()

	storr := NewStorage(ctx, config)
	defer storr.Close()
	calendar := app.New(storr)
	handler := handlers.NewHadnlers(log, calendar)

	router := mux.NewRouter()
	router.Use(mware.LoggingMiddleware(multLogger))

	server := internalhttp.NewServer(config, log, router, handler)
	server.RegisterRoutes()

	// Shutdown
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()

		if err := server.Stop(ctxTimeout); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}
	}()

	// gRPC
	grpcServer := internalgrpc.New(calendar, config)
	grpcServer.Start(ctx, config, log)

	if err := server.Start(); err != nil {
		log.Error("server error", slog.String("error", err.Error()))
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func NewStorage(ctx context.Context, config *config.Config) storage.EventStorage {
	switch config.HTTP.DBType {
	case "memory":
		return memorystorage.New()
	case "postgres":
		pgxConn := database.MustConnectDatabase(ctx, config)
		psqlStorage := sqlstorage.New(pgxConn)
		if err := psqlStorage.Connect(ctx); err != nil {
			log.Fatal("cant connect to db: ", err)
		}
		migrate(ctx, pgxConn.DB, config.DB.MigrationFilepath)
		return psqlStorage
	}
	slog.Warn("storage type not set", slog.String("type", config.HTTP.DBType))
	return memorystorage.New()
}

func migrate(ctx context.Context, db *sql.DB, migrationsPath string) {
	err := goose.UpContext(ctx, db, migrationsPath)
	if err != nil {
		log.Fatal("migration error: %w", err)
	}
	// if err := goose.DownContext(ctx, db, migrationsPath); err != nil {
	// 	log.Fatal("down migration: %w", err)
	// }
}
