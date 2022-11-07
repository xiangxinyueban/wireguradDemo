package center

import "vpn/center/model"

// register -> login -> combo list -> save combo type in user config
// combo list -> create order -> order success ->
// server list(country) -> establish session -> Success page()
// delete session

// mysql struct
// user:
// UserName, Password, UserId, Email
// task:
// ComboType, RemainTraffic, ExpireTime, UserId,
// role?
//

func BootStrap() {
	model.InitDB()

}
