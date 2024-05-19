package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
)

func GetCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received GET /counter request\n")

	counters, err := db.GetAllCounterData()
	if err != nil {
		log.Fatalf("❌ Error retrieving counter data.\n %s", err)
	}

	for _, counter := range counters {
		log.Printf("Current Value: %d, Updated At: %s, Reseted At: %v\n",
			counter.CurrentValue, counter.UpdatedAt, counter.ResetedAt.String)
	}

	// Convert the counters slice to JSON
	jsonData, err := json.Marshal(counters)
	if err != nil {
		log.Fatalf("❌ Error marshaling counter data to JSON.\n %s", err)
		http.Error(w, "Error marshaling counter data to JSON", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func UpdateCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received POST /increment request")
	db.ManualIncrement()
}

func AutoUpdateCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received POST /autoincrement request")
	for {
		err := db.UpsertCounterData()
		if err != nil {
			log.Fatalf("❌ Error incrementing counter.\n %s", err)
		}
	}
}
