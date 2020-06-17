package routine

import (
	"fmt"
	"testing"
	"time"
)

// 退出goroutine，开协成一定要考虑如何以及何时关闭

// for-range，只处理一个chan
// 当channel被关闭时，range就会结束，接着退出for循环
func TestStopRoutine1(t *testing.T) {
	inCh := make(chan int)

	go func(in chan int) {
		for x := range in {
			fmt.Printf("Process %d\n", x)
		}
		// chan关闭后，读取到的是默认值，并不报错
		i, ok := <-in
		fmt.Println("routine finish", i, ok)
		// close后不能写入，否则panic: send on closed channel
		//in <- 1
	}(inCh)

	for i := 1; i < 5; i++ {
		fmt.Println("make i:", i)
		inCh <- i
		time.Sleep(2 * time.Second)
	}

	fmt.Println("main finish")
	close(inCh)
	time.Sleep(2 * time.Second)
}

// for-select1，使用ok判断in是否关闭，关闭了直接退出，其他的case不用管
// select提供了多路复用的能力，所以for-select可以让函数具有持续多路处理多个channel的能力。但select没有感知channel的关闭
// 一定要注意让一方(一般是发送方)进行关闭，因为向关闭chan写数据是panic的
func TestStopRoutine2_1(t *testing.T) {
	in := make(chan int)
	in2 := make(chan int)

	go func() {
		for {
			select {
			case x, ok := <-in:
				if !ok {
					fmt.Println("Process finsh")
					// 去掉后，当main关闭in时，这里会一直符合条件一直执行。。。。
					return
				}
				fmt.Printf("Process %d\n", x)
			case y, ok := <-in2:
				if !ok {
					fmt.Println("Process2 finsh")
					return
				}
				fmt.Printf("process2 %d\n", y)
			}
		}
	}()

	for i := 1; i < 2; i++ {
		fmt.Println("make i:", i)
		in <- i
		in2 <- i
		time.Sleep(2 * time.Second)
	}

	close(in)

	// ch关闭后，上面select中若没有return则fatal error: all goroutines are asleep - deadlock!，因为select一直满足第一个case，第二个不会执行
	//for i := 1; i < 5; i++ {
	//	fmt.Println("make i:", i)
	//	in2 <- i
	//	time.Sleep(2 * time.Second)
	//}

	fmt.Println("main finish")
	time.Sleep(2 * time.Second)
}

// for-select2，in关闭则设置nil，select不会关注nil的chan，其他case继续
func TestStopRoutine2_2(t *testing.T) {
	in := make(chan int)
	in2 := make(chan int)

	go func() {
		for {
			select {
			case x, ok := <-in:
				if !ok {
					fmt.Println("Process finsh set nil")
					in = nil
				}
				fmt.Printf("Process %d\n", x)
			case y, ok := <-in2:
				if !ok {
					fmt.Println("Process2 finsh")
					in2 = nil
				}
				fmt.Printf("process2 %d\n", y)
			}

			// 全部为nil则退出
			if in == nil && in2 == nil {
				fmt.Printf("all process2 finish")
				return
			}
		}
	}()

	for i := 1; i < 2; i++ {
		fmt.Println("make i:", i)
		in <- i
		in2 <- i
		time.Sleep(2 * time.Second)
	}

	close(in)

	for i := 1; i < 5; i++ {
		fmt.Println("make i:", i)
		in2 <- i
		time.Sleep(2 * time.Second)
	}

	fmt.Println("main finish")
	time.Sleep(2 * time.Second)
}

// 使用stopCh主动通知(send<->receive)都行。使用一个专门的通道，发送退出的信号
func TestStopRoutine3(t *testing.T) {
	in := make(chan int)
	stopCh := make(chan struct{})

	go func(stopCh <-chan struct{}) {
		go func() {
			defer fmt.Println("worker exit")
			for {
				select {
				// 关闭时收到信号
				case <-stopCh:
					fmt.Println("Recv stop signal")
					return
				case x, ok := <-in:
					if !ok {
						in = nil
					}
					fmt.Println("Working , ", x)
				}
			}
		}()
		return
	}(stopCh)

	for i := 1; i < 5; i++ {
		fmt.Println("make i:", i)
		in <- i
		time.Sleep(2 * time.Second)
	}

	close(stopCh)
	fmt.Println("main finish")
	time.Sleep(2 * time.Second)
}
