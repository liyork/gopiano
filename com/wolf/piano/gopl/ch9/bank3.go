package main

import "sync"

var (
	mu sync.Mutex
)

// 上锁
func Deposit3(amount int) {
	mu.Lock()
	balance = balance + amount
	mu.Unlock()
}

func Balance3() int {
	mu.Lock()
	b := balance
	mu.Unlock()
	return b
}
