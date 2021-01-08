package main

import "fmt"

// go run failMap.go
func main() {
	aMap := map[string]int{}
	aMap["test"] = 1

	aMap2 := map[string]int{}
	aMap2 = nil
	fmt.Println(aMap2)
	// panic: assignment to entry in nil map
	aMap2["test"] = 1
}
