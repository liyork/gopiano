package concurrent

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

// select是执行选择操作的一个结构，它里面有一组case语句，它会执行其中无阻塞的那一个，如果都阻塞了，那就等待其中一个不阻塞，
// 进而继续执行，它有一个default语句，该语句是永远不会阻塞的
// 每个 case 语句里必须是一个IO操作
// select 可以监听channel 上的数据流动
// select 默认是阻塞的，只有当监听的channel 中有发送或接收可以进行时才会运行，当多个channel 都准备好的时候，select 是随机的选择一个执行的。

//select
//select 先遍历所有case, 所有channel表达式都会被求值、所有被发送的表达式都会被求值。求值顺序：自上而下、从左到右.
//当case没有阻塞则随机执行一个没有阻塞的case就退出select
//当所有case阻塞时, 则一直阻塞直到某个case解除阻塞, 但是如果有default则直接执行default
//也就是一个select最多只执行一次case里的代码
//要一直检测case则必须外层使用for循环包起来

// select在多个通道上进行读或写操作，让函数可以处理多个事情，但1次只处理1个。
//特性
//每次执行select，都会只(随机)执行其中1个case或者执行default语句
//当没有case或者default可以执行时，select则阻塞，等待直到有1个case可以执行
func TestSelectBase(t *testing.T) {
	readCh := make(chan int, 1)
	writeCh := make(chan int, 1)

	y := 1
	select {
	case x := <-readCh:
		fmt.Printf("Read %d\n", x)
	case writeCh <- y:
		fmt.Printf("Write %d\n", y)
	default:
		fmt.Println("Do what you want")
	}
}

func TestNoBufferReadWitSelect(t *testing.T) {
	noBuffChan := make(chan int)
	readWithSelectDefault(noBuffChan)
}

func TestBufferReadWithSelect(t *testing.T) {
	buffChan := make(chan int, 2)
	readWithSelectDefault(buffChan)
}

// case都阻塞则执行default
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

func TestTickBase(t *testing.T) {
	tick := time.Tick(2000 * time.Millisecond)
	after := time.After(5000 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick now", time.Now())
		case <-after:
			fmt.Println("tick finish", time.Now())
			return
		default:
			fmt.Println("tick default", time.Now())
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func TestTimeout(t *testing.T) {
	timeout := time.After(5000 * time.Millisecond)
	bizchan := make(chan int)
	stop := make(chan int)
	go func() {
		for {
			select {
			case <-timeout:
				fmt.Println("timeout...", time.Now())
				stop <- 1
				break
			case <-bizchan:
				fmt.Println("biz chann", time.Now())
			}
		}
	}()
	<-stop
	fmt.Println("main finish")
}

func TestHeartbeat(t *testing.T) {
	heartbeat := time.Tick(3 * time.Second)
	for {
		select {
		case <-heartbeat:
			fmt.Println("hb")
		}
	}
}

func TestDirection(t *testing.T) {
	bothwaychan := make(chan int)

	go func() {
		// 转成单向限制写入
		write := chan<- int(bothwaychan)
		write <- 1
	}()

	read := <-chan int(bothwaychan)
	fmt.Println("read:", read)
}

// break仅能跳出select，不能跳出for
func TestBreakError(t *testing.T) {
	inCh := make(chan int)
	go func(inCh <-chan int) {
		i := 0
		for {
			fmt.Printf("for: %d\n", i)
			select {
			case x, open := <-inCh:
				if !open {
					fmt.Println("xxxxxx")
					break
				}
				fmt.Printf("read: %d\n", x)
			}
			fmt.Printf("end select: %d\n", i)
			i++
		}
	}(inCh)

	fmt.Println("executed consume")
	time.Sleep(2 * time.Second)
	close(inCh)
	time.Sleep(111111 * time.Second)
}

// 使用return结束select
func TestReturn(t *testing.T) {
	inCh := make(chan int)
	go func(inCh <-chan int) {
		i := 0
		for {
			fmt.Printf("for: %d\n", i)
			select {
			case x, open := <-inCh:
				if !open {
					fmt.Println("xxxxxx")
					return
				}
				fmt.Printf("read: %d\n", x)
			}
			fmt.Printf("end select: %d\n", i)
			i++
		}
	}(inCh)

	fmt.Println("executed consume")
	time.Sleep(2 * time.Second)
	close(inCh)
	fmt.Println("executed consume")
	time.Sleep(111111 * time.Second)
}

// 一直阻塞
func TestSelectNoStatement(t *testing.T) {
	select {}
	//与下面等价
	ch := make(chan int)
	<-ch
}

// 理论上即使receive很快，但是也能感知到stop
// 但是这里为什么会卡在ch <- i???而且跑了一会6060也没响应了
func TestSelectAwareStop(t *testing.T) {
	// http://localhost:6060/debug/pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	ch := make(chan int)
	stop := make(chan struct{})

	go func() {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("TestSelectAwareStop error: %v\n", err)
				//debug.PrintStack()
			}
		}()
		for i := 1; i < 1000000; i++ {
			//fmt.Println("send", )
			ch <- i // 太快了会卡在这里?，也并没有遇到panic.，写不进去，那就是没有读了?
		}
	}()

	go func() {
		for {
			select {
			case i := <-ch: //这里并没有打印出stop或者receive，卡在了哪里?
				fmt.Println("receive:", i)
			case <-stop:
				fmt.Println("receive stop event")
				return
			}
		}
	}()

	time.Sleep(1 * time.Second)
	close(stop)

	time.Sleep(55 * time.Second)
}

// 测试case没有body的时候是否走下一个，不会走下一个，自带break
func TestSelectUseOne(t *testing.T) {
	readCh := make(chan int)
	writeCh := make(chan int, 1)

	select {
	case x := <-readCh:
		fmt.Printf("Read %d\n", x)
	case writeCh <- 1:
		//fmt.Printf("write %d\n", 1)
	default:
		fmt.Println("Do what you want")
	}
}
