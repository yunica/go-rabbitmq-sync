package main

import (
	"backend/internal/config"
	"context"
	"github.com/joho/godotenv"
	"log"
	"sync"

	"backend/internal/model"
	"backend/internal/rabbitmq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Starting publisher...")

	producer := rabbitmq.NewProducer(config.QUEUE_NAME)
	defer producer.Close()

	events := []model.PaymentEvent{
		{UserID: 1, PaymentID: 1, DepositAmount: 10},
		{UserID: 1, PaymentID: 2, DepositAmount: 20},
		{UserID: 2, PaymentID: 3, DepositAmount: 20},
		// bug
		{UserID: 1, PaymentID: 1, DepositAmount: 10},
	}

	var wg sync.WaitGroup
	for _, e := range events {
		wg.Add(1)
		go func(event model.PaymentEvent) {
			defer wg.Done()
			ctx := context.Background()
			if err := producer.Publish(ctx, event); err != nil {
				log.Printf("Failed to publish event: %v", err)
			} else {
				log.Printf("Published event: %+v", event)
			}
		}(e)
	}
	wg.Wait()

}
