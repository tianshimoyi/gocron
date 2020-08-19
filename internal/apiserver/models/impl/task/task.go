package task

import (
	"github.com/go-xorm/xorm"
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/pkg/client/database"
)

func New(client *database.Client) models.TaskStore {
	return &taskStore{db: client.DB()}
}

type taskStore struct {
	db *xorm.Engine
}

func (t *taskStore) Create() (insertId int, err error) {
	return 0, nil
}
