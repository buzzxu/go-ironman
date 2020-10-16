package ironman

import (
	"fmt"
	"github.com/buzzxu/ironman/logger"
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
func CreateDB(callback func(db *gorm.DB)) *gorm.DB {
	dbConfig := conf.ServerConf.DataSource
	db, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DB,
	))

	if err != nil {
		logger.Fatalf("Failed to connect to log mysql: %s", err.Error())
	}
	db.DB().SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Hour)
	callback(db)
	db.DB().Ping()
	db.LogMode(dbConfig.Log)
	return db
}

//DataSourceConnect 初始化数据库链接
func DataSourceConnect(callback func(db *gorm.DB)) {
	Db = CreateDB(callback)
}

func DataSourceClose() {
	Db.Close()
}
