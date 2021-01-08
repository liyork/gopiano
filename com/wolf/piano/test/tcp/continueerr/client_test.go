package base

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestClientBase(t *testing.T) {
	fmt.Println("client launch")
	serverAddr := "localhost:8888"
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		fmt.Println("Resolve TCPAddr error", err)
		return
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		fmt.Println("connect server error", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	conn.Write([]byte("hello , I am client"))
	go recv(conn, &wg)
	wg.Wait()
	fmt.Println("finish")
}

// 有粘包问题
func recv(conn net.Conn, wg *sync.WaitGroup) {
	//fmt.Printf("recv wg:%p\n", &wg)
	defer wg.Done()

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err == nil {
			fmt.Println("read message from server:" + string(buffer[:n]))
			fmt.Println("Message len:", n)
		}
	}
}
