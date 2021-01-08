package main

import (
	"fmt"
	"net"
	"os"
)

func lookIP(address string) ([]string, error) {
	hosts, err := net.LookupAddr(address)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func lookHostname(hostname string) ([]string, error) {
	IPs, err := net.LookupHost(hostname)
	if err != nil {
		return nil, err
	}
	return IPs, nil
}

// go run DNS.go 127.0.0.1
//go run DNS.go 10.0.76.211
//go run DNS.go packtpub.com
//go run DNS.go google.com
//go run DNS.go www.google.com
//go run DNS.go cnn.com
// // go run DNS.go 192.168.1.1命令的输出来自/etc/hosts文件，因为在/etc/hosts文件中配置了IP地址192.168.1.1的别名cisco。
////最后一个命令的输出演示了域名（cnn.com）可能有多个公网IP地址映射。请特别注意*公网*这个词，尽管www.google.com有多个IP地址，但是只有地址（216.58.214.36）是公网IP地址。
func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide an argument!")
		return
	}

	input := arguments[1]
	IPaddress := net.ParseIP(input)

	if IPaddress == nil {
		IPs, err := lookHostname(input)
		if err == nil {
			for _, singleIP := range IPs {
				fmt.Println(singleIP)
			}
		}
	} else {
		hosts, err := lookIP(input)
		if err == nil {
			for _, hostname := range hosts {
				fmt.Println(hostname)
			}
		}
	}
}
