package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	connStr := os.Getenv("POSTGRES_DSN") // e.g. "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}
	log.Println("Connected to PostgreSQL")
}

func Stop() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Error Closing Connection")
		} else {
			fmt.Println("PostgreSQL connection closed successfully!")
		}
	}
}
