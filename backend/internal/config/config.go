package config

import (
	"os"
)

var QUEUE_NAME = "payment_events"

func MySQLConnectionString() string {
	return os.Getenv("MYSQL_CONNECTION_STRING")
}

func RabbitMQURL() string {
	return os.Getenv("RABBITMQ_CONNECTION_STRING")
}
