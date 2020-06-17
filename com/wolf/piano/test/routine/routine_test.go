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
