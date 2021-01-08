package bfbczd

import (
	"fmt"
	"net/http"
	"testing"
)

func TestWhereErr(t *testing.T) {
	// 不断请求urls中的url，直到done或者urls没了
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err) //1
					continue
				}
				select {
				case <-done:
					return
				case responses <- resp:
				}
			}
		}()

		return responses
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.baidu.com", "https://badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}
	//这个我们看到goroutine尽其最大努力展示错误信号。但也仅仅是展示出来，它还能做什么？ 它无法传回，让外界知道！ 如果错误种类太多怎么办？ 再请求一遍吗？
}

type Result struct { //1
	Error    error
	Response *http.Response
}

// 向外传递错误
func TestPassErr(t *testing.T) {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result { //2
		results := make(chan Result)
		go func() {
			defer close(results)

			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				// 正、错都返回err
				result = Result{Error: err, Response: resp} //3
				select {
				case <-done:
					return
				case results <- result: //4
				}
			}
		}()

		return results
	}
	done := make(chan interface{})
	defer close(done)

	process1(checkStatus, done)
	//process2(checkStatus, done)
}

func process1(checkStatus func(done <-chan interface{}, urls ...string) <-chan Result, done chan interface{}) {
	urls := []string{"https://www.baidu.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		// 有权利判断是打印还是继续向上
		if result.Error != nil { //5
			fmt.Printf("error: %v", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

func process2(checkStatus func(done <-chan interface{}, urls ...string) <-chan Result, done chan interface{}) {
	errCount := 0
	urls := []string{"a", "https://www.baidu.com", "b", "c", "d"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
