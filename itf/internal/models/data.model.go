package models

import (
	"time"
	"github.com/google/uuid"
)
// QueueManager.QueueMessages
type MessageModel struct {
	Id uuid.UUID `json:"id"`
	System   string `json:"system"`
	Topic   string `json:"topic"`
	Content string `json:"content"`
	Remark string `json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
	Sign int8 `json:"sign"`
}

// OrchestrationManager.ServiceProvider
type ServiceProviderModel struct {
	WorkerGroup string  `json:"worker_group"`
	ServiceName        string    `json:"service_name"`
	ServiceFullName    string    `json:"service_full_name"`
	ServiceDescription string    `json:"service_description"`
	Tps int  `json:"tps"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	Active             int8      `json:"active"`
}

// OrchestrationManager.ServiceProcesses
type ServiceProcessModel struct {
	Topic       string    `json:"topic"`
	Services  []string  `json:"services"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	Sign        int8      `json:"sign"`
	Version     uint16    `json:"version"`
}

type TaskQueueModel struct {
	Id        uuid.UUID           `json:"id"`
	ServiceName  string       `json:"service_name"`
	ServiceVersion  uint16             `json:"service_version"`
	MessageIds []uuid.UUID           `json:"message_ids"`
	CreatedAt time.Time           `json:"created_at"`
	CreatedBy string              `json:"created_by"`
	Sign      int8                `json:"sign"`
}

type ServiceTasksModel struct {
	Name string `json:"name"`
	Tasks []TaskQueueModel `json:"tasks"`
}
type WorkerGroupModel struct {
	Name string `json:"name"`
	Tps int `json:"tps"`
	ServiceTasks []ServiceTasksModel `json:"service_tasks`
}

// ProcessManager.ProcessFlow
type ProcessFlowModel struct {
	Id         uuid.UUID `json:"id"`
	TaskId     uuid.UUID `json:"task_id"`
	SubTaskNo  uint16    `json:"sub_task_no"`
	TaskState  string    `json:"task_state"` // Use an enum type for better type safety
	Remarks    string    `json:"remarks"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	Sign       int8      `json:"sign"`
	Version    uint16    `json:"version"`
}

// ProcessManager.Archive
type ArchiveModel struct {
	Id           uuid.UUID `json:"id"`
	CorrelationId uuid.UUID `json:"correlation_id"`
	TableName    string    `json:"table_name"`
	JsonData     string    `json:"json_data"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
}

// ProcessManager.AccountMapping
type AccountMappingModel struct {
	AccountId           uuid.UUID `json:"accountId"`
	AccountLegacyId int `json:"accountLegacyId"`
}

