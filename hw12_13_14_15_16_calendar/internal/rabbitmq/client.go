package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

// TODO: сейчас это просто заглушка
/*
yaml
rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
  exchange: "events"
  queue: "notifications"
  routing_key: "notify"
scheduler:
  interval_seconds: 60   # как часто сканировать
  delete_older_than_days: 365
*/
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewClient(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &Client{conn: conn, channel: ch}, nil
}

func (c *Client) Setup(exchange, queue, routingKey string) error {
	// объявление exchange и queue, связывание
	// channel.ExchangeDeclare(...)
	// channel.QueueDeclare(...)
	// channel.QueueBind(...)
	return nil
}

func (c *Client) Publish(exchange, routingKey string, body []byte) error {
	return c.channel.Publish(exchange, routingKey, false, false,
		amqp.Publishing{ContentType: "application/json", Body: body})
}

func (c *Client) Consume(queue string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(queue, "", true, false, false, false, nil)
}

func (c *Client) Close() {
	c.channel.Close()
	c.conn.Close()
}
