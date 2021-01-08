package main

import (
	"fmt"
	"net"
	"syscall"
)

func main() {
	// 若是没有这个，则重复运行会有panic: listen unix abc.sock: bind: address already in use
	// 但是有这个unlink之后，第二个运行的程序会接手这个监听，第一个程序就一直运行也不能接收监听
	syscall.Unlink("abc.sock") // abc.sock是删除abc.sock文件，Listen是生成文件。
	l, err := net.Listen("unix", "abc.sock")
	if err != nil {
		panic(err)
	}

	conn, err := l.Accept()
	if err != nil {
		panic(err)
	}

	addr := conn.RemoteAddr()
	fmt.Println("addr:", addr)
}
