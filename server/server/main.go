package main

import (
	// "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"server/db"
	// "server/handlers"
)

func main() {
	db.Connect()
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Print("🏗️  Starting the server...")
	log.Printf("🚀 Listening on %s\n", addr)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	// http.HandleFunc("/ohno", handlers.RecordOhNoEvent)
	// http.HandleFunc("/", handlers.RedirectToCounter)
	// http.HandleFunc("/fine", handlers.RecordFineEvent)
	// http.HandleFunc("/historical/counter", handlers.GetHistoricalCounter)
	// http.HandleFunc("/historical/ohno-counter", handlers.GetHistoricalOhnoCounter)
	// http.HandleFunc("/counter", handlers.GetCounter)
	// http.HandleFunc("/ohno-counter", handlers.GetOhnoCounter)
	// http.HandleFunc("/start-incr", handlers.StartAutoUpdateCounter)
	// http.HandleFunc("/stop-incr", handlers.StopAutoUpdateCounter)
	// http.HandleFunc("/increment", handlers.IncrementCounter)
	// http.HandleFunc("/manual-increment", handlers.SetCounterValue)

	//
	// err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	//
	// if errors.Is(err, http.ErrServerClosed) {
	// 	log.Printf("server closed\n")
	// } else if err != nil {
	// 	log.Fatalf("error starting server: %s\n", err)
	// 	os.Exit(1)
	// }

}
