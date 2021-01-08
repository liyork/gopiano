package main

import (
	"fmt"
	"time"
)

func bigSlowOperation() {
	defer trace("bigSlowOperation")()
	time.Sleep(10 * time.Second)
}

func trace(msg string) func() {
	start := time.Now()
	fmt.Printf("enter %s\n", msg)
	return func() { fmt.Printf("exit %s (%s)\n", msg, time.Since(start)) }
}

func main() {
	bigSlowOperation()
}
