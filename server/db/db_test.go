package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGetCounter(t *testing.T) {
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
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer postgresContainer.Terminate(ctx)

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %s", err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container port: %s", err)
	}

	dsn := fmt.Sprintf("postgres://postgres:password@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %s", err)
	}
	defer db.Close()

	// Create counter table
	_, err = db.Exec(`CREATE TABLE counter (
			current_value INT NOT NULL,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			reseted_at TIMESTAMP NULL DEFAULT NULL
	)`)
	if err != nil {
		t.Fatalf("failed to create table: %s", err)
	}

	// Insert a row into the counter table
	_, err = db.Exec(`INSERT INTO counter (current_value, updated_at, reseted_at) VALUES (42, '2024-05-30 12:34:56', '2024-05-01 12:00:00')`)
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
}
