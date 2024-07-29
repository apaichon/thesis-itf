package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"os"
	"strconv"

	"github.com/apaichon/thesis-itf/itf/internal/models"
	"github.com/apaichon/thesis-itf/itf/internal/repositories"
	"github.com/google/uuid"
)

var (
	messageChannel = make(chan models.MessageModel, 3000) // Buffer size can be adjusted
	batchDuration  = 500 * time.Millisecond
	mu             sync.Mutex
)

// MessageRequest represents the request body for the message
type MessageRequest struct {
	System    string `json:"system"`
	Topic     string `json:"topic"`
	Content   string `json:"content"`
	Remark    string `json:"remark"`
	CreatedBy string `json:"createdBy"`
	Sign      int8   `json:"sign"`
}

// PerformBatchProcessing handles the batch processing of messages
func PerformBatchProcessing(repo *repositories.MessageRepo) {
	var messages []models.MessageModel

	ticker := time.NewTicker(batchDuration)
	defer ticker.Stop()

	for {
		select {
		case message := <-messageChannel:
			messages = append(messages, message)
		case <-ticker.C:
			mu.Lock()
			if len(messages) > 0 {
				log.Println("Processing batch of messages:", len(messages))
				_, err := repo.InsertBatch(messages)
				if err != nil {
					log.Printf("Error processing batch: %v", err)
				} else {
					log.Println("Batch processed successfully")
				}
				messages = nil
			}
			mu.Unlock()
		}
	}
}

// PerformBatchProcessing handles the batch processing of messages
func PerformBatchProcessing2(repo *repositories.MessageRepo) {
	var messages []models.MessageModel

	ticker := time.NewTicker(batchDuration)
	defer ticker.Stop()

	for {
		select {
		case message := <-messageChannel:
			messages = append(messages, message)
		case <-ticker.C:
			mu.Lock()
			if len(messages) > 0 {
				// Copy messages to a new slice and clear the original slice
				messagesToInsert := make([]models.MessageModel, len(messages))
				copy(messagesToInsert, messages)
				messages = nil // Clear the original slice

				// Process the copied messages in a separate goroutine
				go func(messages []models.MessageModel) {
					log.Println("Processing batch of messages:", len(messages))
					_, err := repo.InsertBatch(messages)
					if err != nil {
						log.Printf("Error processing batch: %v", err)
					} else {
						log.Println("Batch processed successfully")
					}
				}(messagesToInsert)
			}
			mu.Unlock()
		}
	}
}

// HandleMessageRequest handles the HTTP POST request for submitting a message
func HandleMessageRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Printf("Error: Invalid request method %s", r.Method)
		return
	}

	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	message := models.MessageModel{
		Id:        uuid.New(),
		System:    req.System,
		Topic:     req.Topic,
		Content:   req.Content,
		Remark:    req.Remark,
		CreatedAt: time.Now(),
		CreatedBy: req.CreatedBy,
		Sign:      req.Sign,
	}

	mu.Lock()
	messageChannel <- message
	mu.Unlock()

	response := map[string]interface{}{
		"messageId": message.Id,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
		return
	}

	log.Printf("Message %s received", message.Id)
}

func main() {
	// cfg := config.NewConfig()
	repo := repositories.NewMessageRepo()
	go PerformBatchProcessing(repo)

	ports := getPortsFromArgs()
	if len(ports) == 0 {
		log.Fatal("No ports specified. Usage: go run main.go <port1> <port2> ...")
	}

	// Create the handler once
	handler := http.HandlerFunc(HandleMessageRequest)

	var wg sync.WaitGroup
	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			startServer(p, handler)
		}(port)
	}

	wg.Wait()
}

func getPortsFromArgs() []int {
	var ports []int
	for _, arg := range os.Args[1:] {
		port, err := strconv.Atoi(arg)
		if err != nil {
			log.Printf("Invalid port number: %s", arg)
			continue
		}
		ports = append(ports, port)
	}
	return ports
}

func startServer(port int, handler http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/api/submit", handler)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server starting on port %d...", port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Printf("Error starting server on port %d: %v", port, err)
	}
}