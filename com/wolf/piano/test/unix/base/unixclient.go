package main

import (
	"net"
	"time"
)

func main() {
	unixConn, err := net.DialTimeout("unix", "abc.sock", 1*time.Second)
	if err != nil {
		panic(err)
	}

	uc := unixConn.(*net.UnixConn)
	uc.Write([]byte{1})

	uc.Close()
}
