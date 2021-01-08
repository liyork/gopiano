package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

// go run clientsleep.go common.go
func main() {
	file, err := os.Open("./temp")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fdnum := file.Fd()
	fmt.Printf("fd: %b %b %b %b\n", byte(fdnum), byte(fdnum>>8), byte(fdnum>>16), byte(fdnum>>24))
	fmt.Printf("ready to send fd: %d\n", fdnum)
	// 编码fd编译传送给其他进程
	data := syscall.UnixRights(int(fdnum))
	raddr, err := net.ResolveUnixAddr("unix", socksPath)
	if err != nil {
		panic(err)
	}
	// 连接UnixSock
	conn, err := net.DialUnix("unix", nil, raddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("has dial, time:%v\n", time.Now())
	time.Sleep(10 * time.Second)
	// 发送msg
	n, oobn, err := conn.WriteMsgUnix(nil, data, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("has write, time:%v\n", time.Now())
	time.Sleep(10 * time.Second)
	fmt.Printf("WriteMsgUnix = n:%d, oobn:%d; want 1, %d\n", n, oobn, len(data))
	fmt.Printf("write %d data success\n", n)
}
