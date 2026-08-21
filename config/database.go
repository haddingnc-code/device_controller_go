package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

// DB holds the global connection pool instance shared across application layers.
var DB *pgxpool.Pool

// ConnectDatabase initializes a high-performance connection pool with PostgreSQL.
// It reads settings from environment variables with safe fallback defaults.
func ConnectDatabase() {
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "devices_db")

	// Construct the standard PostgreSQL connection URI string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse database configuration URI: %v\n", err)
	}

	// Optimize pool settings for heavy production traffic workloads
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnIdleTime = 30 * time.Minute

	// Establish connection pool context with a strict timeout constraint
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	DB, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Failed to establish PostgreSQL connection pool: %v\n", err)
	}

	// Ping the server to guarantee the database engine is online and accessible
	if err := DB.Ping(ctx); err != nil {
		log.Fatalf("Database ping verification failed: %v\n", err)
	}

	log.Println("Successfully connected to the PostgreSQL database pool.")
}

// getEnv captures system environment properties or returns a fallback string value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
