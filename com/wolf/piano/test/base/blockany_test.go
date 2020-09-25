package base

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
)

// 各种阻塞

//使用sync.WaitGroup,一直等待直到WaitGroup等于0(需要有done)
func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

//空select{}是一个没有任何case的select，它会一直阻塞
func TestSelect(t *testing.T) {
	select {}
}

// 死循环,虽然能阻塞，但会100%占用一个cpu。不建议使用
func TestForNil(t *testing.T) {
	//for {}
}

// 用sync.Mutex,一个已经锁了的锁，再锁一次会一直阻塞，这个不建议使用，不可重用
func TestLockRepeat(t *testing.T) {
	var m sync.Mutex
	m.Lock()
	m.Lock()
}

// os.Signal,系统信号量，在go里面也是个channel，在收到特定的消息之前一直阻塞
func TestSignal(t *testing.T) {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
}

// 空channel或者nil channel
func TestWaitForChan(t *testing.T) {
	// 正常channel会一直阻塞直到收到消息
	c := make(chan struct{})
	<-c

	//nil channel永远阻塞。
	var c2 chan struct{}
	<-c2

}
