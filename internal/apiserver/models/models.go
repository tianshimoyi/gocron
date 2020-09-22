package models

import "time"

const DefaultTimeFormat = "2006-01-02 15:04:05"

type Model struct {
	Id        uint       `json:"id,omitempty" xorm:"pk autoincr index"`
	CreatedAt *time.Time `json:"create_time,omitempty" xorm:"created index"`
	UpdatedAt *time.Time `json:"update_time,omitempty" xorm:"updated index"`
	//DeletedAt *time.Time `xorm:"deleted"`
}
