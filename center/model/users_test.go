package model

import (
	"fmt"
	"testing"
)

func TestInitDB(t *testing.T) {
	InitDB()
	var usr User
	DB.Model(&User{}).Create(&User{
		UserName: "harris",
		Password: "bossis42",
		Email:    "1315909600@qq.com",
		Token:    "xxxxxxxxxxxx",
		Status:   1,
	})
	DB.Model(&User{}).Where("user_name=?", "harris").First(&usr)
	usr.Password = "qwert12345"
	DB.Model(&User{}).Where("user_name=?", "harris").Update("password", usr.Password)
	fmt.Printf("%+v", usr)
}
