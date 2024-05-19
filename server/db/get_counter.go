package db

import (
	"fmt"
)

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
