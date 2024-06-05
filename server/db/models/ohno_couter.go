package models

import (
	"database/sql"
	"fmt"
	"log"
	"server/utils"
)

func CreateOhnoCounterTableIfNotExists(db *sql.DB) {
	// Create the counter table for ohno period
	tableName := utils.TableInstance.OhnoCounter

	rawCreateTableQuery := `
		CREATE TABLE IF NOT EXISTS %s (
			current_value INT NOT NULL,
			is_locked BOOLEAN NOT NULL DEFAULT TRUE,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			reseted_at TIMESTAMP NULL DEFAULT NULL
		);
	`
	createTableQuery := fmt.Sprintf(rawCreateTableQuery, tableName)
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("❌ Error creating %s table.\n %s", tableName, err)
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
		log.Fatalf("❌ Error creating trigger function for %s table.\n %s", tableName, err)
	}

	// Create the trigger conditionally
	rawCreateTriggerQuery := `
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 
				FROM pg_trigger 
				WHERE tgname = 'update_ohno_counter_updated_at_column'
			) THEN
				CREATE TRIGGER update_ohno_counter_updated_at_column
				BEFORE UPDATE ON %s
				-- FOR EACH ROW
				EXECUTE FUNCTION update_ohno_counter_updated_at_column();
			END IF;
		END $$;
	`
	createTriggerQuery := fmt.Sprintf(rawCreateTriggerQuery, tableName)
	_, err = db.Exec(createTriggerQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger for %s counter table.\n %s", tableName, err)
	}

	log.Printf("✅ Ensured %s table and trigger exist.", tableName)
}
