package db

import (
	"fmt"
	"github.com/google/uuid"
	"log"
)

func CreateHistoricalCounter(lastValue int) error {
	if lastValue <= 0 {
		message := fmt.Sprintf("❌ Error creating new historical counter. Value must be greater than 0. Received: %d", lastValue)
		log.Printf(message)
		return nil
	}
	newCounterId := uuid.New().String()
	_, err := db.Exec("INSERT INTO historical_counters (counter_id, value) VALUES (?,?)", newCounterId, lastValue)
	if err != nil {
		message := fmt.Sprintf("❌ Error inserting new historical counter row.\n %s", err)
		log.Fatalf(message)
		return fmt.Errorf(message)
	}
	return nil
}
