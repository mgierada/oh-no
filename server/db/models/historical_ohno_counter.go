package models

import (
	"database/sql"
	"fmt"
	"log"
	"server/utils"
)

func CreateHistoricalOhnoCountersTableIfNotExists(db *sql.DB) {
	// Create the historical_ohno_counters table
	tableName := utils.TableInstance.HistoricalOhnoCounter
	rawCreateTableQuery := `
		CREATE TABLE IF NOT EXISTS %s (
			counter_id UUID PRIMARY KEY NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			value INT NOT NULL
		);
	`
	createTableQuery := fmt.Sprintf(rawCreateTableQuery, tableName)
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("❌ Error creating %s table.\n %s", tableName, err)
	}

	// Create or replace the trigger function
	createTriggerFunctionQuery := `
		CREATE OR REPLACE FUNCTION update_historical_ohno_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = now();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`
	_, err = db.Exec(createTriggerFunctionQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger function for historical_ohno_counters table.\n %s", err)
	}

	// Create the trigger conditionally
	rawCreateTriggerQuery := `
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 
				FROM pg_trigger 
				WHERE tgname = 'update_historical_ohno_updated_at'
			) THEN
				CREATE TRIGGER update_historical_ohno_updated_at
				BEFORE UPDATE ON %s
				FOR EACH ROW
				EXECUTE FUNCTION update_historical_ohno_updated_at_column();
			END IF;
		END $$;
	`
	createTriggerQuery := fmt.Sprintf(rawCreateTriggerQuery, tableName)
	_, err = db.Exec(createTriggerQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger for %s table.\n %s", tableName, err)
	}

	log.Printf("✅ Ensured %s table and trigger exist.", tableName)
}
