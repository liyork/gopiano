package base

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

func recv(conn net.Conn, wg *sync.WaitGroup) {
	//fmt.Printf("recv wg:%p\n", &wg)
	defer wg.Done()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err == nil {
		fmt.Println("read message from server:" + string(buffer[:n]))
		fmt.Println("Message len:", n)
	}
}

func TestClientReadNonblock(t *testing.T) {
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

	//fmt.Printf("recv wg:%p\n", &wg)
	defer wg.Done()

	buffer := make([]byte, 1024)
	// 一直卡住，直到服务端有响应数据、关闭
	n, err := conn.Read(buffer)
	if err == nil {
		fmt.Println("read message from server:" + string(buffer[:n]))
		fmt.Println("Message len:", n)
	} else {
		fmt.Println("read message from server err, err:", err)
	}

	wg.Wait()
	fmt.Println("finish")
}

// 读取数据，需要server端进行close产生EOF，这里会进行累加数据，直到退出
func TestClientRead1(t *testing.T) {
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

	//fmt.Printf("recv wg:%p\n", &wg)
	defer wg.Done()

	//buf := make([]byte,0, 1000)// 这里若是，len=0，那么下面读取的数据永远是0长度。。。,看来是依据len进行最大读取的
	buf := make([]byte, 1000)
	len := 0
	for {
		n, err := conn.Read(buf[len:])
		if n > 0 {
			len += n
		}
		if err != nil {
			if err != io.EOF {
				//Error Handler
			}
			break
		}
	}

	fmt.Println("buf[:len]:", string(buf[:len]))
}

// 同样需要server进行close，不过这里用了buffer进行copy
func TestClientRead2(t *testing.T) {
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

	//fmt.Printf("recv wg:%p\n", &wg)
	defer wg.Done()

	defer conn.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, conn)
	if err != nil {
		// Error handler
		fmt.Println("err:", err)
	}

	fmt.Println("buf[:len]:", buf.String())
}

// 也是EOF或者close或者错误才能返回
func TestClientRead3(t *testing.T) {
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

	//fmt.Printf("recv wg:%p\n", &wg)
	defer wg.Done()

	defer conn.Close()

	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		// Error Handler
		fmt.Println("connect server error", err)
		return
	}

	fmt.Println("buf[:len]:", string(buf[:]))
}
