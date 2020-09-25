package routine

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/utils"
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
