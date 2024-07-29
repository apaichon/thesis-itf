package process

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/apaichon/thesis-itf/itf/internal/models"
	"github.com/apaichon/thesis-itf/itf/internal/repositories"
)

type ProcessManager struct {
	banking360Service  *Banking360Service
	coreBankingService *CoreBankingService
}

func NewProcessManager() *ProcessManager {
	repo := repositories.NewProcessManagerRepo()
	return &ProcessManager{
		banking360Service:  &Banking360Service{processManagerRepo: repo},
		coreBankingService: &CoreBankingService{processManagerRepo: repo},
	}
}

type Banking360Service struct {
	processManagerRepo *repositories.ProcessManagerRepo
}

func (b *Banking360Service) GetTaskQueue(limit int) ([]models.TaskQueueModel, error) {

	tasks, err := b.processManagerRepo.GetTaskQueuesByWorkerGroup("Banking360", limit)

	if err != nil {
		return nil, err
	}
	return tasks, err
}

func (b *Banking360Service) Process() {
	for {
		tasks, err := b.GetTaskQueue(1000)
		if err != nil {
			log.Printf("Error getting task queue: %v", err)
			continue
		}
		for _, task := range tasks {
			// Process the task as needed

			TransferIntraBankRestApiFlow(task)

			log.Println("Processing Banking360 task:", task.ServiceName)
		}
		time.Sleep(3 * time.Second)
	}
}

type CoreBankingService struct {
	processManagerRepo *repositories.ProcessManagerRepo
}

func (c *CoreBankingService) GetTaskQueue(limit int) ([]models.TaskQueueModel, error) {

	tasks, err := c.processManagerRepo.GetTaskQueuesByWorkerGroup("Legacy", limit)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

func (c *CoreBankingService) Process() {
	for {
		tasks, err := c.GetTaskQueue(1000)
		if err != nil {
			log.Printf("Error getting task queue: %v", err)
			continue
		}
		for _, task := range tasks {
			// Process the task as needed

			switch task.ServiceName {
			case "Legacy.DepositService":
				DepositCoreBankFlow(task)
			case "Legacy.WithdrawService":
				WithdrawalCoreBankFlow(task)
			default:
				fmt.Println("[error] Service not found")
			}
			log.Println("Processing CoreBanking task:", task.ServiceName)
		}
		time.Sleep(3 * time.Second)
	}
}

func (pm *ProcessManager) Start() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		pm.banking360Service.Process()
	}()

	go func() {
		defer wg.Done()
		pm.coreBankingService.Process()
	}()

	wg.Wait()
}
