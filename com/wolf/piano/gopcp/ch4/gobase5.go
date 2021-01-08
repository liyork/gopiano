package main

import (
	"fmt"
	"time"
)

//若在go函数中使用的name值不会受到外部变量变化的影响，就既可以保证go函数的独立执行，又不用担心他们的正确性受到破坏。
//若这个假设被实现，该go函数就可以称作“可重入函数”
//潜在原因是，name类型string是一个非引用类型，在把一个值作为参数传递给函数时，该值会被复制。对于外部修改不会对函数内部产生影响。
//对于引用类型(如切片或字典类型)的值，由于它类似于指向真正数据的指针，所以即使被复制了，之后在外部对该值的修改也会反映到该函数的内部。
func main() {
	names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		// 值被传入go函数，传入过程中，该值会被复制并在go函数中由参数who指代。
		// go函数内不再使用外部变量name，使用参数who
		go func(who string) {
			fmt.Printf("hello, %s!\n", who)
		}(name)
	}
	time.Sleep(time.Millisecond)
}
