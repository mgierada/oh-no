package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Counter represents a row in the counter table
type Counter struct {
	CurrentValue int
	UpdatedAt    string
	ResetedAt    sql.NullString
}

var cancelFunc context.CancelFunc
var taskRunning bool

// UpsertCounterData upserts the counter data, increasing current_value by one
func UpsertCounterData() error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("❌ Error starting transaction.\n %s", err)
	}

	var counter Counter
	err = tx.QueryRow("SELECT current_value, updated_at, reseted_at FROM counter LIMIT 1 FOR UPDATE").Scan(&counter.CurrentValue, &counter.UpdatedAt, &counter.ResetedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec("INSERT INTO counter (current_value, updated_at) VALUES (1, NOW())")
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("❌ Error inserting new counter row.\n %s", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("❌ Error querying counter table.\n %s", err)
		}
	} else {
		lastUpdated, err := time.Parse("2006-01-02 15:04:05", counter.UpdatedAt)
		// if time.Since(lastUpdated) >= 24*time.Hour {
		if time.Since(lastUpdated) >= time.Second {
			_, err = tx.Exec("UPDATE counter SET current_value = current_value + 1, updated_at = NOW()")
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("❌ Error updating counter row.\n %s", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("❌ Error committing transaction.\n %s", err)
	}

	return nil
}

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func getEnvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}

func getEnvInt(key string) (int, error) {
	s, err := getEnvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func runBackgroundTask(ctx context.Context) {
	incrementFrequencyInHours, ok := getEnvInt("COUNTER_INCREMENT_FREQUENCY_IN_HOURS")
	if ok != nil {
		log.Println("❌ Error getting COUNTER_INCREMENT_FREQUENCY_IN_HOURS")
		return
	}
	ticker := time.NewTicker(time.Duration(incrementFrequencyInHours) * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := UpsertCounterData()
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			log.Println("🛑 Background task stopped")
			return
		}
	}
}

func RunBackgroundTask() {
	if taskRunning {
		log.Println("⚠️ Background task is already running")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancelFunc = cancel
	taskRunning = true
	go func() {
		runBackgroundTask(ctx)
		taskRunning = false
	}()
}

func StopBackgroundTask() {
	if cancelFunc != nil {
		cancelFunc()
	}
}
