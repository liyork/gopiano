package main

import "fmt"

func f11() int {
	fmt.Println("Entering f11()")
	//return -10
	fmt.Println("Exiting f11()")
	return -1
}

func f22() int {
	if true {
		return 10
	}
	fmt.Println("Exiting f22()")
	return 0
}

// go vet cannotReach.go
func main() {
	fmt.Println(f11())
	fmt.Println("Exiting program...")
}
