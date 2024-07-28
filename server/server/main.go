package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"server/db"
	"server/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Print("🏗️  Starting the server...")

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/counter", handlers.GetCounter)
	r.GET("/", handlers.RedirectToCounter)
	r.POST("/ohno", handlers.RecordOhNoEvent)
	r.POST("/fine", handlers.RecordFineEvent)
	// http.HandleFunc("/fine", handlers.RecordFineEvent)
	// http.HandleFunc("/historical/counter", handlers.GetHistoricalCounter)
	// http.HandleFunc("/historical/ohno-counter", handlers.GetHistoricalOhnoCounter)
	// http.HandleFunc("/counter", handlers.GetCounter)
	// http.HandleFunc("/ohno-counter", handlers.GetOhnoCounter)
	// http.HandleFunc("/start-incr", handlers.StartAutoUpdateCounter)
	// http.HandleFunc("/stop-incr", handlers.StopAutoUpdateCounter)
	// http.HandleFunc("/increment", handlers.IncrementCounter)
	// http.HandleFunc("/manual-increment", handlers.SetCounterValue)

	r.Run(addr)
	log.Printf("🚀 Listening on %s\n", addr)

}
