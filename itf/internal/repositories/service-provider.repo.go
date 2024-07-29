package repositories

import (
	"fmt"
	"log"
	"time"

	"github.com/apaichon/thesis-itf/itf/config"
	"github.com/apaichon/thesis-itf/itf/internal/data"
	"github.com/apaichon/thesis-itf/itf/internal/models"
)

type ServiceProviderRepo struct {
	DBPool *data.DBPool
}

func NewServiceProviderRepo() *ServiceProviderRepo {
	cfg := config.NewConfig()
	dsn := cfg.GetDns()
	pool, err := data.GetDBPool(dsn, 10, 5, 5*time.Minute)
	if err != nil {
		log.Fatalf("Failed to get database pool: %v", err)
	}
	return &ServiceProviderRepo{DBPool: pool}
}

func (repo *ServiceProviderRepo) Insert(serviceInfo models.ServiceProviderModel) (int64, error) {

	command := "INSERT INTO OrchestrationManager.ServiceProvider (ServiceName, ServiceFullName, ServiceDescription, CreatedAt, CreatedBy, Active) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := repo.DBPool.Insert(command, serviceInfo.ServiceName, serviceInfo.ServiceFullName, serviceInfo.ServiceDescription, serviceInfo.CreatedAt, serviceInfo.CreatedBy, serviceInfo.Active)
	if err != nil {
		log.Printf("Insert error: %v", err)
	}
	return result, nil
}

func (repo *ServiceProviderRepo) GetAllServices() ([]models.ServiceProviderModel, error) {

	query := `SELECT *
   FROM OrchestrationManager.ServiceProvider final`

	rows, err := repo.DBPool.Query(query)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var services []models.ServiceProviderModel
	for rows.Next() {
		var service models.ServiceProviderModel
		if err := rows.Scan(
			&service.WorkerGroup,
			&service.ServiceName,
			&service.ServiceFullName,
			&service.ServiceDescription,
			&service.Tps,
			&service.CreatedAt,
			&service.CreatedBy,
			&service.Active,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		services = append(services, service)
	}

	return services, nil
}
