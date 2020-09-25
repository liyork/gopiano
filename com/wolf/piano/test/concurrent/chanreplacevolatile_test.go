package concurrent

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// go中没有C和Java中volatile的强行刷缓存的概念，所以只能用锁和chan去实现同步，而不能依赖于共享全局变量

var inited = false

func Setup() {
	time.Sleep(time.Second)
	inited = true
}

// go test -race -run TestRace chanreplacevolatile_test.go
// -race 用于检测代码中可能存在的 data race。
// 这里的告警信息写得很清楚，一个 goroutine 的写操作和另一个 goroutine 的读操作可能发生 data race，有风险。
func TestRace(t *testing.T) {
	go Setup()

	for {
		if inited {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("setup succeed")
}

var lock sync.Mutex

func Setup1() {
	time.Sleep(time.Second)
	// 加锁写
	lock.Lock()
	inited = true
	lock.Unlock()
}

// 当某个变量存在并发访问的可能时，请一定记得加锁：
// go test -race -run TestResolveRaceByLock  chanreplacevolatile_test.go
func TestResolveRaceByLock(t *testing.T) {
	go Setup1()

	for {
		// 加锁读
		lock.Lock()
		b := inited
		lock.Unlock()
		if b {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("setup succeed")
}

var c = make(chan bool)

func Setup2() {
	time.Sleep(time.Second * 3)
	c <- true
}

// 在 c/c++以及 java 中，习惯设置共享对象或基本类型，然后在不同的线程中对它进行读写操作从而达到线程间通信的目的。
// 这样做的好处就是操作简单，坏处也很明显：容易造成 data race，而为了避免 data race 又不得不加锁，加锁又降低了并发性能。
// golang 提出了一个新的理念： Do not communicate by sharing memory; instead, share memory by communicating.
// 无锁的版本
//事实上使用 chan 的好处不仅在于解决了并发访问的 data race 和锁的问题，
// 而且还提高了代码的运行效率——chan 的接收端 goroutine 会被挂起 直到 chan 中有值可读，相比于传统方法的循环检测共享变量，这种方式效率明显要高不少，且优雅。
// go test -race -run TestResolveRaceByChan  chanreplacevolatile_test.go
func TestResolveRaceByChan(t *testing.T) {
	go Setup2()
	if <-c {
		fmt.Println("setup succeed")
	}
}
