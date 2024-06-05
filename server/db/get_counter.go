package db

import (
	"database/sql"
	"fmt"
	"server/utils"
)

func GetCounter(tableName string) (Counter, error) {
	var counter Counter
	query := fmt.Sprintf(`
		SELECT
			current_value, is_locked, updated_at, reseted_at 
		FROM %s 
		LIMIT 1
		`, tableName)
	row := db.QueryRow(query)
	err := row.Scan(&counter.CurrentValue, &counter.IsLocked, &counter.UpdatedAt, &counter.ResetedAt)

	if err != nil {
		defaultIsLocked := false
		if tableName == utils.TableInstance.OhnoCounter {
			defaultIsLocked = true
		}
		if err == sql.ErrNoRows {
			emptyCounter := Counter{
				CurrentValue: 0,
				IsLocked:     defaultIsLocked,
				UpdatedAt:    "",
				ResetedAt:    sql.NullString{},
			}
			return emptyCounter, nil
		}
		return Counter{}, fmt.Errorf("❌ Error querying %s table.\n %s", tableName, err)
	}
	return counter, nil
}
