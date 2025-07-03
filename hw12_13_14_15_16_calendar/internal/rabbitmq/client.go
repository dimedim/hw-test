package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client interface {
	Setup(ctx context.Context, exch, exchType, queue, key string) error
	Publish(ctx context.Context, exch, key string, body []byte) error
	Consume(ctx context.Context, queue string) (<-chan amqp.Delivery, error)
	Close() error
}

type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func New(url string) (Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitClient{conn: conn, ch: ch}, nil
}

func (c *RabbitClient) Setup(_ context.Context, exch, exchType, queue, key string) error {
	if err := c.ch.ExchangeDeclare(
		exch,     // имя
		exchType, // тип ("direct", "fanout", "topic")
		true,     // durable «Сохраняй эту очередь на диске, чтобы при перезапуске RabbitMQ она осталась».
		false,    // autoDelete Не удалять очередь автоматически, когда из неё перестанут потреблять
		false,    // internal Разрешаю публиковать в этот exchange из клиентов».
		// Если выставить true, можно использовать exchange только для промежуточного маршрута,
		//  но не для внешних публикаций.
		false, // noWait
		nil,   // args
	); err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}

	q, err := c.ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive Эта очередь может использоваться разными соединениями и клиентами
		false, // no-wait Жду от сервера подтверждения, что очередь объявлена
		// Если true — клиент просто шлёт команду и не дожидается ответа (асинхронно).
		nil, // arguments
	)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	if err := c.ch.QueueBind(
		q.Name,
		key,  // routing key
		exch, // exchange
		false,
		nil,
	); err != nil {
		return fmt.Errorf("queue bind: %w", err)
	}
	return nil
}

func (c *RabbitClient) Publish(ctx context.Context, exch, key string, body []byte) error {
	return c.ch.PublishWithContext(
		ctx,
		exch,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (c *RabbitClient) Consume(ctx context.Context, queue string) (<-chan amqp.Delivery, error) {
	return c.ch.ConsumeWithContext(
		ctx,
		queue,
		"", // consumerTag — это просто строковое имя для «потребителя» (твоего процесса).
		// Если оставить его пустым (""), RabbitMQ сам сгенерирует уникальный тег.
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args

	)
}

func (c *RabbitClient) Close() error {
	err := c.ch.Close()
	if err != nil {
		c.conn.Close()
		return err
	}
	return c.conn.Close()
}
