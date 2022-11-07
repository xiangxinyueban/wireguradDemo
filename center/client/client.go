package client

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"vpn/center/snow"
)

const BASE_URI = "http://43.142.50.173:8080/"

type CreateResponse struct {
	Code int64         `json:"code"`
	Msg  string        `json:"msg"`
	Data []interface{} `json:"data,omitempty"`
}

type CreateParam struct {
	PayId int64
	Type  int
	Price float32
	Param string
	Sign  string
}

var snowFlake *snow.SNOW

func init() {
	var err error
	snowFlake, err = snow.NewSnow(1)
	if err != nil {
		log.Fatalln("init failed, err:", err)
	}

}

func CreateOrder(url string) {
	key := "194ba281abba943113d7c43cb2f990ce"
	payId := strconv.FormatInt(snowFlake.GetID(), 10)
	typ := strconv.Itoa(2)
	price := strconv.FormatFloat(20, 'f', 2, 32)
	param := ""
	params := map[string][]string{
		"payId": {payId},
		"type":  {typ}, // wechat:1, alipay:2
		"price": {price},
		"param": {param},
	}
	params["sign"] = []string{
		md5Encode(payId + param + typ + price + key),
	}
	resp, err := http.PostForm(url, params)
	ErrPrint(err)
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	//decoder := json.NewDecoder(resp.Body)
	//var createRsp CreateResponse
	//err = decoder.Decode(&createRsp)
	//ErrPrint(err)
	//if createRsp.Code != 1 {
	//	fmt.Println(createRsp.Msg)
	//}
	fmt.Printf("%s\r\n", bytes)
}

func ErrPrint(err error) {
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func md5Encode(str string) string {
	md := md5.New()
	md.Write([]byte(str))
	return hex.EncodeToString(md.Sum(nil))
}
