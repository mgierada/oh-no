package models

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateCounterTableIfNotExists(db *sql.DB, tableName string) error {
	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			current_value INT NOT NULL,
			max_value INT NULL DEFAULT 0,
			is_locked BOOLEAN NOT NULL DEFAULT FALSE,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			reseted_at TIMESTAMP NULL DEFAULT NULL
		);`, tableName)
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("❌ Error creating %s table.\n %s", tableName, err)
	}
	log.Printf("✅ Ensured table %s exist.", tableName)
	return nil
}

func CreateOrReplaceTrigger(db *sql.DB, tableName string) error {
	var err error
	triggerFunctionName := fmt.Sprintf("%s_update", tableName)
	triggerName := fmt.Sprintf("%s_update", tableName)

	createTriggerFunctionQuery := fmt.Sprintf(`
		CREATE OR REPLACE FUNCTION %s()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = now();
			IF NEW.current_value > NEW.max_value THEN
				NEW.max_value = NEW.current_value;
			END IF;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`, triggerFunctionName)
	_, err = db.Exec(createTriggerFunctionQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger function for %s table.\n %s", tableName, err)
	}

	createTriggerQuery := fmt.Sprintf(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 
				FROM pg_trigger 
				WHERE tgname = '%s'
			) THEN
				CREATE TRIGGER %s
				BEFORE UPDATE ON %s
				FOR EACH ROW
				EXECUTE FUNCTION %s();
			END IF;
		END $$;
	`, triggerName, triggerName, tableName, triggerFunctionName)
	_, err = db.Exec(createTriggerQuery)
	if err != nil {
		log.Fatalf("❌ Error creating trigger %s for %s table.\n %s", triggerName, tableName, err)
	}
	log.Printf("✅ Ensured trigger for table %s exist.", tableName)
	return nil
}
