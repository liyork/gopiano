package base

import (
	"fmt"
	"testing"
	"time"
)

// defer的原意是推迟、延期。它的思想类似与C++的析构函数，不过go语言中的析构的不是对象，而是函数，defer就是用来添加函数结束时执行的语句。
//注意这里强调的是添加，而不是指定，即可以有多个。因为不同于C++中的析构函数是静态的，Go中的defer是动态的。
func TestDeferBase(t *testing.T) {
	fmt.Println("result1:", Defer1())
	fmt.Println("result2:", Defer2())
}

// 返回之前，调用defer，进行+1
// 应该是返回0，但是赋值给result了，然后defer中++就更新了然后返回
func Defer1() (result int) {
	defer func() {
		result++
	}()
	return 0
}

// 应该是先执行return 0，然后defer中再对result进行操作不影响结果了
func Defer2() int {
	var result = 0
	defer func() {
		result++
	}()
	return result
}

// 一直等到包含defer语句的函数执行完毕时，延迟函数（defer后的函数）才会被执行，
// 而不管函数是通过return的正常结束，还是由于panic导致的异常结束。
func TestDefer(t *testing.T) {

	i := return1()
	fmt.Println("i:", i)
}

func return1() int {
	defer println(1111)
	fmt.Println("222")
	return 1
}

// defer声明和执行顺序相反，defer栈,后进先出
func TestDeferSeq(t *testing.T) {
	defer println(1111)
	defer println(222)
	fmt.Println("3333")
}

type passValue struct {
	a int
}

func TestPassValue(t *testing.T) {
	value := passValue{1}
	// 对闭包的参数传递是值传递
	defer func(pv passValue) {
		fmt.Println("passValue1", pv)
	}(value)

	// 闭包内使用外部变量是引用传递
	defer func() {
		fmt.Println("passValue2", value)
	}()

	value.a = 2
	fmt.Println("value:", value)
}

// 得先定义defer才能最后被调用
func TestDeferUnder(t *testing.T) {

	fmt.Println("1111")
	testP()
	defer println(2222)
}

func testP() {
	panic(33333)
}

// panic: send on closed channel
// defer在所属方法退出时进行调用，若有goroutine执行另一个方法产生panic，那么本调用方法defer不会执行，
// 所以defer要放在目标方法中
func TestNested(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			println("xx Start error, err:%v", err)
		}
	}()

	ints := make(chan int)
	go func() {
		time.Sleep(5 * time.Second)
		close(ints)
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				println("xx Start error, err:%v", err)
			}
		}()
		for {
			ints <- 1
			time.Sleep(2 * time.Second)
		}
	}()

	time.Sleep(111 * time.Second)
}
