package main

import (
	"log"
	"time"

	"github.com/apaichon/thesis-itf/itf/config"
	qmanager "github.com/apaichon/thesis-itf/itf/internal/queue"
	"github.com/apaichon/thesis-itf/itf/internal/repositories"
)

func main() {
	cfg := config.NewConfig()
	repo := repositories.NewMessageRepo()
	pq := qmanager.NewPersistencyQueue(cfg.TotalMemQueues, cfg.PersistencyBatchSize, repo)

	go func() {
		ticker := time.NewTicker(time.Duration(cfg.QueueManagerInterval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Collection and Insert.")
			pq.CollectAndInsert()
		}
	}()

	// Block main goroutine
	select {}
}
