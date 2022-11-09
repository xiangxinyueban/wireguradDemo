package model

import (
	"golang.org/x/crypto/bcrypt"
	//"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Password string
	Email    string
	Token    string
	Status   int
}

var DB *gorm.DB

func InitDB() {
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 使用用彩色打印
		})
	db, err := gorm.Open(mysql.Open("root:12345@tcp/demo?charset=utf8&parseTime=True&loc=Local"),
		&gorm.Config{
			Logger: newLogger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	DB = db
	err = DB.AutoMigrate(&User{}, &Orders{}, &Session{}, &Blade{})
	if err != nil {
		log.Fatal(err)
	}
}

// User:
// id
// name
// password
// token
// status
// flux

const (
	Normal       = iota
	UnNormal     = iota
	PasswordCost = 12
)

func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
