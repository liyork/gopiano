package main

import "fmt"

// 1,2,3,4得到3个数
func main() {
	for i := 1; i <= 4; i++ {
		for j := 1; j <= 4; j++ {
			if j == i {
				continue
			}
			for q := 1; q <= 4; q++ {
				if q == i || q == j {
					continue
				}
				fmt.Println(i, j, q)
			}
		}
	}
}
