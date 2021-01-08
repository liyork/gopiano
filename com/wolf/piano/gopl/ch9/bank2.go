package main

var (
	sema = make(chan struct{}, 1)
)

// 获取token
func Deposit2(amount int) {
	sema <- struct{}{}
	balance = balance + amount
	<-sema
}

func Balance2() int {
	sema <- struct{}{}
	b := balance
	<-sema
	return b
}
