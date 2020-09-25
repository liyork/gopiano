package concurrent

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"testing"
	"time"
)

var a1 string
var once sync.Once
var wg1 sync.WaitGroup

func setup() {
	time.Sleep(time.Second * 2) //1
	a1 = "hello, world"
	fmt.Println("setup over") //2
}

func doprint() {
	once.Do(setup)         //3
	fmt.Println("a1:", a1) //4
	wg1.Done()
}

func twoprint() {
	go doprint()
	go doprint()
}

func TestRightSync(t *testing.T) {
	wg1.Add(2)
	twoprint()
	wg1.Wait()
}

var done bool

func setup2() {
	a1 = "hello, world"
	done = true
}

func doprint2() {
	if !done {
		fmt.Println("doprint2...")
		once.Do(setup2)
	}
	println(a1)
	wg1.Done()
}

func twoprint2() {
	go doprint2()
	go doprint2()
}

//并不能保证一定输出hello, world，而可能输出空字符串，
// 因为在doPrint2函数内即使可以能够看到setup中对done变量的写操作，也不能保证在doPrint里面看到对变量a的写操作。因为两者并没有happenbefore同步机制
// 与TestRightSync不同在，这里可能另一个协成看到done进而不阻塞在once.Do，也就导致直接打印a1
func TestErrorSync1(t *testing.T) {
	// http://localhost:6060/debug/pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	wg1.Add(2)
	twoprint2()
	wg1.Wait()
}

func setup3() {
	a1 = "hello, world"
	done = true
}

// 并不能确保在TestErrorSync2函数内既看到对变量done的写操作，也可以看到对变量a的操作，所以TestErrorSync2函数还是可能会输出空串。
// 更糟糕的是由于两个goroutine没有对变量done做同步措施，TestErrorSync2函数所在goroutine可能看不到对done的写操作，从而导致main函数所在goroutine一直运行在for循环出。java中使用的是volatile
func TestErrorSync2(t *testing.T) {
	go setup3()
	for !done {
	}
	println(a)
}

type T struct {
	msg string
}

var g *T

func setup4() {
	t := new(T)
	t.msg = "hello, world"
	g = t
}

// 即使TestErrorSync3函数内可以看到setup函数内对g的赋值，从而让main函数退出，
// 但是也不能保证TestErrorSync3函数可以看到对 g.msg的赋值，也就是可能输出空串
func TestErrorSync3(t *testing.T) {
	go setup4()
	for g == nil {
	}
	println(g.msg)
}

// 总结，对于单个内存，有多个routine进行读写时记得要保证原子、可见性，正确的同步是需要的
