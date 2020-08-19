package schema

type TaskRequest struct {
	Level            int8   `json:"level" description:"task level" enum:"1|2" validate:"required,oneof=1 2"`                          // 任务等级 1: 主任务 2: 依赖任务
	DependencyStatus int8   `json:"dependency_status" description:"task dependency status" enum:"1|2" validate:"required,oneof=1 2"`  // 依赖关系 1:强依赖 主任务执行成功, 依赖任务才会被执行 2:弱依赖
	DependencyTaskId string `json:"dependency_task_id" description:"dependency task id" optional:"true"`                              // 依赖任务ID,多个ID逗号分隔
	Name             string `json:"name" description:"task name" validate:"required"`                                                 // 任务名称
	Spec             string `json:"spec" description:"crontab expression" validate:"required"`                                        // crontab
	Protocol         int8   `json:"protocol" description:"task protocol type,shell or http" enum:"1|2" validate:"required,oneof=1 2"` // 协议 1:http 2:系统命令
	Command          string `json:"command" description:"shell command or http url" validate:"required"`
	HttpMethod       int8   `json:"http_method" description:"get or post" enum:"1|2" optional:"true"`                         // http请求方法
	Timeout          int    `json:"timeout" description:"task timeout, second, 0 is not limit" optional:"true"`               // 任务执行超时时间(单位秒),0不限制
	Multi            int8   `json:"multi" description:"allow multi task run the same time, default is false" optional:"true"` // 是否允许多实例运行
	RetryTimes       int8   `json:"retry_times" description:"task retry times when task run failed" optional:"true"`          // 重试次数
	RetryInterval    int16  `json:"retry_interval" description:"retry interval, second" optional:"true"`                      // 重试间隔时间
	HostId           string `json:"hosts" description:"host ids" validate:"true"`                                             //任务运行HOST,多个ID逗号分隔
	Tag              string `json:"tag" description:"task tag" optional:"true"`
	Remark           string `json:"remark" description:"task remark" optional:"true"`                  // 备注
	NotifyStatus     int8   `json:"notify_status" description:"notify status" optional:"true"`         // 任务执行结束是否通知 0: 不通知 1: 失败通知 2: 执行结束通知 3: 任务执行结果关键字匹配通知
	NotifyType       int8   `json:"notify_type"description:"notify type" optional:"true"`              // 通知类型 1: 邮件 2: slack 3: webhook
	NotifyReceiverId string `json:"notify_receiver_id" description:"notify receivers" optional:"true"` // 通知接受者ID, setting表主键ID，多个ID逗号分隔
	NotifyKeyword    string `json:"notify_keyword" description:"notify keyword" optional:"true"`
}
