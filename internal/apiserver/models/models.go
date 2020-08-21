package models

import "time"

type Status int8
type CommonMap map[string]interface{}

const (
	Disabled Status = 0 // 禁用
	Failure  Status = 0 // 失败
	Enabled  Status = 1 // 启用
	Running  Status = 1 // 运行中
	Finish   Status = 2 // 完成
	Cancel   Status = 3 // 取消
)

const (
	Page        = 1    // 当前页数
	PageSize    = 20   // 每页多少条数据
	MaxPageSize = 1000 // 每次最多取多少条
)

const DefaultTimeFormat = "2006-01-02 15:04:05"

const (
	dbPingInterval = 90 * time.Second
	dbMaxLiftTime  = 2 * time.Hour
)

type Model struct {
	ID        uint      `xorm:"pk autoincr index"`
	CreatedAt time.Time `xorm:"created index"`
	UpdatedAt time.Time `xorm:"updated index"`
	//DeletedAt *time.Time `xorm:"deleted"`
}

type BaseModel struct {
	Page     int `xorm:"-"`
	PageSize int `xorm:"-"`
}

func (model *BaseModel) parsePageAndPageSize(params CommonMap) {
	page, ok := params["Page"]
	if ok {
		model.Page = page.(int)
	}
	pageSize, ok := params["PageSize"]
	if ok {
		model.PageSize = pageSize.(int)
	}
	if model.Page <= 0 {
		model.Page = Page
	}
	if model.PageSize <= 0 {
		model.PageSize = MaxPageSize
	}
}

func (model *BaseModel) pageLimitOffset() int {
	return (model.Page - 1) * model.PageSize
}
