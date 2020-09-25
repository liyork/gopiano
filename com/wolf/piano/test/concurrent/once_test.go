package concurrent

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestOnceBase(t *testing.T) {
	var once sync.Once
	go func() {
		fmt.Println("before 11111111")
		once.Do(func() {
			fmt.Println("11111111")
			time.Sleep(5 * time.Second)
		})
		fmt.Println("after 11111111")
	}()

	go func() {
		fmt.Println("before 222222")
		// 会阻塞，直到其他协成的Do执行完
		once.Do(func() {
			fmt.Println("222222")
			time.Sleep(5 * time.Second)
		})
		fmt.Println("after 222222")
	}()
	time.Sleep(55 * time.Second)
}

type wrapOnce struct {
	once sync.Once
}

func TestOnceReplicate(t *testing.T) {
	wrapOnce1 := &wrapOnce{}
	wrapOnce2 := &wrapOnce{}
	go func() {
		fmt.Println("before 11111111")
		wrapOnce1.once.Do(func() {
			fmt.Println("11111111")
			time.Sleep(5 * time.Second)
		})
		fmt.Println("after 11111111")
	}()

	go func() {
		fmt.Println("before 222222")
		// 会阻塞，直到其他协成的Do执行完
		wrapOnce1.once.Do(func() {
			fmt.Println("222222")
			time.Sleep(5 * time.Second)
		})
		fmt.Println("after 222222")
	}()

	// 不阻塞
	go func() {
		fmt.Println("before 333333")
		wrapOnce2.once.Do(func() {
			fmt.Println("333333")
			time.Sleep(5 * time.Second)
		})
		fmt.Println("after 333333")
	}()

	time.Sleep(10 * time.Second)
	// 重新执行，不再执行once
	go func() {
		fmt.Println("before 222222______")
		wrapOnce1.once.Do(func() {
			fmt.Println("222222______")
			time.Sleep(5 * time.Second)
		})
		fmt.Println("after 222222______")
	}()

	time.Sleep(55 * time.Second)
}
