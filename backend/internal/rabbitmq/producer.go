package rabbitmq

import (
	"backend/internal/config"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewProducer(queueName string) *Producer {
	connStr := config.RabbitMQURL()
	if connStr == "" {
		log.Fatal("RABBITMQ_CONNECTION_STRING not set")
	}

	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect : %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	return &Producer{
		conn:    conn,
		channel: ch,
		queue:   q,
	}
}

func (p *Producer) Publish(ctx context.Context, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = p.channel.PublishWithContext(ctx,
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}

func (p *Producer) Close() {
	p.channel.Close()
	p.conn.Close()
}
