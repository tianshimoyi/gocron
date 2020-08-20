package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"k8s.io/klog/v2"
	"xorm.io/core"
)

type Client struct {
	db *xorm.Engine
}

func NewDatabaseClient(options *Options, stopCh <-chan struct{}) (*Client, error) {
	engine, err := xorm.NewEngine(options.Type, options.GetDSN())
	if err != nil {
		klog.Error("unable to connect to database", err)
		return nil, err
	}
	engine.SetMaxIdleConns(options.MaxIdleConnections)
	engine.SetMaxOpenConns(options.MaxOpenConnections)
	engine.SetConnMaxLifetime(options.MaxConnectionLifeTime)

	mapper := core.NewPrefixMapper(core.SnakeMapper{}, defaultTablePrefix)
	engine.SetTableMapper(mapper)

	if options.Debug {
		engine.ShowSQL(true)
		engine.Logger().SetLevel(core.LOG_DEBUG)
	}

	go func() {
		<-stopCh
		if err := engine.Close(); err != nil {
			klog.Warning("error happened during closing database connection", err)
		}
	}()

	return &Client{db: engine}, nil
}

func NewDatabaseClientOrDie(options *Options, stopCh <-chan struct{}) *Client {
	c, err := NewDatabaseClient(options, stopCh)
	if err != nil {
		klog.Error("unable to init database", err)
		panic(err)
	}
	return c
}

func (c *Client) DB() *xorm.Engine {
	if c == nil {
		klog.Warning("database client is nil")
		return nil
	}
	return c.db
}
