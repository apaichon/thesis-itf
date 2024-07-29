package repositories

import (
	"log"
	"strings"
	"time"

	"github.com/apaichon/thesis-itf/itf/config"
	"github.com/apaichon/thesis-itf/itf/internal/data"
	"github.com/apaichon/thesis-itf/itf/internal/models"
)

type ServiceProcessRepo struct {
	DBPool *data.DBPool
}

func NewServiceProcessRepo() *ServiceProcessRepo {
	cfg := config.NewConfig()
	dsn := cfg.GetDns()
	pool, err := data.GetDBPool(dsn, 10, 5, 5*time.Minute)
	if err != nil {
		log.Fatalf("Failed to get database pool: %v", err)
	}
	return &ServiceProcessRepo{DBPool: pool}
}

func (repo *ServiceProcessRepo) GetProcessesByTopics(topics []string) ([]models.ServiceProcessModel, error) {
	query := `Select * from OrchestrationManager.ServiceProcesses final where Topic in (?` + strings.Repeat(", ?", len(topics)-1) + ")"

	// Convert topics slice to interface slice
	args := make([]interface{}, len(topics))
	for i, topic := range topics {
		args[i] = topic
	}

	// log.Printf("query services: %v", query)

	rows, err := repo.DBPool.Query(query, args...)
	if err != nil {
		log.Printf("query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var serviceProcesses []models.ServiceProcessModel
	for rows.Next() {
		// log.Printf("got rows:")
		var sp models.ServiceProcessModel
		// var services []string
		if err := rows.Scan(&sp.Topic, &sp.Services, &sp.Description, &sp.CreatedAt, &sp.CreatedBy, &sp.Sign, &sp.Version); err != nil {
			log.Printf("error:%v", err)
			return nil, err
		}

		log.Printf("sp:%v", sp)
		// sp.Services = strings.Split(services, ",")
		serviceProcesses = append(serviceProcesses, sp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return serviceProcesses, nil
}
