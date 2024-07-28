package handlers

import (
	"fmt"
	"log"
	"net/http"
	"server/db"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type ServerResponse struct {
	Message string `json:"message"`
}

func RedirectToCounter(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/counter")
}

func recordEvent(c *gin.Context, tableToResetAndLock string, tableToUnlock string, historicalTable string, serverResponseOkMessage string) {
	switch c.Request.Method {
	case "POST":

		last_value, err := db.ResetCounter(tableToResetAndLock)
		if err != nil {
			log.Printf("❌ Error resetting %s.\n %s", tableToResetAndLock, err)
			c.JSON(http.StatusInternalServerError, ServerResponse{Message: fmt.Sprintf("Error resetting %s.", tableToResetAndLock)})
			return
		}

		_, err = db.UnlockCounter(tableToUnlock)
		log.Printf("🔓 Unlocking %s...", tableToUnlock)
		if err != nil {
			log.Printf("❌ Error unlocking %s .\n %s", tableToUnlock, err)
			c.JSON(http.StatusInternalServerError, ServerResponse{Message: fmt.Sprintf("Error unlocking %s.", tableToUnlock)})
			return
		}

		_, err = db.LockCounter(tableToResetAndLock)
		log.Printf("🔒 Locking %s...", tableToResetAndLock)
		if err != nil {
			log.Printf("❌ Error locking %s .\n %s", tableToResetAndLock, err)
			c.JSON(http.StatusInternalServerError, ServerResponse{Message: fmt.Sprintf("Error locking %s.", tableToResetAndLock)})
			return
		}

		err = db.CreateHistoricalCounter(historicalTable, last_value)
		if err != nil {
			log.Printf("❌ Error creating %s.\n %s", historicalTable, err)
			c.JSON(http.StatusInternalServerError, ServerResponse{Message: fmt.Sprintf("Error creating %s.", historicalTable)})
			return
		}

		response := ServerResponse{Message: serverResponseOkMessage}
		c.JSON(http.StatusOK, response)
		log.Printf("🟢 %s", serverResponseOkMessage)

	default:
		log.Printf("❌ Only POST method is allowed")
		errResponse := ServerResponse{Message: "Only POST method is allowed"}
		c.JSON(http.StatusMethodNotAllowed, errResponse)
		return
	}
}

func RecordOhNoEvent(c *gin.Context) {
	log.Printf("🔗 received /ohno request of type %s", c.Request.Method)
	serverResponseOkMessage := "Oh No! Event recorded"
	recordEvent(c, utils.TableInstance.Counter, utils.TableInstance.OhnoCounter, utils.TableInstance.HistoricalCounter, serverResponseOkMessage)
}

func RecordFineEvent(c *gin.Context) {
	log.Printf("🔗 received /fine request of type %s", c.Request.Method)
	serverResponseOkMessage := "It's all good now! Event recorded"
	recordEvent(c, utils.TableInstance.OhnoCounter, utils.TableInstance.Counter, utils.TableInstance.HistoricalOhnoCounter, serverResponseOkMessage)
}
