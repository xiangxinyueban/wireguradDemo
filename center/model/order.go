package model

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	User          User  `gorm:"ForeignKey:Uid"`
	Uid           uint  `gorm:"not null"`
	ComboType     int   `gorm:"default:1"`
	RemainTraffic int64 `gorm:"default:0"`
	StartTime     time.Time
	EndTime       time.Time
}
