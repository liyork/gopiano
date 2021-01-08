package bfbczd

import (
	"fmt"
	"sync"
	"testing"
)

func TestJoinPoint(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	sayHello := func() {
		defer wg.Done()
		fmt.Println("hello")
	}

	go sayHello()
	wg.Wait() //1 在这里加入连接点
}

//闭包围绕它们创建的词法范围，从而捕捉变量。
// 如果在goroutine中使用闭包，闭包是否在这些变量还是原始引用的副本上运行？
//事实证明，goroutine在它创建的同一地址空间内执行，因此我们的程序打印出“welcome”
func TestClosureErr1(t *testing.T) {
	var wg sync.WaitGroup
	salutation := "hello"
	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "welcome" // 1
	}()
	wg.Wait()
	fmt.Println(salutation)
}

// goroutine正在运行一个已经关闭迭代变量salutation的闭包，它有一个字符串类型。
// 当我们的循环迭代时，salutation被分配给切片中的下一个字符串值。
// 由于运行时调度器安排的goroutine可能会在将来的任何时间点运行，因此不确定在goroutine内将打印哪些值。
// 在goroutines开始之前，循环很可能会退出。 这意味着salutation变量超出了范围。 然后会发生什么？ goroutines仍然可以引用已经超出范围的东西吗？ 这个goroutine会访问可能已经被回收的内存吗？
//
// Go运行时足够敏锐地知道对salutation变量的引用仍然保留，因此会将内存传输到堆中，以便goroutine可以继续访问它。
//
// 循环在任何goroutines开始运行之前退出，所以salutation转移到堆中，并保存对字符串切片“good day”中最后一个值的引用。所以会看到“good day”打印三次
func TestClosureErr2(t *testing.T) {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation) // 测试打印字符串切片创建的循环变量salutation。
		}()
	}
	wg.Wait()
}

// 正确方法是将salutation的副本传递给闭包，以便在运行goroutine时，它将对参数进行操作：
func TestClosureRight(t *testing.T) {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(a string) { // 1
			defer wg.Done()
			fmt.Println(a)
		}(salutation) // 2 一个字符串的副本被创建
	}
	wg.Wait()
}
