package main

import (
	"backend/internal/config"
	"backend/internal/model"
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"backend/internal/db"
	"backend/internal/rabbitmq"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.Init()

	log.Println("Starting consumer...")

	consumer := rabbitmq.NewConsumer(config.QUEUE_NAME)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// listenning
	handleMessage := func(event model.PaymentEvent) error {
		err := db.InsertPaymentEvent(event)

		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				log.Printf("Duplicate detected for payment_id %d, storing in skipped_message", event.PaymentID)
				insertErr := db.InsertSkippedMessage(event)
				if insertErr != nil {
					log.Printf("Failed to insert into skipped_messages: %v", insertErr)
				}
			}
		} else {
			log.Println("Inserted Payment Event", event)
		}
		return nil
	}

	// go routine
	go func() {
		if err := consumer.Consume(ctx, handleMessage); err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
	}()

	// live consumer
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down consumer...")
}
