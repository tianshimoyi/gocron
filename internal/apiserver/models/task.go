package models

import "time"

type TaskProtocol int8

const (
	TaskHTTP TaskProtocol = iota + 1 // HTTP协议
	TaskRPC                          // RPC方式执行命令
)

type TaskLevel int8

const (
	TaskLevelParent TaskLevel = 1 // 父任务
	TaskLevelChild  TaskLevel = 2 // 子任务(依赖任务)
)

type TaskDependencyStatus int8

const (
	TaskDependencyStatusStrong TaskDependencyStatus = 1 // 强依赖
	TaskDependencyStatusWeak   TaskDependencyStatus = 2 // 弱依赖
)

type TaskHTTPMethod int8

const (
	TaskHTTPMethodGet  TaskHTTPMethod = 1
	TaskHttpMethodPost TaskHTTPMethod = 2
)

// 任务
type Task struct {
	Id               int                  `json:"id" xorm:"int pk autoincr"`
	Name             string               `json:"name" xorm:"varchar(32) notnull"`                            // 任务名称
	Level            TaskLevel            `json:"level" xorm:"tinyint notnull index default 1"`               // 任务等级 1: 主任务 2: 依赖任务
	DependencyTaskId string               `json:"dependency_task_id" xorm:"varchar(64) notnull default ''"`   // 依赖任务ID,多个ID逗号分隔
	DependencyStatus TaskDependencyStatus `json:"dependency_status" xorm:"tinyint notnull default 1"`         // 依赖关系 1:强依赖 主任务执行成功, 依赖任务才会被执行 2:弱依赖
	Spec             string               `json:"spec" xorm:"varchar(64) notnull"`                            // crontab
	Protocol         TaskProtocol         `json:"protocol" xorm:"tinyint notnull index"`                      // 协议 1:http 2:系统命令
	Command          string               `json:"command" xorm:"varchar(256) notnull"`                        // URL地址或shell命令
	HttpMethod       TaskHTTPMethod       `json:"http_method" xorm:"tinyint notnull default 1"`               // http请求方法
	Timeout          int                  `json:"timeout" xorm:"mediumint notnull default 0"`                 // 任务执行超时时间(单位秒),0不限制
	Multi            int8                 `json:"multi" xorm:"tinyint notnull default 1"`                     // 是否允许多实例运行
	RetryTimes       int8                 `json:"retry_times" xorm:"tinyint notnull default 0"`               // 重试次数
	RetryInterval    int16                `json:"retry_interval" xorm:"smallint notnull default 0"`           // 重试间隔时间
	NotifyStatus     int8                 `json:"notify_status" xorm:"tinyint notnull default 1"`             // 任务执行结束是否通知 0: 不通知 1: 失败通知 2: 执行结束通知 3: 任务执行结果关键字匹配通知
	NotifyType       int8                 `json:"notify_type" xorm:"tinyint notnull default 0"`               // 通知类型 1: 邮件 2: slack 3: webhook
	NotifyReceiverId string               `json:"notify_receiver_id" xorm:"varchar(256) notnull default '' "` // 通知接受者ID, setting表主键ID，多个ID逗号分隔
	NotifyKeyword    string               `json:"notify_keyword" xorm:"varchar(128) notnull default '' "`
	Tag              string               `json:"tag" xorm:"varchar(32) notnull default ''"`
	Remark           string               `json:"remark" xorm:"varchar(100) notnull default ''"` // 备注
	Status           Status               `json:"status" xorm:"tinyint notnull index default 0"` // 状态 1:正常 0:停止
	Created          time.Time            `json:"created" xorm:"datetime notnull created"`       // 创建时间
	Deleted          time.Time            `json:"deleted" xorm:"datetime deleted"`               // 删除时间
	BaseModel        `json:"-" xorm:"-"`
	Hosts            []TaskHostDetail `json:"hosts" xorm:"-"`
	NextRunTime      time.Time        `json:"next_run_time" xorm:"-"`
}

type TaskStore interface {
	Create() (insertId int, err error)
	//UpdateBean(id int) (int64, error)
	//Update(id int, data CommonMap) (int64, error)
	//Delete(id int) (int64, error)
	//Disable(id int) (int64, error)
	//Enable(id int) (int64, error)
	//ActiveList(page, pageSize int) ([]Task, error)
	//ActiveListByHostId(hostId int16) ([]Task, error)
}
