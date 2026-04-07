package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	conn := "host=localhost port=5432 user=books password=password dbname=books sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL database")
}
