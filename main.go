package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/the-code-genin/golang_integration_testing/http"
	"github.com/the-code-genin/golang_integration_testing/repository"
	"github.com/the-code-genin/golang_integration_testing/service"
)

// Config holds environment variables
type Config struct {
	PostgresUser     string `envconfig:"POSTGRES_USER" default:"postgres"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:"password"`
	PostgresHost     string `envconfig:"POSTGRES_HOST" default:"localhost"`
	PostgresPort     string `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresDB       string `envconfig:"POSTGRES_DB" default:"postgres"`
	ServerPort       int    `envconfig:"SERVER_PORT" default:"8080"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	// Build Postgres connection string
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB,
	)

	// Connect to Postgres using pgxpool
	connPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer connPool.Close()

	// Initialize repository and service
	repo := repository.NewRepository(connPool)
	svc := service.NewService(repo)

	// Start the HTTP server
	server := http.NewServer(svc)
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("starting server on %s...", addr)
	if err := server.Start(addr); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
