package ironman

import (
	"fmt"
	"log"
	"time"

	"github.com/buzzxu/ironman/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DbConfig 数据库配置
type DbConfig struct {
	Host            string
	Port            int16
	User            string
	Password        string
	DBName          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

// Db 数据库操作
var Db *gorm.DB

func init() {

}

// CreateDB 创建数据库链接
func CreateDB(dbConfig *DbConfig) *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	))

	if err != nil {
		log.Panic(fmt.Errorf("Failed to connect to log mysql: %s", err))
	}
	db.DB().SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Hour)
	db.DB().Ping()
	db.LogMode(true)
	return db
}

//DataSourceConnect 初始化数据库链接
func DataSourceConnect() {
	dbConfig := &DbConfig{
		Host:            conf.ServerConf.DataSource.Host,
		Port:            conf.ServerConf.DataSource.Port,
		User:            conf.ServerConf.DataSource.User,
		Password:        conf.ServerConf.DataSource.Password,
		DBName:          conf.ServerConf.DataSource.DB,
		MaxIdleConns:    conf.ServerConf.DataSource.MaxIdleConns,
		MaxOpenConns:    conf.ServerConf.DataSource.MaxOpenConns,
		ConnMaxLifetime: conf.ServerConf.DataSource.ConnMaxLifetime,
	}
	Db = CreateDB(dbConfig)
}
