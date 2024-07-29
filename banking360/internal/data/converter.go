package data

import (
	"fmt"
	"banking360/internal/data/models"
)
func StringToTransactionStatus(status string) (models.TransactionStatus, error) {
    switch status {
    case "pending":
        return models.Pending, nil
    case "completed":
        return models.Completed, nil
    case "failed":
        return models.Failed, nil
    default:
        return 0, fmt.Errorf("unknown transaction status: %s", status)
    }
}