package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

// 1.pprof可以做什么
//CPU Profiling：CPU 分析，按照一定的频率采集所监听的应用程序 CPU（含寄存器）的使用情况，可确定应用程序在主动消耗 CPU 周期时花费时间的位置
//Memory Profiling：内存分析，在应用程序进行堆分配时记录堆栈跟踪，用于监视当前和历史内存使用情况，以及检查内存泄漏
//Block Profiling：阻塞分析，记录 goroutine 阻塞等待同步（包括定时器通道）的位置
//Mutex Profiling：互斥锁分析，报告互斥锁的竞争情况

// 2.http://127.0.0.1:6060/debug/pprof/
// 可以看到
//cpu（CPU Profiling）: $HOST/debug/pprof/profile，默认进行 30s 的 CPU Profiling，得到一个分析用的 profile 文件
//block（Block Profiling）：$HOST/debug/pprof/block，查看导致阻塞同步的堆栈跟踪
//goroutine：$HOST/debug/pprof/goroutine，查看当前所有运行的 goroutines 堆栈跟踪
//heap（Memory Profiling）: $HOST/debug/pprof/heap，查看活动对象的内存分配情况
//mutex（Mutex Profiling）：$HOST/debug/pprof/mutex，查看导致互斥锁的竞争持有者的堆栈跟踪
//threadcreate：$HOST/debug/pprof/threadcreate，查看创建新OS线程的堆栈跟踪

// 3.通过命令方式下载cpu占用// pprof help
// go tool pprof http://localhost:6060/debug/pprof/profile\?seconds\=60
// (pprof) top
//Showing nodes accounting for 7.15s, 96.10% of 7.44s total
//Dropped 57 nodes (cum <= 0.04s)
//Showing top 10 nodes out of 43
//      flat  flat%   sum%        cum   cum%
//     4.93s 66.26% 66.26%      5.30s 71.24%  syscall.syscall
//     0.51s  6.85% 73.12%      0.51s  6.85%  runtime.nanotime
//     0.48s  6.45% 79.57%      0.48s  6.45%  runtime.usleep
//     0.37s  4.97% 84.54%      0.37s  4.97%  runtime.exitsyscallfast
//     0.32s  4.30% 88.84%      0.60s  8.06%  runtime.notetsleep
//     0.29s  3.90% 92.74%      0.29s  3.90%  runtime.pthread_cond_timedwait_relative_np
//     0.09s  1.21% 93.95%      0.09s  1.21%  runtime.pthread_cond_signal
//     0.08s  1.08% 95.03%      0.08s  1.08%  runtime.memmove
//     0.04s  0.54% 95.56%      0.04s  0.54%  runtime.findObject
//     0.04s  0.54% 96.10%      0.04s  0.54%  runtime.madvise
//解释
//flat：给定函数上运行耗时
//flat%：同上的 CPU 运行耗时总比例
//sum%：给定函数累积使用 CPU 总比例
//cum：当前函数加上它之上的调用运行总耗时
//cum%：同上的 CPU 运行耗时总比例

//查看内存
//go tool pprof http://localhost:6060/debug/pprof/heap
//(pprof) top
//Showing nodes accounting for 148.75MB, 100% of 148.75MB total
//      flat  flat%   sum%        cum   cum%
//  148.75MB   100%   100%   148.75MB   100%  main.Add
//         0     0%   100%   148.75MB   100%  main.main.func1
//
//-inuse_space：分析应用程序的常驻内存占用情况
//-alloc_objects：分析应用程序的内存临时分配情况
//
//go tool pprof http://localhost:6060/debug/pprof/block
//go tool pprof http://localhost:6060/debug/pprof/mutex

// 4.通过浏览器查看cpu等信息
// go tool pprof -http=:8080 profile
// 如果出现 Could not execute dot; may need to install graphviz.，就是提示你要安装 graphviz 了
// 安装graphviz：
// brew install graphviz

// http://localhost:8080/ui/
//http://localhost:8080/ui/top
//框越大，线越粗代表它占用的时间越大
//http://localhost:8080/ui/peek
//http://localhost:8080/ui/source
// http://localhost:8080/ui/flamegraph
// 调用顺序由上到下（A -> B -> C -> D），每一块代表一个函数，越大代表占用 CPU 的时间更长。同时它也支持点击块深入进行分析！
func main() {
	go func() {
		for {
			log.Println(Add("https://github.com/EDDYCJY"))
		}
	}()

	http.ListenAndServe("0.0.0.0:6060", nil)
}
