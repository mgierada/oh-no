package db

import (
	"database/sql"
	"fmt"
	"log"
	"server/utils"
	"strings"
	"time"
)

type Counter struct {
	CurrentValue int
	MaxValue     int
	UpdatedAt    string
	ResetedAt    sql.NullString
	IsLocked     bool
}

type UpdateCounterType struct {
	CurrentValue *int
	UpdatedAt    *string
	ResetedAt    *sql.NullString
	IsLocked     *bool
}

func upsertCounterData(tableName string) (bool, error) {
	if tableName == "" {
		return false, fmt.Errorf("❌ Error upserting counter data. Table name cannot be empty.")
	}

	tx, err := db.Begin()
	if err != nil {
		return false, fmt.Errorf("❌ Error starting transaction.\n %s", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("❌ Error committing transaction.\n %s", err)
			}
		}
	}()

	var counter Counter

	upsertCounterQuery := fmt.Sprintf(`
		SELECT 
			current_value, is_locked, updated_at, reseted_at 
		FROM 
			%s
		LIMIT 1 
		FOR UPDATE
	`, tableName)

	err = tx.QueryRow(upsertCounterQuery).Scan(&counter.CurrentValue, &counter.IsLocked, &counter.UpdatedAt, &counter.ResetedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found in counter table. Inserting new row.")

			insertCounterQuery := fmt.Sprintf(`
				INSERT INTO %s (current_value, updated_at)
				VALUES (1, NOW())
			`, tableName)

			_, err = tx.Exec(insertCounterQuery)

			if err != nil {
				return false, fmt.Errorf("❌ Error inserting new counter row.\n %s", err)
			}
			return true, nil

		} else {
			return false, fmt.Errorf("❌ Error querying %s table.\n %s", tableName, err)
		}

	} else {
		lastUpdated, err := time.Parse(time.RFC3339Nano, counter.UpdatedAt)
		if err != nil {
			return false, fmt.Errorf("❌ Error parsing updated_at timestamp.\n %s", err)
		}

		if counter.ResetedAt.Valid {
			lastReseted, err := time.Parse(time.RFC3339Nano, counter.ResetedAt.String)
			if err != nil {
				return false, fmt.Errorf("❌ Error parsing reseted_at timestamp.\n %s", err)
			}

			if lastReseted.After(lastUpdated) || lastReseted.Equal(lastUpdated) {
				log.Println("Counter was reseted. lastReseted <= lastUpdated")
				updateQuery := fmt.Sprintf(`
					UPDATE %s 
					SET 
						current_value = 1, updated_at = NOW()
				`, tableName)
				_, err = tx.Exec(updateQuery)
				if err != nil {
					return false, fmt.Errorf("❌ Error updating counter row.\n %s", err)
				}
			}
		}

		updateIntervalInt, err := utils.GetEnvInt("UPDATE_INTERVAL_IN_HOURS")
		if err != nil {
			return false, fmt.Errorf("❌ Error getting UPDATE_INTERVAL_IN_HOURS environment variable.\n %s", err)
		}

		updateInterval := time.Duration(updateIntervalInt)

		// if time.Since(lastUpdated) < updateInterval*time.Hour {
		if time.Since(lastUpdated) < updateInterval*time.Second {
			log.Printf("🙅 %d hours have not passed since the last update. Counter not increased...", updateIntervalInt)
			return false, nil
		}

		updateQuery := fmt.Sprintf(`
			UPDATE %s 
			SET current_value = current_value + 1, updated_at = NOW()
		`, tableName)

		_, err = tx.Exec(updateQuery)

		if err != nil {
			return false, fmt.Errorf("❌ Error updating counter row.\n %s", err)
		}
	}

	return true, err
}

