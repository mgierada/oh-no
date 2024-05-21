package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
)

type ServerResponse struct {
	Message string `json:"message"`
}

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
	log.Printf("🔗 received /ohno request")
	switch r.Method {
	case "POST":
		last_value, ok := db.ResetCounter()
		if ok != nil {
			log.Printf("❌ Error resetting counter.\n %s", ok)
			http.Error(w, "Error resetting counter.", http.StatusInternalServerError)
			return
		}
		ok = db.CreateHistoricalCounter(last_value)
		if ok != nil {
			log.Printf("❌ Error creating historical counter.\n %s", ok)
			http.Error(w, "Error creating historical counter.", http.StatusInternalServerError)
			return
		}
		response := ServerResponse{Message: "Oh No! Event recorded"}
		marshalJson(w, http.StatusOK, response)
		log.Println("🟢 Oh No! Event recorded")
	default:
		log.Printf("❌ Only POST method is allowed")
		errResponse := ServerResponse{Message: "Only POST method is allowed"}
		marshalJson(w, http.StatusMethodNotAllowed, errResponse)
		return
	}
}

func marshalJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("❌ Error marshaling data to JSON.\n %s", err)
		http.Error(w, "Error marshaling data to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
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
