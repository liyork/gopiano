package main

import (
	"fmt"
	"sort"
)

type aStructure struct {
	person string
	height int
	weight int
}

// go run sortSlice.go
func main() {
	mySlice := []aStructure{{"Mihalis", 180, 90}, {"Bill", 134, 45}, {"Merietta", 155, 45}, {"Epifanios", 144, 50}, {"Athina", 134, 40}}

	// 升序
	sort.Slice(mySlice, func(i, j int) bool {
		return mySlice[i].height < mySlice[j].height
	})
	fmt.Println("<:", mySlice)

	sort.Slice(mySlice, func(i, j int) bool {
		return mySlice[i].height > mySlice[j].height
	})
	fmt.Println(">:", mySlice)
}
