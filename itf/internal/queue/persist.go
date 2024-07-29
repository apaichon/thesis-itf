package qmanager

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/apaichon/thesis-itf/itf/internal/models"
	"github.com/apaichon/thesis-itf/itf/internal/repositories"
)

type PersistencyQueue struct {
	queue     chan models.MessageModel
	batchSize int
	repo      *repositories.MessageRepo
}

func NewPersistencyQueue(size int, batchSize int, repo *repositories.MessageRepo) *PersistencyQueue {
	return &PersistencyQueue{
		queue:     make(chan models.MessageModel, size),
		batchSize: batchSize,
		repo:      repo,
	}
}

func (p *PersistencyQueue) CollectAndInsert() {
	var messages []models.MessageModel

	// Collect messages from both endpoints
	// , "http://localhost:8082/api/messages"
	// , "http://localhost:8083/api/messages"
	for _, url := range []string{"http://localhost:8081/api/messages",
		"http://localhost:8082/api/messages",
		"http://localhost:8083/api/messages"} {
		data, err := p.getDataFromAPI(url)
		if err != nil {
			log.Printf("Error getting data from %s: %v", url, err)
			continue
		}
		messages = append(messages, data.Messages...)
	}

	if len(messages) > 0 {
		p.insertIntoClickhouse(messages)

	}

}

func (p *PersistencyQueue) getDataFromAPI(url string) (struct {
	Count    int                   `json:"count"`
	Messages []models.MessageModel `json:"messages"`
}, error) {
	var result struct {
		Count    int                   `json:"count"`
		Messages []models.MessageModel `json:"messages"`
	}

	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	return result, err
}

func (p *PersistencyQueue) insertIntoClickhouse(messages []models.MessageModel) {
	for len(messages) > 0 {
		batchSize := min(len(messages), p.batchSize)
		batch := messages[:batchSize]
		messages = messages[batchSize:]
		log.Printf("batchSize: %v", batchSize)

		_, err := p.repo.InsertBatch(batch)
		if err != nil {
			log.Printf("Error inserting batch into Clickhouse: %v", err)
			go p.writeErrorToFile(batch, err)
		}
	}
}

func (p *PersistencyQueue) writeErrorToFile(messages []models.MessageModel, err error) {
	file, fileErr := os.OpenFile("error_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		log.Printf("Error opening log file: %v", fileErr)
		return
	}
	defer file.Close()

	for _, msg := range messages {
		logEntry := fmt.Sprintf("Error: %v, Message: %+v\n", err, msg)
		if _, fileErr := file.WriteString(logEntry); fileErr != nil {
			log.Printf("Error writing to log file: %v", fileErr)
		}
	}
}

func (p *PersistencyQueue) Print() {
	log.Println("Queue Run")
}

func (p *PersistencyQueue) GetMessageQueues() ([]models.MessageModel, error) {
	messages := []models.MessageModel{}

	return messages, nil
}
