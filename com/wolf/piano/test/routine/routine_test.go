package routine

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/utils"
	"runtime"
	"testing"
	"time"
)

//操作系统并不知道goroutine的存在，goroutine的调度是由Golang运行时进行管理的
//Golang有自己的调度器，许多goroutine的数据都是共享的，因此goroutine之间的切换会快很多，启动goroutine所耗费的资源也很少
//goroutine是Go语言中的轻量级线程实现，由Go运行时（runtime）管理.goroutine 比thread 更易用、更高效、更轻便
//Go语言提供的是另一种通信模型，即以消息机制而非共享内存作为通信方式
//Go语言提供的消息通信机制被称为channel
//不要通过共享内存来通信，而应该通过通信来共享内存

// goroutine到底什么时候能执行到呢？golang的调度器是非抢占式的，在GPM的架构里，Waiting态的goroutine（在LRQ或者GRQ中等待执行）
// 必须要等Executing态的goroutine主动退出执行，才能绑定到M上执行。主动退出的方式有很多，比如time.Sleep()，比如runtime.Gosched()。

// Go语言从语言层面上就支持了并发
//
//Go语言中有个概念叫做goroutine, 这类似我们熟知的线程，但是更轻。

// 程序默认不等所有 goroutine 都执行完才退出

func TestGetGoRoutineId(t *testing.T) {
	go testg1()
	go testg1()
	go testg1()
	time.Sleep(time.Second * 5)
}

func testg1() {
	fmt.Println("getgid:", utils.GetGID())
	time.Sleep(time.Second * 50)
}

type testv struct {
	age int
}

func (t *testv) test() {
	fmt.Println("age:", t.age)
}

// 由于闭包只是绑定到这个value变量上，并没有被保存到goroutine栈中，代码极有可能运行的结构都输出为切片的最后一个元素。
// 因为这样写会导致for循环结束后才执行goroutine多线程操作，这时候value值只指向了最后一个元素。
// go多协程编程，共享变量是一种极度不可靠的方式
func TestGoroutineErr(t *testing.T) {
	values := []*testv{{1}, {2}, {3}}
	for _, value := range values {
		go func() {
			fmt.Println(value)
		}()
	}
}

// 将 val 作为一个参数传入 goroutine 中，每个 val 都会被独立计算并保存到 goroutine 的栈中
func TestGoroutineCorrect1(t *testing.T) {
	values := []*testv{{1}, {2}, {3}}
	for _, value := range values {
		go func(val interface{}) {
			fmt.Println(val)
		}(value)
	}
}

// 在循环内定义新的变量，由于在循环内定义的变量在循环遍历的过程中是不共享的
func TestGoroutineCorrect2(t *testing.T) {
	values := []*testv{{1}, {2}, {3}}
	for _, value := range values {
		val := value
		go func() {
			fmt.Println(val)
		}()
	}
}

//协程通常称为coroutine，在golang中称为goroutine。
//协程本质上是一种用户态线程，它不需要操作系统来进行抢占式调度，在实际实现中寄存在线程之中。
//协程系统开销极小，可有效提高单个线程的任务并发性，从而避免使用多线程。而且使用协程编程简单，结构清晰。缺点就是需要编程语言的支持，如果不支持，则需要用户在程序中自行实现调度器。

//goroutine是golang中的轻量级线程实现，由go runtime管理
//启动一个goroutine使用go关键字，go是golang中最重要的关键字，因此这个语言也是以这个关键字命名。
//在一个函数前加上go关键字调用，这次调用就会在一个新的goroutine中并发执行，开启goroutine的线程将继续执行。
//当被go调用的函数返回时，这个goroutine也自动结束了。如果这个函数有返回值，那么这个返回值会被丢弃。
//golang程序从main()函数开始执行，当main()函数返回时，程序结束且不等待其他goroutine结束。

// 并发与并行
// 两个队列，一个Coffee机器，那是并发
// 两个队列，一个Coffee机器，那是并发

// golang并发与并行
//并发执行
//如果在单核cpu情况下，golang所有的goroutine只能在一个线程里跑
//如果当前goroutine不发生阻塞，它是不会让出cpu时间给其他goroutine，除非调用runtime.Gosched()主动让出时间片
//如果当前goroutine发生阻塞，它会主动让出cpu时间给其他goroutine执行
//golang的runtime包是goroutine的调度器，其中使用runtime.GOMAXPROCS(n)可以控制使用cpu核数
func loop() {
	fmt.Println("loop starts!")
	for i := 0; i < 101; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()
}

func TestConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(1) // 强制使用1个cpu

	go loop()
	go loop()

	time.Sleep(3 * time.Second)
}

// 并行执行
// 默认情况下，golang是开启多核的，因此一开始的示例中的输出就是无序的
//也可以告诉golang我们允许同时最多使用核数
func TestParallel(t *testing.T) {
	runtime.GOMAXPROCS(2)

	go loop()
	go loop()

	time.Sleep(3 * time.Second)
}
