package coroutines

import (
	"context"
	"log"
	"server/db"
	"server/utils"
	"time"
)

var cancelFunc context.CancelFunc
var taskRunning bool

func runBackgroundTask(ctx context.Context) {
	incrementFrequencyInHours, ok := utils.GetEnvInt("COUNTER_INCREMENT_FREQUENCY_IN_HOURS")
	if ok != nil {
		log.Println("❌ Error getting COUNTER_INCREMENT_FREQUENCY_IN_HOURS")
		return
	}
	ticker := time.NewTicker(time.Duration(incrementFrequencyInHours) * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := db.UpdateCounter()
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			log.Println("🛑 Background task stopped")
			return
		}
	}
}

func RunBackgroundTask() {
	if taskRunning {
		log.Println("⚠️ Background task is already running")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancelFunc = cancel
	taskRunning = true
	go func() {
		runBackgroundTask(ctx)
		taskRunning = false
	}()
}

func StopBackgroundTask() {
	if cancelFunc != nil {
		cancelFunc()
	}
}
