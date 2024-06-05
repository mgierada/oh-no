package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"server/db/models"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
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

	models.CreateCounterTableIfNotExists(db)
	models.CreateHistoricalCountersTableIfNotExists(db)
	models.CreateOhnoCounterTableIfNotExists(db)
	models.CreateHistoricalOhnoCountersTableIfNotExists(db)
}
