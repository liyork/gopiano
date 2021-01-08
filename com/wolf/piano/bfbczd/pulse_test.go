package bfbczd

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestPulseBase(t *testing.T) {
	// 会产生心跳的goroutine
	var doWork = func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{}) //1
		results := make(chan time.Time)

		go func() {
			// done时关闭chan
			defer close(heartbeat)
			defer close(results)

			// 定期心跳，表明自己ok
			pulse := time.Tick(pulseInterval) //2
			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default: //4
				}
			}

			//注意，由于我们可能在等待输入/发送结果时发送多个pulse所有select语句都需要在for循环内
			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return
					case <-pulse: //5
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			workGen := time.Tick(2 * pulseInterval) //3
			for {
				select {
				case <-done:
					return
				case <-pulse: //5
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	} //1.在这里，我们设置了一个发送心跳信号的通道。doWork会返回该通道。
	//2.我们按传入的pulseInterval值定时发送心跳，每次心跳都意味着可以从该通道上读取到内容。
	//3.这只是用来模拟进入的工作的另一处代码。我们选择一个比pulseInterval更长的持续时间，以便我们可以看到来自goroutine的心跳。
	//4.请注意，我们包含一个default子句。我们必须考虑如果没有人接受到心跳的情况。从goroutine发出的结果是至关重要的，但心跳不是。
	//5.就像done通道，无论何时执行发送或接收，你都需要考虑心跳发送的情况。

	// using，主要接收结果，可以顺便通过心跳监测目标工作状况，若是超时则返回
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) }) //1

	const timeoutInterval = 2 * time.Second               //2
	heartbeat, results := doWork(done, timeoutInterval/2) //3
	for {
		select {
		case _, ok := <-heartbeat: //4
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results: //5
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeoutInterval): //6
			fmt.Println("timeout ...")
			return
		}
	} //1.我们设置done通道并在10秒后关闭它。
	//2.我们在这里设定超时时间 我们将用它将心跳间隔与超时时间相耦合。
	//3.我们向dowork传入超时时间的一半。--心跳
	//4.我们将hearbeat的读取放入select语句中。每间隔 timeout/2 获取一次来自心跳通道的消息。如果我们没有收到消息，那就说明该goroutine存在问题。
	//5.我们从result通道获取数据。
	//6.如果我们没有收到心跳或result，程序就会超时结束。
}

// 将在两次迭代后停止goroutine来模拟循环中断
func TestPulseExample(t *testing.T) {
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		go func() {
			// 不关闭任何一个通道

			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}
			sendResult := func(r time.Time) {
				for {
					select {
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for i := 0; i < 2; i++ { //1
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}

	// using
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r)
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
	//1.这里我们简单模拟循环中断。TestPulseBase中，未收到通知会无限循环。而这里我们只循环两次。
	//在两秒钟之内，我们的系统意识到goroutine未能正确读取，并且打破了for-select循环。通过使用心跳，我们已经成功地避免了死锁，并且不必通过依赖较长的超时而保持稳定性
	//请注意，心跳会帮助处理相反的情况：它让我们知道长时间运行的goroutine依然存在，但花了一段时间才产生一个值并发送至通道。
}

// 场景：工作单元开始时产生心跳。这对测试非常有用
func TestPulseExample2(t *testing.T) {
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		heartbeatStream := make(chan interface{}, 1) //1
		workStream := make(chan int)
		go func() {
			defer close(heartbeatStream)
			defer close(workStream)

			for i := 0; i < 10; i++ {
				select { //2
				case heartbeatStream <- struct{}{}:
				default: //3
				}

				select {
				case <-done:
					return
				case workStream <- rand.Intn(10):
				}
			}
		}()

		return heartbeatStream, workStream
	}

	done := make(chan interface{})
	defer close(done)

	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
	//1.这里我们用一个缓冲区创建心跳通道。这确保即使没有人及时监听发送，也总会发送至少一个pulse。表明开始
	//2.在这里，我们为心跳设置了一个单独的select块。我们不希望将它与发送结果一起包含在同一个select块中，因为若在一起时，如果接收器未准备好，它们将接收到一个pulse，而result的当前值将会丢失。我们也没有为心跳提供done通道提供case语句，因为我们有一个default可以处理这种情况。
	//3.我们再次处理如果没有人监听到心跳。因为我们的心跳通道是用缓冲区创建的，如果有人在监听，但没有及时处理第一个心跳，仍会被通知。
	//如预期一致，每个结果都会有一个心跳。
}

// 测试的编写
// 简单的生成器，它将传入的数字转换为它返回通道上的数据流
func DoWork(done <-chan interface{}, nums ...int) (<-chan interface{}, <-chan int) {

	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)

	go func() {

		defer close(heartbeat)
		defer close(intStream)
		time.Sleep(2 * time.Second) // 1

		for _, n := range nums {
			select {
			case heartbeat <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()

	return heartbeat, intStream
} //1.我们在goroutine开始工作之前模拟延迟。在实践中，延迟可以由各种各样的原因导致，例如CPU负载，磁盘争用，网络延迟和bug。

// 测试的反例：因为它的不确定性。
// 如果移除time.Sleep情况会变得更糟：这个测试会有时通过，有时失败。
// 我们之前提到过程中的外部因素可能会导致goroutine花费更长的时间才能完成第一次迭代。关键在于我们不能保证在超时之前第一次迭代会完成，
// 所以我们开始考虑：这时候超时会有多大意义？我们可以增加超时时间，但这意味着测试时失败也需要很长时间，从而减慢我们的测试效率。
func TestDoWork_GeneratesAllNumbersErr(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	_, results := DoWork(done, intSlice...)

	for i, expected := range intSlice {
		select {
		case r := <-results:
			if r != expected {
				t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
			}
		case <-time.After(1 * time.Second): // 1
			t.Fatal("test timed out")
		}
	}
} //1.在这里，我们设置超时，以防止goroutine出现问题导致死锁。

// 正确的测试：
func TestDoWork_GeneratesAllNumbersRight(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	heartbeat, results := DoWork(done, intSlice...)

	<-heartbeat //1

	i := 0
	for r := range results {
		if expected := intSlice[i]; r != expected {
			t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
		}
		i++
	}
} //1.在这里，我们等待goroutine发出信号表示它正在开始处理迭代

// 使用基于间隔的心跳的测试示例

func DoWork2(done <-chan interface{}, pulseInterval time.Duration, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(intStream)
		time.Sleep(2 * time.Second)
		pulse := time.Tick(pulseInterval)
	numLoop: //2
		for _, n := range nums {
			for { //1
				select {
				case <-done:
					return

				case <-pulse:
					select {
					case heartbeat <- struct{}{}:
					default:
					}
				case intStream <- n:
					// 跳到numLoop继续for
					continue numLoop //3
				}
			}
		}
	}()

	return heartbeat, intStream
}

func TestDoWork_GeneratesAllNumbers(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 2 * time.Second
	heartbeat, results := DoWork2(done, timeout/2, intSlice...)

	<-heartbeat //4

	i := 0
	for {
		select {
		case r, ok := <-results:
			if ok == false {
				return
			} else if expected := intSlice[i]; r != expected {
				t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
			}
			i++
		case <-heartbeat: //5
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	}
} //1.我们需要两个循环：一个用来覆盖我们的数字列表，并且这个内部循环会运行直到intStream上的数字成功发送。
//2.我们在这里使用一个标签来使内部循环继续更简单一些。
//3.这里我们继续执行外部循环。
//4.我们仍然等待第一次心跳出现，表明我们已经进入了goroutine的循环。
//5.我们在这里获取心跳以实现超时。
