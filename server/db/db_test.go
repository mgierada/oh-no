package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS counter (
		current_value INT NOT NULL,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		reseted_at TIMESTAMP NULL DEFAULT NULL
	)`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Create historical table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS historical_counters (
			counter_id UUID PRIMARY KEY NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			value INT NOT NULL
		);
	`)
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

// NOTE: Test covers a situation where there are no existing rows in the counter table and we simply need to create one with CurrentValue=1 and UpdatedAt=NOW(). ResetedAt should be a sql null string
func TestUpsertCounterData(t *testing.T) {
	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Call the UpsertCounterData function
	err := UpsertCounterData()
	if err != nil {
		t.Fatalf("failed to upsert counter data: %s", err)
	}

	// Test GetCounter function
	counter, err := GetCounter()
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	// Check current value
	if counter.CurrentValue != 1 {
		t.Errorf("expected current_value to be 1, got %d", counter.CurrentValue)
	}

	// Check updated_at (should be close to current time)
	expectedTime := time.Now().UTC()
	parsedUpdatedAt, err := time.Parse(time.RFC3339, counter.UpdatedAt)
	if err != nil {
		t.Fatalf("failed to parse updated_at: %s", err)
	}

	// Allow for a small time difference (e.g., 5 seconds)
	if expectedTime.Sub(parsedUpdatedAt).Seconds() > 5 {
		t.Errorf("expected updated_at to be close to '%s', got '%s'", expectedTime, counter.UpdatedAt)
	}

	// Check reseted_at (should be null)
	if counter.ResetedAt.Valid {
		t.Errorf("expected reseted_at to be null, got %v", counter.ResetedAt)
	}

	// Cleanup table after test
	if _, err = db.Exec(`DELETE FROM counter`); err != nil {
		t.Fatalf("failed to clean up table: %s", err)
	}
}

// NOTE: Test covers a typical situation where there are some existing rows in the counter table and we simply need to update the counter.
// It is expected to update a increment the counter by one and update the updated_at field to NOW().
func TestUpsertCounterDataTypicalCase(t *testing.T) {

	// Insert a row into the counter table
	_, ierr := db.Exec(`INSERT INTO counter (current_value, updated_at) VALUES (42, '2024-05-28 12:34:56')`)
	if ierr != nil {
		t.Fatalf("failed to insert into table: %s", ierr)
	}

	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Test UpsertCounterData
	err := UpsertCounterData()
	if err != nil {
		t.Fatalf("failed to upsert counter data: %s", err)
	}

	// Test GetCounter function
	counter, err := GetCounter()
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 43 {
		t.Errorf("expected current_value to be 43, got %d", counter.CurrentValue)
	}

	// Check updated_at (should be close to current time)
	expectedTime := time.Now().UTC()
	parsedUpdatedAt, err := time.Parse(time.RFC3339, counter.UpdatedAt)
	if err != nil {
		t.Fatalf("failed to parse updated_at: %s", err)
	}

	// Allow for a small time difference - 1 seconds
	if expectedTime.Sub(parsedUpdatedAt).Seconds() > 1 {
		t.Errorf("expected updated_at to be close to '%s', got '%s'", expectedTime, counter.UpdatedAt)
	}

	if counter.ResetedAt.Valid {
		t.Errorf("expected reseted_at to be null, got %v", counter.ResetedAt)
	}

	// Cleanup table after test
	if _, err = db.Exec(`DELETE FROM counter`); err != nil {
		t.Fatalf("failed to clean up table: %s", err)
	}
}

// NOTE: Test covers a situation where there are some existing rows in the counter table and we
// cannot update the counter because 24h did not pass since the last update.
// It is expected that no counter is updated.
func TestUpsertCounterDataTimeDidNotPass(t *testing.T) {
	// Check updated_at (should be close to current time)
	updatedLessThan24hAgo := time.Now().UTC().Add(-23 * time.Hour)
	parsedUpdatedLessThan24hAgo := updatedLessThan24hAgo.Format(time.RFC3339)

	// Insert a row into the counter table
	_, ierr := db.Exec(`INSERT INTO counter (current_value, updated_at) VALUES (42, $1)`, parsedUpdatedLessThan24hAgo)
	if ierr != nil {
		t.Fatalf("failed to insert into table: %s", ierr)
	}

	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Test UpsertCounterData
	err := UpsertCounterData()
	if err != nil {
		t.Fatalf("failed to upsert counter data: %s", err)
	}

	// Test GetCounter function
	counter, err := GetCounter()
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 42 {
		t.Errorf("expected current_value to be 42, got %d", counter.CurrentValue)
	}

	// Check updated_at (should be intact)
	if counter.UpdatedAt != parsedUpdatedLessThan24hAgo {
		t.Errorf("expected updated_at: '%s' to not change, got '%s'", parsedUpdatedLessThan24hAgo, counter.UpdatedAt)
	}

	if counter.ResetedAt.Valid {
		t.Errorf("expected reseted_at to be null, got %v", counter.ResetedAt)
	}

	// Cleanup table after test
	if _, err = db.Exec(`DELETE FROM counter`); err != nil {
		t.Fatalf("failed to clean up table: %s", err)
	}
}
