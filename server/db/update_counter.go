package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Counter struct {
	CurrentValue int
	UpdatedAt    string
	ResetedAt    sql.NullString
}

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

		lastUpdated, err := time.Parse(time.RFC3339Nano, counter.UpdatedAt)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("❌ Error parsing updated_at timestamp.\n %s", err)
		}

		lastReseted, err := time.Parse(time.RFC3339Nano, counter.ResetedAt.String)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("❌ Error parsing updated_at timestamp.\n %s", err)
		}

		if lastReseted.Before(lastUpdated) || lastReseted.Equal(lastUpdated) {
			_, err = tx.Exec("UPDATE counter SET current_value = 1, updated_at = NOW()")
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("❌ Error updating counter row.\n %s", err)
			}
		}

		if time.Since(lastUpdated) >= 24*time.Hour {
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
