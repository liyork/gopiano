package test

import (
	"fmt"
	"testing"
	"time"
)

// Timer经过指定的时间后触发一个事件，只执行一次，这个事件通过其本身提供的channel进行通知

// 基本使用，延迟2s后执行
func TestTimer(t *testing.T) { // Test还一定要大写而且在前面
	// 一经创建便开始计。实际上，创建Timer意味着把一个计时任务交给系统守护协程，该协程管理着所有的Timer，
	// 当Timer的时间到达后向Timer的管道中发送当前的时间作为事件
	timer := time.NewTimer(2 * time.Second) //

	// 延时到则会写入chan，下面得以执行
	a := <-timer.C // 得到当前时间
	fmt.Println("time arrive,receive:", a)
}

// 超时设置
func TestWaitChan(t *testing.T) {
	readChan := make(chan int)
	timer := time.NewTimer(2 * time.Second)

	select { // 轮询readChan和timer.C两个管道，有任意可用则退出
	case <-readChan:
		fmt.Println("receive chan")
	case <-timer.C:
		fmt.Println("timeout")
	}
}

// 退出执行后，延迟执行某个方法
func TestDelayRun(t *testing.T) {
	timer := time.NewTimer(2 * time.Second)
	defer func() {
		<-timer.C
		fmt.Println("delay run")
	}()

	fmt.Println("end method TestDelayRun")
}

func TestStop(t *testing.T) {
	timer := time.NewTimer(2 * time.Second)
	stopSucc := timer.Stop()
	fmt.Println("stopSucc:", stopSucc) // true: 定时器超时前停止，后续不会再有事件发送

	stopSucc = timer.Stop()
	fmt.Println("stopSucc:", stopSucc) // false: 定时器超时后停止
}

// 重置的动作实质上是先停掉定时器，再启动。其返回值也即停掉计时器的返回值。
func TestReset(t *testing.T) {
	timer := time.NewTimer(1 * time.Second)
	timer.Reset(5)

	fmt.Println("reset")
}

// 时间后出发，便捷方法，没有需求提前停止定时器，也没有需求复用该定时器，那么可以使用匿名的定时器
func TestAfter(t *testing.T) {
	fmt.Println(time.Now())
	<-time.After(2 * time.Second)
	fmt.Println(time.Now())
}

// 便捷延迟方法，异步执行
func TestAfterFunc(t *testing.T) {
	stopChan := make(chan int)

	time.AfterFunc(2*time.Second, func() {
		fmt.Println("after 2 second func...")
		stopChan <- 1
	})

	<-stopChan // 由于AfterFunc是异步，所以等待协成退出，不然单元测试直接退出。
}

// 由于被time.Sleep阻塞，所以下次timer时间就是这个5s
func TestDelay(t *testing.T) {
	duration := 2 * time.Second
	timer := time.NewTimer(duration)

	for {
		select {
		case <-timer.C:
			fmt.Println("1111,", time.Now())
			timer.Reset(duration)
			time.Sleep(5 * time.Second)
		}
	}
}

func TestTimerImediate(t *testing.T) {
	duration := 5 * time.Second
	timer := time.NewTimer(1) //
	for {
		select {
		case <-timer.C:
			fmt.Println("1111111111")
			timer.Reset(duration)
		}
	}
}
