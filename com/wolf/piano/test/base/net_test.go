package base

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestResolve(t *testing.T) {
	time.Sleep(10 * time.Second)
	//os.Setenv("GODEBUG", "netdns=cgo")// 默认
	//os.Setenv("GODEBUG", "netdns=go")// 快
	addr, err := net.ResolveIPAddr("ip", "xx.local")
	if err != nil {
		fmt.Println("Resolvtion error", err.Error())
		os.Exit(1)
	}
	fmt.Println(addr.String())
	ips, err := net.LookupIP("xx.local")
	fmt.Println("ips:", ips)
}
