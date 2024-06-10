package handlers

import (
	"fmt"
	"log"
	"net/http"
	"server/db"
	"server/utils"
)

type ServerResponse struct {
	Message string `json:"message"`
}

func recordEvent(w http.ResponseWriter, r *http.Request, tableToResetAndLock string, tableToUnlock string, historicalTable string, serverResponseOkMessage string) {
	switch r.Method {
	case "POST":

		last_value, err := db.ResetCounter(tableToResetAndLock)
		if err != nil {
			log.Printf("❌ Error resetting %s.\n %s", tableToResetAndLock, err)
			http.Error(w, fmt.Sprintf("Error resetting %s.", tableToResetAndLock), http.StatusInternalServerError)
			return
		}

		_, err = db.UnlockCounter(tableToUnlock)
		log.Printf("🔓 Unlocking %s...", tableToUnlock)
		if err != nil {
			log.Printf("❌ Error unlocking %s .\n %s", tableToUnlock, err)
			http.Error(w, fmt.Sprintf("Error unlocking %s.", tableToUnlock), http.StatusInternalServerError)
			return
		}

		_, err = db.LockCounter(tableToResetAndLock)
		log.Printf("🔒 Locking %s...", tableToResetAndLock)
		if err != nil {
			log.Printf("❌ Error locking %s .\n %s", tableToResetAndLock, err)
			http.Error(w, fmt.Sprintf("Error locking %s.", tableToResetAndLock), http.StatusInternalServerError)
			return
		}

		err = db.CreateHistoricalCounter(historicalTable, last_value)
		if err != nil {
			log.Printf("❌ Error creating %s.\n %s", historicalTable, err)
			http.Error(w, fmt.Sprintf("Error creating %s.", historicalTable), http.StatusInternalServerError)
			return
		}

		response := ServerResponse{Message: serverResponseOkMessage}
		MarshalJson(w, http.StatusOK, response)
		log.Printf("🟢 %s", serverResponseOkMessage)

	default:
		log.Printf("❌ Only POST method is allowed")
		errResponse := ServerResponse{Message: "Only POST method is allowed"}
		MarshalJson(w, http.StatusMethodNotAllowed, errResponse)
		return
	}
}

func RecordOhNoEvent(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received /ohno request")
	serverResponseOkMessage := "Oh No! Event recorded"
	recordEvent(w, r, utils.TableInstance.Counter, utils.TableInstance.OhnoCounter, utils.TableInstance.HistoricalCounter, serverResponseOkMessage)

}

func RecordFineEvent(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔗 received /fine request")
	serverResponseOkMessage := "It's all good now! Event recorded"
	recordEvent(w, r, utils.TableInstance.OhnoCounter, utils.TableInstance.Counter, utils.TableInstance.HistoricalOhnoCounter, serverResponseOkMessage)
}
