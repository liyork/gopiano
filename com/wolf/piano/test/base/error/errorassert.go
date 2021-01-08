package main

import (
	"fmt"
	"net"
	"os"
)

//断言底层结构类型并从结构字段获取更多信息
func main() {
	openFile()
	//connect()
}

// 使用类型断言获得错误接口的基本值
func openFile() {
	f, err := os.Open("/test.txt")
	if err, ok := err.(*os.PathError); ok {
		fmt.Println("File at path", err.Path, "failed to open")
		return
	}
	fmt.Println(f.Name(), "opened successfully")
}

// 正在尝试获取一个无效域名的ip地址，这是一个无效的域名。http://golangbot123.com。通过类型断言为*net.DNSError来获得错误的潜在价值。
func connect() {
	addr, err := net.LookupHost("golangbot123.com")
	if err, ok := err.(*net.DNSError); ok {
		if err.Timeout() {
			fmt.Println("operation timed out")
		} else if err.Temporary() {
			fmt.Println("temporary error")
		} else {
			fmt.Println("generic error: ", err)
		}
		return
	}
	fmt.Println(addr)
}
