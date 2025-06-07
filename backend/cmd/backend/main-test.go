package main

import (
	"backend/internal/db"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.Init()
	// read db
	//payments, err := db.GetAllPaymentEvents()

	//if err != nil {
	//	log.Fatalf("Error getting payments: %v", err)
	//}

	// insert db
	//for i := 100; i < 120; i++ {
	//	payment := model.PaymentEvent{12, i, i * 2}
	//	err_ := db.InsertPaymentEvent(payment)
	//	if err_ != nil {
	//		log.Fatal(err_)
	//	}
	//}
}
