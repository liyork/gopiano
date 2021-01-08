package main

import (
	"fmt"
	"os"
	"syscall"
)

// sudo go run syscallNet.go
// ping 127.0.0.1
func main() {
	// `syscall.AF_INET使用 IPv4。
	//`syscall.SOCK_RAW参数使生成的 socket 成为原始 socket。
	// syscall.IPPROTO_ICMP`，告诉 `syscall.Socket()` 您只对 ICMP 通信感兴趣。
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		fmt.Println("Error in syscall.Socket:", err)
		return
	}
	f := os.NewFile(uintptr(fd), "captureICMP")
	if f == nil {
		fmt.Println("Error is os.NewFile:", err)
		return
	}

	// 设置 socket 的接收 buffer 大小为 256。`syscall.SOL_SOCKET` 参数是为了说明您想要在 socket 层级上工作。
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, 256)
	if err != nil {
		fmt.Println("Error in syscall.Socket:", err)
		return
	}

	for {
		buf := make([]byte, 1024)
		numRead, err := f.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("% X\n", buf[:numRead])
	}
}
