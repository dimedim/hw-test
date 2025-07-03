package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/config"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/rabbitmq"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/pkg/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func main() {
	flag.Parse()
	cfg := config.LoadRabbitCfg(configFile)
	logger := logger.New(cfg.Logger.Level, os.Stdout)
	client, err := rabbitmq.New(cfg.Rabbit.URL)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer client.Close()

	ctx := context.Background()

	err = client.Setup(ctx, cfg.Rabbit.ExchName, cfg.Rabbit.ExchType, cfg.Rabbit.Queue, cfg.Rabbit.Key)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	sigCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	ch, err := client.Consume(ctx, cfg.Rabbit.Queue)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	go func() {
		for {
			select {
			case <-sigCtx.Done():
				logger.Info("Shutting down…")
				client.Close()
				return
			default:
				for msg := range ch {
					var data models.Notification

					if err := json.Unmarshal(msg.Body, &data); err != nil {
						logger.Error("invalid message format", slog.String("err", err.Error()))
						msg.Nack(false, false)
						continue
					}
					fmt.Fprintf(os.Stdout, "GOT MESSAGE: %+v\n", data)
				}
			}
		}
	}()
	logger.Info("Wait messages...")
	<-sigCtx.Done()
	logger.Info("Shutting DOWN…")
	client.Close()
}
