package models

import (
	//"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type User struct {
	gorm.Model
	Name     string
	Password string
	Token    string
	Status   int
	Flux     uint64
	ExpireAt time.Time
}

func init() {
}

func InitDB() *gorm.DB {
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 使用用彩色打印
		})
	db, err := gorm.Open(mysql.Open("root:12345@tcp/demo?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// User:
// id
// name
// password
// token
// status
// flux

const (
	Normal   = iota
	UnNormal = iota
)
