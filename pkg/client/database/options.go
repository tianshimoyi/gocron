package database

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/x893675/gocron/pkg/utils/reflectutils"
	"time"
)

const defaultTablePrefix = "g_"

// database options
type Options struct {
	Host                  string        `json:"host,omitempty" yaml:"host" description:"db service host address"`
	Username              string        `json:"username,omitempty" yaml:"username"`
	Password              string        `json:"-" yaml:"password"`
	Type                  string        `json:"type" yaml:"type"`
	DBName                string        `json:"dbName" yaml:"dbName"`
	Debug                 bool          `json:"debug" yaml:"debug"`
	Port                  string        `json:"port" yaml:"port"`
	MaxIdleConnections    int           `json:"maxIdleConnections,omitempty" yaml:"maxIdleConnections"`
	MaxOpenConnections    int           `json:"maxOpenConnections,omitempty" yaml:"maxOpenConnections"`
	MaxConnectionLifeTime time.Duration `json:"maxConnectionLifeTime,omitempty" yaml:"maxConnectionLifeTime"`
}

// NewMySQLOptions create a `zero` value instance
func NewDatabaseOptions() *Options {
	return &Options{
		Host:                  "",
		Username:              "",
		Password:              "",
		DBName:                "",
		Type:                  "postgres",
		Port:                  "5432",
		Debug:                 false,
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
	}
}

func (m *Options) Validate() []error {
	var errors []error

	return errors
}

func (m *Options) ApplyTo(options *Options) {
	reflectutils.Override(options, m)
}

func (m *Options) AddFlags(fs *pflag.FlagSet, c *Options) {

	fs.StringVar(&m.Host, "db-host", c.Host, ""+
		"Database service host address. If left blank, the following related db options will be ignored.")

	fs.StringVar(&m.Port, "db-port", c.Port, ""+
		"Database service port. If left blank, the following related db options will be ignored.")

	fs.StringVar(&m.Username, "db-username", c.Username, ""+
		"Username for access to db service.")

	fs.StringVar(&m.Password, "db-password", c.Password, ""+
		"Password for access to db, should be used pair with password.")

	fs.StringVar(&m.Type, "db-type", c.Type, ""+
		"Database type. eg: sqlite3,mysql,postgres, default is postgres")

	fs.IntVar(&m.MaxIdleConnections, "db-max-idle-connections", c.MaxOpenConnections, ""+
		"Maximum idle connections allowed to connect to db.")

	fs.IntVar(&m.MaxOpenConnections, "db-max-open-connections", c.MaxOpenConnections, ""+
		"Maximum open connections allowed to connect to db.")

	fs.DurationVar(&m.MaxConnectionLifeTime, "db-max-connection-life-time", c.MaxConnectionLifeTime, ""+
		"Maximum connection life time allowed to connecto to db.")

	fs.BoolVar(&m.Debug, "db-debug", c.Debug, ""+
		"enable / disable database log")
}

func (m *Options) GetDSN() string {
	switch m.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true", m.Username, m.Password, m.Host, m.Port, m.DBName)
	case "sqlite3":
		return m.DBName
	case "postgres":
		fallthrough
	default:
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", m.Host, m.Port, m.Username, m.DBName, m.Password)
	}
}
