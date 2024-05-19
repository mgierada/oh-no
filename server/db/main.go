package db

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *sql.DB

func Connect() {
	e := godotenv.Overload("../.env")
	if e != nil {
		log.Fatalf("❌ Error loading .env file.\n %s", e)
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbProtocol := os.Getenv("DB_PROTOCOL")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbAddress := fmt.Sprintf("%s:%s", dbHost, dbPort)

	if dbUser == "" || dbPassword == "" || dbName == "" || dbProtocol == "" || dbPort == "" || dbHost == "" {
		log.Fatal("❌ One or more environment variables are missing")
	}

	cfg := mysql.Config{
		User:   dbUser,
		Passwd: dbPassword,
		Net:    dbProtocol,
		Addr:   dbAddress,
		DBName: dbName,
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalf("❌ Error getting a database handle.\n, %s", err)
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
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS counter (
			current_value INT NOT NULL,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			reseted_at TIMESTAMP NULL DEFAULT NULL
		);
	`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("❌ Error creating counter table.\n %s", err)
	}
	log.Println("✅ Ensured counter table exists.")
}

func createHistoricalCountersTableIfNotExists() {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS historical_counters (
			counter_id CHAR(36) PRIMARY KEY NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			value INT NOT NULL
		);
	`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("❌ Error creating histiorical_counter table.\n %s", err)
	}
	log.Println("✅ Ensured histiorical_counter table exists.")
}
