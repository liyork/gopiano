package main

import (
	"fmt"
	"net"
)

// 可以看到有的网络接口没有显示网络地址，可能的主要原因是接口是关闭的，或者该接口没有进行配置。
// 并非所有列出的网络接口都关联了真正的硬件网络设备。最典型的例子是lo0接口，它是环回设备。环回设备是一种特殊的虚拟网络接口，主机可以通过该接口与自身通信。

// go run netConfig.go
func main() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, i := range interfaces {
		fmt.Printf("Interface: %v\n", i.Name)
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			fmt.Println(err)
		}

		addresses, err := byName.Addrs()
		for k, v := range addresses {
			fmt.Printf("Interface Address #%v: %v\n", k, v.String())
		}
		fmt.Println()
	}
}
