package orchestrator

import (
	"fmt"
	"log"
	"time"

	"github.com/apaichon/thesis-itf/itf/internal/models"
	"github.com/apaichon/thesis-itf/itf/internal/repositories"
	"github.com/google/uuid"
)

type OrchestrationManager struct {
	messageRepo        *repositories.MessageRepo
	serviceProcessRepo *repositories.ServiceProcessRepo
	processManagerRepo *repositories.ProcessManagerRepo
}

func NewOrchestrationManager() *OrchestrationManager {
	return &OrchestrationManager{
		messageRepo:        repositories.NewMessageRepo(),
		serviceProcessRepo: repositories.NewServiceProcessRepo(),
		processManagerRepo: repositories.NewProcessManagerRepo(),
	}
}

func (o *OrchestrationManager) Produce(totalMessage int) {
	// get messages => done
	messages, err := o.messageRepo.GetMessages(totalMessage, 0)

	if err != nil {
		log.Printf("[error] om step 1-get messages: %v", err)
		return
	}

	if len(messages) > 0 {
		_, err = o.processManagerRepo.InsertMessages(messages)

		if err != nil {
			log.Printf("[error] om step 2-insert: %v", err)
			return
		}

		topics := o.ExtractUniqueTopics(messages)

		// get service Process by topics => done
		log.Printf("Topics: %v", topics)
		serviceProcesses, err := o.serviceProcessRepo.GetProcessesByTopics(topics)
		log.Printf("ServiceProcesses: %v", serviceProcesses)

		chunkSize := 1000

		for i := 0; i < len(messages); i += chunkSize {
			end := i + chunkSize
			if end > len(messages) {
				end = len(messages)
			}
			chunk := messages[i:end]
			fmt.Printf("Processing chunk from %d to %d\n", i, end)

			// Call your function with the chunk
			tq := o.ConvertToTaskQueueModels(serviceProcesses, chunk)
			err = o.registerTasksWithChunkSize(tq)
			if err != nil {
				log.Printf("[error] om step 3-register task queue: %v", err)
			}
		}

		// flag delete message queues => done
		err = o.flagDeleteWithChunkSize(messages)
		if err != nil {
			log.Printf("[error] om step 4-flag delete: %v", err)
			return
		}

		log.Println("[info] om register task completed.")

	}

}

func (o *OrchestrationManager) registerTasksWithChunkSize(tasks []models.TaskQueueModel) error {
	const chunkSize = 2000

	// Process messages in chunks
	for i := 0; i < len(tasks); i += chunkSize {
		end := i + chunkSize
		if end > len(tasks) {
			end = len(tasks)
		}
		chunk := tasks[i:end]
		_, err := o.processManagerRepo.RegisterTasks(chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *OrchestrationManager) flagDeleteWithChunkSize(messages []models.MessageModel) error {
	const chunkSize = 2000

	// Process messages in chunks
	for i := 0; i < len(messages); i += chunkSize {
		end := i + chunkSize
		if end > len(messages) {
			end = len(messages)
		}
		chunk := messages[i:end]
		err := o.messageRepo.FlagDeleteQueueMessages(chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

// Function to extract unique topics from a slice of MessageModel
func (o *OrchestrationManager) ExtractUniqueTopics(messages []models.MessageModel) []string {
	topicMap := make(map[string]struct{})
	for _, message := range messages {
		topicMap[message.Topic] = struct{}{}
	}

	uniqueTopics := make([]string, 0, len(topicMap))
	for topic := range topicMap {
		uniqueTopics = append(uniqueTopics, topic)
	}

	return uniqueTopics
}

func (o *OrchestrationManager) ConvertToTaskQueueModels(serviceProcesses []models.ServiceProcessModel, messages []models.MessageModel) []models.TaskQueueModel {
	taskQueueMap := make(map[string]*models.TaskQueueModel)
	for _, serviceProcess := range serviceProcesses {
		for _, service := range serviceProcess.Services {
			key := fmt.Sprintf("%s:%s", service, serviceProcess.Topic)
			if _, exists := taskQueueMap[key]; !exists {
				taskQueueMap[key] = &models.TaskQueueModel{
					Id:             uuid.New(),
					ServiceName:    service,
					ServiceVersion: serviceProcess.Version,
					MessageIds:     []uuid.UUID{},
					CreatedAt:      time.Now(),
					CreatedBy:      serviceProcess.CreatedBy,
					Sign:           serviceProcess.Sign,
				}
			}
		}
	}

	for _, message := range messages {
		for _, serviceProcess := range serviceProcesses {
			if message.Topic == serviceProcess.Topic {
				for _, service := range serviceProcess.Services {
					key := fmt.Sprintf("%s:%s", service, serviceProcess.Topic)
					if taskQueue, exists := taskQueueMap[key]; exists {
						taskQueue.MessageIds = append(taskQueue.MessageIds, message.Id)
					}
				}
			}
		}
	}

	// Convert map to slice
	taskQueues := make([]models.TaskQueueModel, 0, len(taskQueueMap))
	for _, taskQueue := range taskQueueMap {
		taskQueues = append(taskQueues, *taskQueue)
	}

	return taskQueues
}
