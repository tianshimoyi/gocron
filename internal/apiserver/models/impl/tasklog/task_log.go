package tasklog

import (
	"github.com/go-xorm/xorm"
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/pkg/client/database"
)

func New(client *database.Client) models.TaskLogStore {
	return &taskLogStore{db: client.DB()}
}

type taskLogStore struct {
	db *xorm.Engine
}

func (t *taskLogStore) Create() (insertId int64, err error) {
	return 0, nil
}
