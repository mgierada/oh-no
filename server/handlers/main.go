package handlers

import (
	"log"
	"net/http"
	"server/db"
)

type ServerResponse struct {
	Message string `json:"message"`
}

func GetCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received /counter request\n")

	counter, err := db.GetCounter()
	if err != nil {
		log.Fatalf("❌ Error retrieving counter data.\n %s", err)
	}

	MarshalJson(w, http.StatusOK, counter)
}

func StartAutoUpdateCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received POST /start-incr request")
	db.RunBackgroundTask()
	response := ServerResponse{Message: "Background task stared."}
	MarshalJson(w, http.StatusOK, response)
	log.Println("🟢 Background task started")
}

func StopAutoUpdateCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received POST /stop_incr request")
	db.StopBackgroundTask()
	response := ServerResponse{Message: "Background task stopped."}
	MarshalJson(w, http.StatusOK, response)
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
		MarshalJson(w, http.StatusOK, response)
		log.Println("🟢 Oh No! Event recorded")
	default:
		log.Printf("❌ Only POST method is allowed")
		errResponse := ServerResponse{Message: "Only POST method is allowed"}
		MarshalJson(w, http.StatusMethodNotAllowed, errResponse)
		return
	}
}

func GetHistoricalCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received GET /historical request\n")

	hCounters, err := db.GetHistoricalCounters()
	if err != nil {
		log.Fatalf("❌ Error retrieving historical_counter data.\n %s", err)
	}

	MarshalJson(w, http.StatusOK, hCounters)
}
