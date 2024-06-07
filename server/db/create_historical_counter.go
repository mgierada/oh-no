package db

import (
	"fmt"
	"github.com/google/uuid"
	"log"
)

func CreateHistoricalCounter(tableName string, lastValue int) error {
	if lastValue <= 0 {
		message := fmt.Sprintf("❌ Error creating new historical counter. Value must be greater than 0. Received: %d", lastValue)
		log.Printf(message)
		return nil
	}
	newCounterId := uuid.New().String()
	rawInsertQuery := `
		INSERT INTO %s (counter_id, value)
		VALUES ('%s', %d);
	`
	insertQuery := fmt.Sprintf(rawInsertQuery, tableName, newCounterId, lastValue)
	_, err := db.Exec(insertQuery)
	if err != nil {
		message := fmt.Sprintf("❌ Error inserting new %s row.\n %s", tableName, err)
		log.Fatalf(message)
		return fmt.Errorf(message)
	}
	return nil
}
