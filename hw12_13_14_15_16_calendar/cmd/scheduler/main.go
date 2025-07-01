package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/rabbitmq"
)

// import "fmt"

// func main() {
// 	fmt.Println("sheduler producer")
// }

// TODO: набросок реббит
func main() {
	// HelloExample()

	//? Логгер

	//? Конфиг
	cfgPath := flag.String("config", "configs/rabbit.yaml", "path to config file")
	flag.Parse()

	// cfg, err := config.Load(*cfgPath)
	// if err != nil {
	// 	log.Fatalf("config load: %v", err)
	// }
	_ = cfgPath

	//? Подключиться к БД

	//? Подключиться к КроликуМэКу
	client, err := rabbitmq.New("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// ? объявить exchange/queue и биндинг.

	ctx := context.Background()
	err = client.Setup(ctx, "", "", "", "")
	if err != nil {
		log.Fatal(err)
	}

	//? Обработка сигнала прерывания
	sigCtx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	//? ticker
	// TODO: cfg.Scheduler.IntervalSeconds
	ticker := time.NewTicker(time.Duration(1) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigCtx.Done():
			client.Close()
		case now := <-ticker.C:
			ProcessMsg(now)
		}
	}
}

// TODO:
func ProcessMsg(now time.Time) {
	//? ListEventsToNotify(now) — выбрать все события, у которых notify_before

	// Сформировать объекты Notification и сериализовать их в JSON.

	// Опубликовать в очередь (Publish(exchange, routingKey, body)).

	// Удалить события старее года (DeleteOlderThan(now.AddDate(-1,0,0))).

	// логировать
}

// func runIteration(
//   ctx context.Context,
//   db store.Store,
//   rmq rabbitmq.Client,
//   cfg *config.Config,
//   now time.Time,
// ) {
//   // 1) Выбираем события, нуждающиеся в напоминании
//   events, err := db.ListEventsToNotify(ctx, now)
//   if err != nil {
//     log.Printf("db.ListEventsToNotify error: %v", err)
//     return
//   }
//   log.Printf("→ Found %d events to notify", len(events))

//   // 2) Публикуем каждое в очередь
//   for _, ev := range events {
//     notif := models.Notification{
//       EventID:   ev.ID,
//       UserID:    ev.UserID,
//       Title:     ev.Title,
//       EventTime: ev.Time,
//     }
//     data, _ := json.Marshal(notif)
//     if err := rmq.Publish(ctx, cfg.RabbitMQ.Exchange, cfg.RabbitMQ.RoutingKey, data); err != nil {
//       log.Printf("publish error: %v", err)
//     } else {
//       log.Printf("→ Published notification for event %d", ev.ID)
//     }
//   }

//   // 3) Удаляем старые события
//   cutoff := now.AddDate(0, 0, -cfg.Scheduler.DeleteOlderThanDays)
//   delCount, err := db.DeleteOlderThan(ctx, cutoff)
//   if err != nil {
//     log.Printf("db.DeleteOlderThan error: %v", err)
//   } else if delCount > 0 {
//     log.Printf("→ Deleted %d old events", delCount)
//   }
// }
