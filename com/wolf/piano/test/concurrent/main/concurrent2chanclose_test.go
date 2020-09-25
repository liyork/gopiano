package main

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

// 这里closed the chan字后还有一堆数字打印之后才是all finish是有缓存，还是打印的慢?
// 再测试，还是有可能打印不出closed the chan，也就是没有关闭，但是程序读写都不运行了。。6060端口不知道怎么停了，开始时还是有的,
// o，是不是test用例的问题?...因为执行完了所以6060结束，但是加上sleep在最后也有可能不能停止，应该是chan的问题。
// /Users/lichao30/goland_workspace/gopiano/com/wolf/piano/test/concurrent/closechan_test.go:27  --http.ListenAndServe
// /Users/lichao30/goland_workspace/gopiano/com/wolf/piano/test/concurrent/closechan_test.go:100  --time.Sleep(2222 * time.Second)
// 看样子应该是完成，但是为什么没有打印close也没有打印finish，若是把所有打印去掉，那么没有任何问题....???
// 同样的代码，在concurrent/main/testconcurrent2chan.go就能执行完成，应该是log或者testunit的问题！
func TestOneSenderMultiReceiverErrForTest(t *testing.T) {
	// http://localhost:6060/debug/pprof
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recover has err,", err)
			}
		}()
		defer func() {
			fmt.Println("listen return")
		}()
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	defer func() {
		fmt.Println("main return")
	}()
	//go func() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover has err,", err)
		}
	}()
	rand.Seed(time.Now().UnixNano())
	// ...
	const MaxRandomNumber = 100000
	const NumReceivers = 100
	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)
	// ...
	dataCh := make(chan int, 100)
	// the sender
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recover has err,", err)
			}
		}()
		defer func() {
			fmt.Println("sender return")
		}()
		for {
			intn := rand.Intn(MaxRandomNumber)
			fmt.Println("intn:", intn)
			if value := intn; value == 0 {
				fmt.Println("closed the chan")
				// the only sender can close the channel safely.
				close(dataCh)
				return
			} else {
				dataCh <- value
			}
		}
	}()
	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("recover has err,", err)
				}
			}()
			defer func() {
				fmt.Println("receiver return")
			}()
			//defer wgReceivers.Done()
			// receive values until dataCh is closed and
			// the value buffer queue of dataCh is empty.
			for value := range dataCh {
				fmt.Print(value)
				//if value != 0 {
				//}
			}
			wgReceivers.Done()
			//fmt.Println("sub finish", )
		}()
	}
	wgReceivers.Wait()
	fmt.Println("all finish ")
	//}()
	time.Sleep(2222 * time.Second)
}
