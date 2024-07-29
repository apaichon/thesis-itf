package repositories

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/apaichon/thesis-itf/itf/config"
	"github.com/apaichon/thesis-itf/itf/internal/data"
	"github.com/apaichon/thesis-itf/itf/internal/models"
)

type MessageRepo struct {
	DBPool *data.DBPool
}

// NewMessageRepo creates a new instance of MessageRepo
func NewMessageRepo() *MessageRepo {
	cfg := config.NewConfig()
	dsn := cfg.GetDns()
	pool, err := data.GetDBPool(dsn, cfg.MaxConnection, cfg.MaxIdleConnection, time.Duration(cfg.ConnectionLifeTime)*time.Minute)
	if err != nil {
		log.Fatalf("Failed to get database pool: %v", err)
	}
	return &MessageRepo{DBPool: pool}
}

// InsertBatchTopicA inserts a batch of messages into the database
func (msg *MessageRepo) InsertBatch(messages []models.MessageModel) (int64, error) {
	if len(messages) == 0 {
		return 0, nil
	}

	// Construct the base command
	baseCommand := "INSERT INTO QueueManager.MessageQueues (Id, SystemOwner, Topic, Content, Remarks, CreatedAt, CreatedBy, Sign) VALUES "

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
	result, err := msg.DBPool.Insert(command, args...)
	if err != nil {
		log.Printf("Insert error: %v", err)
		return 0, err
	}

	return result, nil
}

// Insert Message inserts a new role into the database
func (msg *MessageRepo) Insert(message models.MessageModel) (int64, error) {
	// Example usage
	command := "INSERT INTO QueueManager.MessageQueues (Id, SystemOwner, Topic, Content, Remarks, CreatedAt, CreatedBy, Sign) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := msg.DBPool.Insert(command, message.Id, message.System, message.Topic, message.Content, message.Remark, message.CreatedAt, message.CreatedBy, message.Sign)
	if err != nil {
		log.Printf("Insert error: %v", err)
	}
	return result, nil
}

// getMessages fetches data from QueueManager.MessageQueues table with limit and offset
func (msg *MessageRepo) GetMessages(limit, offset int) ([]models.MessageModel, error) {
	// Calculate the timestamp for "now - 2 hours"
	twoHoursAgo := time.Now().Add(-2 * time.Hour).Format("2006-01-02 15:04:05")

	// Query to get data with limit, offset and CreatedAt condition
	query := `SELECT 
		Id, 
		SystemOwner, 
		Topic, 
		Content, 
		Remarks, 
		CreatedAt, 
		CreatedBy, 
		Sign 
	FROM QueueManager.MessageQueues final 
	WHERE CreatedAt >= ?
	LIMIT ? OFFSET ?`

	rows, err := msg.DBPool.Query(query, twoHoursAgo, limit, offset)

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

func (msg *MessageRepo) GetTopics(limit int) ([]string, error) {
	query := `select distinct Topic from ( Select Topic from  QueueManager.MessageQueues final limit ?)`
	rows, err := msg.DBPool.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	var topics []string
	for rows.Next() {
		topic := ""
		if err := rows.Scan(
			&topic,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		topics = append(topics, topic)
	}
	return topics, nil
}

func transformMessageToArchive(message models.MessageModel) (*models.ArchiveModel, error) {
	id, error := uuid.NewV7()
	if error != nil {
		id = uuid.New()
	}

	var archive models.ArchiveModel
	jsonData, error := json.Marshal(message)
	if error != nil {
		return nil, error
	}
	archive = models.ArchiveModel{
		Id:            id,         // New unique ID for Archive record
		CorrelationId: message.Id, // Using MessageQueue Id as CorrelationId
		TableName:     "MessageQueues",
		JsonData:      string(jsonData),
		CreatedAt:     message.CreatedAt,
		CreatedBy:     message.CreatedBy,
	}
	return &archive, nil
}

func (msg *MessageRepo) InsertToArchive(messages []models.MessageModel) (int64, error) {
	var archiveRecords []models.ArchiveModel
	for _, message := range messages {
		archiveRecord, err := transformMessageToArchive(message)
		if err != nil {
			fmt.Printf("Error transform message. %s", err)
			continue
		}
		archiveRecords = append(archiveRecords, *archiveRecord)
	}
	if len(archiveRecords) == 0 {
		return -1, fmt.Errorf("error no archive data")
	}

	// Construct the base command
	command := `INSERT INTO QueueManager.Archive 
		( Id,CorrelationId,TableName,JsonData, CreatedAt,CreatedBy) VALUES `

	// Prepare placeholders and arguments
	var placeholders []string
	var args []interface{}

	for _, archive := range archiveRecords {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?)")
		args = append(args, archive.Id, archive.CorrelationId, archive.TableName, archive.CreatedAt, archive.CreatedBy)
	}

	// Join placeholders and form the final command
	command = command + strings.Join(placeholders, ", ")

	// Execute the command
	result, err := msg.DBPool.Insert(command, args...)
	if err != nil {
		log.Printf("Insert error: %v", err)
		return 0, err
	}
	return result, nil

}

func (msg *MessageRepo) ClearMessagesFromQueue(messageIds []string) error {
	query := fmt.Sprintf(`ALTER TABLE QueueManager.MessageQueues DELETE WHERE Id IN (%s)`, strings.Join(messageIds, ","))
	err := msg.DBPool.Update(query, nil)
	return err
}

func (msg *MessageRepo) FlagDeleteQueueMessages(messages []models.MessageModel) error {
	// Convert the UUIDs to strings
	uuidStrings := make([]string, len(messages))
	for i, message := range messages {
		uuidStrings[i] = fmt.Sprintf("'%s'", message.Id.String())
	}

	// Create the parameterized query
	query := fmt.Sprintf(`
		INSERT INTO QueueManager.MessageQueues
		(Id, SystemOwner, Topic, Content, Remarks, CreatedAt, CreatedBy, Sign)
		SELECT Id, SystemOwner, Topic, Content, Remarks, CreatedAt, CreatedBy, -1 as Sign
		FROM QueueManager.MessageQueues
		WHERE Id IN (%s)
	`, strings.Join(uuidStrings, ", "))

	// Execute the query
	err := msg.DBPool.Delete(query)
	if err != nil {
		log.Printf("Failed to execute insert-select query: %v", err)
		return err
	}

	return nil
}
