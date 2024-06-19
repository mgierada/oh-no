package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"server/utils"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var postgresContainer testcontainers.Container

/*
Setup
*/
func TestMain(m *testing.M) {
	// Setup before running tests
	if err := setup(); err != nil {
		log.Fatalf("Could not set up test container: %v", err)
	}

	// Ensure db is properly initialized before running tests
	if db == nil {
		log.Fatalf("Database connection is not initialized")
	}

	// Set relevant environment variables
	os.Setenv("UPDATE_INTERVAL_IN_HOURS", "24")

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

	createCounterTableForTest(db, utils.TableInstance.Counter)
	createCounterTableForTest(db, utils.TableInstance.OhnoCounter)
	createTriggerForCounterTableForTest(db, utils.TableInstance.Counter)
	createTriggerForCounterTableForTest(db, utils.TableInstance.OhnoCounter)
	createHistoricalCounterForTest(db, utils.TableInstance.HistoricalCounter)
	createHistoricalCounterForTest(db, utils.TableInstance.HistoricalOhnoCounter)

	return nil
}

func createTriggerForCounterTableForTest(db *sql.DB, tableName string) error {
	var err error
	createTriggerFunctionQuery := `
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = now();
			IF NEW.current_value > NEW.max_value THEN
				NEW.max_value = NEW.current_value;
			END IF;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`
	_, err = db.Exec(createTriggerFunctionQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger function for %s table.\n %s", tableName, err)
	}

	// Create the trigger conditionally
	rawCreateTriggerQuery := `
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 
				FROM pg_trigger 
				WHERE tgname = 'update_updated_at'
			) THEN
				CREATE TRIGGER update_updated_at
				BEFORE UPDATE ON %s
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();
			END IF;
		END $$;
	`
	createTriggerQuery := fmt.Sprintf(rawCreateTriggerQuery, tableName)
	_, err = db.Exec(createTriggerQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger for %s table.\n %s", tableName, err)
	}
	return nil
}

func createCounterTableForTest(db *sql.DB, tableName string) error {
	var err error
	rawCreateQuery := `CREATE TABLE IF NOT EXISTS %s (
		current_value INT NOT NULL,
		max_value INT NULL DEFAULT 0,
		is_locked BOOLEAN NOT NULL DEFAULT FALSE,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		reseted_at TIMESTAMP NULL DEFAULT NULL
	);`
	createQuery := fmt.Sprintf(rawCreateQuery, tableName)
	_, err = db.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("failed to create table %s: %w", tableName, err)
	}
	return nil
}

func createHistoricalCounterForTest(db *sql.DB, tableName string) error {
	var err error
	rawCreateQuery := `CREATE TABLE IF NOT EXISTS %s (
			counter_id UUID PRIMARY KEY NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			value INT NOT NULL
	);`
	createQuery := fmt.Sprintf(rawCreateQuery, tableName)
	_, err = db.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("failed to create table %s: %w", tableName, err)
	}
	return nil
}

func cleanupTable(t *testing.T, tableName string) {
	var err error
	rawDeleteQuery := `
		DELETE FROM %s;
	`
	deleteQuery := fmt.Sprintf(rawDeleteQuery, tableName)
	if _, err = db.Exec(deleteQuery); err != nil {
		t.Fatalf("failed to clean up table: %s", err)
	}
}

func teardown() error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return postgresContainer.Terminate(context.Background())
}

