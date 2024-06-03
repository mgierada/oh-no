package models

import (
	"database/sql"
	"log"
)

func CreateCounterTableIfNotExists(db *sql.DB) {
	// Create the counter table
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS counter (
			current_value INT NOT NULL,
			is_locked BOOLEAN NOT NULL DEFAULT FALSE,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			reseted_at TIMESTAMP NULL DEFAULT NULL
		);
	`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("❌ Error creating counter table.\n %s", err)
	}

	// Create or replace the trigger function
	createTriggerFunctionQuery := `
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = now();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`
	_, err = db.Exec(createTriggerFunctionQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger function for counter table.\n %s", err)
	}

	// Create the trigger conditionally
	createTriggerQuery := `
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 
				FROM pg_trigger 
				WHERE tgname = 'update_updated_at'
			) THEN
				CREATE TRIGGER update_updated_at
				BEFORE UPDATE ON counter
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();
			END IF;
		END $$;
	`
	_, err = db.Exec(createTriggerQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger for counter table.\n %s", err)
	}

	log.Println("✅ Ensured counter table and trigger exist.")
}
