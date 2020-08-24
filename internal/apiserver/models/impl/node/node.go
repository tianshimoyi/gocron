package node

import (
	"context"
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

func (t *hostStore) Create(ctx context.Context, host *models.Host) error {
	_, err := t.db.Insert(host)
	if err != nil {
		return err
	}
	return nil
}

func (t *hostStore) Update(ctx context.Context, host *models.Host) error {
	_, err := t.db.ID(host.Id).Update(host)
	if err != nil {
		return err
	}
	return nil
}

func (t *hostStore) Delete(ctx context.Context, param models.DeleteParam) error {
	s := t.db.Table(new(models.Host))
	if param.ID > 0 {
		s = s.ID(param.ID)
	}
	if param.Name != "" {
		s = s.Where("name = ?", param.Name)
	}
	_, err := s.Delete(new(models.Host))
	if err != nil {
		return err
	}
	return nil
}

func (t *hostStore) List(ctx context.Context, param models.ListHostParam) ([]*models.Host, int64, error) {
	s := t.db.Table(new(models.Host))
	if param.ID > 0 {
		s = s.Where("id = ?", param.ID)
	}
	if param.Name != "" {
		s = s.And("name LIKE ?", param.Name)
	}
	if param.SortKey != "" {
		s = s.OrderBy(param.SortKey)
	}
	if param.Reverse {
		s = s.Asc("id")

	} else {
		s = s.Desc("id")
	}
	list := make([]*models.Host, 0)
	cnt, err := s.Limit(param.Limit, param.Offset).FindAndCount(&list)
	if err != nil {
		return nil, 0, err
	}
	return list, cnt, nil
}

func (t *hostStore) Get(ctx context.Context, param models.GetParam) (*models.Host, error) {
	s := t.db.Table(new(models.Host))
	if param.ID > 0 {
		s = s.Where("id = ?", param.ID)
	}
	if param.Name != "" {
		s = s.And("name = ?", param.Name)
	}
	host := models.Host{}
	_, err := s.Get(&host)
	if err != nil {
		return &host, err
	}
	return &host, nil
}

func (t *hostStore) Exist(ctx context.Context, param models.GetParam) (bool, error) {
	s := t.db.Table(new(models.Host))
	if param.ID > 0 {
		s = s.Where("id = ?", param.ID)
	}
	if param.Name != "" {
		s = s.And("name = ?", param.Name)
	}
	return s.Exist(new(models.Host))
}
