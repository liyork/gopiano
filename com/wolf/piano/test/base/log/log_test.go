package log

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestLogBase(t *testing.T) {

	err := os.MkdirAll("/export/server", 0755)
	fmt.Println("err:", err)
	// 按照所需读写权限创建文件
	f, err := os.OpenFile("/export/server/test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	// 完成后延迟关闭
	defer f.Close()
	//设置日志输出到 f
	log.SetOutput(f)
	//测试用例
	log.Println("check to make sure it works!!!!")
}

// Ldate 日期: 2009/01/23
// Ltime 时间: 01:23:23
// Lmicroseconds 毫秒级时间: 01:23:23.123123。该设置会覆盖 Ltime 标志
// Llongfile 完整路径的文件名和行号: /a/b/c/d.go:23
// Lshortfile 最终的文件名元素和行号: d.go:23,覆盖 Llongfile
// LstdFlags = Ldate | Ltime 标准日志记录器的初始值

// log 包这些日志记录器是多 goroutine 安全的。这意味着在多个goroutine 可以同时调用来自同一个日志记录器的这些函数，而不会有彼此间的写冲突

func TestLogFlag(t *testing.T) {

	log.SetPrefix("logprefix:")
	//log.SetFlags(log.LstdFlags|log.Lshortfile)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

	log.Println("message")
	log.Fatalln("fatal message")
	log.Panicln("panic message")
}

func TestFatal(t *testing.T) {
	log.Fatal("Fatal level log: log entry\n") // 输出信息后，程序终止执行
	log.Println("Nomal level log: log entry")
}
