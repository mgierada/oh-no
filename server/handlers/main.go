package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
)

func GetCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received GET /counter request\n")

	counter, err := db.GetCounter()
	if err != nil {
		log.Fatalf("❌ Error retrieving counter data.\n %s", err)
	}

	log.Printf("Current Value: %d, Updated At: %s, Reseted At: %v\n",
		counter.CurrentValue, counter.UpdatedAt, counter.ResetedAt.String)

	// Convert the counters slice to JSON
	jsonData, err := json.Marshal(counter)
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

func StartAutoUpdateCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received POST /start-incr request")
	db.RunBackgroundTask()
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Background task stared"}
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("❌ Error marshaling counter data to JSON.\n %s", err)
		http.Error(w, "Error marshaling counter data to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
	log.Println("🟢 Background task started")
}

func StopAutoUpdateCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received POST /stop_incr request")
	db.StopBackgroundTask()
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Background task stopped"}
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("❌ Error marshaling counter data to JSON.\n %s", err)
		http.Error(w, "Error marshaling counter data to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
	log.Println("🔴 Background task stopped")
}

func RecordOhNoEvent(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received POST /ohno request")
	last_value, err := db.ResetCounter()
	if err != nil {
		log.Fatalf("❌ Error resetting counter.\n %s", err)
		http.Error(w, "Error resetting counter.", http.StatusInternalServerError)
		return
	}
	db.CreateHistoricalCounter(last_value)
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Oh No! Event recorded"}
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("❌ Error marshaling counter data to JSON.\n %s", err)
		http.Error(w, "Error marshaling counter data to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
	log.Println("🟢 Oh No! Event recorded")
}

func GetHistoricalCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received GET /historical request\n")

	hCounters, err := db.GetHistoricalCounters()
	if err != nil {
		log.Fatalf("❌ Error retrieving historical_counter data.\n %s", err)
	}

	for _, hCounter := range hCounters {
		log.Printf("CounterId: %s, Updated At: %s, Created_At: %v Value: %d,\n",
			hCounter.CounterID, hCounter.CreatedAt, hCounter.UpdatedAt, hCounter.Value)
	}

	jsonData, err := json.Marshal(hCounters)
	if err != nil {
		log.Fatalf("❌ Error marshaling historical_counter data to JSON.\n %s", err)
		http.Error(w, "Error marshaling historical_counter data to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
