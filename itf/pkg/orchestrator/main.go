package main

import (
	"log"
	"time"

	"github.com/apaichon/thesis-itf/itf/config"
	"github.com/apaichon/thesis-itf/itf/internal/orchestrator"
)

func main() {

	cfg := config.NewConfig()
	om := orchestrator.NewOrchestrationManager()

	go func() {
		ticker := time.NewTicker(time.Duration(cfg.OrchestrationManagerInterval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Starting Produce Tasks")
			om.Produce(cfg.PersistencyBatchSize)
		}
	}()

	// Block main goroutine
	select {}
}
