package node

import (
	"github.com/go-xorm/xorm"
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/pkg/client/database"
)

func New(client *database.Client) models.HostStore {
	return &hostStore{db: client.DB()}
}

type hostStore struct {
	db *xorm.Engine
}

func (t *hostStore) Create() (insertId int16, err error) {
	return 0, nil
}
