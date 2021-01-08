package main

import (
	"fmt"
	"net"
)

// 网络包采用二进制格式，这要求您查找特定类型的网络数据包，而不仅仅是任何类型的网络数据包

// 为了产生 ICMP 流量，可以使用ping/traceroute
// ping -c 5 localhost
// traceroute localhost

// sudo go run lowlevel.go
func main() {
	netaddr, err := net.ResolveIPAddr("ip4", "127.0.0.1")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenIP("ip4:icmp", netaddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFrom(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("% X\n", buffer[0:n])
}
