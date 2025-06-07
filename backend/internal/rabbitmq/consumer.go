package rabbitmq

import (
	"backend/internal/config"
	"backend/internal/model"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewConsumer(queueName string) *Consumer {
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

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   q,
	}
}

func (c *Consumer) Consume(ctx context.Context, process func(model.PaymentEvent) error) error {
	msgs, err := c.channel.Consume(
		c.queue.Name,
		"",    // consumer tag
		false, // autoAck false
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				// Channel closed
				return nil
			}

			var event model.PaymentEvent

			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error decoding message: %v", err)
				// Reject message
				d.Nack(false, false)
				continue
			}

			if err := process(event); err != nil {
				log.Printf("Processing error: %v", err)
				// Reject and requeue
				d.Nack(false, true)
				continue
			}

			d.Ack(false)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.conn.Close()
}
