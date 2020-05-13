package base

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"testing"
	"time"
)

func m1() {
	m2()
}
func m2() {
	m3()
}
func m3() {
	panic("panic from m3")
}
func a() {
	time.Sleep(time.Hour)
}

// 若输出所有goroutine信息,设置GOTRACEBACK=1，
// GOTRACEBACK=1 go run p.go
// 0、1、2，分别代表none、all、system
func TestOccurPanicPrintStack(t *testing.T) {
	go a()
	m1()
}

func m1_1() {
	m2_1()
}
func m2_1() {
	m3_1()
}
func m3_1() {
	time.Sleep(time.Hour)
}
func a_1() {
	time.Sleep(time.Hour)
}

// kill -SIGQUIT <pid> 杀死这个程序，程序在退出的时候输出strack trace
func TestOccurSIGQUITPrintStack(t *testing.T) {
	go a_1()
	m1_1()
}

func m1_2() {
	m2_2()
}
func m2_2() {
	m3_2()
}
func m3_2() {
	debug.PrintStack()
	time.Sleep(time.Hour)
}
func a_2() {
	time.Sleep(time.Hour)
}

// debug.PrintStack()方法可以将当前所在的goroutine的stack trace打印出来
func TestOccurDebugPrintStack(t *testing.T) {
	go a_2()
	m1_2()
}

func m1_3() {
	m2_3()
}
func m2_3() {
	m3_3()
}
func m3_3() {
	time.Sleep(time.Hour)
}
func a_3() {
	time.Sleep(time.Hour)
}
func setupSigusr1Trap() {
	c := make(chan os.Signal, 1)
	fmt.Println("1111")
	//signal.Notify(c, syscall.SIGUSR1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	fmt.Println("222")
	go func() {
		x := <-c
		fmt.Println("x:", x)
		for range c {
			DumpStacks()
		}
	}()
	fmt.Println("333")
}
func DumpStacks() {
	buf := make([]byte, 16384)
	buf = buf[:runtime.Stack(buf, true)]
	fmt.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===", buf)
}

// todo 没成功
// 使用runtime.Stack得到所有的goroutine的stack trace信息，其实debug.PrintStack()也是通过这个方法获得的。
func TestPrintAllRoutine2(t *testing.T) {
	setupSigusr1Trap()
	go a_3()
	m1_3()
}
