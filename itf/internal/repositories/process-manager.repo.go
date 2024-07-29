package repositories

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/apaichon/thesis-itf/itf/config"
	"github.com/apaichon/thesis-itf/itf/internal/data"
	"github.com/apaichon/thesis-itf/itf/internal/models"
	"github.com/google/uuid"
)

type ProcessManagerRepo struct {
	DBPool *data.DBPool
}

func NewProcessManagerRepo() *ProcessManagerRepo {
	cfg := config.NewConfig()
	dsn := cfg.GetDns()
	pool, err := data.GetDBPool(dsn, 10, 5, 5*time.Minute)
	if err != nil {
		log.Fatalf("Failed to get database pool: %v", err)
	}
	return &ProcessManagerRepo{DBPool: pool}
}

// InsertBatchTopicA inserts a batch of messages into the database
func (repo *ProcessManagerRepo) InsertMessages(messages []models.MessageModel) (int64, error) {
	if len(messages) == 0 {
		return 0, nil
	}

	// Construct the base command
	baseCommand := "INSERT INTO ProcessManager.MessageQueues (Id, SystemOwner, Topic, Content, Remarks, CreatedAt, CreatedBy, Sign) VALUES "

	// Prepare placeholders and arguments
	var placeholders []string
	var args []interface{}

	for _, message := range messages {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args, message.Id, message.System, message.Topic, message.Content, message.Remark, message.CreatedAt, message.CreatedBy, message.Sign)
	}

	// Join placeholders and form the final command
	command := baseCommand + strings.Join(placeholders, ", ")

	// Execute the command
	result, err := repo.DBPool.Insert(command, args...)
	if err != nil {
		log.Printf("Insert error: %v", err)
		return 0, err
	}

	return result, nil
}

func (repo *ProcessManagerRepo) RegisterTasks(tasks []models.TaskQueueModel) (int64, error) {

	if len(tasks) == 0 {
		return 0, nil
	}

	// Construct the base command
	baseCommand := "INSERT INTO ProcessManager.TaskQueues (Id,  MessageIds , ServiceName, ServiceVersion, CreatedAt, CreatedBy, Sign) VALUES "

	// Prepare placeholders and arguments
	var placeholders []string
	var args []interface{}

	for _, task := range tasks {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?)")
		args = append(args, task.Id, task.MessageIds, task.ServiceName, task.ServiceVersion, task.CreatedAt, task.CreatedBy, task.Sign)
	}

	// Join placeholders and form the final command
	command := baseCommand + strings.Join(placeholders, ", ")

	// Execute the command
	result, err := repo.DBPool.Insert(command, args...)
	if err != nil {
		log.Printf("Insert error: %v", err)
		return 0, err
	}

	return result, nil
}

