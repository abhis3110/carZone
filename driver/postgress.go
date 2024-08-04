package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

var db *sql.DB

func InitDB() {
	connStr := fmt.Sprintf("host = %s port = %s user = %s password = %s dbname = %s sslmode = disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME "))

	fmt.Println("Connecting to PostgreSQL...")
	time.Sleep(5 * time.Second)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening in database", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	fmt.Println("Successfully connected to database")
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("Error in closing the Database : %v", err)
	}
}
