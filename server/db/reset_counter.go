package db

import (
	"database/sql"
	"fmt"
)

func ResetCounter() (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return -1, fmt.Errorf("❌ Error starting transaction.\n %s", err)
	}

	var counter Counter
	var lastValue int

	err = tx.QueryRow("SELECT current_value, updated_at, reseted_at FROM counter LIMIT 1 FOR UPDATE").Scan(&counter.CurrentValue, &counter.UpdatedAt, &counter.ResetedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec("INSERT INTO counter (current_value, updated_at, reseted_at) VALUES (0, NOW(), NOW())")
			if err != nil {
				tx.Rollback()
				return -1, fmt.Errorf("❌ Error inserting new counter row.\n %s", err)
			}
		} else {
			tx.Rollback()
			return -1, fmt.Errorf("❌ Error querying counter table.\n %s", err)
		}
	} else {
		lastValue = counter.CurrentValue
		_, err = tx.Exec("UPDATE counter SET current_value = 0, updated_at = NOW(), reseted_at = NOW()")
		if err != nil {
			tx.Rollback()
			return -1, fmt.Errorf("❌ Error updating counter row.\n %s", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return -1, fmt.Errorf("❌ Error committing transaction.\n %s", err)
	}

	return lastValue, err
}
