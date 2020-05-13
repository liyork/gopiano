package main

import (
	"fmt"
	"github.com/gohouse/gorose"
	"github.com/liyork/gopiano/com/wolf/piano/testmod/utils" //依赖自己工程中内容(moduleName+path)
)

func main() {
	fmt.Println("hello")
	gorose.NewLogger()
	utils.PrintText("abc")
}
