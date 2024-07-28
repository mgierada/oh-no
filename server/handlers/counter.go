package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/db"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// func GetCounter(w http.ResponseWriter, r *http.Request) {
// 	log.Printf("🔗 received /counter request\n")
//
// 	counter, err := db.GetCounter("counter")
// 	if err != nil {
// 		log.Fatalf("❌ Error retrieving counter data.\n %s", err)
// 	}
// EnableCors(c)
// c.JSON(http.StatusOK, counter)
// }

func GetCounter(c *gin.Context) {
	log.Printf("🔗 received /counter request\n")

	counter, err := db.GetCounter("counter")
	if err != nil {
		log.Printf("❌ Error retrieving counter data.\n %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving counter data"})
		return
	}
	c.JSON(http.StatusOK, counter)

}

func GetOhnoCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received /ohno-counter request\n")

	counter, err := db.GetCounter("ohno_counter")
	if err != nil {
		log.Fatalf("❌ Error retrieving ohno_counter data.\n %s", err)
	}

	MarshalJson(&w, http.StatusOK, counter)
}

func getHistoricalCounterEntries(w http.ResponseWriter, tableName string) {
	hCounters, err := db.GetHistoricalCounters(tableName)
	if err != nil {
		log.Fatalf("❌ Error retrieving %s data.\n %s", tableName, err)
	}
	MarshalJson(&w, http.StatusOK, hCounters)
}

func GetHistoricalCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received GET /historical/counter request\n")
	getHistoricalCounterEntries(w, utils.TableInstance.HistoricalCounter)
}

func GetHistoricalOhnoCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received GET /historical/ohno-counter request\n")
	getHistoricalCounterEntries(w, utils.TableInstance.HistoricalOhnoCounter)
}

func IncrementCounter(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received /increment request")
	switch r.Method {
	case "POST":

		if !IsCounterLocked() && !IsOhnoCounterLocked() {
			log.Printf("🤔 Both counters are unlocked. Something went wrong.")
			errResponse := ServerResponse{Message: "Both counters are unlocked. Something went wrong."}
			MarshalJson(&w, http.StatusInternalServerError, errResponse)
			log.Println("⁉️ Both counters are unlocked. Something went wrong.")
			return
		}

		if IsCounterLocked() && IsOhnoCounterLocked() {
			log.Printf("🤔 Both counters are locked. Something went wrong.")
			errResponse := ServerResponse{Message: "Both counters are locked. Something went wrong."}
			log.Println("⁉️ Both counters are locked. Something went wrong.")
			MarshalJson(&w, http.StatusInternalServerError, errResponse)
			return
		}

		if IsOhnoCounterLocked() {
			log.Printf("😀 Ohno Counter is locked. Proceeding with incrementing counter. Another happy day.")
			isUpdated := db.UpdateCounter()

			if !isUpdated {
				errResponse := ServerResponse{Message: "Counter not incremented. Conditions not met."}
				MarshalJson(&w, http.StatusOK, errResponse)
				return
			}

			response := ServerResponse{Message: "Counter incremented successfully"}
			MarshalJson(&w, http.StatusOK, response)
			log.Println("🟢 Counter incremented successfully")
		}

		if IsCounterLocked() {
			log.Printf("🤮 Counter is locked. Proceeding with incrementing ohno counter. Illness continues.")
			isUpdated := db.UpdateOhnoCounter()

			if !isUpdated {
				errResponse := ServerResponse{Message: "Counter not incremented. Conditions not met."}
				MarshalJson(&w, http.StatusOK, errResponse)
				return
			}

			response := ServerResponse{Message: "Ohno counter incremented successfully"}
			MarshalJson(&w, http.StatusOK, response)
			log.Println("🟢 Ohno counter incremented successfully")
		}

	default:
		log.Printf("❌ Only POST method is allowed")
		errResponse := ServerResponse{Message: "Only POST method is allowed"}
		MarshalJson(&w, http.StatusMethodNotAllowed, errResponse)
		return
	}
}

type ManualCouterIncrementRequest struct {
	Value int `json:"value"`
}

func SetCounterValue(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received /manual-increment request")

	switch r.Method {
	case "POST":
		var body ManualCouterIncrementRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Printf("❌ Error decoding request body.\n %s", err)
			errResponse := ServerResponse{Message: "Error decoding request body"}
			MarshalJson(&w, http.StatusBadRequest, errResponse)
			return
		}
		db.SetCounter(body.Value)
		response := ServerResponse{Message: "Counter incremented successfully"}
		MarshalJson(&w, http.StatusOK, response)
		log.Println("🟢 Counter incremented successfully")
	default:
		log.Printf("❌ Only POST method is allowed")
		errResponse := ServerResponse{Message: "Only POST method is allowed"}
		MarshalJson(&w, http.StatusMethodNotAllowed, errResponse)
		return
	}
}
