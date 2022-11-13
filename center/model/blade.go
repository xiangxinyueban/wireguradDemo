package model

import (
	"gorm.io/gorm"
	"net"
)

type Blade struct {
	gorm.Model
	Traffic  uint64 //traffic used
	Country  string
	Status   byte
	Address  net.IP
	Password string
	UserName string
}
