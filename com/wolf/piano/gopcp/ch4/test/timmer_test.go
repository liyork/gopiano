package test

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/utils"
	"testing"
	"time"
)

func TestTimerLate(t *testing.T) {
	timer := time.NewTimer(2 * time.Second)
	fmt.Printf("Present time: %v.\n", time.Now())

	// 即使晚了，再取，也可以取出
	time.Sleep(3 * time.Second)
	expirationTime := <-timer.C
	fmt.Printf("Expiration time: %v.\n", expirationTime)
	fmt.Printf("Stop timer: %v.\n", timer.Stop())
}

// 若定时器到期了，但由于某种原因未能及时地从它的C字段中接收元素值，那么该字段就会一直缓冲着那个元素值，即使在该定时器重置之后也如此。
//由于C的容量总是1，因此会影响重置后的定时器再次发送到期通知。后面的通知会被直接丢掉。因此，若想要复用定时器，就应该确保旧的通知已被接收。
func TestResetLate(t *testing.T) {
	timer := time.NewTimer(2 * time.Second)
	fmt.Printf("Present time: %v.\n", time.Now())

	// 缓冲中内容未被取出，即使reset之后，再取出还是之前的time
	time.Sleep(3 * time.Second)
	timer.Reset(2 * time.Second)
	fmt.Printf("now1:%v\n", time.Now())
	expirationTime := <-timer.C
	fmt.Printf("now2:%v\n", time.Now())
	fmt.Printf("Expiration time: %v.\n", expirationTime)
}

// time.AfterFunc，同样会返回新建的定时器，不过这样的定时器在到期时，不会向它的通知通道发送元素值，而是新启用一个goroutine执行调用方传入的函数，
//无论它是否被重置以及被重置多少次都这样
func TestAfterFunc(t *testing.T) {
	gid := utils.GetGID()
	fmt.Printf("current gid: %v\n", gid)
	timer := time.AfterFunc(2*time.Second, func() {
		gid1 := utils.GetGID()
		fmt.Printf("in AfterFunc gid: %v\n", gid1)
		fmt.Println("11111")
	})

	timer.Reset(2 * time.Second)
	//timer.Stop()

	gid2 := utils.GetGID()
	fmt.Printf("current gid2: %v\n", gid2)
	time.Sleep(6 * time.Second)
}
