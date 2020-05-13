package init

import "fmt"

// 全局变量
var Age int
var Name string

func init() {
	fmt.Println("init2_test.go init")
	Age = 100
	Name = "Tao"
}
