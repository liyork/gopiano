package bfbczd

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

//在这个例子中，整个进程的生命周期很短，但是在一个真正的程序中，goroutines可以很容易地在一个长期生命的程序开始时启动，导致内存利用率下降。
func TestGoLeak(t *testing.T) {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings { // 只要strings进行了close，这个routine就退出
				fmt.Println(s)
			}
		}()
		return completed
	}

	//看到doWork被传递了一个nil通道。所以strings通道永远无法读取到其承载的内容，而且包含doWork的goroutine将在这个过程的整个生命周期中保留在内存中
	doWork(nil)
	fmt.Println("Done.")
}

func TestGoLeakSolve(t *testing.T) {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} { //1
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for { // 一定用for不断检查
				select {
				case s := <-strings:
					// Do something interesting
					fmt.Println(s)
				case <-done: //2
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() { //3
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated //4
	fmt.Println("Done.")

	//1.这里我们传递done通道给doWork函数。作为惯例，这个通道被作为首个参数。
	//2.这里我们看到使用了for-select的使用模式之一。我们的目的是检查done通道有没有发出信号。如果有的话，我们退出当前goroutine。
	//3.在这里我们创建另一个goroutine，一秒后就会取消doWork中产生的goroutine。
	//4.这是我们在main goroutine中调用doWork函数返回结果的地方。
}

func TestGoLeadWrite(t *testing.T) {
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.") // 1
			defer close(randStream)
			for {
				randStream <- rand.Int()
			}
		}()

		return randStream
	}

	randStream := newRandStream()
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	//1.当goroutine成功退出时我们打印一行消息。
	//可以看到注释1所在的打印语句并未执行。在main循环的第三次迭代之后，goroutine块试图将下一个随机整数发送到不再被读取的通道。
}

func TestGoLeadWriteSolve(t *testing.T) {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)

			for {
				select {
				case randStream <- rand.Int():
				case <-done: // 感知到关闭则退出
					return
				}
			}

		}()

		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}

	close(done)
	//模拟正在进行的工作
	time.Sleep(1 * time.Second)
}
