package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/sftp"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"time"
	"vpn/center/model"
	"vpn/center/serializer"
)

type BladeService struct {
	Address  string `form:"address" json:"address" binding:"required" example:"Anonymous"`
	Username string `form:"username" json:"username" binding:"required" example:"Anonymous"`
	Password string `form:"password" json:"password" binding:"required" example:"Anonymous"`
	Country  string `form:"country" json:"country" binding:"-" example:"Anonymous"`
	Traffic  int    `form:"traffic" json:"traffic" binding:"required" example:"Anonymous"`
	Vendor   string `form:"vendor" json:"vendor" binding:"-" example:"Anonymous"`
	Center   center `form:"center" json:"center" binding:"-" example:"Anonymous"`
	Local    local  `form:"local" json:"local" binding:"-" example:"Anonymous"`
}

type center struct {
	rpcport int    `form:"rpc_port" json:"rpc_port" binding:"-" example:"Anonymous"`
	address string `form:"address" json:"address" binding:"-" example:"Anonymous"`
}

type local struct {
	rpcport       int    `form:"rpc_port" json:"rpc_port" binding:"-" example:"Anonymous"`
	bootstrapPeer string `form:"bootstrap_peer" json:"bootstrap_peer" binding:"-" example:"Anonymous"`
}

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

func (b *BladeService) Register() serializer.Response {
	var blade model.Blade

	if err := model.DB.Where("hostname = ?", b.Address).First(&blade).Error; err != nil {
		//如果查询不到，返回相应的错误
		if gorm.IsRecordNotFoundError(err) {
			return serializer.Response{
				Code:  -1,
				Error: "服务器已存在",
			}
		}
		if gorm.IsRecordNotFoundError(err) {
			log.Error().Err(err).Msgf("添加服务器出错")
			return serializer.Response{
				Code:  -1,
				Error: "添加服务器出错",
			}
		}

	}
	blade.Hostname = b.Address
	blade.Address = net.ParseIP(b.Address)
	blade.UserName = b.Username
	blade.Password = b.Password
	blade.Vendor = b.Vendor
	blade.Traffic = uint64(b.Traffic * 1024 * 1024 * 1024)
	blade.Country = b.Country
	blade.Status = 1

	addr := fmt.Sprintf("%s:%v", b.Address, 22)
	config := ssh.ClientConfig{
		User: b.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(b.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	client, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: fmt.Sprintf("ssh Connect blade failed, err: %s", err),
		}
	}

	if sftpClient, err := sftp.NewClient(client); err != nil { //创建客户端
		fmt.Println("创建客户端失败", err)
		return serializer.Response{
			Code:  -1,
			Data:  nil,
			Msg:   "",
			Error: fmt.Sprintf("ssh Connect blade failed, err: %s", err),
		}
	} else {
		//execute(client, "sudo mkdir -p /etc/vpn/config")
		//sftpClient.MkdirAll("/etc/vpn/config/")
		execute(client, "sudo killall vpn")
		upload(sftpClient, "/usr/local/vpn", "/home/ubuntu/vpn")
		upload(sftpClient, "/etc/vpn/config/server_temp.ini", "/home/ubuntu/server.ini")
	}
	//execute(client, "sudo -i")
	execute(client, "chmod u+x /home/ubuntu/vpn")
	execute(client, "id")
	execute(client, "sudo nohup /home/ubuntu/vpn > /home/ubuntu/vpn.log 2>&1 &")
	//model.DB.Create(&blade)
	return serializer.Response{
		Code:  1,
		Data:  nil,
		Msg:   "服务器启动成功",
		Error: "",
	}
}

func execute(client *ssh.Client, command string) string {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Sprintf("client session create failed, err:%s", err)
	}
	var outputBytes []byte
	outputBytes, err = session.CombinedOutput(command)
	if err != nil {
		return fmt.Sprintf("client session create failed, err:%s", err)
	}
	log.Debug().Msgf("executed: %s, result: %s", command, string(outputBytes))
	return string(outputBytes)
}

func upload(client *sftp.Client, srcPath, dstPath string) {
	srcFile, _ := os.Open(srcPath)       //本地
	dstFile, _ := client.Create(dstPath) //远程
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()
	buf := make([]byte, 1024*1024*10)
	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal().Err(err)
			} else {
				break
			}
		}
		_, _ = dstFile.Write(buf[:n])
	}
}
