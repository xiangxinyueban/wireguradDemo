package model

import (
	"fmt"
	"testing"
)

func TestInitDB(t *testing.T) {
	db := InitDB()
	//db.Create(&[]User{
	//	{
	//		Name:     "harris",
	//		ExpireAt: time.Now(),
	//		Status:   Normal,
	//		Password: "bossis42",
	//		Token:    "xxxxxxxxx",
	//		Flux:     uint64(55),
	//	},
	//})
	var usr User
	db.First(&usr, 1)
	fmt.Printf("%+v", usr)
}
