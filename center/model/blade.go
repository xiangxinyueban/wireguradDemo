package model

import (
	"gorm.io/gorm"
	"net"
)

type Blade struct {
	gorm.Model
	Hostname    string `gorm:"unique"`
	Traffic     uint64 //traffic used
	Country     string
	Status      byte
	Address     net.IP
	Password    string
	UserName    string
	Vendor      string
	Users       int
	BootStrapID string
}
