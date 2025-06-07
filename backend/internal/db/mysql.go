package db

import (
	_ "github.com/go-sql-driver/mysql"

	"backend/internal/config"
	"database/sql"
	"log"
)

var DB *sql.DB

func Init() {
	var err error

	DB, err = sql.Open("mysql", config.MySQLConnectionString())
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Error pinging DB: %v", err)
	}
}
