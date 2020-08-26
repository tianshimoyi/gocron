package models

import (
	"context"
	"github.com/x893675/gocron/internal/apiserver/schema"
)

// 主机
type Host struct {
	Model  `xorm:"extends"`
	Name   string `json:"name,omitempty" xorm:"varchar(64) notnull"`
	Alias  string `json:"alias,omitempty" xorm:"varchar(32) notnull default '' "`   // 主机别名
	Port   int    `json:"port,omitempty" xorm:"notnull default 5921"`               // 主机端口
	Remark string `json:"remark,omitempty" xorm:"varchar(100) notnull default '' "` // 备注
	Addr   string `json:"addr,omitempty" xorm:"varchar(100) notnull default '' "`
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

type SchemaHost schema.NodeRequest

func (s SchemaHost) ToModelHost() *Host {
	return &Host{
		Name:   s.Name,
		Alias:  s.Alias,
		Port:   s.Port,
		Remark: s.Remark,
		Addr:   s.Address,
	}
}
