package iputil

import (
	"fmt"
	"testing"
)

func TestGetPublicIP(t *testing.T) {
	public := GetPublicIP()
	fmt.Println(public)
}
