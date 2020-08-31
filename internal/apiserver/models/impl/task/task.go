package task

import (
	"context"
	"github.com/go-xorm/xorm"
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/pkg/client/database"
	"k8s.io/klog/v2"
	"sync"
	"time"
)

func New(client *database.Client) models.TaskStore {
	return &taskStore{db: client.DB()}
}

type taskStore struct {
	db *xorm.Engine
}

func (t *taskStore) Create(ctx context.Context, tsk models.SchemaTask) (uint, error) {
	session := t.db.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return 0, err
	}
	item := tsk.ToModelTask()
	if _, err := session.Insert(item); err != nil {
		return 0, err
	}
	klog.V(2).Infof("create task id is %d", item.Id)
	if item.Protocol == models.TaskShell {
		itemHost := tsk.ToModelTaskHosts(int(item.Id))
		if _, err := session.Insert(&itemHost); err != nil {
			return 0, err
		}
	}
	if err := session.Commit(); err != nil {
		return 0, err
	}
	return item.Id, session.Commit()
}

func (t *taskStore) Delete(ctx context.Context, param models.DeleteParam) error {
	session := t.db.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	if _, err := session.ID(param.ID).Delete(new(models.Task)); err != nil {
		return err
	}
	if _, err := session.Where("task_id = ?", param.ID).Delete(new(models.TaskHost)); err != nil {
		return err
	}
	return session.Commit()
}

func (t *taskStore) List(ctx context.Context, param models.ListTaskParam) ([]*models.Task, int64, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	defer close(errChan)
	var total int64
	var err error
	wg.Add(1)
	go func() {
		defer wg.Done()
		total, err = t.total(param)
		if err != nil {
			errChan <- err
			return
		}
	}()
	list := make([]*models.Task, 0)
	session := t.db.Alias("t").Join("LEFT", []string{"g_task_host", "th"}, "t.id = th.task_id")
	t.parseListCondition(session, param)
	session.GroupBy("t.id")
	if param.Reverse {
		session.Asc("t.id").Cols("t.*").Limit(param.Limit, param.Offset)
	} else {
		session.Desc("t.id").Cols("t.*").Limit(param.Limit, param.Offset)
	}
	err = session.Find(&list)
	if err != nil {
		return nil, 0, err
	}
	list, err = t.setHostsForTasks(ctx, list)
	if err != nil {
		return nil, 0, err
	}
	wg.Wait()
	if len(errChan) != 0 {
		return nil, 0, <-errChan
	}
	return list, total, nil
}

func (t *taskStore) total(param models.ListTaskParam) (int64, error) {
	list := make([]*models.Task, 0)
	session := t.db.Alias("t").Join("LEFT", []string{"g_task_host", "th"}, "t.id = th.task_id")
	t.parseListCondition(session, param)
	session.GroupBy("t.id")
	if param.Reverse {
		session.Asc("t.id").Cols("t.*").Limit(param.Limit, param.Offset)
	} else {
		session.Desc("t.id").Cols("t.*").Limit(param.Limit, param.Offset)
	}
	err := session.Find(&list)
	return int64(len(list)), err
}

func (t *taskStore) parseListCondition(session *xorm.Session, param models.ListTaskParam) {
	if param.ID > 0 {
		session.And("t.id = ?", param.ID)
	}
	if param.HostID > 0 {
		session.And("th.host_id = ?", param.HostID)
	}
	if param.Name != "" {
		session.And("t.name LIKE ?", "%"+param.Name+"%")
	}
	if param.Protocol != "" {
		session.And("protocol = ?", param.Protocol)
	}
	if param.Status != "" {
		session.And("status = ?", param.Status)
	}
	if param.Tag != "" {
		session.And("tag = ?", param.Tag)
	}
	if param.Level != "" {
		session.And("level = ?", param.Level)
	}
	if param.Type != "" {
		session.And("type = ?", param.Type)
	}
	if param.RunAtInterval != 0 {
		t := time.Now().UTC()
		//tAfter := t.Add(param.RunAtInterval
		session.And("run_at >= ?", t.Format(models.DefaultTimeFormat))
		session.And("run_at <= ?", t.Add(param.RunAtInterval).Format(models.DefaultTimeFormat))
	}
}

func (t *taskStore) setHostsForTasks(ctx context.Context, tasks []*models.Task) ([]*models.Task, error) {
	var err error
	for i, value := range tasks {
		taskHostDetails, err := t.GetTaskHostByTaskID(ctx, value.Id)
		if err != nil {
			return nil, err
		}
		tasks[i].Hosts = taskHostDetails
	}
	return tasks, err
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
	fields := "th.id,th.host_id,h.alias,h.name,h.port,h.addr"
	err := t.db.Alias("th").
		Join("LEFT", []string{"g_host", "h"}, "th.host_id=h.id").
		Where("th.task_id = ?", taskID).
		Cols(fields).
		Find(&list)
	klog.V(2).Infof("task host detail is %v", list)
	return list, err
}
