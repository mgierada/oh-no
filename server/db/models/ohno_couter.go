package models

import (
	"database/sql"
	"log"
)

func CreateOhnoCounterTableIfNotExists(db *sql.DB) {
	// Create the counter table for ohno period
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS ohno_counter (
			current_value INT NOT NULL,
			is_locked BOOLEAN NOT NULL DEFAULT TRUE,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			reseted_at TIMESTAMP NULL DEFAULT NULL
		);
	`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("❌ Error creating ohno counter table.\n %s", err)
	}

	// Create or replace the trigger function
	createTriggerFunctionQuery := `
		CREATE OR REPLACE FUNCTION update_ohno_counter_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = now();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`
	_, err = db.Exec(createTriggerFunctionQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger function for ohno_counter table.\n %s", err)
	}

	// Create the trigger conditionally
	createTriggerQuery := `
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 
				FROM pg_trigger 
				WHERE tgname = 'update_ohno_counter_updated_at_column'
			) THEN
				CREATE TRIGGER update_ohno_counter_updated_at_column
				BEFORE UPDATE ON ohno_counter
				-- FOR EACH ROW
				EXECUTE FUNCTION update_ohno_counter_updated_at_column();
			END IF;
		END $$;
	`
	_, err = db.Exec(createTriggerQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger for ohno counter table.\n %s", err)
	}

	log.Println("✅ Ensured ohno counter table and trigger exist.")
}
