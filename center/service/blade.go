package service

//
//import (
//	"fmt"
//	"golang.org/x/crypto/ssh"
//	"golang.org/x/net/http2"
//	"net"
//	"time"
//	"vpn/center/serializer"
//)
//
//type BladeService struct {
//	UserName     string `json:"user_name"`
//	Password     string `json:"password"`
//	Traffic      uint64 `json:"traffic,omitempty"` //traffic used
//	TotalTraffic uint64 `json:"total_traffic"`
//	Country      string `json:"country"`
//	Status       byte   `json:"status"`
//	Address      string `json:"address"`
//	Config       string `json:"config"`
//}
//
//func (b *BladeService) Register() serializer.Response {
//	addr := fmt.Sprintf("%s:%s", b.Address, 22)
//	config := ssh.ClientConfig{
//		User: b.UserName,
//		Auth: []ssh.AuthMethod{
//			ssh.Password(b.Password),
//		},
//		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
//			return nil
//		},
//		Timeout: 10 * time.Second,
//	}
//	client, err := ssh.Dial("tcp", addr, &config)
//	if err != nil {
//		return serializer.Response{
//			Code:  -1,
//			Data:  nil,
//			Msg:   "",
//			Error: fmt.Sprintf("ssh Connect blade failed, err: %s", err),
//		}
//	}
//	execute(client, "")
//}
//
//func execute(client ssh.Client, command string) string {
//	session, err := client.NewSession()
//	if err != nil {
//		return fmt.Sprintf("client session create failed, err:%s", err)
//	}
//	var outputBytes []byte
//	outputBytes, err = session.CombinedOutput(command)
//	if err != nil {
//		return fmt.Sprintf("client session create failed, err:%s", err)
//	}
//	http2.Server{}
//	http2.ServeConnOpts
//	return string(outputBytes)
//}
