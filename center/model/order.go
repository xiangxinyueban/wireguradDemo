package model

import (
	"gorm.io/gorm"
	"time"
)

type Orders struct {
	gorm.Model
	User          User `gorm:"ForeignKey:Uid"`
	Uid           uint `gorm:"not null"`
	ComboType     byte
	RemainTraffic int64
	StartTime     time.Time
	EndTime       time.Time
}