func (repo *ProcessManagerRepo) GetTaskQueues(limit int) ([]models.TaskQueueModel, error) {
	query := `SELECT 
                Id,
                ServiceName,
                ServiceVersion,
                MessageIds,
                CreatedAt,
                CreatedBy,
                Sign
              FROM 
                ProcessManager.TaskQueues final limit ?`

	rows, err := repo.DBPool.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var taskQueues []models.TaskQueueModel
	for rows.Next() {
		var taskQueue models.TaskQueueModel

		if err := rows.Scan(
			&taskQueue.Id,
			&taskQueue.ServiceName,
			&taskQueue.ServiceVersion,
			&taskQueue.MessageIds,
			&taskQueue.CreatedAt,
			&taskQueue.CreatedBy,
			&taskQueue.Sign,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		taskQueues = append(taskQueues, taskQueue)
	}

	return taskQueues, nil
}

func (repo *ProcessManagerRepo) GetTaskQueuesByWorkerGroup(workerGroup string, limit int) ([]models.TaskQueueModel, error) {
	query := `SELECT 
                Id,
                ServiceName,
                ServiceVersion,
                MessageIds,
                CreatedAt,
                CreatedBy,
                Sign
              FROM 
                ProcessManager.TaskQueues final where ServiceName in (Select ServiceFullName from OrchestrationManager.ServiceProvider where WorkerGroup = ? ) limit ? `
	rows, err := repo.DBPool.Query(query, workerGroup, limit)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var taskQueues []models.TaskQueueModel
	for rows.Next() {
		var taskQueue models.TaskQueueModel

		if err := rows.Scan(
			&taskQueue.Id,
			&taskQueue.ServiceName,
			&taskQueue.ServiceVersion,
			&taskQueue.MessageIds,
			&taskQueue.CreatedAt,
			&taskQueue.CreatedBy,
			&taskQueue.Sign,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		taskQueues = append(taskQueues, taskQueue)
	}

	return taskQueues, nil
}

func (repo *ProcessManagerRepo) FlagDeleteTaskQueue(id uuid.UUID) error {

	// Create the parameterized query
	query := `INSERT INTO ProcessManager.TaskQueues
		(Id, ServiceName , ServiceVersion, MessageIds, CreatedAt, CreatedBy, Sign)
		select Id, ServiceName , ServiceVersion, MessageIds, CreatedAt, CreatedBy, -1 as sign
		from ProcessManager.TaskQueues final
		WHERE Id =?`

	// Execute the query
	err := repo.DBPool.Delete(query, id)
	if err != nil {
		log.Printf("Failed to execute insert-select query: %v", err)
		return err
	}

	return nil
}

func (repo *ProcessManagerRepo) FlagDeleteTaskQueues(ids []uuid.UUID) error {
	// Convert the UUIDs to strings
	uuidStrings := make([]string, len(ids))
	for i, id := range ids {
		uuidStrings[i] = fmt.Sprintf("'%s'", id.String())
	}

	// Create the parameterized query
	query := fmt.Sprintf(`
		INSERT INTO ProcessManager.TaskQueues
		(Id, ServiceName , serviceVersion, MessageIds, CreatedAt, CreatedBy, Sign)
		select Id, ServiceName , serviceVersion, MessageIds, CreatedAt, CreatedBy, -1 as sign
		from ProcessManager.TaskQueues final
		WHERE Id IN (%s)
	`, strings.Join(uuidStrings, ", "))

	// Execute the query
	err := repo.DBPool.Delete(query)
	if err != nil {
		log.Printf("Failed to execute insert-select query: %v", err)
		return err
	}

	return nil
}

// getMessages fetches data from QueueManager.MessageQueues table with limit and offset
func (repo *ProcessManagerRepo) GetMessagesByIds(ids []uuid.UUID) ([]models.MessageModel, error) {

	uuidStrings := make([]string, len(ids))
	for i, id := range ids {
		uuidStrings[i] = fmt.Sprintf("'%s'", id.String())
	}
	// Query to get data with limit, offset and CreatedAt condition
	query := fmt.Sprintf(`SELECT 
	   Id, 
	   SystemOwner, 
	   Topic, 
	   Content, 
	   Remarks, 
	   CreatedAt, 
	   CreatedBy, 
	   Sign 
   FROM ProcessManager.MessageQueues final 
   WHERE Id IN (%s)`, strings.Join(uuidStrings, ", "))

	rows, err := repo.DBPool.Query(query)

	// fmt.Printf("query:%s time:%v" ,query, twoHoursAgo)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var messages []models.MessageModel
	for rows.Next() {
		var msg models.MessageModel
		if err := rows.Scan(
			&msg.Id,
			&msg.System,
			&msg.Topic,
			&msg.Content,
			&msg.Remark,
			&msg.CreatedAt,
			&msg.CreatedBy,
			&msg.Sign,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (repo *ProcessManagerRepo) GetAccountMapping(accountId uuid.UUID) (*models.AccountMappingModel, error) {
	query := `SELECT 
                *
              FROM 
                ProcessManager.AccountMapping Where AccountId= ?`

	rows, err := repo.DBPool.Query(query, accountId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var accountMappings []models.AccountMappingModel
	for rows.Next() {
		var accountMapping models.AccountMappingModel

		if err := rows.Scan(
			&accountMapping.AccountId,
			&accountMapping.AccountLegacyId,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		accountMappings = append(accountMappings, accountMapping)
	}

	if len(accountMappings) > 0 {
		return ToPtr(accountMappings[0]), nil
	}

	return nil, fmt.Errorf("data not found: %v", accountId)
}

// parseMessageIds converts a comma-separated string of UUIDs to a slice of UUIDs
func parseMessageIds(messageIds string) ([]uuid.UUID, error) {
	idStrings := strings.Split(messageIds, ",")
	uuids := make([]uuid.UUID, len(idStrings))

	for i, idStr := range idStrings {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		uuids[i] = id
	}

	return uuids, nil
}

func ToPtr[T any](value T) *T {
	return &value
}
