package db

import (
	"database/sql"
	"fmt"
)

func getCounter(tableName string) (Counter, error) {
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
		if tableName == "ohno_counter" {
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

func GetCounter() (Counter, error) {
	return getCounter("counter")
}

func GetOhnoCounter() (Counter, error) {
	return getCounter("ohno_counter")
}

func GetCounterLocked() bool {
	counter, err := GetCounter()
	if err != nil {
		return false
	}
	return counter.IsLocked
}

func GetOhnoCounterLocked() bool {
	counter, err := GetOhnoCounter()
	if err != nil {
		return false
	}
	return counter.IsLocked
}
