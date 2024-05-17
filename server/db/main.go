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
		// Addr:   "127.0.0.1:3306",
		Addr:   dbAddress,
		DBName: dbName,
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Printf("✅ Connected to database %s on %s", dbName, dbAddress)
}