/*
Test Cases
*/
func TestGetCounter(t *testing.T) {
	// Insert a row into the counter table
	tableName := utils.TableInstance.Counter
	rawInsertQuery := `
		INSERT INTO %s (current_value, is_locked, updated_at, reseted_at) 
		VALUES 
			(42, false, '2024-05-30 12:34:56', '2024-05-01 12:00:00');`
	insertQuery := fmt.Sprintf(rawInsertQuery, tableName)

	_, err := db.Exec(insertQuery)
	if err != nil {
		t.Fatalf("failed to insert into table: %s, err: %s", tableName, err)
	}

	counter, err := GetCounter(tableName)
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 42 {
		t.Errorf("expected current_value to be 42, got %d", counter.CurrentValue)
	}
	if counter.MaxValue != 0 {
		t.Errorf("expected max_value to be 0, got %d", counter.MaxValue)
	}
	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be true, got %v", counter.IsLocked)
	}
	if counter.UpdatedAt != "2024-05-30T12:34:56Z" {
		t.Errorf("expected updated_at to be '2024-05-30T12:34:56Z', got %s", counter.UpdatedAt)
	}
	if counter.ResetedAt.String != "2024-05-01T12:00:00Z" {
		t.Errorf("expected reseted_at to be '2024-05-01T12:00:00Z', got %s", counter.ResetedAt.String)
	}

	// Cleanup table after test
	cleanupTable(t, tableName)
}

func TestGetCounterEmpty(t *testing.T) {
	tableName := utils.TableInstance.Counter

	// Ensure table is empty
	rawDeleteQuery := `
		DELETE FROM %s;
	`
	deleteQuery := fmt.Sprintf(rawDeleteQuery, tableName)
	_, err := db.Exec(deleteQuery)
	if err != nil {
		t.Fatalf("failed to clean up table: %s", err)
	}

	counter, err := GetCounter(tableName)
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 0 {
		t.Errorf("expected current_value to be 0, got %d", counter.CurrentValue)
	}
	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be true, got %v", counter.IsLocked)
	}
	if counter.UpdatedAt != "" {
		t.Errorf("expected updated_at to be '', got %s", counter.UpdatedAt)
	}
	if counter.ResetedAt.Valid {
		t.Errorf("expected reseted_at to be invalid, got valid")
	}
}

// NOTE: Test covers a situation where there are no existing rows in the counter table and we
// simply need to create one with CurrentValue=1 and UpdatedAt=NOW(). ResetedAt should be a sql
// null string
func TestUpdateCounter(t *testing.T) {
	tableName := utils.TableInstance.Counter
	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Call the UpdateCounter function
	isUpdated := UpdateCounter()

	if !isUpdated {
		t.Errorf("expected isUpdated to be true, got %v", isUpdated)
	}

	// Test GetCounter function
	counter, err := GetCounter(tableName)
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

	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be true, got %v", counter.IsLocked)
	}

	// Cleanup table after test
	cleanupTable(t, tableName)
}

// NOTE: Test covers a typical situation where there are some existing rows in the counter table
// and we simply need to update the counter. It is expected to increment the counter by
// one and update the updated_at field to NOW().
func TestUpdateCounterTypicalCase(t *testing.T) {
	tableName := utils.TableInstance.Counter

	// Insert a row into the counter table
	rawInsertQuery := `
		INSERT INTO %s (current_value, max_value, is_locked, updated_at) 
		VALUES 
			(42, 42, false, '2024-05-30 12:34:56');`
	insertQuery := fmt.Sprintf(rawInsertQuery, tableName)

	_, err := db.Exec(insertQuery)
	if err != nil {
		t.Fatalf("failed to insert into table: %s, err: %s", tableName, err)
	}

	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Test UpdateCounter
	isUpdated := UpdateCounter()

	if !isUpdated {
		t.Errorf("expected isUpdated to be true, got %v", isUpdated)
	}

	counter, err := GetCounter(tableName)
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 43 {
		t.Errorf("expected current_value to be 43, got %d", counter.CurrentValue)
	}
	if counter.MaxValue != 43 {
		t.Errorf("expected max_value to be 43, got %d", counter.MaxValue)
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

	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be true, got %v", counter.IsLocked)
	}

	// Cleanup table after test
	cleanupTable(t, tableName)
}

