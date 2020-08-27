package models

type TaskHost struct {
	Model  `xorm:"extends"`
	TaskId int `json:"task_id,omitempty" xorm:"not null index"`
	HostId int `json:"host_id,omitempty" xorm:"not null index"`
}

type TaskHostDetail struct {
	TaskHost `xorm:"extends"`
	Name     string `json:"name,omitempty"`
	Port     int    `json:"port,omitempty"`
	Alias    string `json:"alias,omitempty"`
	Addr     string `json:"addr,omitempty"`
}

func (TaskHostDetail) TableName() string {
	return "g_" + "task_host"
}
