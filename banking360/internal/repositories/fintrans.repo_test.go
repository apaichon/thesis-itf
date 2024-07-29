package repositories

import (
	"fmt"
	"testing"
	"time"

	"banking360/internal/data/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeposit(t *testing.T) {
	repo := NewFinancialTransactionRepo()

	destinationInstitution := "x"

	// Create a sample FinancialTransaction
	sampleTrans := models.FinancialTransaction{
		TransactionID:          uuid.New(),
		TransactionDate:        time.Now(),
		Amount:                 999.00,
		Currency:               "THB",
		AccountID:              uuid.New(),
		DestinationInstitution: &destinationInstitution,
		ActBy:                  uuid.New(),
		ActAt:                  time.Now(),
		Description:            "Test deposit",
		CreatedBy:              uuid.New(),
	}
	_, err := repo.Deposit(sampleTrans)
	// fmt.Printf("exec %v", iexec)
	assert.Equal(t, err, nil)
}

func TestGetDeposits(t *testing.T) {
	repo := NewFinancialTransactionRepo()

	results, paginated, err := repo.GetDepositsByTextSearchPagination("", 1, 50)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	fmt.Printf("paginated %v", paginated.TotalItems)
	assert.Greater(t, len(results), 0)
}
