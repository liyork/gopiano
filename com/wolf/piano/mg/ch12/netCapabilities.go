package main

import (
	"fmt"
	"net"
)

// 打印UNIX操作系统的主机上每个网络接口。

// go run netCapabilities.go
func main() {
	interfaces, err := net.Interfaces()

	if err != nil {
		fmt.Print(err)
		return
	}

	for _, i := range interfaces {
		fmt.Printf("Name: %v\n", i.Name)
		fmt.Println("Interface Flags:", i.Flags.String())
		fmt.Println("Interface MTU:", i.MTU)
		fmt.Println("Interface Hardware Address:", i.HardwareAddr)
		fmt.Println()
	}
}

// 查找主机的默认网关
// netstat -nr
