package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"banking360/config"
	"banking360/internal/data/models"
	"banking360/internal/repositories"
)

var (
	transactionChannel chan models.FinancialTransaction
	batchDuration      = 500 * time.Millisecond
	mu                 sync.Mutex
)

type TransferRequest struct {
	SenderAccountID    uuid.UUID `json:"senderAccountId"`
	ReceiverAccountID  uuid.UUID `json:"receiverAccountId"`
	ActBy              uuid.UUID `json:"actBy"`
	CreatedBy          uuid.UUID `json:"createdBy"`
	Amount             float64   `json:"amount"`
	TransactionDate    time.Time `json:"transactionDate"`
	Description        string    `json:"description"`
}

func PerformIntraBankTransfers(ctx context.Context, repo *repositories.FinancialTransactionRepo) {
	ticker := time.NewTicker(batchDuration)
	defer ticker.Stop()

	var transactions []models.FinancialTransaction

	for {
		select {
		case <-ctx.Done():
			return
		case transaction := <-transactionChannel:
			transactions = append(transactions, transaction)
		case <-ticker.C:
			if len(transactions) > 0 {
				processBatch(repo, transactions)
				transactions = transactions[:0] // Clear slice while preserving capacity
			}
		}
	}
}

func processBatch(repo *repositories.FinancialTransactionRepo, transactions []models.FinancialTransaction) {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("Processing batch of transactions: %d", len(transactions))
	_, err := repo.PerformIntraBankTransfers(transactions)
	if err != nil {
		log.Printf("Error processing batch: %v", err)
	} else {
		log.Println("Batch processed successfully")
	}
}

func PerformIntraBankTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction := models.FinancialTransaction{
		TransactionDate:    req.TransactionDate,
		Amount:             req.Amount,
		Currency:           "THB",
		AccountID:          req.SenderAccountID,
		RecipientAccountID: &req.ReceiverAccountID,
		Status:             models.Pending,
		Description:        req.Description,
		ActBy:              req.ActBy,
		ActAt:              time.Now(),
		CreatedBy:          req.CreatedBy,
		CreatedAt:          time.Now(),
	}

	transactionID := uuid.New()
	select {
	case transactionChannel <- transaction:
		response := map[string]interface{}{
			"transactionId": transactionID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		log.Printf("Transaction %s received", transactionID)
	default:
		http.Error(w, "Server is busy, please try again later", http.StatusServiceUnavailable)
	}
}

func main() {
	cfg := config.NewConfig()
	repo := repositories.NewFinancialTransactionRepo()

	// Use a buffered channel with a reasonable size
	transactionChannel = make(chan models.FinancialTransaction, 3000)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go PerformIntraBankTransfers(ctx, repo)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.RestApiPort),
		Handler: http.HandlerFunc(PerformIntraBankTransfer),
	}

	go func() {
		log.Printf("Server starting on port %v...", cfg.RestApiPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	cancel() // Stop the PerformIntraBankTransfers goroutine
	log.Println("Server exiting")
}