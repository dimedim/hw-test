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
		exch,
		exchType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}

	q, err := c.ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	if err := c.ch.QueueBind(
		q.Name,
		key,
		exch,
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
		"",
		true,
		false,
		false,
		false,
		nil,
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
