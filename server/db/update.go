package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Counter represents a row in the counter table
type Counter struct {
	CurrentValue int
	UpdatedAt    string
	ResetedAt    sql.NullString
}

// UpsertCounterData upserts the counter data, increasing current_value by one
func UpsertCounterData() error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("❌ Error starting transaction.\n %s", err)
	}

	var counter Counter
	err = tx.QueryRow("SELECT current_value, updated_at, reseted_at FROM counter LIMIT 1 FOR UPDATE").Scan(&counter.CurrentValue, &counter.UpdatedAt, &counter.ResetedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec("INSERT INTO counter (current_value, updated_at) VALUES (1, NOW())")
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("❌ Error inserting new counter row.\n %s", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("❌ Error querying counter table.\n %s", err)
		}
	} else {
		lastUpdated, err := time.Parse("2006-01-02 15:04:05", counter.UpdatedAt)
		// if time.Since(lastUpdated) >= 24*time.Hour {
		if time.Since(lastUpdated) >= time.Second {
			_, err = tx.Exec("UPDATE counter SET current_value = current_value + 1, updated_at = NOW()")
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("❌ Error updating counter row.\n %s", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("❌ Error committing transaction.\n %s", err)
	}

	return nil
}

func runBackgroundTask() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := UpsertCounterData()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func RunBackgroundTask() {
	go runBackgroundTask()
}
