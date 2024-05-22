package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *sql.DB

func Connect() {
	e := godotenv.Overload("../.env")
	if e != nil {
		log.Printf("❌ Error loading .env file.\n %s", e)
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbAddress := fmt.Sprintf("%s:%s", dbHost, dbPort)
	dbDriver := os.Getenv("DB_DRIVER")

	dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", dbUser, dbPassword, dbAddress, dbName)

	if dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" || dbHost == "" || dbDriver == "" {
		log.Fatalf("❌ One or more environment variables are missing")
	}

	// Get a database handle.
	var err error
	db, err = sql.Open(dbDriver, dbConnectionString)
	if err != nil {
		log.Printf("❌ Error getting a database handle.\n, %s", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatalf("❌ Error connecting to database.\n %s", pingErr)
	}
	log.Printf("✅ Connected to database %s on %s", dbName, dbAddress)

	createCounterTableIfNotExists()
	createHistoricalCountersTableIfNotExists()
}

func createCounterTableIfNotExists() {
	// Create the counter table
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS counter (
			current_value INT NOT NULL,
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

func createHistoricalCountersTableIfNotExists() {
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
