package models

import (
	"context"
	"github.com/x893675/gocron/internal/apiserver/schema"
	"strconv"
	"strings"
	"time"
)

const (
	ParentLevelTask            = "parent"
	ChildLevelTask             = "child"
	TaskDependencyStatusStrong = "strong" // 强依赖
	TaskDependencyStatusWeak   = "weak"   // 弱依赖
	TaskHTTP                   = "http"   // HTTP协议
	TaskShell                  = "shell"  // shell方式执行命令
	TaskHTTPMethodGet          = "get"
	TaskHttpMethodPost         = "post"
	TaskStatusDisabled         = "disabled"
	TaskStatusEnabled          = "enabled"
	TaskLogStatusFailure       = "failure"
	TaskLogStatusRunning       = "running"
	TaskLogStatusFinish        = "finish"
	TaskLogStatusCancel        = "cancel"
	TaskTypeJob                = "job"
	TaskTypeCronJob            = "cronjob"
	TaskTypePlanJob            = "planjob"
)

type Task struct {
	Model            `xorm:"extends"`
	Name             string           `json:"name,omitempty" xorm:"varchar(32) notnull"`                               // 任务名称
	Level            string           `json:"level,omitempty" xorm:"varchar(32) notnull default 'parent'"`             // 任务等级 parent: 主任务, child: 依赖任务
	DependencyTaskId string           `json:"dependency_task_id,omitempty" xorm:"varchar(64) notnull default ''"`      // 依赖任务ID,多个ID逗号分隔
	DependencyStatus string           `json:"dependency_status,omitempty" xorm:"varchar(32) notnull default 'strong'"` // 依赖关系 strong:强依赖 主任务执行成功, 依赖任务才会被执行, weak:弱依赖
	Spec             string           `json:"spec,omitempty" xorm:"varchar(64) notnull"`                               // crontab 表达式
	Protocol         string           `json:"protocol,omitempty" xorm:"varchar(32) notnull index"`                     // 协议 1:http 2:系统命令
	Command          string           `json:"command,omitempty" xorm:"varchar(1024) notnull"`                          // URL地址或shell命令
	HttpMethod       string           `json:"http_method,omitempty" xorm:"varchar(32) notnull default 'get'"`          // http请求方法
	Timeout          int              `json:"timeout,omitempty" xorm:"mediumint notnull default 0"`                    // 任务执行超时时间(单位秒),0不限制
	Multi            bool             `json:"multi,omitempty" xorm:"notnull default true"`                             // 是否允许多实例运行
	RetryTimes       int8             `json:"retry_times,omitempty" xorm:"tinyint notnull default 0"`                  // 重试次数
	RetryInterval    int16            `json:"retry_interval,omitempty" xorm:"smallint notnull default 0"`              // 重试间隔时间
	NotifyStatus     int8             `json:"notify_status,omitempty" xorm:"tinyint notnull default 1"`                // 任务执行结束是否通知 0: 不通知 1: 失败通知 2: 执行结束通知 3: 任务执行结果关键字匹配通知
	NotifyType       int8             `json:"notify_type,omitempty" xorm:"tinyint notnull default 0"`                  // 通知类型 1: 邮件 2: slack 3: webhook
	NotifyReceiverId string           `json:"notify_receiver_id,omitempty" xorm:"varchar(256) notnull default '' "`    // 通知接受者ID, setting表主键ID，多个ID逗号分隔
	NotifyKeyword    string           `json:"notify_keyword,omitempty" xorm:"varchar(128) notnull default '' "`
	Tag              string           `json:"tag,omitempty" xorm:"varchar(32) notnull default ''"`
	Type             string           `json:"type,omitempty" xorm:"varchar(32) notnull default 'cronjob'"`
	Remark           string           `json:"remark,omitempty" xorm:"varchar(100) notnull default ''"` // 备注
	Status           string           `json:"status,omitempty" xorm:"varchar(32) notnull index"`       // 状态 1:正常 0:停止
	RunAt            *time.Time       `json:"run_at,omitempty" xorm:"index"`
	Creator          string           `json:"creator,omitempty" xorm:"index"`
	NextRunTime      *time.Time       `json:"next_run_time,omitempty" xorm:"-"`
	Hosts            []TaskHostDetail `json:"hosts,omitempty" xorm:"-"`
}

type TaskStore interface {
	Create(context.Context, SchemaTask) (uint, error)
	//Update(context.Context, SchemaTask) error
	Delete(context.Context, DeleteParam) error
	List(context.Context, ListTaskParam) ([]*Task, int64, error)
	Get(context.Context, GetParam) (*Task, error)
	Exist(context.Context, GetParam) (bool, error)
	GetTaskHostByTaskID(context.Context, uint) ([]TaskHostDetail, error)
	UpdateTaskStatus(context.Context, int, string) error
}

type SchemaTask schema.TaskRequest

type ListTaskParam struct {
	BaseListParam
	GetParam
	Status        string
	Level         string
	HostID        int
	Protocol      string
	Tag           string
	Type          string
	Creator       string
	RunAtInterval time.Duration
}

func (s SchemaTask) ToModelTask() *Task {
	t := &Task{
		Name:             s.Name,
		Level:            s.Level,
		DependencyTaskId: s.DependencyTaskId,
		DependencyStatus: s.DependencyStatus,
		Spec:             s.Spec,
		Protocol:         s.Protocol,
		Command:          s.Command,
		HttpMethod:       s.HttpMethod,
		Timeout:          s.Timeout,
		Multi:            s.Multi,
		RetryTimes:       s.RetryTimes,
		RetryInterval:    s.RetryInterval,
		NotifyStatus:     s.NotifyStatus,
		NotifyType:       s.NotifyType,
		NotifyReceiverId: s.NotifyReceiverId,
		NotifyKeyword:    s.NotifyKeyword,
		Tag:              s.Tag,
		Type:             s.Type,
		Remark:           s.Remark,
		Creator:          s.Creator,
		Status:           TaskStatusEnabled,
	}
	if s.RunAt != nil {
		ts := s.RunAt.Time
		t.RunAt = &ts
	}
	return t
}

func (s SchemaTask) ToModelTaskHosts(taskId int) []TaskHost {

	hostIdStrList := strings.Split(s.HostId, ",")
	hostIds := make([]TaskHost, len(hostIdStrList))
	for i, value := range hostIdStrList {
		hostIds[i].TaskId = taskId
		hostIds[i].HostId, _ = strconv.Atoi(value)
	}
	return hostIds
}
