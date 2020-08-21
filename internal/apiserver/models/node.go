package models

import "context"

// 主机
type Host struct {
	Model  `xorm:"extends"`
	Name   string `json:"name" xorm:"varchar(64) notnull"`
	Alias  string `json:"alias" xorm:"varchar(32) notnull default '' "`   // 主机别名
	Port   int    `json:"port" xorm:"notnull default 5921"`               // 主机端口
	Remark string `json:"remark" xorm:"varchar(100) notnull default '' "` // 备注
}

type HostStore interface {
	Create(context.Context, *Host) error
	Update(context.Context, *Host) error
	Delete(context.Context, DeleteParam) error
	List(context.Context, ListHostParam) ([]*Host, int64, error)
	Get(context.Context, GetParam) (*Host, error)
	Exist(context.Context, GetParam) (bool, error)
}

type BaseListParam struct {
	Reverse bool
	SortKey string
	Offset  int
	Limit   int
}

type ListHostParam struct {
	BaseListParam
	GetParam
}

type GetParam struct {
	ID   int
	Name string
}

type DeleteParam GetParam
