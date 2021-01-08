package main

import "fmt"

// go run loops.go
func main() {
	for i := 0; i < 100; i++ {
		// for循环的例子以及关键字continue和break的使用
		if i%20 == 0 {
			continue
		}
		if i == 95 {
			break
		}
		fmt.Print(i, " ")

		// 使用了关键字break来退出for循环
		fmt.Println()
		i := 10
		for {
			if i < 0 {
				break
			}
			fmt.Print(i, " ")
			i--
		}
		fmt.Println()

		// for循环来进行do...while循环的工作的操作
		i = 0
		anExpression := true
		for ok := true; ok; ok = anExpression {
			if i > 10 {
				anExpression = false
			}
			fmt.Print(i, " ")
			i++
		}
		fmt.Println()

		// range
		anArray := [5]int{0, 1, -1, 2, -2}
		for i, value := range anArray {
			fmt.Println("index:", i, "value: ", value)
		}
	}
}
