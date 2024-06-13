package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"server/utils"
)

func MarshalJson(w *http.ResponseWriter, statusCode int, data interface{}) {
	utils.EnableCors(w)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(statusCode)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("❌ Error marshaling data to JSON.\n %s", err)
		http.Error(*w, "Error marshaling data to JSON", http.StatusInternalServerError)
		return
	}
	(*w).Write(jsonData)
}

func IsCounterLocked() bool {
	counter, err := db.GetCounter("counter")

	if err != nil {
		return false
	}
	return counter.IsLocked
}

func IsOhnoCounterLocked() bool {
	counter, err := db.GetCounter("ohno_counter")
	if err != nil {
		return false
	}
	return counter.IsLocked
}
