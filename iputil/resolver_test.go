package iputil

import (
	"fmt"
	"testing"
	"vpn/requests"
)

func TestPublicIp(t *testing.T) {
	resolver := NewResolver(requests.NewHTTPClient("0.0.0.0", requests.DefaultTimeout), "0.0.0.0", "http://www.baidu.com", []string{})
	fmt.Println(resolver.GetPublicIP())
}
