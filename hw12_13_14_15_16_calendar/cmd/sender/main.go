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

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/config"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/rabbitmq"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/pkg/logger"
)

// import (
// 	"fmt"
// )

// func main() {
// 	fmt.Println("sender consumer")
// }

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/rabbit.yaml", "Path to configuration file")
}
func main() {
	// HelloExample()

	//TODO:
	/*
			Задачи рассыльщика:

		Прочитать тот же rabbitmq-блок конфига (только url и queue).

		Подключиться к RabbitMQ и убедиться, что очередь существует.

		Вызвать Consume(queue) и в цикле читать Delivery.

		Для каждого сообщения распарсить в Notification и просто вывести в STDOUT или лог:

		При SIGINT/TERM корректно закрыть соединение.
	*/

	cfg := config.LoadRabbitCfg(configFile)
	logger := logger.New(cfg.Logger.Level, os.Stdout)
	client, err := rabbitmq.New("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer client.Close()

	ctx := context.Background()

	err = client.Setup(ctx, cfg.Rabbit.ExchName, cfg.Rabbit.ExchType, cfg.Rabbit.Queue, cfg.Rabbit.Key)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sigCtx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	ch, err := client.Consume(ctx, cfg.Rabbit.Queue)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
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
						log.Printf("invalid message format: %v", err)
						// отклоняем без переотправки
						msg.Nack(false, false)
						continue
					}
					fmt.Printf("GET MESSAGE!!: %v\n", data)
					// if err := msg.Ack(false); err != nil {
					// 	log.Printf("failed to ack message: %v", err)
					// }
				}
			}
		}
	}()
	logger.Info("Wait messages...")
	<-sigCtx.Done()
	logger.Info("Shutting DOWN…")
	client.Close()
}
