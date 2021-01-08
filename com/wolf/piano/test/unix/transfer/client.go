package main

import (
	"fmt"
	"net"
	"time"
)

// go run client.go common.go  启动两个
// 发送ping，接收pong
func main() {
	var count uint32
	conn, err := net.Dial("tcp", "127.0.0.1:7000")
	if err != nil {
		panic(err)
	}
	defer func() {
		conn.Close()
	}()
	for {
		time.Sleep(time.Second)
		err := conn.SetWriteDeadline(time.Now().Add(writeTimeout))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println("send ping")
		_, err = conn.Write([]byte(`ping`))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		err = conn.SetReadDeadline(time.Now().Add(readTimeout))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		rBuf := make([]byte, 4)
		_, err = conn.Read(rBuf)
		if err != nil {
			fmt.Println(err.Error())
		}
		count++
		fmt.Println("count:", count, "recv "+string(rBuf))
	}
}
