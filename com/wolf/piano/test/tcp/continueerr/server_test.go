package base

import (
	"fmt"
	"net"
	"testing"
)

func TestServerBase(t *testing.T) {
	fmt.Println("hello world")

	lner, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("listener creat error", err)
	}
	fmt.Println("waiting for client")
	for {
		// 阻塞
		conn, err := lner.Accept()
		if err != nil {
			fmt.Println("accept error", err)
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	fmt.Println("connection success")
	fmt.Println("client address: ", conn.RemoteAddr())

	buffer := make([]byte, 1024)
	recvLen, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Read error", err)
		return
	}
	strBuffer := string(buffer[:recvLen])
	fmt.Println("Message: ", strBuffer)
	fmt.Println("Message len :", recvLen)

	// 持续写
	for {
		buffer = make([]byte, 1024)
		//time.Sleep(time.Second * 1) //等一秒钟，可以看出client里面的read函数有阻塞效果
		sendLen, err := conn.Write([]byte("I am server, you message :" + strBuffer))
		if err != nil {
			fmt.Println("send message error", err)
		}
		fmt.Println("send message len；", sendLen)

		//time.Sleep(2 * time.Second)
		fmt.Println("")
	}
}
