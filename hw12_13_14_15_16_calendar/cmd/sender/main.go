package main

// import (
// 	"fmt"
// )

// func main() {
// 	fmt.Println("sender consumer")
// }

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue. Error: %s", err)
	}

	messages, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer. Error: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()

	go func() {
		for message := range messages {
			log.Printf("GET MESSAGE!!!!!: %s\n", message.Body)
		}
	}()

	log.Println("wait msges")
	<-ctx.Done()
}
