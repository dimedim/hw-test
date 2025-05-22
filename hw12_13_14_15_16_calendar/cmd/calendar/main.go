package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
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
	logg.Info("awdawdawd")

	stor := memorystorage.New()

	for {

		testEvent := &storage.Event{ID: "123", CreatedAt: time.Now().UTC()}
		stor.CreateEvent(context.Background(), testEvent)

		fmt.Println(stor.DB["123"])
		time.Sleep(time.Minute)
	}
	// calendar := app.New(logg, storage)

	// server := internalhttp.NewServer(logg, calendar)

	// ctx, cancel := signal.NotifyContext(context.Background(),
	// 	syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	// defer cancel()

	// go func() {
	// 	<-ctx.Done()

	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// 	defer cancel()

	// 	if err := server.Stop(ctx); err != nil {
	// 		logg.Error("failed to stop http server: " + err.Error())
	// 	}
	// }()

	// logg.Info("calendar is running...")

	// if err := server.Start(ctx); err != nil {
	// 	logg.Error("failed to start http server: " + err.Error())
	// 	cancel()
	// 	os.Exit(1) //nolint:gocritic
	// }
}
