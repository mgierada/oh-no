package handlers

import (
	"log"
	"net/http"
	"server/db"
)

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