// NOTE: Test covers a typical situation where there are some existing rows in the counter table
// and we simply need to update the counter. It is expected to increment the counter by
// one and update the updated_at field to NOW(). The max_value should not be updated because it is
// lower than the current_value.
func TestUpdateCounterTypicalCaseMaxValueNotReached(t *testing.T) {
	tableName := utils.TableInstance.Counter

	// Insert a row into the counter table
	rawInsertQuery := `
		INSERT INTO %s (current_value, max_value, is_locked, updated_at) 
		VALUES 
			(42, 100, false, '2024-05-30 12:34:56');`
	insertQuery := fmt.Sprintf(rawInsertQuery, tableName)

	_, err := db.Exec(insertQuery)
	if err != nil {
		t.Fatalf("failed to insert into table: %s, err: %s", tableName, err)
	}

	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Test UpdateCounter
	isUpdated := UpdateCounter()

	if !isUpdated {
		t.Errorf("expected isUpdated to be true, got %v", isUpdated)
	}

	counter, err := GetCounter(tableName)
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 43 {
		t.Errorf("expected current_value to be 43, got %d", counter.CurrentValue)
	}
	if counter.MaxValue != 100 {
		t.Errorf("expected max_value to be 100, got %d", counter.MaxValue)
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

	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be true, got %v", counter.IsLocked)
	}

	// Cleanup table after test
	cleanupTable(t, tableName)
}

// NOTE: Test covers a situation where there are some existing rows in the counter table and we
// cannot update the counter because 24h did not pass since the last update.
// It is expected that no counter is updated.
func TestUpdateCounterTimeDidNotPass(t *testing.T) {
	tableName := utils.TableInstance.Counter
	// Check updated_at (should be close to current time)
	updatedLessThan24hAgo := time.Now().UTC().Add(-23 * time.Hour)
	parsedUpdatedLessThan24hAgo := updatedLessThan24hAgo.Format(time.RFC3339)

	// Insert a row into the counter table
	rawInsertQuery := `
		INSERT INTO %s (current_value, is_locked, updated_at) 
		VALUES 
			(42, false, '%s');`
	insertQuery := fmt.Sprintf(rawInsertQuery, tableName, parsedUpdatedLessThan24hAgo)

	_, err := db.Exec(insertQuery)
	if err != nil {
		t.Fatalf("failed to insert into table: %s, err: %s", tableName, err)
	}

	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Test UpdateCounter
	isUpdated := UpdateCounter()

	if isUpdated {
		t.Errorf("expected isUpdated to be false, got %v", isUpdated)
	}

	counter, err := GetCounter(tableName)
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 42 {
		t.Errorf("expected current_value to be 42, got %d", counter.CurrentValue)
	}

	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be true, got %v", counter.IsLocked)
	}

	// Check updated_at (should be intact)
	if counter.UpdatedAt != parsedUpdatedLessThan24hAgo {
		t.Errorf("expected updated_at: '%s' to not change, got '%s'", parsedUpdatedLessThan24hAgo, counter.UpdatedAt)
	}

	if counter.ResetedAt.Valid {
		t.Errorf("expected reseted_at to be null, got %v", counter.ResetedAt)
	}

	// Cleanup table after test
	cleanupTable(t, tableName)
}

// NOTE: Resetting counter. Counter has no values - no entry exists, calling ResetCounter
// should create a new counter element
func TestResetCounterNoData(t *testing.T) {
	tableName := utils.TableInstance.Counter

	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Test ResetCounter
	lastValue, err := ResetCounter(tableName)
	if err != nil {
		t.Fatalf("failed to reset %s data: %s", tableName, err)
	}

	if lastValue != 0 {
		t.Errorf("expected last value to be 1, got %d", lastValue)
	}

	counter, err := GetCounter(tableName)
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 1 {
		t.Errorf("expected current_value to be 1, got %d", counter.CurrentValue)
	}

	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be false, got %v", counter.IsLocked)
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

	if !counter.ResetedAt.Valid {
		t.Errorf("expected reseted_at to be a valid sql null string, got %v", counter.ResetedAt)
	}

	// Cleanup table after test
	cleanupTable(t, tableName)
}

