package concurrent

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// context可以用来实现一对多的goroutine协作。
// 上下文。当一个请求来时，会产生一个goroutine，但是这个goroutine往往要衍生出许多额外的goroutine去处理操作，例如链接database、请求rpc请求。。等等，这些衍生的goroutine和主goroutine有很多公用数据的，例如同一个请求生命周期、用户认证信息、token等，当这个请求超时或者被取消的时候，这里所有的goroutine都应该结束。context就可以帮助我们达到这个效果。
//主goroutine和衍生的所有子goroutine之间形成了一颗树结构。我们的context可以从根节点遍布整棵树，当然，是线程安全的
//线程之间的基本是这样的：
//func DoSomething(ctx context.Context, arg Arg) error {
//    // ... use ctx ...
//}
//有两个根context:background和todo；这两个根都是contenxt空的，没有值的。两者也没啥太本质的区别，Background是最常用的，作为Context这个树结构的最顶层的Context，它不能被取消。当不知道用啥context的时候就可以用TODO。

// // 要想结束所有线程，就调用ctx, cancel := context.WithCancel(context.Background())函数返回的cancel函数即可，当撤销函数被调用之后，对应的Context值会先关闭它内部的接收通道，也就是它的Done方法返回的通道。
// 可撤销的Context
func TestCancelContext(t *testing.T) {
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done(): //只有撤销函数被调用后，才会触发
					return
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //调用返回的cancel方法来让 context声明周期结束

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}

//两者唯一区别是WithTimeout表示从现在开始xxx超时，而WithDeadline的时间可以是指定的任意的时间(包括以前的时间)
// WithDeadline
func TestWithDeadline(t *testing.T) {
	d := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	// Even though ctx will be expired, it is good practice to call its
	// cancelation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	defer cancel() //时间超时会自动调用，或者主动调用

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}

func TestWithTimeout(t *testing.T) {
	// Pass a context with a timeout to tell a blocking function that it
	// should abandon its work after the timeout elapses.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel() //时间超时会自动调用

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // prints "context deadline exceeded"
	}
}

// 不可撤销的context，传递值
func TestNoCancelContext(t *testing.T) {
	type favContextKey string
	k := favContextKey("language")
	k1 := favContextKey("Chinese")

	ctx := context.WithValue(context.Background(), k, "Go")
	ctx1 := context.WithValue(ctx, k1, "Go1")

	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	f(ctx1, k1)
	f(ctx1, k)
}
