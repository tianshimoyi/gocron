package task

import (
	"context"
	"github.com/go-xorm/xorm"
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/pkg/client/database"
	"k8s.io/klog/v2"
)

func New(client *database.Client) models.TaskStore {
	return &taskStore{db: client.DB()}
}

type taskStore struct {
	db *xorm.Engine
}

func (t *taskStore) Create(ctx context.Context, tsk models.SchemaTask) error {
	session := t.db.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	item := tsk.ToModelTask()
	if _, err := session.Insert(item); err != nil {
		return err
	}
	klog.V(2).Infof("create task id is %d", item.Id)
	if item.Protocol == models.TaskShell {
		itemHost := tsk.ToModelTaskHosts(int(item.Id))
		if _, err := session.Insert(&itemHost); err != nil {
			return err
		}
	}
	return session.Commit()
}

func (t *taskStore) Delete(ctx context.Context, param models.DeleteParam) error {
	session := t.db.NewSession()
	defer session.Close()
	if _, err := session.ID(param.ID).Delete(new(models.Task)); err != nil {
		return err
	}
	if _, err := session.Where("task_id = ?", param.ID).Delete(new(models.TaskHost)); err != nil {
		return err
	}
	return session.Commit()
}

func (t *taskStore) List(ctx context.Context, param models.ListTaskParam) ([]*models.Task, int64, error) {
	return nil, 0, nil
}

func (t *taskStore) Get(ctx context.Context, param models.GetParam) (*models.Task, error) {
	s := t.db.Table(new(models.Task))
	if param.ID > 0 {
		s = s.Where("id = ?", param.ID)
	}
	if param.Name != "" {
		s = s.And("name = ?", param.Name)
	}
	item := models.Task{}
	_, err := s.Get(&item)
	if err != nil {
		return &item, err
	}
	item.Hosts, err = t.GetTaskHostByTaskID(ctx, item.Id)
	return &item, err
}

func (t *taskStore) Exist(ctx context.Context, param models.GetParam) (bool, error) {
	s := t.db.Table(new(models.Task))
	if param.ID > 0 {
		s = s.Where("id = ?", param.ID)
	}
	if param.Name != "" {
		s = s.And("name = ?", param.Name)
	}
	return s.Exist(new(models.Task))
}

func (t *taskStore) UpdateTaskStatus(ctx context.Context, id int, status string) error {
	if _, err := t.db.Table(new(models.Task)).ID(id).Update(map[string]interface{}{"status": status}); err != nil {
		return err
	}
	return nil
}

func (t *taskStore) GetTaskHostByTaskID(ctx context.Context, taskID uint) ([]models.TaskHostDetail, error) {
	list := make([]models.TaskHostDetail, 0)
	fields := "th.id,th.host_id,h.alias,h.name,h.port"
	err := t.db.Alias("th").
		Join("LEFT", []string{"g_host", "h"}, "th.host_id=h.id").
		Where("th.task_id = ?", taskID).
		Cols(fields).
		Find(&list)
	klog.V(2).Infof("task host detail is %v", list)
	return list, err
}
