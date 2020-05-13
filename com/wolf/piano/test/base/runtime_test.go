package base

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func doSomething() {
	var i int
	for {
		i++
		time.Sleep(2 * time.Second)
	}
}

func TestRuntimeBase(t *testing.T) {
	cpu := runtime.NumCPU()
	fmt.Println("本机cpu核数:", cpu)

	//设置核数
	runtime.GOMAXPROCS(12)

	cpu = runtime.NumCPU()
	fmt.Println("本机cpu核数:", cpu)
}
