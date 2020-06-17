package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func convertChan1(writeChan chan<- int) {

}

// 只读
func convertChan2(readChan <-chan int) {
	for a := range readChan {
		fmt.Println(a)
	}
	// 编译报错,不允许对只读chan写
	//readChan<-1
}

func testConvert() {
	bichan := make(chan int)
	convertChan1(bichan) // 转成只读
	convertChan2(bichan) // 转成只写

	bichan2 := make(chan int, 2) // capacity
	fmt.Println(bichan2)
}

// 一个chan多个读取只能有一个成功
func TestMultiReceive(t *testing.T) {
	pubchan := make(chan int)

	go func() {
		for {
			a := <-pubchan
			fmt.Println("receive a1 succss, a2:", a)
		}
	}()

	go func() {
		for {
			a := <-pubchan
			fmt.Println("receive a2 succss, a2:", a)
		}
	}()

	for i := 0; i < 10; i++ {
		//time.Sleep(1 * time.Second)
		pubchan <- i
	}
}

// 无缓冲channel则写入/读取阻塞直到另一端有人读/写
func TestNoBufferBlockRead(t *testing.T) {
	noBuffChan := make(chan int)
	go writeNoBuffer(noBuffChan)
	time.Sleep(time.Second * 2)
	noBuffChan <- 2
	fmt.Println("main finish:")
}

func writeNoBuffer(noBuffChan chan int) {
	a := <-noBuffChan
	// 阻塞，直到有人写入
	fmt.Println("a:", a)
}

func TestBufferBlockRead(t *testing.T) {
	buffChan := make(chan int, 2)
	a := <-buffChan // 阻塞，因为缓冲中没有数据
	fmt.Println("a:", a)
}

// 缓冲区满之前都不阻塞
func TestBufferBlockWrite(t *testing.T) {
	buffChan := make(chan int, 2)
	buffChan <- 1
	buffChan <- 2
	buffChan <- 3 // 阻塞，缓冲被写满
	fmt.Println("a:")
}

func TestChannelClose(t *testing.T) {
	chann := make(chan int)
	go processWithChanel(chann)
	fmt.Println("wait processWithChanel")
	<-chann
	fmt.Println("finish main")
}

func processWithChanel(chann chan int) {
	fmt.Println("business")
	time.Sleep(time.Second * 3)
	// there is error
	close(chann)
}

func TestChanTimeout(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		time.Sleep(122 * time.Second)
		ch <- 1
	}()

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	select {
	case <-ch:
		fmt.Println("read from ch")
	case <-timeout:
		fmt.Println("read from ch timeout")
	}
}

// 只要一个结果即可
func Query(conns []string, query string) string {
	ch := make(chan string)
	for _, conn := range conns {
		go func(c string) {
			select {
			case ch <- DoQuery(query):
			default: // If the send cannot go through immediately the default case will be selected.
			}
		}(conn)
	}
	return <-ch
}

func DoQuery(s string) string {
	return ""
}
