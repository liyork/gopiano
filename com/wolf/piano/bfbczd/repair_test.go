package bfbczd

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

// 启动协程的函数，返回心跳chan
type startGoroutineFn func(done <-chan interface{}, pulseInterval time.Duration) (heartbeat <-chan interface{}) //1

var newSteward = func(timeout time.Duration, startGoroutine startGoroutineFn) startGoroutineFn { //2
	return func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
		heartbeat := make(chan interface{})

		go func() {
			defer close(heartbeat)

			var wardDone chan interface{}
			var wardHeartbeat <-chan interface{}
			startWard := func() { //3
				wardDone = make(chan interface{})                             //4
				wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2) //5
			}
			startWard()
			pulse := time.Tick(pulseInterval)

		monitorLoop:
			for { //6
				timeoutSignal := time.After(timeout)
				for {
					select {
					case <-pulse: // 本身自己的心跳，发送出去
						select {
						case heartbeat <- struct{}{}:
						default:
						}
					case <-wardHeartbeat: //7，目标协程的心跳，跳转continue重新计时
						continue monitorLoop
					case <-timeoutSignal: //8，超时
						log.Println("steward: ward unhealthy; restarting")
						close(wardDone)
						startWard()
						continue monitorLoop
					case <-done:
						return
					}
				}
			}
		}()
		return heartbeat
	}
} //1.这里我们定义一个可以监控和重新启动的goroutine的函数签名。
//2.在这里我们设置了超时时间，并使用函数startGoroutine来启动它正在监控的goroutine。有趣的是，监控器本身返回一个startGoroutineFn，表示监控器自身也是可监控的。
//3.在这里我们定义一个闭包，它以同样的的方式来启动我们正在监视的goroutine。
//4.这是我们创建一个新通道，我们会将其传递给监控通道，以响应发出的停止信号。
//5.在这里，我们开启对目标goroutine的监控。如果监控器停止工作，或者监控器想要停止被监控区域，我们希望监控者也停止，因此我们将两个done通道都包含在逻辑中。我们传入的心跳间隔是超时时间的一半，这可以调整。
//6.这是我们的内部循环，它确保监控者可以发出自己的心跳。
//7.在这里我们如果接收到监控者的心跳，就会知道它还处于正常工作状态，程序会继续监测循环。
//8.这里如果我们发现监控者超时，我们要求监控者停下来，并开始一个新的goroutine。然后开始新的监测。

func or(wardDone chan interface{}, done <-chan interface{}) <-chan interface{} {
	return nil
}

// 如果监控一个行为异常的goroutine，会发生什么：
func TestRepair(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWork := func(done <-chan interface{}, _ time.Duration) <-chan interface{} {
		log.Println("ward: Hello, I'm irresponsible!")
		go func() {
			<-done // 1
			log.Println("ward: I am halting.")
		}()
		return nil
	}

	doWorkWithSteward := newSteward(4*time.Second, doWork) // 2

	done := make(chan interface{})
	time.AfterFunc(9*time.Second, func() { // 3
		log.Println("main: halting steward and ward.")
		close(done)
	})

	for range doWorkWithSteward(done, 4*time.Second) { // 4
	}

	log.Println("Done")
	//1.可以看到这个goroutine什么都没干，持续阻塞等待被取消，它同样不会发出任何表明自己正常信号。
	//2.这里开始建立被监控的例程，其4秒后会超时。
	//3.这里我们9秒后向done通道发出信号停止整个程序。
	//4.最后，我们启动监控器并在其心跳范围内防止示例停止。
}

// 用闭包强化一下：
var doWorkFn = func(done <-chan interface{}, intList ...int) (startGoroutineFn, <-chan interface{}) { //1
	intChanStream := make(chan (<-chan interface{})) //2
	intStream := bridge(done, intChanStream)

	doWork := func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} { //3
		intStream := make(chan interface{}) //4
		heartbeat := make(chan interface{})

		go func() {
			defer close(intStream)
			select {
			case intChanStream <- intStream: //5
			case <-done:
				return
			}

			pulse := time.Tick(pulseInterval)

			for {
			valueLoop:
				for _, intVal := range intList {
					if intVal < 0 {
						log.Printf("negative value: %v\n", intVal) //6
						return
					}

					for {
						select {
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						case intStream <- intVal:
							continue valueLoop
						case <-done:
							return
						}
					}
				}
			}
		}()
		return heartbeat
	}
	return doWork, intStream
} //1.我们将监控器关闭的内容放入返回值，并返回所有监控器用来交流数据的通道。
//2.我们建立通道的通道。
//3.这里我们建立闭包控制监控器的启动和关闭。
//4.这是各通道与监控器交互数据的实例。
//5.这里我们向起数据交互作用的通道传入数据。
//6.这里我们返回负数并从goroutine返回以模拟不正常的工作状态。

func TestRepair2(t *testing.T) {
	log.SetFlags(log.Ltime | log.LUTC)
	log.SetOutput(os.Stdout)

	done := make(chan interface{})
	defer close(done)

	doWork, intStream := doWorkFn(done, 1, 2, -1, 3, 4, 5)      //1
	doWorkWithSteward := newSteward(1*time.Millisecond, doWork) //2
	doWorkWithSteward(done, 1*time.Hour)                        //3

	for intVal := range take(done, intStream, 6) { //4
		fmt.Printf("Received: %v\n", intVal)
	}
	//1.这里我们调用该函数，它会将传入的不定长整数参数转换为可通信的流。
	//2.在这里，我们创建了一个检查doWork关闭的监视器。我们预计这里会极快的进入失败流程，所以将监控时间设置为一毫秒。
	//3.我们通知 steward 开启监测。
	//4.最后，我们使用该管道，并从intStream中取出前六个值。
}
