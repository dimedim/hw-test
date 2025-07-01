package main

// import "fmt"

// func main() {
// 	fmt.Println("sheduler producer")
// }

import (
	"context"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// TODO: набросок реббит
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
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body1 := "Hello!!!"

	i := 0
	for {
		i++
		body := body1 + strconv.Itoa(i)
		err = ch.PublishWithContext(
			ctx,
			"",
			queue.Name,
			false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)

		if err != nil {
			log.Fatal(err)
		}
		log.Printf(" [x] Sent %s\n", body)
		time.Sleep(time.Second * 3)
	}

}
