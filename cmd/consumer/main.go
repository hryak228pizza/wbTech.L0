package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	kafka "github.com/hryak228pizza/wbTech.L0/internal/kafka/consumer"
)

func main() {
	dsn := "user=postgres password=123 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	kafka.Consumer(db)
}