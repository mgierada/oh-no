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
		last_value, err := db.ResetCounter()
		if err != nil {
			log.Printf("❌ Error resetting counter.\n %s", err)
			http.Error(w, "Error resetting counter.", http.StatusInternalServerError)
			return
		}

		_, err = db.LockCounter("counter")
		log.Printf("🔓 Locking counter")

		if err != nil {
			log.Printf("❌ Error locking counter.\n %s", err)
			http.Error(w, "Error locking counter.", http.StatusInternalServerError)
			return
		}

		_, err = db.UnlockCounter("ohno_counter")
		log.Printf("🔓 Unlocking ohno counter")
		if err != nil {
			log.Printf("❌ Error unlocking ohno counter.\n %s", err)
			http.Error(w, "Error unlocking ohno counter.", http.StatusInternalServerError)
			return
		}

		err = db.CreateHistoricalCounter(last_value)
		if err != nil {
			log.Printf("❌ Error creating historical counter.\n %s", err)
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
