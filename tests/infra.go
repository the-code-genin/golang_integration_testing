package tests

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresImage    = "postgres:16.11"
	postgresUser     = "postgres"
	postgresPassword = "password"
	postgresDB       = "postgres"
)

type CleanupFunc func() error

// Spins up a PostgreSQL container, applies the up migrations,
// and returns a pgx connection and a cleanup function which will run the downMigrations and terminate the container.
func SetupPostgresDB(
	ctx context.Context, upMigrations, downMigrations *bytes.Buffer,
) (container *testcontainers.DockerContainer, conn *pgx.Conn, cleanupFunc CleanupFunc, err error) {
	// Spin up a postgres database
	container, err = testcontainers.Run(
		ctx, postgresImage,
		testcontainers.WithExposedPorts("5432/tcp"),
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDB,
		}),
		testcontainers.WithFiles(
			testcontainers.ContainerFile{
				Reader:            upMigrations,
				ContainerFilePath: "/docker-entrypoint-initdb.d/migrations.up.sql",
				FileMode:          0o644,
			},
		),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to start postgres container: %w", err)
	}

	// Get the host and mapped port
	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	// Connect to the DB
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		postgresUser, postgresPassword, host, port.Port(), postgresDB,
	)
	conn, err = pgx.Connect(ctx, connStr)
	if err != nil {
		container.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Cleanup function to run down migrations, close connection and terminate container
	cleanupFunc = func() error {
		_, execErr := conn.Exec(ctx, downMigrations.String())
		if execErr != nil {
			return fmt.Errorf("warning: failed to apply down migrations: %w", execErr)
		}

		if err := conn.Close(ctx); err != nil {
			return err
		}

		if err := container.Terminate(ctx); err != nil {
			return err
		}

		return nil
	}

	return container, conn, cleanupFunc, nil
}
