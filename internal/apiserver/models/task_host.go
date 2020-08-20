package models

type TaskHost struct {
	Id     int   `json:"id" xorm:"int pk autoincr"`
	TaskId int   `json:"task_id" xorm:"int not null index"`
	HostId int16 `json:"host_id" xorm:"smallint not null index"`
}

type TaskHostDetail struct {
	TaskHost `xorm:"extends"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Alias    string `json:"alias"`
}
