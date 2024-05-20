package db

import (
	"fmt"
)

type HistoricalCounter struct {
	CounterID string
	CreatedAt string
	UpdatedAt string
	Value     int
}

func GetHistoricalCounters() ([]HistoricalCounter, error) {

	query := "SELECT counter_id, updated_at, created_at, value FROM historical_counters"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("❌ Error querying historical_counters table.\n %s", err)
	}
	defer rows.Close()

	var historicalCounters []HistoricalCounter

	for rows.Next() {
		var historicalCounter HistoricalCounter
		err := rows.Scan(&historicalCounter.CounterID, &historicalCounter.CreatedAt, &historicalCounter.UpdatedAt, &historicalCounter.Value)
		if err != nil {
			return nil, fmt.Errorf("❌ Error scanning row.\n %s", err)
		}
		historicalCounters = append(historicalCounters, historicalCounter)

	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("❌ Row iteration error.\n %s", err)
	}

	return historicalCounters, nil
}
