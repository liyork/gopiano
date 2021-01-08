package main

import (
	"container/ring"
	"fmt"
)

var size = 10

// go run conRing.go
func main() {
	myRing := ring.New(size + 1)
	fmt.Println("Empty ring:", *myRing)

	// 0~9
	for i := 0; i < myRing.Len()-1; i++ {
		myRing.Value = i
		myRing = myRing.Next()
	}

	// 10
	myRing.Value = 2

	sum := 0
	myRing.Do(func(x interface{}) {
		t := x.(int)
		sum = sum + t
	})
	fmt.Println("Sum:", sum)

	// 多出来2就循环了
	for i := 0; i < myRing.Len()+2; i++ {
		myRing = myRing.Next()
		fmt.Print(myRing.Value, " ")
	}
	fmt.Println()
}
