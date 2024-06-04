package handlers

import (
	"log"
	"net/http"
	"server/db"
)

type ServerResponse struct {
	Message string `json:"message"`
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
