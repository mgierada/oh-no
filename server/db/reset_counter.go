package db

import (
	"database/sql"
	"fmt"
	"log"
)

func ResetCounter(tableName string) (int, error) {
	var counter Counter
	var lastValue int

	tx, err := db.Begin()
	if err != nil {
		return -1, fmt.Errorf("❌ Error starting transaction.\n %s", err)
	}

	rawQuery := `
		SELECT 
			current_value, is_locked, updated_at, reseted_at 
		FROM 
		%s	
		LIMIT 1 FOR UPDATE;
	`

	query := fmt.Sprintf(rawQuery, tableName)
	err = tx.QueryRow(query).Scan(&counter.CurrentValue, &counter.IsLocked, &counter.UpdatedAt, &counter.ResetedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("❌ No %s, initializing one", tableName)

			rawInsertQuery := `
				INSERT INTO %s (current_value, is_locked, updated_at, reseted_at)
				VALUES (1, false, NOW(), NOW());
			`

			insertQuery := fmt.Sprintf(rawInsertQuery, tableName)
			_, err = tx.Exec(insertQuery)
			if err != nil {
				tx.Rollback()
				return -1, fmt.Errorf("❌ Error inserting new $s row.\n %s", tableName, err)
			}

		} else {
			tx.Rollback()
			return -1, fmt.Errorf("❌ Error querying %s table.\n %s", tableName, err)
		}

	} else {
		lastValue = counter.CurrentValue

		rawUpdateQuery := (`
			UPDATE
				%s	
			SET 
				current_value = 1, updated_at = NOW(), reseted_at = NOW()
			`)
		updateQuery := fmt.Sprintf(rawUpdateQuery, tableName)

		_, err = tx.Exec(updateQuery)
		if err != nil {
			tx.Rollback()
			return -1, fmt.Errorf("❌ Error updating %s row.\n %s", tableName, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return -1, fmt.Errorf("❌ Error committing transaction.\n %s", err)
	}

	return lastValue, err
}
