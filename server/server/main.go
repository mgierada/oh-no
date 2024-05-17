package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)
import (
	"server/db"
)

func getCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received GET /counter request\n")

	counters, err := db.GetAllCounterData()
	if err != nil {
		log.Fatalf("❌ Error retrieving counter data.\n %s", err)
	}

	for _, counter := range counters {
		log.Printf("Current Value: %d, Updated At: %s, Reseted At: %v\n",
			counter.CurrentValue, counter.UpdatedAt, counter.ResetedAt.String)
	}

	// Convert the counters slice to JSON
	jsonData, err := json.Marshal(counters)
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

func updateCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received POST /increment request")
	db.ManualIncrement()
}

func getHello(w http.ResponseWriter, r *http.Request) {
	log.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func main() {
	db.Connect()
	http.HandleFunc("/counter", getCounter)
	http.HandleFunc("/hello", getHello)
	http.HandleFunc("/increment", updateCounter)
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
