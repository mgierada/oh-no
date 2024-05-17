package main

import (
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
	log.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	log.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func main() {
	db.Connect()
	http.HandleFunc("/", getCounter)
	http.HandleFunc("/hello", getHello)
	var port int = 3333
	var host string = "localhost"
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("🚀 Starting server on %s...\n", addr)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Fatalf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
