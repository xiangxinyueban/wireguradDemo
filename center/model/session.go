package model

import "gorm.io/gorm"

const (
	ACTIVE   = iota
	INACTIVE = iota
)

type Session struct {
	gorm.Model
	User      User   `gorm:"ForeignKey:Uid"`
	Uid       uint   `gorm:"not null"`
	SessionID string `gorm:"PrimaryKey"`
	Traffic   uint64 //traffic used
	Country   string
	Status    byte
}
