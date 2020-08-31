package tasklog

import (
	"context"
	"github.com/go-xorm/xorm"
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/pkg/client/database"
	"time"
)

func New(client *database.Client) models.TaskLogStore {
	return &taskLogStore{db: client.DB()}
}

type taskLogStore struct {
	db *xorm.Engine
}

func (t *taskLogStore) Create(ctx context.Context, taskLog *models.TaskLog) (insertId int64, err error) {
	_, err = t.db.Insert(taskLog)
	if err == nil {
		insertId = int64(taskLog.Id)
	}
	return
}

func (t *taskLogStore) Update(ctx context.Context, taskLog *models.TaskLog) error {
	_, err := t.db.ID(taskLog.Id).Update(taskLog)
	if err != nil {
		return err
	}
	return nil
}

func (t *taskLogStore) List(ctx context.Context, param models.ListTaskLogParam) ([]*models.TaskLog, int64, error) {
	session := t.db.Table(new(models.TaskLog))
	t.parseListCondition(session, param)
	session.Limit(param.Limit, param.Offset)
	if param.Reverse {
		session.Asc("id").Limit(param.Limit, param.Offset)
	} else {
		session.Desc("id").Limit(param.Limit, param.Offset)
	}
	list := make([]*models.TaskLog, 0)
	total, err := session.FindAndCount(&list)
	if err != nil {
		return nil, 0, err
	}
	if len(list) > 0 {
		for i, item := range list {
			endTime := item.EndTime
			if item.Status == models.TaskLogStatusRunning {
				ts := time.Now()
				endTime = &ts
			}
			execSeconds := endTime.Sub(*item.StartTime).Seconds()
			list[i].TotalTime = int(execSeconds)
		}
	}
	return list, total, nil
}

func (t *taskLogStore) parseListCondition(session *xorm.Session, param models.ListTaskLogParam) {
	if param.Protocol != "" {
		session.And("protocol = ?", param.Protocol)
	}
	if param.Status != "" {
		session.And("status = ?", param.Status)
	}
	if param.TaskID != 0 {
		session.And("task_id = ?", param.TaskID)
	}
}

func (t *taskLogStore) Delete(ctx context.Context, param models.DeleteTaskLogParam) error {
	session := t.db.Where("task_id = ?", param.TaskID)
	if param.Status != "" {
		session.And("status = ?", param.Status)
	}
	if param.Mon != 0 {
		t := time.Now().UTC().AddDate(0, -param.Mon, 0)
		session.And("start_time <= ?", t.Format(models.DefaultTimeFormat))
	}
	_, err := session.Delete(new(models.TaskLog))
	if err != nil {
		return err
	}
	return nil
}
