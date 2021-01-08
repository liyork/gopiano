package main

import (
	"fmt"
	"time"
)

// go run usingTime.go
func main() {
	// 返回UNIX时间（UNIX时间是计算了从00:00:00 UTC，1970年1月1日以来的秒数）。
	fmt.Println("Epoch Time:", time.Now().Unix())
	t := time.Now()
	// 将`time`类型的变量转换成其他格式，例如`RFC3339`格式。
	fmt.Println(t, " | ", t.Format(time.RFC3339))
	fmt.Println(t.Weekday(), t.Day(), t.Month(), t.Year())

	time.Sleep(time.Second)
	t1 := time.Now()
	fmt.Println("Time difference:", t1.Sub(t))

	formatT := t.Format("01 January 2006")
	fmt.Println(formatT)
	loc, _ := time.LoadLocation("Europe/Paris")
	LondonTime := t.In(loc)
	fmt.Println("Paris:", LondonTime)
}
