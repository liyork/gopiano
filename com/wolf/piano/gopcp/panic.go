package main

import "fmt"

func main() {
	outerFunc()
}

func outerFunc() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("Recovered panic: %s\n", p)
		}
	}()
	innerFunc()
}

func innerFunc() {
	//panic(errors.New("An intened fatal error!"))
	myIndex := 4
	ia := [3]int{1, 2, 3}
	_ = ia[myIndex]

	fmt.Println("1111111")
}
