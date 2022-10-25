package iputil

import (
	"fmt"
	"net"
	"testing"
)

func TestAllocateIPNet(t *testing.T) {
	subnet := net.IPNet{
		IP:   net.ParseIP("10.182.0.0").To4(),
		Mask: net.IPv4Mask(255, 255, 0, 0),
	}
	allocator := NewAllocator(subnet)
	fmt.Println(allocator.AllocateIPNet())
	fmt.Println(allocator.AllocateIPNet())
	fmt.Println(allocator.AllocateIPNet())
	fmt.Println(allocator.AllocateIPNet())
	fmt.Println(allocator.AllocateIPNet())
	fmt.Println(allocator.AllocateIPNet())
	fmt.Println(allocator.AllocateIPNet())
	ipNet, _ := allocator.AllocateIPNet()
	fmt.Println(FirstIP(ipNet))
}