// NOTE: Counter has some value, resetting, should be zero now
func TestResetCounter(t *testing.T) {
	tableName := utils.TableInstance.Counter
	_, err := db.Exec(`
		INSERT INTO counter (current_value, is_locked, updated_at, reseted_at) 
		VALUES (42, false, '2024-05-30 12:34:56', '2024-05-01 12:00:00')
	`)
	if err != nil {
		t.Fatalf("failed to insert into table: %s", err)
	}

	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Test ResetCounter
	lastValue, err := ResetCounter(tableName)
	if err != nil {
		t.Fatalf("failed resetting %s data: %s", tableName, err)
	}

	if lastValue != 42 {
		t.Errorf("expected last value to be 42, got %d", lastValue)
	}

	counter, err := GetCounter(tableName)
	if err != nil {
		t.Fatalf("failed to get %s: %s", tableName, err)
	}

	if counter.CurrentValue != 1 {
		t.Errorf("expected current_value to be 1, got %d", counter.CurrentValue)
	}

	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be , got %v", counter.IsLocked)
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

	if !counter.ResetedAt.Valid {
		t.Errorf("expected reseted_at to be a valid sql null string, got %v", counter.ResetedAt)
	}

	// Cleanup table after test
	cleanupTable(t, tableName)
}

func TestGetOhnoCounter(t *testing.T) {
	// Insert a row into the ohno_counter table
	tableName := utils.TableInstance.OhnoCounter
	rawInsertQuery := `
		INSERT INTO %s (current_value, is_locked, updated_at, reseted_at) 
		VALUES 
			(42, false, '2024-05-30 12:34:56', '2024-05-01 12:00:00');`
	insertQuery := fmt.Sprintf(rawInsertQuery, tableName)

	_, err := db.Exec(insertQuery)
	if err != nil {
		t.Fatalf("failed to insert into table: %s, err: %s", tableName, err)
	}

	counter, err := GetCounter(tableName)
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 42 {
		t.Errorf("expected current_value to be 42, got %d", counter.CurrentValue)
	}
	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be true, got %v", counter.IsLocked)
	}
	if counter.UpdatedAt != "2024-05-30T12:34:56Z" {
		t.Errorf("expected updated_at to be '2024-05-30T12:34:56Z', got %s", counter.UpdatedAt)
	}
	if counter.ResetedAt.String != "2024-05-01T12:00:00Z" {
		t.Errorf("expected reseted_at to be '2024-05-01T12:00:00Z', got %s", counter.ResetedAt.String)
	}

	// Cleanup table after test
	cleanupTable(t, tableName)
}

func TestGetOhnoCounterEmpty(t *testing.T) {
	tableName := utils.TableInstance.OhnoCounter

	// Ensure table is empty
	rawDeleteQuery := `
		DELETE FROM %s;
	`
	deleteQuery := fmt.Sprintf(rawDeleteQuery, tableName)
	_, err := db.Exec(deleteQuery)
	if err != nil {
		t.Fatalf("failed to clean up table: %s", err)
	}

	counter, err := GetCounter(tableName)
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 0 {
		t.Errorf("expected current_value to be 0, got %d", counter.CurrentValue)
	}
	if counter.IsLocked != true {
		t.Errorf("expected isLocked to be true, got %v", counter.IsLocked)
	}
	if counter.UpdatedAt != "" {
		t.Errorf("expected updated_at to be '', got %s", counter.UpdatedAt)
	}
	if counter.ResetedAt.Valid {
		t.Errorf("expected reseted_at to be invalid, got valid")
	}
}

// NOTE: Test covers a situation where there are no existing rows in the ohno_counter table and we
// simply need to create one with CurrentValue=1 and UpdatedAt=NOW(). ResetedAt should be a sql
// null string
func TestUpdateOhnoCounter(t *testing.T) {
	tableName := utils.TableInstance.OhnoCounter
	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Call the UpdateOhnoCounter function
	isUpdated := UpdateOhnoCounter()

	if !isUpdated {
		t.Errorf("expected isUpdated to be true, got %v", isUpdated)
	}

	// Test GetCounter function
	counter, err := GetCounter(tableName)
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

	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be true, got %v", counter.IsLocked)
	}

	// Cleanup table after test
	cleanupTable(t, tableName)
}

