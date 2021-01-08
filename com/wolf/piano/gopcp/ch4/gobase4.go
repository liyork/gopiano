package main

import (
	"fmt"
	"time"
)

// 随着迭代的进行，每次获取到的迭代值(每个名字)都会被赋给这个迭代变量name，
//这里并发执行的5个go函数(由5个G分别封装的同一个函数)中，name值都是Mark，因为他们都在for语句执行完毕之后才执行，name那时已经是mark了
//不要对go函数的执行时机作任何假设，除非确实作出让这种假设成为绝对事实的保证

func main() {
	names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		go func() {
			fmt.Printf("hello, %s!\n", name)
		}()
	}
	time.Sleep(time.Millisecond)
}
