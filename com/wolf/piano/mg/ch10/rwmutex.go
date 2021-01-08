package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var Password = secret{password: "myPassword"}

type secret struct {
	RWM      sync.RWMutex
	M        sync.Mutex
	password string
}

// `Change()` 函数修改共享变量，意味着需要使用一个排他锁，使用 `Lock()` 和 `Unlock()`
func Change(c *secret, pass string) {
	c.RWM.Lock()
	fmt.Println("LChange")
	time.Sleep(200 * time.Millisecond)
	c.password = pass
	c.RWM.Unlock()
}

// `show()` 函数使用 `RLock()` 和 `RUnlock()` 函数是因为它的关键部分是用来读取共享变量的
func show(c *secret) string {
	c.RWM.RLock()
	fmt.Println("show")
	time.Sleep(200 * time.Millisecond)
	defer c.RWM.RUnlock()
	return c.password
}

// showWithLock()` 函数为了读操作使用了排他锁，这意味着只有一个 `showWithLock()` 函数能读取 `secret` 结构体的 `password` 字段。
func showWithLock(c *secret) string {
	c.M.Lock()
	fmt.Println("showWithLock")
	time.Sleep(200 * time.Millisecond)
	defer c.M.Unlock()
	return c.password
}

// time go run rwMutex.go 10 >/dev/null
// time go run rwMutex.go >/dev/null
// > /dev/null是为了忽略这两个命令的输出。因此，使用 `sync.RWMutex` 互斥体的版本要比使用 `sync.Mutex` 的版本快很多。
func main() {
	var showFunction = func(c *secret) string { return "" }
	if len(os.Args) != 2 {
		fmt.Println("Using sync.RWMutex!")
		showFunction = show
	} else {
		fmt.Println("Using sync.Mutex!")
		showFunction = showWithLock
	}

	var waitGroup sync.WaitGroup
	fmt.Println("Pass:", showFunction(&Password))

	for i := 0; i < 10; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			fmt.Println("Go Pass:", showFunction(&Password))
		}()

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			Change(&Password, "123456")
		}()
		waitGroup.Wait()
		fmt.Println("Pass:", showFunction(&Password))
	}
}
