package repositories

import (
    "testing"
	"time"
    "github.com/google/uuid"
	
)


// Mock models for testing
type TaskQueueModel struct {
    Id           uuid.UUID `json:"id"`
    ServiceName  string    `json:"service_name"`
    ServiceVersion string    `json:"service_version"`
    MessageIds   string    `json:"message_ids"`
    CreatedAt    time.Time `json:"created_at"`
    CreatedBy    uuid.UUID `json:"created_by"`
    Sign         int       `json:"sign"`
}

func TestGetTaskQueuesByWorkerGroup(t *testing.T) {
    repo := NewProcessManagerRepo()
    workerGroup := "Banking360"
    limit := 10

    taskQueues, err := repo.GetTaskQueuesByWorkerGroup(workerGroup, limit)
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    if len(taskQueues) != 0 { // Adjust based on expected results from the mock DBPool
        t.Error("Expected zero task queues but got a different number")
    }
}
