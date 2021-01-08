package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

// go run server.go common.go
// kill -USR2 pid
func main() {
	// unlink删除已存在的unixSock文件
	syscall.Unlink(socksPath)
	laddr, err := net.ResolveUnixAddr("unix", socksPath)
	if err != nil {
		panic(err)
	}
	l, err := net.ListenUnix("unix", laddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("waiting for conn from unix socks\n")
	// 等待建联，阻塞
	conn, err := l.AcceptUnix()
	if err != nil {
		panic(err)
	}

	fmt.Printf("has accept...time:%v\n", time.Now())
	// msg分为两部分数据
	buf := make([]byte, 32)
	oob := make([]byte, 32)
	// 一直读取，直到读取到符合结构或者close或者异常
	n, oobn, _, _, err := conn.ReadMsgUnix(buf, oob)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read from conn...time:%v\n", time.Now())
	fmt.Println("buf:", buf[:n])
	// 解出SocketControlMessage数组
	scms, err := syscall.ParseSocketControlMessage(oob[:oobn])
	if err != nil {
		panic(err)
	}
	if len(scms) > 0 {
		// 从SocketControlMessage中得到UnixRights
		fds, err := syscall.ParseUnixRights(&(scms[0]))
		if err != nil {
			panic(err)
		}
		fmt.Printf("parse %d fds: %v \n", len(fds), fds)
		// os.NewFile()将文件描述符fd转为 *os.File对象, 并不创建新文件
		f := os.NewFile(uintptr(fds[0]), "")
		// 下面操作后最后close这file
		defer f.Close()
		// 从文件中读取文本内容
		buf := make([]byte, 1024)
		n, err := f.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("read len:%d, data:%s from file success\n", n, string(buf[:n]))
		return
	}
	err = conn.Close()
	if err != nil {
		panic(err)
	}
}
