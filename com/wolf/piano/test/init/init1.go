package main

import "fmt"

var Age int
var Name string

// Age Nane 是全局变量

func init() {
	fmt.Println("init 包的初始化 init（）。。。")
	Age = 100
	Name = "Tao"

}
