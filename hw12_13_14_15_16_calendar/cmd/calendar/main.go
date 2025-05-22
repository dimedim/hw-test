package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/database"
	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

// ./calendar --config=/path/to/config.yaml

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		PrintVersion()
		return
	}

	config := config.LoadConfig(configFile)
	fmt.Println("config -->: ", config)
	logg := logger.New(config.Logger.Level)

	dummyCheck(logg)

	var storr app.Storage
	switch config.App.DBType {
	case "memory":
		storr = memorystorage.New()
	case "postgres":
		pgxConn := database.InitDatabase(context.Background(), config)
		psqlStorage := sqlstorage.New(pgxConn)
		// storr = psqlStorage
		_ = psqlStorage
	default:
		storr = memorystorage.New()
	}

	calendar := app.New(logg, storr)
	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func dummyCheck(logg logger.LoggerTodo) {
	logg.Info("awdawdawd")

	stor := memorystorage.New()

	for {

		testEvent := &storage.Event{ID: "123", CreatedAt: time.Now().UTC()}
		stor.CreateEvent(context.Background(), testEvent)

		fmt.Println(stor.DB["123"])
		time.Sleep(time.Minute)
	}
}
