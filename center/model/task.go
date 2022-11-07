package model

type User struct {
	gorm.Model
	ComboType     int
	RemainTraffic int64
	ExpireTime    int64
	UserId
}
