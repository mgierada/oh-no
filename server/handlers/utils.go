package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func MarshalJson(w http.ResponseWriter, statusCode int, data interface{}) {
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
