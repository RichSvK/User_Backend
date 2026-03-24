package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
)

func DatabaseConfig() *sql.DB {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	db.SetMaxIdleConns(28)
	db.SetMaxOpenConns(60)
	db.SetConnMaxIdleTime(3 * time.Minute)
	db.SetConnMaxLifetime(10 * time.Minute)
	return db
}
