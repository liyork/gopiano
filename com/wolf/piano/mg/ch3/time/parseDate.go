package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// go run parseDate.go "12 January 2019"
func main() {
	var myDate string

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s date\n",
			filepath.Base(os.Args[0]))
		os.Exit(0)
	}

	myDate = os.Args[1]

	// Go解析日期的常量是：`Jan`用来解析月份（英文月份简写）/January得到月份的英文全称，`2006`用来解析年，`02`用来解析天，`Mon`用来解析周几（如果是`Monday`，那就是周几的英文全称），
	// 02 January-2006
	d, err := time.Parse("02 January 2006", myDate)
	if err == nil {
		fmt.Println("Full", d)
		fmt.Println("Time", d.Day(), d.Month(), d.Year())
	} else {
		fmt.Println(err)
	}
}
