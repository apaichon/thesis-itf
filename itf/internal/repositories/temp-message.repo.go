package repositories

import (
	"github.com/apaichon/thesis-itf/itf/internal/data"
	"github.com/apaichon/thesis-itf/itf/internal/models"
)

type TempMessageRepo struct {
	DB *data.SqliteDB
}

// NewMessageRepo creates a new instance of MessageRepo
func NewTempMessageRepo() *TempMessageRepo {
	db := data.NewSqliteDB()
	return &TempMessageRepo{DB: db}
}

func (repo *TempMessageRepo) Insert(msg models.MessageModel) (int64, error) {
	formattedTime := msg.CreatedAt.Format("2006-01-02 15:04:05")
	result, err := repo.DB.Insert(
		`INSERT INTO messages (id, topic, content, created_at, sign)
		VALUES (?, ?, ?, ?, ?)`, msg.Id, msg.Topic, msg.Content, formattedTime, msg.Sign)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
