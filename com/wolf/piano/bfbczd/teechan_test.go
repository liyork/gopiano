package bfbczd

import (
	"fmt"
	"testing"
)

func TestTeeChan(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}

// 1分2
var tee = func(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {

		defer close(out1)
		defer close(out2)

		for val := range orDone(done, in) {
			var out1, out2 = out1, out2 //1
			for i := 0; i < 2; i++ {    //2
				select {
				case <-done:
				case out1 <- val:
					out1 = nil //3
				case out2 <- val:
					out2 = nil //3
				}
			}
		}
	}()

	return out1, out2
}

//1.我们希望使用使用本地的变量，所以建立了他们的副本。
//2.我们将使用一条select语句，以便写入out1和out2不会彼此阻塞。 为了确保两者都顺利写入，我们将执行select语句的两个迭代。
//3.一旦我们写入了通道，我们将其副本设置为零，这样继续写入将阻塞，而另一个通道可以继续执行。
