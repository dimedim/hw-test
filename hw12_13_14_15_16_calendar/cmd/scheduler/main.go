package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
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

// import "fmt"

// func main() {
// 	fmt.Println("sheduler producer")
// }

// TODO: набросок реббит
// TODO: flag parse
var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/rabbit.yaml", "Path to configuration file")
}

func main() {
	// HelloExample()

	flag.Parse()

	//? Конфиг
	// cfgPath := flag.String("config", "configs/rabbit.yaml", "path to config file")
	// flag.Parse()

	cfg := config.LoadRabbitCfg(configFile)

	//? Логгер
	logger := logger.New(cfg.Logger.Level, os.Stdout)

	// cfg, err := config.Load(*cfgPath)
	// if err != nil {
	// 	log.Fatalf("config load: %v", err)
	// }
	// _ = cfgPath

	//? Подключиться к БД
	ctx := context.Background()
	stor := storage.NewStorage(ctx, cfg.DB.Type, cfg.DB.DSN, "")
	defer stor.Close()

	//? Подключиться к КроликуМэКу
	client, err := rabbitmq.New("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer client.Close()

	// ? объявить exchange/queue и биндинг.

	err = client.Setup(ctx, cfg.Rabbit.ExchName, cfg.Rabbit.ExchType, cfg.Rabbit.Queue, cfg.Rabbit.Key)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	//? Обработка сигнала прерывания
	sigCtx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	//? ticker
	// TODO: cfg.Scheduler.IntervalSeconds
	ticker := time.NewTicker(cfg.Scheduler.Interval)
	defer ticker.Stop()

	logger.Info("Start scheduler")
	fmt.Println(cfg)

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

// TODO:
func ProcessMsg(
	ctx context.Context,
	stor storage.EventStorage,
	client rabbitmq.Client,
	cfg *config.RabbitConfig,
	logger logger.Logger,
	now time.Time,
) {
	//? ListEventsToNotify(now) — выбрать все события, у которых notify_before

	// Сформировать объекты Notification и сериализовать их в JSON.

	// Опубликовать в очередь (Publish(exchange, routingKey, body)).

	// Удалить события старее года (DeleteOlderThan(now.AddDate(-1,0,0))).

	// логировать
	events, err := stor.ListEventsToNotify(ctx, now)
	if err != nil {
		log.Printf("ListEventsToNotify error: %v", err)
		return
	}
	log.Printf("→ Found %d events to notify", len(events))

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
			log.Printf("publish error: %v", err)
		} else {
			log.Printf("→ Published notification for event %s", ev.ID)
		}
	}

	expire := now.AddDate(0, 0, -cfg.Scheduler.DeleteDays)
	delCount, err := stor.DeleteOlderThan(ctx, expire)
	if err != nil {
		log.Printf("db.DeleteOlderThan error: %v", err)
	} else if delCount > 0 {
		log.Printf("→ Deleted %d old events", delCount)
	}

	//! TODO: DELETE
	notif := models.Notification{
		EventID: "123123",
		Title:   "TUOWAODWJWADJ",
	}
	data, _ := json.Marshal(notif)
	client.Publish(ctx, cfg.Rabbit.ExchName, cfg.Rabbit.Key, data)
	fmt.Printf("Send msg: %v\n", notif)
}
