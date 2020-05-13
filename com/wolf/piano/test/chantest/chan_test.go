package chantest

import (
	"fmt"
	"testing"
	"time"
)

func convertChan1(writeChan chan<- int) {

}

func convertChan2(readChan <-chan int) {
	for a := range readChan {
		fmt.Println(a)
	}
}

func testConvert() {
	bichan := make(chan int)
	convertChan1(bichan) // 转成只读
	convertChan2(bichan) // 转成只写

	bichan2 := make(chan int, 2) // capacity
	fmt.Println(bichan2)
}

// 一个chan多个读取只能有一个成功，似乎第二次再来的读的goroutine会等待前者读完
func TestMultiReceive(t *testing.T) {
	pubchan := make(chan int)

	go func() {
		for {
			a := <-pubchan
			fmt.Println("a1:", a)
		}
	}()

	go func() {
		for {
			a := <-pubchan
			fmt.Println("a2:", a)
		}
	}()

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		pubchan <- i
	}
}

func testNoBufferBlockRead() {
	noBuffChan := make(chan int)
	a := <-noBuffChan // 阻塞，直到有人写入
	fmt.Println("a:", a)
}

func testNoBufferBlockWrite() {
	noBuffChan := make(chan int)
	noBuffChan <- 2 // 阻塞，直到有人读取
	fmt.Println("xxx")
}

func TestBufferBlockRead(t *testing.T) {
	buffChan := make(chan int, 2)
	a := <-buffChan // 阻塞，因为缓冲中没有数据
	fmt.Println("a:", a)
}

func TestBufferBlockWrite(t *testing.T) {
	buffChan := make(chan int, 2)
	buffChan <- 1
	buffChan <- 2
	buffChan <- 3 // 阻塞，缓冲被写满
	fmt.Println("a:")
}
