package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// go run parseTime.go 12:10
func main() {
	var myTime string

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s string\n", filepath.Base(os.Args[0]))
		os.Exit(0)
	}

	myTime = os.Args[1]

	// `15`代表解析小时，`04`代表解析分钟，`05`解析秒，同时你可以使用`PM`将字符串中的字母转为大写，`pm`转为小写。
	d, err := time.Parse("15:04", myTime)
	if err == nil {
		fmt.Println("Full", d)
		fmt.Println("Time", d.Hour(), d.Minute())
	} else {
		fmt.Println(err)
	}
}
