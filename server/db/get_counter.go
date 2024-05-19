package db

import (
	"database/sql"
	"fmt"
)

func GetCounter() (Counter, error) {
	var counter Counter
	query := "SELECT current_value, updated_at, reseted_at FROM counter LIMIT 1"
	row := db.QueryRow(query)
	err := row.Scan(&counter.CurrentValue, &counter.UpdatedAt, &counter.ResetedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			emptyCounter := Counter{
				CurrentValue: 0,
				UpdatedAt:    "",
				ResetedAt:    sql.NullString{},
			}
			return emptyCounter, nil
		}
		return Counter{}, fmt.Errorf("❌ Error querying counter table.\n %s", err)
	}
	return counter, nil
}
