package client

import (
	"github.com/skip2/go-qrcode"
	"log"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	CreateOrder(BASE_URI + "createOrder")

	err := qrcode.WriteFile("https://qr.alipay.com/fkx15582jpnx1nctrv9zd5f", qrcode.Medium, 256, "qr.png")
	if err != nil {
		log.Println(err)
	}
}
