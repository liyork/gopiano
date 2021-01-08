package main

import (
	"fmt"
	"time"
)

// time.NewTimer参数含义是，从定时器被初始化那一刻起，距到期时间需要多少纳秒(ns)。
//Reset用于重置定时器，返回值与档次重置操作是否成功无关，一旦Reset调用完成，该定时器就已被重置。
//Stop用于停止计时器，返回结果含义：为false，说明该定时器早已到期(或已过期)或已被停止，为true说明该定时器由于被调用而被停止。
//接受操作<-timer.C一直阻塞，直到定时器到期。
func main() {
	// 2s后到期
	timer := time.NewTimer(2 * time.Second)
	fmt.Printf("Present time: %v.\n", time.Now())
	expirationTime := <-timer.C
	fmt.Printf("Expiration time: %v.\n", expirationTime)
	fmt.Printf("Stop timer: %v.\n", timer.Stop())
}
