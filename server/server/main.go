package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/db"
	"server/handlers"
)

func main() {
	db.Connect()
	http.HandleFunc("/counter", handlers.GetCounter)
	http.HandleFunc("/start-incr", handlers.StartAutoUpdateCounter)
	http.HandleFunc("/stop-incr", handlers.StopAutoUpdateCounter)
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Print("🏗️  Starting the server...")
	log.Printf("🚀 Listening on %s\n", addr)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Fatalf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