// NOTE: Test covers a typical situation where there are some existing rows in the ohno_counter
// table and we simply need to update the counter. It is expected to increment the counter by
// one and update the updated_at field to NOW().
func TestUpdateOhnoCounterTypicalCase(t *testing.T) {
	tableName := utils.TableInstance.OhnoCounter

	// Insert a row into the counter table
	rawInsertQuery := `
		INSERT INTO %s (current_value, is_locked, updated_at) 
		VALUES 
			(42, false, '2024-05-30 12:34:56');`
	insertQuery := fmt.Sprintf(rawInsertQuery, tableName)

	_, err := db.Exec(insertQuery)

	if err != nil {
		t.Fatalf("failed to insert into table: %s, err: %s", tableName, err)
	}

	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Test UpdateOhnoCounter
	isUpdated := UpdateOhnoCounter()

	if !isUpdated {
		t.Errorf("expected isUpdated to be true, got %v", isUpdated)
	}

	counter, err := GetCounter(tableName)
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

	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be true, got %v", counter.IsLocked)
	}

	// Cleanup table after test
	cleanupTable(t, tableName)
}

// NOTE: Test covers a situation where there are some existing rows in the ohno_counter table and
// we cannot update the counter because 24h did not pass since the last update.
// It is expected that no counter is updated.
func TestUpdateOhnoCounterTimeDidNotPass(t *testing.T) {
	tableName := utils.TableInstance.OhnoCounter
	// Check updated_at (should be close to current time)
	updatedLessThan24hAgo := time.Now().UTC().Add(-23 * time.Hour)
	parsedUpdatedLessThan24hAgo := updatedLessThan24hAgo.Format(time.RFC3339)

	// Insert a row into the counter table
	rawInsertQuery := `
		INSERT INTO %s (current_value, is_locked, updated_at) 
		VALUES 
			(42, false, '%s');`
	insertQuery := fmt.Sprintf(rawInsertQuery, tableName, parsedUpdatedLessThan24hAgo)

	_, err := db.Exec(insertQuery)
	if err != nil {
		t.Fatalf("failed to insert into table: %s, err: %s", tableName, err)
	}

	// Ensure db_test is not nil
	if db == nil {
		t.Fatal("db is nil")
	}

	// Test UpdateOhnoCounter
	isUpdated := UpdateOhnoCounter()

	if isUpdated {
		t.Errorf("expected isUpdated to be false, got %v", isUpdated)
	}

	counter, err := GetCounter(tableName)
	if err != nil {
		t.Fatalf("failed to get counter: %s", err)
	}

	if counter.CurrentValue != 42 {
		t.Errorf("expected current_value to be 42, got %d", counter.CurrentValue)
	}

	if counter.IsLocked != false {
		t.Errorf("expected isLocked to be true, got %v", counter.IsLocked)
	}

	// Check updated_at (should be intact)
	if counter.UpdatedAt != parsedUpdatedLessThan24hAgo {
		t.Errorf("expected updated_at: '%s' to not change, got '%s'", parsedUpdatedLessThan24hAgo, counter.UpdatedAt)
	}

	if counter.ResetedAt.Valid {
		t.Errorf("expected reseted_at to be null, got %v", counter.ResetedAt)
	}

	// Cleanup table after test
	cleanupTable(t, tableName)
}

