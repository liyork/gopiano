package ch3

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8085"
	DELIMITER      = '\t'
)

// 为了遵循开放封闭原则，该变量是包级私有的
var wg sync.WaitGroup

func TestServer(t *testing.T) {
	wg.Add(2)
	go serverGo()
	time.Sleep(500 * time.Millisecond)
	go clientGo(1)
	wg.Wait()
}

// 遵循单一职责原则，也利于并发运行
func serverGo() {
	var listener net.Listener
	listener, err := net.Listen(SERVER_NETWORK, SERVER_ADDRESS)
	if err != nil {
		printServerLog("Listen Error: %s", err)
		return
	}
	defer listener.Close()
	// 隔离将来可能发生的日志记录方式的变化，并能够避免散弹式修改
	printServerLog("Got listener for the server. (local address: %s)", listener.Addr())

	for {
		conn, err := listener.Accept() // 阻塞直到新连接到来
		if err != nil {
			printServerLog("Accept Error: %s", err)
			return
		}
		printServerLog("Established a connection with a client application. (remote address: %s)", conn.RemoteAddr())
		// 为了快速、独立地处理已经建立的每一个连接，应该尽量让这些处理过程并发执行
		go handleConn(conn)
	}
}

// go语句携带的函数(或称go函数)向外传递结果值的方式决定了handleConn函数没有结果声明
func handleConn(conn net.Conn) {
	defer wg.Done()
	defer conn.Close() // 执行结束时应把连接关闭。执行结束可能是由于主题已经执行结束，也可能是某些代码引起panic
	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second)) // 超时闲置
		strReq, err := read(conn)
		if err != nil {
			if err == io.EOF { // 连接已关闭
				printServerLog("The connection is closed by another side.")
			} else {
				printServerLog("Read Error: %s", err)
			}
			break
		}
		printServerLog("Received request: %s.", strReq)

		intReq, err := strToInt32(strReq)
		if err != nil {
			n, err := write(conn, err.Error())
			printServerLog("Sent error message (written %d bytes): %s.", n, err)
			continue
		}
		floatResp := cbrt(intReq)
		respMsg := fmt.Sprintf("The cube root of %d is %d.", intReq, floatResp)
		n, err := write(conn, respMsg)
		if err != nil {
			printServerLog("Write Error: %s", err)
		}
		printServerLog("Sent response (written %d bytes): %s.", n, respMsg)
	}
}

func cbrt(i int32) int32 {
	return i * i * i
}

func strToInt32(s string) (int32, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// 从conn中读取符合分隔符的一块内容
func read(conn net.Conn) (string, error) {
	readBytes := make([]byte, 1) // 每次读一字节，防止从连接值中读出多余的数据从而对后续的读取操作造成影响
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return "", err
		}
		readByte := readBytes[0]
		if readByte == DELIMITER {
			break
		}
		buffer.WriteByte(readByte)
	}
	return buffer.String(), nil
}

func write(conn net.Conn, content string) (int, error) {
	var buffer bytes.Buffer
	buffer.WriteString(content)
	buffer.WriteByte(DELIMITER) // 追加数据分界符
	return conn.Write(buffer.Bytes())
}

func printLog(role string, sn int, format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Printf("%s[%d]: %s", role, sn, fmt.Sprintf(format, args...))
}
func printServerLog(format string, args ...interface{}) {
	printLog("Server", 0, format, args...)
}

func printClientLog(sn int, format string, args ...interface{}) {
	printLog("Client", sn, format, args...)
}

func clientGo(id int) {
	defer wg.Done()
	conn, err := net.DialTimeout(SERVER_NETWORK, SERVER_ADDRESS, 2*time.Second)
	if err != nil {
		printClientLog(id, "Dial Error: %s", err)
		return
	}
	defer conn.Close()
	printClientLog(id, "Connected to server. (remote address:%s, local address: %s)", conn.RemoteAddr(), conn.LocalAddr())
	time.Sleep(200 * time.Millisecond)

	requestNumber := 5
	conn.SetDeadline(time.Now().Add(5 * time.Millisecond))
	for i := 0; i < requestNumber; i++ {
		req := rand.Int31()
		n, err := write(conn, fmt.Sprintf("%d", req))
		if err != nil {
			printClientLog(id, "Write Error: %s", err)
			continue
		}
		printClientLog(id, "Sent request (written %d bytes): %d.", n, req)
	}

	for j := 0; j < requestNumber; j++ {
		strResp, err := read(conn)
		if err != nil {
			if err == io.EOF {
				printClientLog(id, "The connection is closed by another side.")
			} else {
				printClientLog(id, "Read Error: %s", err)
			}
			break
		}
		printClientLog(id, "Received response: %s.", strResp)
	}
}
