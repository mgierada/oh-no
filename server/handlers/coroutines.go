package handlers

import (
	"log"
	"net/http"
	"server/coroutines"
)

func StartAutoUpdateCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received POST /start-incr request")
	coroutines.RunBackgroundTask()
	response := ServerResponse{Message: "Background task stared."}
	MarshalJson(w, http.StatusOK, response)
	log.Println("🟢 Background task started")
}

func StopAutoUpdateCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received POST /stop_incr request")
	response := ServerResponse{Message: "Background task stopped."}
	coroutines.StopBackgroundTask()
	MarshalJson(w, http.StatusOK, response)
	log.Println("🔴 Background task stopped")
}
