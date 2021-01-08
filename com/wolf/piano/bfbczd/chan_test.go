package bfbczd

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
)

func Test111(t *testing.T) {
	testBufChan()
}

// 这是一个在正确条件下可以使用的优化示例：
// 如果写入通道的goroutine明确知道将会写入多少信息，则创建相对应的缓冲通道容量会很有用，就可以尽可能快地进行读取
func testBufChan() {
	var stdoutBuff bytes.Buffer         //1
	defer stdoutBuff.WriteTo(os.Stdout) //2

	intStream := make(chan int, 4) //3
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ { // 一共5各元素
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %v.\n", integer)
	}
	//1.这里我们创建一个内存缓冲区来帮助缓解输出的不确定性。 它不会给带来我们任何保证，但比直接写stdout要快一些。
	//2.在这里，我们确保在进程退出之前将缓冲区内容写入标准输出。
	//3.这里我们创建一个容量为4的缓冲通道。
}

func TestChanInClosure(t *testing.T) {
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5) //1
		go func() {                       //2
			defer close(resultStream) //3
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream //4
	}

	resultStream := chanOwner()
	for result := range resultStream { //5
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
	//1.这里我们实例化一个缓冲通道。 由于我们知道我们将产生六个结果，因此我们创建了五个缓冲通道，以便该goroutine可以尽快完成操作，但是需要外界接收一个才能结束。
	//2.在这里，我们启动一个匿名的goroutine，它在resultStream上执行写操作。 请注意，我们是如果创建goroutines的， 它现在被封装在函数中。
	//3.这里我们确保resultStream在操作完成后关闭。作为通道所有者，这是我们的责任。
	//4.我们在这里返回通道。由于返回值被声明为只读通道，resultStream将隐式转换为只读的。
	//5.这里我们消费了resultStream。 作为消费者，我们只关心阻塞和通道的关闭。
}

//这使得goroutine可以在等待另一个goroutine报告结果的同时取得进展
//有一个循环正在做某种工作，偶尔检查它是否应该停止
func TestSelect(t *testing.T) {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		// Simulate work
		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
}
