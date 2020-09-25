package base

import (
	"errors"
	"fmt"
	"runtime/debug"
	"testing"
)

// Go 语言官方表态不要轻易使用 panic recover，除非你真的无法预料中间可能会发生的错误，或者它能非常显著地简化你的代码。
// 在一个常见的 Web 应用中，不能因为个别 URL 处理器抛出异常而导致整个程序崩溃，就需要在每个 URL 处理器外面包括一层 recover() 来恢复异常。
// 在 json 序列化过程中，逻辑上需要递归处理 json 内部的各种类型，每一种容器类型内部都可能会遇到不能序列化的类型。如果对每个函数都使用返回错误的方式来编写代码，会显得非常繁琐。所以在内置的 json 包里也使用了 panic，然后在调用的最外层包裹了 recover 函数来进行恢复，最终统一返回一个 error 类型。

//Recover内建函数，可以让进入宕机流程中的 goroutine 恢复过来，recover 仅在延迟函数 defer 中有效，在正常的执行过程中，调用 recover 会返回 nil 并且没有其他任何效果，

func TestPanic(t *testing.T) {
	testOccurPanic()
	fmt.Println("111")
}

// Golang中引入两个内置函数panic和recover来触发和终止异常处理流程
// 当程序运行时，如果遇到引用空指针、下标越界或显式调用panic函数等情况，则先触发panic函数的执行，然后调用延迟函数。
// 继续向上传递panic，因此该过程一直在调用栈中重复发生：函数停止执行，调用延迟执行函数等。
// 如果一路在延迟函数中没有recover函数的调用，则会到达该携程的起点，该携程结束，然后终止其他所有携程，包括主携程。
func TestRecover(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("testOccurPanic error: %v\n", err)
			//debug.PrintStack()
		}
	}() // 对匿名函数进行了调用

	testOccurPanic()
	// 异常点后面的逻辑也不会再继续执行了
	fmt.Println("111")
}

func testOccurPanic() {
	panic(22222)
}

func TestPassError(t *testing.T) {
	err := funcA()
	if err == nil {
		fmt.Printf("err is nil\\n")
	} else {
		fmt.Printf("err is %v\\n", err)
	}
}

// panic异常处理机制不会自动将错误信息传递给error，所以要在funcA函数中进行显式的传递，代码如下所示：
func funcA() (err error) {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic recover! p:", p)
			str, ok := p.(string)
			if ok {
				err = errors.New(str)
			} else { // 或者再抛出 panic(err)  // rethrow
				err = errors.New("panic")
			}
			debug.PrintStack()
		}
	}()
	return funcB()
}

func funcB() error {
	// simulation
	panic("foo")
	return errors.New("success")
}

//但不是所有panic都能通过recover()捕捉到的，比如：并发操作map实例。
//不能单纯依靠recover()解决函数内部所有panic异常，应该做到以下几点：
//a) 通过编写代码校验，防止能预期到的panic，比如：空指针引用的指针判断。
//b) 对于无法预期的panic，使用recover()捕获并加以处理。
//c) 使用map时，必须要考虑是否存在并发读写场景，存在时，应使用ConcurrentMap组件或自己加sync.RWMutex进行加锁保护。
func TestConcurrentIsNotRecover(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("TestConcurrentIsNotRecover error: %v\n", err)
		}
	}()

	m := make(map[int]int)
	go func() {
		for {
			m[0] = 0
		}
	}()
	for {
		fmt.Println("m[0]:", m[0])
	}
}