// TODO: Refactor this function to improve error handling and readability
func SetCounter(value int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("❌ Error starting transaction.\n %s", err)
	}

	var counter Counter

	err = tx.QueryRow("SELECT current_value, updated_at, reseted_at FROM counter LIMIT 1 FOR UPDATE").Scan(&counter.CurrentValue, &counter.UpdatedAt, &counter.ResetedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec("INSERT INTO counter (current_value, updated_at) VALUES (?, NOW())", value)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("❌ Error inserting new counter row.\n %s", err)
			}

		} else {
			tx.Rollback()
			return fmt.Errorf("❌ Error querying counter table.\n %s", err)
		}

	} else {
		_, err = tx.Exec("UPDATE counter SET current_value = $1, updated_at = NOW()", value)
		if err != nil {
			log.Printf("Error updating counter.\n %s", err)
			tx.Rollback()
			// BUG Why this does not throw error when update counter fails?
			return fmt.Errorf("❌ Error updating counter row.\n %s", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("❌ Error committing transaction.\n %s", err)
	}
	log.Println("✅ Transaction committed successfully")
	return nil
}

func UpdateCounter() bool {
	isUpdated, err := upsertCounterData(utils.TableInstance.Counter)

	if err != nil {
		log.Printf("❌ Error updating counter.\n %s", err)
	}

	if !isUpdated {
		log.Printf("❌ Counter not incremented. Conditions not met.")
	}

	return isUpdated
}

func UpdateOhnoCounter() bool {
	isUpdated, err := upsertCounterData(utils.TableInstance.OhnoCounter)

	if err != nil {
		log.Printf("❌ Error updating counter.\n %s", err)
	}

	if !isUpdated {
		log.Printf("❌ Counter not incremented. Conditions not met.")
	}

	return isUpdated
}

func updateCounter(tableName string, properties UpdateCounterType) (bool, error) {
	if tableName == "" {
		return false, fmt.Errorf("❌ Error updating counter data. Table name cannot be empty.")
	}

	tx, err := db.Begin()
	if err != nil {
		return false, fmt.Errorf("❌ Error starting transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("❌ Error committing transaction: %s", err)
			}
		}
	}()

	upsertCounterQuery := fmt.Sprintf(`
		SELECT 
			current_value, is_locked, updated_at, reseted_at 
		FROM 
			%s
		LIMIT 1 
		FOR UPDATE
	`, tableName)

	var counter Counter

	err = tx.QueryRow(upsertCounterQuery).Scan(&counter.CurrentValue, &counter.IsLocked, &counter.UpdatedAt, &counter.ResetedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No rows found in %s table. Inserting new row.", tableName)

			insertCounterQuery := fmt.Sprintf(`
				INSERT INTO %s (current_value, updated_at, is_locked)
				VALUES (1, NOW(), false)
			`, tableName)

			_, err = tx.Exec(insertCounterQuery)

			if err != nil {
				return false, fmt.Errorf("❌ Error inserting new counter row.\n %s", err)
			}
			return true, nil

		} else {
			return false, fmt.Errorf("❌ Error querying %s table.\n %s", tableName, err)
		}
	}

	var setClauses []string
	var args []interface{}
	argIndex := 1

	if properties.CurrentValue != nil {
		setClauses = append(setClauses, fmt.Sprintf("current_value = $%d", argIndex))
		args = append(args, *properties.CurrentValue)
		argIndex++
	}
	if properties.UpdatedAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
		args = append(args, *properties.UpdatedAt)
		argIndex++
	}
	if properties.ResetedAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("reseted_at = $%d", argIndex))
		args = append(args, *properties.ResetedAt)
		argIndex++
	}
	if properties.IsLocked != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_locked = $%d", argIndex))
		args = append(args, *properties.IsLocked)
		argIndex++
	}

	if len(setClauses) == 0 {
		return false, fmt.Errorf("❌ No properties provided for update")
	}

	query := fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(setClauses, ", "))

	_, err = tx.Exec(query, args...)
	if err != nil {
		log.Printf("❌ Error updating table: %s: %s", tableName, err)
		return false, err
	}

	return true, nil
}

func LockCounter(tableName string) (bool, error) {
	isLocked := true
	return updateCounter(tableName, UpdateCounterType{IsLocked: &isLocked})
}

func UnlockCounter(tableName string) (bool, error) {
	isLocked := false
	return updateCounter(tableName, UpdateCounterType{IsLocked: &isLocked})
}
