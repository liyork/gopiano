package main

import (
	"fmt"
	"os"
	"strconv"
)

// man ascii
// ASCII 字符表中可打印的字符一共有 94 个[0~94)
// go run generatePassword.go common.go
// go run generatePassword.go common.go
// go run generatePassword.go common.go 20
func main() {
	MIN := 0
	MAX := 94
	var LENGTH int64 = 8

	arguments := os.Args
	switch len(arguments) {
	case 2:
		LENGTH, _ = strconv.ParseInt(os.Args[1], 10, 64)
	default:
		fmt.Println("Using default values!")
	}

	startChar := "!"
	var i int64 = 1
	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		fmt.Print(newChar)
		if i == LENGTH {
			break
		}
		i++
	}
	fmt.Println()
}
