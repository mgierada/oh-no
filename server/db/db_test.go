package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var postgresContainer testcontainers.Container

func TestMain(m *testing.M) {
	// Setup before running tests
	if err := setup(); err != nil {
		log.Fatalf("Could not set up test container: %v", err)
	}

	// Ensure db is properly initialized before running tests
	if db == nil {
		log.Fatalf("Database connection is not initialized")
	}

	// Run tests
	code := m.Run()

	// Teardown after tests
	if err := teardown(); err != nil {
		log.Fatalf("Could not tear down test container: %v", err)
	}

	os.Exit(code)
}

func setup() error {
	ctx := context.Background()

	// Start a PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	var err error
	postgresContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return fmt.Errorf("failed to get container host: %w", err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return fmt.Errorf("failed to get container port: %w", err)
	}

	dsn := fmt.Sprintf("postgres://postgres:password@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create counter table
	_, err = db.Exec(`CREATE TABLE counter (
		current_value INT NOT NULL,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		reseted_at TIMESTAMP NULL DEFAULT NULL
	)`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func teardown() error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return postgresContainer.Terminate(context.Background())
}

func TestGetCounter(t *testing.T) {
	// Insert a row into the counter table
	_, err := db.Exec(`INSERT INTO counter (current_value, updated_at, reseted_at) VALUES (42, '2024-05-30 12:34:56', '2024-05-01 12:00:00')`)
	if err != nil {
		t.Fatalf("failed to insert into table: %s", err)
	}

	// Test GetCounter function
	counter, err := GetCounter()
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 42 {
		t.Errorf("expected current_value to be 42, got %d", counter.CurrentValue)
	}
	if counter.UpdatedAt != "2024-05-30T12:34:56Z" {
		t.Errorf("expected updated_at to be '2024-05-30T12:34:56Z', got %s", counter.UpdatedAt)
	}
	if counter.ResetedAt.String != "2024-05-01T12:00:00Z" {
		t.Errorf("expected reseted_at to be '2024-05-01T12:00:00Z', got %s", counter.ResetedAt.String)
	}

	// Cleanup table after test
	if _, err = db.Exec(`DELETE FROM counter`); err != nil {
		t.Fatalf("failed to clean up table: %s", err)
	}
}

func TestGetCounterEmpty(t *testing.T) {
	// Ensure table is empty
	_, err := db.Exec(`DELETE FROM counter`)
	if err != nil {
		t.Fatalf("failed to clean up table: %s", err)
	}

	// Test GetCounter function
	counter, err := GetCounter()
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 0 {
		t.Errorf("expected current_value to be 0, got %d", counter.CurrentValue)
	}
	if counter.UpdatedAt != "" {
		t.Errorf("expected updated_at to be '', got %s", counter.UpdatedAt)
	}
	if counter.ResetedAt.Valid {
		t.Errorf("expected reseted_at to be invalid, got valid")
	}
}
