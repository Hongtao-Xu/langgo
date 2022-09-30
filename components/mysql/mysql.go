package mysql

import (
	"time"

	"github.com/Hongtao-Xu/langgo/core"
	"github.com/Hongtao-Xu/langgo/core/log"

	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const name = "mysql" //模块名称

type Instance struct {
}

var instance *Instance

var connections = make(map[string]*gorm.DB)

//mysql连接配置信息
type item struct {
	Dsn                       string              `yaml:"dsn"`
	MaxIdleConns              int                 `yaml:"max_idle_conns"`
	MaxOpenConns              int                 `yaml:"max_open_conns"`
	ConnMaxLifetime           time.Duration       `yaml:"conn_max_lifetime"`
	SlowThreshold             time.Duration       `yaml:"slow_threshold"`
	Colorful                  bool                `yaml:"colorful"`
	IgnoreRecordNotFoundError bool                `yaml:"ignore_record_not_found_error"`
	LogLevel                  gormLogger.LogLevel `yaml:"log_level"`
}

func (i *Instance) GetName() string {
	return name
}

func (i *Instance) Load() error {
	instance = i
	//配置信息
	items := make(map[string]item)

	//加载yml配置信息
	core.GetComponentConfiguration(name, &items)

	for k, c := range items {
		zl := log.Logger("mysql", k)
		//构造sql日志配置类
		l := New(*zl, gormLogger.Config{
			SlowThreshold:             c.SlowThreshold,
			Colorful:                  c.Colorful,
			IgnoreRecordNotFoundError: c.IgnoreRecordNotFoundError,
			LogLevel:                  c.LogLevel,
		})
		//默认关闭事务
		conn, err := gorm.Open(gormMysql.Open(c.Dsn), &gorm.Config{Logger: l, SkipDefaultTransaction: true})
		if err != nil {
			log.Logger("component", "mysql").Warn().Err(err).Send()
			continue
		}

		sqlDB, err := conn.DB()
		if err != nil {
			log.Logger("component", "mysql").Warn().Err(err).Send()
			continue
		}
		sqlDB.SetMaxIdleConns(c.MaxIdleConns)
		sqlDB.SetMaxOpenConns(c.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime)
		//放入连接池
		connections[k] = conn
	}
	return nil
}

//Main 获取名为"main"的连接
func Main() *gorm.DB {
	return Get("main")
}

func Get(name string) *gorm.DB {
	conn, ok := connections[name]
	if ok {
		return conn
	} else {
		return nil
	}
}
