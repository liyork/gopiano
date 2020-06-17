package concurrent

import (
	"fmt"
	"testing"
	"time"
)

// 问题：前9个协程调用Wait()阻塞，第10个调用后10个协程全部唤醒
// 约束：不能用任何同步相关的操作，但可以用chan，前提是无缓冲模式的。

// 分析：
// chan没有缓冲的生产和消费必须同时调用，二者都会被唤醒，否则任意一方都会被阻塞
// 但是总数如何统计，chan是无缓冲，不能用cap。因为goroutine执行，而不用同步操作，count势必会产生竞态条件。只能从chan中读取计数

// 抽象一个栅栏
type Barrier interface {
	Wait()
}

// 栅栏的实现类
type barrier struct {
	chCount chan struct{} // 写入，负责记数
	chSync  chan struct{} // 阻塞、唤醒
	count   int
}

// 创建栅栏对象
func NewBarrier(n int) Barrier {
	b := &barrier{count: n, chCount: make(chan struct{}), chSync: make(chan struct{})}
	go b.Sync()
	return b
}
func (b *barrier) Wait() {
	b.chCount <- struct{}{}
	<-b.chSync // 阻塞
}

func (b *barrier) Sync() {
	count := 0
	for range b.chCount { // 有新数据则唤醒
		fmt.Println("count:", count, ",b.count:", b.count)
		count++
		if count >= b.count {
			close(b.chSync) // close这个chan所有阻塞协程都会被激活
			break
		}
	}
}

func TestBarrier(t *testing.T) {
	newBarrier := NewBarrier(10)
	for x := 0; x < 10; x++ {
		time.Sleep(2 * time.Second)
		go newBarrier.Wait()
	}
}
