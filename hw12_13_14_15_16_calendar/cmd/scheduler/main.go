package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/config"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/rabbitmq"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/storage"
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

	ctx := context.Background()
	stor := storage.NewStorage(ctx, cfg.DB.Type, cfg.DB.DSN, "")
	defer stor.Close()

	client, err := rabbitmq.New(cfg.Rabbit.URL)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer client.Close()

	err = client.Setup(ctx, cfg.Rabbit.ExchName, cfg.Rabbit.ExchType, cfg.Rabbit.Queue, cfg.Rabbit.Key)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	sigCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	ticker := time.NewTicker(cfg.Scheduler.Interval)
	defer ticker.Stop()

	logger.Info("Start scheduler")

	for {
		select {
		case <-sigCtx.Done():
			logger.Info("Shutting down…")
			client.Close()
			return
		case now := <-ticker.C:
			ProcessMsg(ctx, stor, client, cfg, logger, now)
		}
	}
}

func ProcessMsg(
	ctx context.Context,
	stor storage.EventStorage,
	client rabbitmq.Client,
	cfg *config.RabbitConfig,
	logger logger.Logger,
	now time.Time,
) {
	events, err := stor.ListEventsToNotify(ctx, now)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info("Found events to notify", slog.Int("num of events", len(events)))

	for _, ev := range events {
		notif := models.Notification{
			EventID:  ev.ID,
			UserID:   ev.UserID,
			Title:    ev.Title,
			StartsAt: ev.StartsAt,
		}
		data, err := json.Marshal(notif)
		if err != nil {
			logger.Error(err.Error())
		}
		if err := client.Publish(ctx, cfg.Rabbit.ExchName, cfg.Rabbit.Key, data); err != nil {
			logger.Error("publish", slog.String("err", err.Error()))
		} else {
			logger.Info("Published notification for event", slog.String("eventID", ev.ID))
		}
	}

	expire := now.AddDate(0, 0, -cfg.Scheduler.DeleteDays)
	delCount, err := stor.DeleteOlderThan(ctx, expire)
	if err != nil {
		logger.Error("DeleteOlderThan", slog.String("err", err.Error()))
	} else if delCount > 0 {
		logger.Info("Deleted", slog.Int("Num of deleted events", delCount))
	}
}
