package qmanager

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/apaichon/thesis-itf/itf/internal/models"
	"github.com/apaichon/thesis-itf/itf/internal/repositories"
)

type InMemQueue struct {
	queue chan models.MessageModel
	temp  repositories.TempMessageRepo
	mu    sync.Mutex
}

func NewInMemQueue(size int, tempRepo repositories.TempMessageRepo) *InMemQueue {
	return &InMemQueue{
		queue: make(chan models.MessageModel, size),
		temp:  tempRepo,
	}
}

func (q *InMemQueue) AddMessage(msg models.MessageModel) {
	msg.CreatedAt = time.Now()
	msg.Sign = 1
	select {
	case q.queue <- msg:
		log.Printf("Message submitted to queue %s: %+v", msg.Topic, msg)
	default:
		log.Printf("Error: Queue is full")
		q.temp.Insert(msg)
	}
}

func (p *InMemQueue) MessageHandler(w http.ResponseWriter, r *http.Request) {
	var msg models.MessageModel
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuid := uuid.New()
	msg.Id = uuid
	go p.AddMessage(msg)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Message:" + uuid.String() + " submitted successfully"})
}

func (q *InMemQueue) GetMessages() []models.MessageModel {
	q.mu.Lock()
	defer q.mu.Unlock()

	var messages []models.MessageModel
	timeout := time.After(2 * time.Second) // Set a timeout to prevent blocking indefinitely

	// Collect messages from the channel
collectLoop:
	for {
		select {
		case msg, ok := <-q.queue:
			if !ok {
				// Channel is closed
				break collectLoop
			}
			messages = append(messages, msg)
		case <-timeout:
			// Timeout reached
			break collectLoop
		default:
			// No more messages available without blocking
			break collectLoop
		}
	}

	return messages
}

func (p *InMemQueue) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages := p.GetMessages()

	// Prepare the response
	response := struct {
		Count    int                   `json:"count"`
		Messages []models.MessageModel `json:"messages"`
	}{
		Count:    len(messages),
		Messages: messages,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("Retrieved and cleared %d messages", len(messages))
}
