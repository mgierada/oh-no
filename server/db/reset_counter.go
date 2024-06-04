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

	query := (`
		SELECT 
			current_value, is_locked, updated_at, reseted_at 
		FROM 
			counter 
		LIMIT 1 FOR UPDATE
	`)

	err = tx.QueryRow(query).Scan(&counter.CurrentValue, &counter.IsLocked, &counter.UpdatedAt, &counter.ResetedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			insertQuery := (`
				INSERT INTO counter 
					(current_value, is_locked, updated_at, reseted_at) 
				VALUES 
					(0, false, NOW(), NOW())
			`)

			_, err = tx.Exec(insertQuery)

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
		updateQuery := (`
			UPDATE
				counter 
			SET 
				current_value = 0, updated_at = NOW(), reseted_at = NOW()
		`)
		_, err = tx.Exec(updateQuery)
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
