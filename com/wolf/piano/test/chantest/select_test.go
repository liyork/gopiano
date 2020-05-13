package chantest

import (
	"fmt"
	"testing"
	"time"
)

// select是执行选择操作的一个结构，它里面有一组case语句，它会执行其中无阻塞的那一个，如果都阻塞了，那就等待其中一个不阻塞，
// 进而继续执行，它有一个default语句，该语句是永远不会阻塞的

func TestNoBufferReadWitSelect(t *testing.T) {
	noBuffChan := make(chan int)
	readWithSelectDefault(noBuffChan)
}

func TestBufferReadWithSelect(t *testing.T) {
	buffChan := make(chan int, 2)
	readWithSelectDefault(buffChan)
}

func readWithSelectDefault(tmpChan chan int) {
	select {
	case a := <-tmpChan:
		fmt.Println("a:", a)
	default:
		fmt.Println("read no data in chan")
	}
}

func TestNoBufferReadWitSelectTimer(t *testing.T) {
	noBuffChan := make(chan int)
	readWithSelectTimer(noBuffChan)
}

// 相对于default，timer不会立即返回而是等待几秒
func readWithSelectTimer(tmpChan chan int) {
	timer := time.NewTimer(2 * time.Second)

	select {
	case a := <-tmpChan:
		fmt.Println("a:", a)
	case <-timer.C:
		fmt.Println("read no data in chan")
	}
}

func TestNoBufferWriteWithSelect(t *testing.T) {
	noBuffChan := make(chan int)
	writeWithSelect(noBuffChan, 2)
}

func TestNoBufferWriteWithSelectTimer(t *testing.T) {
	noBuffChan := make(chan int)
	writeWithSelectTimer(noBuffChan, 2)
}

func TestBufferWriteWithSelect(t *testing.T) {
	buffChan := make(chan int, 2)
	writeWithSelect(buffChan, 1)
	writeWithSelect(buffChan, 2)
	writeWithSelect(buffChan, 3)
}

func writeWithSelect(tmpChan chan int, n int) {
	select {
	case tmpChan <- n:
		fmt.Println("write to chan success")
	default:
		fmt.Println("write to chan fail")
	}
}

func writeWithSelectTimer(tmpChan chan int, n int) {
	timer := time.NewTimer(2000 * time.Millisecond)

	select {
	case tmpChan <- n:
		fmt.Println("write to chan success")
	case <-timer.C:
		fmt.Println("write to chan fail")
	}
}
