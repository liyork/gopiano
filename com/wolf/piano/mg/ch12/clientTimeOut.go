package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var timeout2 = time.Duration(time.Second)

// 自己控制发送请求的时间
func Timeout(network, host string) (net.Conn, error) {
	conn, err := net.DialTimeout(network, host, timeout2)
	if err != nil {
		return nil, err
	}
	// 设置给定网络连接的读写截止时间。由于 `SetDeadline()` 函数的工作方式，您需要在任何读写操作前调用 `SetDeadline()`。
	// 请记住，Go 使用截止日期来实现超时，因此您不需要在应用程序每次接收或发送任何数据时重置它。
	conn.SetDeadline(time.Now().Add(timeout2))
	return conn, nil
}

// go run clientTimeOut.go http://localhost:8001
// go run clientTimeOut.go http://localhost:8001 15
func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s URL TIMEOUT\n", filepath.Base(os.Args[0]))
		return
	}

	if len(os.Args) == 3 {
		temp, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Using Default Timeout!")
		} else {
			timeout2 = time.Duration(time.Duration(temp) * time.Second)
		}
	}

	URL := os.Args[1]
	t := http.Transport{
		Dial: Timeout,
	}

	client := http.Client{
		Transport: &t,
	}

	data, err := client.Get(URL)
	if err != nil {
		fmt.Println("err:", err)
		return
	} else {
		defer data.Body.Close()
		_, err := io.Copy(os.Stdout, data.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