// NOTE: Test covers a situation when counter has some values and we create a historical counter
// by recording ohno event
func TestCreatHistoricalCounter(t *testing.T) {
	// Insert a row into the counter table
	counterTableName := utils.TableInstance.Counter
	historicalCounterTableName := utils.TableInstance.HistoricalCounter
	rawInsertQuery := `
		INSERT INTO %s (current_value, is_locked, updated_at, reseted_at) 
		VALUES 
			(42, false, '2024-05-30 12:34:56', '2024-05-01 12:00:00');`
	insertQuery := fmt.Sprintf(rawInsertQuery, counterTableName)

	_, err := db.Exec(insertQuery)
	if err != nil {
		t.Fatalf("failed to insert into table: %s, err: %s", counterTableName, err)
	}

	// Reset the counter
	lastValue, err := ResetCounter(counterTableName)
	if err != nil {
		t.Fatalf("failed to reset counter data: %s", err)
	}

	// Create historical counter entry
	CreateHistoricalCounter(historicalCounterTableName, lastValue)

	// Get historical counter
	historicalCounter, err := GetHistoricalCounters(historicalCounterTableName)
	if err != nil {
		t.Fatalf("failed to get %s: %s", historicalCounterTableName, err)
	}

	// Check updated_at and created_at should be close to current time
	expectedTime := time.Now().UTC()
	parsedUpdatedAtTime, err := time.Parse(time.RFC3339, historicalCounter[0].UpdatedAt)
	if err != nil {
		t.Fatalf("failed to parse updated_at: %s", err)
	}
	parsedCreatedAtTime, err := time.Parse(time.RFC3339, historicalCounter[0].CreatedAt)
	if err != nil {
		t.Fatalf("failed to parse created_at: %s", err)
	}

	// Assert historical counter was created and it is a list
	if len(historicalCounter) != 1 {
		t.Fatalf("expected 1 historical counter, got %d", len(historicalCounter))
	}
	if historicalCounter[0].Value != 42 {
		t.Fatalf("expected value to be 42, got %d", historicalCounter[0].Value)
	}
	// Allow for a small time difference (e.g., 5 seconds)
	if expectedTime.Sub(parsedUpdatedAtTime).Seconds() > 5 {
		t.Errorf("expected updated_at to be close to '%s', got '%s'", expectedTime, historicalCounter[0].UpdatedAt)
	}
	// Allow for a small time difference (e.g., 5 seconds)
	if expectedTime.Sub(parsedCreatedAtTime).Seconds() > 5 {
		t.Errorf("expected updated_at to be close to '%s', got '%s'", expectedTime, historicalCounter[0].UpdatedAt)
	}

	// Cleanup table after test
	cleanupTable(t, counterTableName)
	cleanupTable(t, historicalCounterTableName)
}

// NOTE: Test retrieves historical counter items
func TestGetHistoricalCouters(t *testing.T) {
	// Insert a row into the counter table
	historicalCounterTableName := utils.TableInstance.HistoricalCounter
	rawInsertQuery := `
		INSERT INTO %s (counter_id,  created_at, updated_at, value)
		VALUES ('%s', '%s', '%s', %d);
`

	entries := []HistoricalCounter{
		{"123e4567-e89b-12d3-a456-426614174000", "2024-05-30 12:34:56", "2024-07-01 12:00:00", 42},
		{"223e4567-e89b-12d3-a456-426614174001", "2024-05-31 13:34:56", "2024-07-02 13:00:00", 43},
		{"323e4567-e89b-12d3-a456-426614174002", "2024-06-01 14:34:56", "2024-07-03 14:00:00", 44},
	}

	for _, entry := range entries {
		insertQuery := fmt.Sprintf(rawInsertQuery, historicalCounterTableName, entry.CounterID, entry.CreatedAt, entry.UpdatedAt, entry.Value)
		_, err := db.Exec(insertQuery)
		if err != nil {
			t.Fatalf("failed to insert into table: %s, err: %s", historicalCounterTableName, err)
		}
	}

	// Get historical counters
	historicalCounters, err := GetHistoricalCounters(historicalCounterTableName)
	if err != nil {
		t.Fatalf("failed to get %s: %s", historicalCounterTableName, err)
	}

	// Assert historical_counter was created and it is a list
	if len(historicalCounters) != 3 {
		t.Fatalf("expected 3 historical counters, got %d", len(historicalCounters))
	}
	expectedValues := []int{42, 43, 44}
	for counter_idx, expectedValue := range expectedValues {
		if historicalCounters[counter_idx].Value != expectedValue {
			t.Fatalf("expected value of the counter %d to be %d, got %d", counter_idx+1, expectedValue, historicalCounters[counter_idx].Value)
		}
	}

	// Cleanup table after test
	cleanupTable(t, historicalCounterTableName)
}
