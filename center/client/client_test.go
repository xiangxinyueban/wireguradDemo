package client

import (
	"testing"
)

func TestCreateOrder(t *testing.T) {
	CreateOrder(BASE_URI + "createOrder")

	//err := qrcode.WriteFile("https://qr.alipay.com/fkx17776qje1i0zys2zww42", qrcode.Medium, 256, "qr.png")
	//if err != nil {
	//	log.Println(err)
	//}
}
