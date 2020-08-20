package models

// 主机
type Host struct {
	Id        int16  `json:"id" xorm:"smallint pk autoincr"`
	Name      string `json:"name" xorm:"varchar(64) notnull"`                // 主机名称
	Alias     string `json:"alias" xorm:"varchar(32) notnull default '' "`   // 主机别名
	Port      int    `json:"port" xorm:"notnull default 5921"`               // 主机端口
	Remark    string `json:"remark" xorm:"varchar(100) notnull default '' "` // 备注
	BaseModel `json:"-" xorm:"-"`
	Selected  bool `json:"-" xorm:"-"`
}

type HostStore interface {
	Create() (insertId int16, err error)
}
