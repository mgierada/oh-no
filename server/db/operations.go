package db

import (
	"database/sql"
	"fmt"
)

// ManualIncrement increments the counter value manually
func ManualIncrement() error {
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
		_, err = tx.Exec("UPDATE counter SET current_value = current_value + 1, updated_at = NOW()")
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("❌ Error updating counter row.\n %s", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("❌ Error committing transaction.\n %s", err)
	}

	return nil
}

// GetAllCounterData retrieves all rows from the counter table
func GetAllCounterData() ([]Counter, error) {
	query := "SELECT current_value, updated_at, reseted_at FROM counter"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("❌ Error querying counter table.\n %s", err)
	}
	defer rows.Close()

	var counters []Counter
	for rows.Next() {
		var counter Counter
		err := rows.Scan(&counter.CurrentValue, &counter.UpdatedAt, &counter.ResetedAt)
		if err != nil {
			return nil, fmt.Errorf("❌ Error scanning row.\n %s", err)
		}
		counters = append(counters, counter)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("❌ Row iteration error.\n %s", err)
	}

	return counters, nil
}
