package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var timeout = time.Duration(time.Second)

// go run ../ch10/http/slowWWW.go
// go run anotherTimeOut.go http://localhost:8001
// go run anotherTimeOut.go http://localhost:8001 15
func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please provide a URL")
		return
	}

	if len(os.Args) == 3 {
		temp, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Using Default Timeout!")
		} else {
			timeout = time.Duration(time.Duration(temp) * time.Second)
		}
	}

	URL := os.Args[1]

	client := http.Client{
		// 指定超时针对请求，包含连接、重定向、读取返回body
		Timeout: timeout,
	}
	client.Get(URL)

	data, err := client.Get(URL)
	if err != nil {
		fmt.Println(err)
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
