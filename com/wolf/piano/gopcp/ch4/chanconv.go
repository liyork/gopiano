package main

import "fmt"

// go不允许程序关闭接收通道。
// 单向通道通常由双向通道转换而来，但是单向通道不可以转换回双向通道。通道允许的数据传递方向是其类型的一部分。
// 对于通道，数据传递方向的不同意味着他们的类型不同。
// 利用函数声明将双向通道转换为单向通道的做法，只是go语言的一个语法糖。
// 不能利用函数声明把单向通道转换成双向通道。
func main() {
	// 所有类型转换表达式的结果都为false
	var ok bool
	ch := make(chan int, 1)
	_, ok = interface{}(ch).(<-chan int)
	fmt.Println("chan int => <-chan int:", ok)
	_, ok = interface{}(ch).(chan<- int)
	fmt.Println("chan int => chan<- int:", ok)

	sch := make(chan<- int, 1)
	_, ok = interface{}(sch).(chan int)
	fmt.Println("chan<- int => chan int:", ok)

	rch := make(<-chan int, 1)
	_, ok = interface{}(rch).(chan int)
	fmt.Println("<-chan int => chan int:", ok)
}
