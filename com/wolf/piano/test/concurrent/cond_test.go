package concurrent

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCondBase(t *testing.T) {
	mutex := sync.Mutex{}
	cond := sync.NewCond(&mutex)

	go func() {
		mutex.Lock()
		time.Sleep(3 * time.Second)
		cond.Signal()
		mutex.Unlock()
	}()

	mutex.Lock()
	cond.Wait()
	mutex.Unlock()
	fmt.Println("finish")
}
