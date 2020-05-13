package base

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func TestGetGoRoutineId(t *testing.T) {
	go testg1()
	go testg1()
	go testg1()
	time.Sleep(time.Second * 5)
}

func testg1() {
	fmt.Println("getgid:", GetGID())
	time.Sleep(time.Second * 50)
}
