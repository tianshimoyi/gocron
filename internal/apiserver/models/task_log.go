package models

import (
	"context"
	"time"
)

type TaskLogStore interface {
	Create(ctx context.Context, taskLog *TaskLog) (insertId int64, err error)
	Update(ctx context.Context, taskLog *TaskLog) error
	List(ctx context.Context, param ListTaskLogParam) ([]*TaskLog, int64, error)
	Delete(ctx context.Context, param DeleteTaskLogParam) error
}

type TaskType int8

// 任务执行日志
type TaskLog struct {
	Id         uint       `json:"id,omitempty" xorm:"pk autoincr index"`
	TaskId     int        `json:"task_id" xorm:"int notnull index default 0"`                // 任务id
	Name       string     `json:"name" xorm:"varchar(32) notnull"`                           // 任务名称
	Spec       string     `json:"spec" xorm:"varchar(64) notnull"`                           // crontab
	Protocol   string     `json:"protocol" xorm:"varchar(32) notnull index"`                 // 协议 1:http 2:RPC
	Command    string     `json:"command" xorm:"varchar(2048) notnull"`                      // URL地址或shell命令
	Timeout    int        `json:"timeout" xorm:"mediumint notnull default 0"`                // 任务执行超时时间(单位秒),0不限制
	RetryTimes int8       `json:"retry_times" xorm:"tinyint notnull default 0"`              // 任务重试次数
	Hostname   string     `json:"hostname" xorm:"varchar(128) notnull default '' "`          // RPC主机名，逗号分隔
	StartTime  *time.Time `json:"start_time" xorm:"datetime created"`                        // 开始执行时间
	EndTime    *time.Time `json:"end_time" xorm:"datetime updated"`                          // 执行完成（失败）时间
	Status     string     `json:"status" xorm:"varchar(32) notnull index default 'running'"` // 状态 0:执行失败 1:执行中  2:执行完毕 3:任务取消(上次任务未执行完成) 4:异步执行
	Result     string     `json:"result" xorm:"mediumtext notnull "`                         // 执行结果
	TotalTime  int        `json:"total_time" xorm:"-"`                                       // 执行总时长
}

type ListTaskLogParam struct {
	BaseListParam
	Status   string
	Protocol string
	TaskID   int
}

type DeleteTaskLogParam struct {
	TaskID int
	Status string
	Mon    int
}
