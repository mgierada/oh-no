package models

import (
	"database/sql"
	"log"
)

func CreateHistoricalCountersTableIfNotExists(db *sql.DB) {
	// Create the historical_counters table
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS historical_counters (
			counter_id UUID PRIMARY KEY NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			value INT NOT NULL
		);
	`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("❌ Error creating historical_counters table.\n %s", err)
	}

	// Create or replace the trigger function
	createTriggerFunctionQuery := `
		CREATE OR REPLACE FUNCTION update_historical_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = now();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`
	_, err = db.Exec(createTriggerFunctionQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger function for historical_counters table.\n %s", err)
	}

	// Create the trigger conditionally
	createTriggerQuery := `
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 
				FROM pg_trigger 
				WHERE tgname = 'update_historical_updated_at'
			) THEN
				CREATE TRIGGER update_historical_updated_at
				BEFORE UPDATE ON historical_counters
				FOR EACH ROW
				EXECUTE FUNCTION update_historical_updated_at_column();
			END IF;
		END $$;
	`
	_, err = db.Exec(createTriggerQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger for historical_counters table.\n %s", err)
	}

	log.Println("✅ Ensured historical_counters table and trigger exist.")
}
