package service

import "testing"

func TestBladeRegister(t *testing.T) {
	var blade = &BladeService{
		Password: "#!1234wsjcl1234!#",
		Username: "ubuntu",
		Address:  "43.142.50.173",
		Country:  "CN",
		Traffic:  500 * 1024 * 1024 * 1024,
	}
	blade.Register()
}
