package benchmark

import (
	"fmt"
	"testing"
	"time"
)

//基准测试可以测试一段程序的运行性能及耗费 CPU 的程度。
//测试代码需要保证函数可重入性及无状态
//测试代码不使用全局变量等带有记忆性质的数据结构。避免多次运行同一段代码时的环境不一致，不能假设 N 值范围。

//-bench=.表示运行 benchmark_test.go 文件里的所有基准测试
//1000000000 表示测试的次数，也就是 testing.B 结构中提供给程序使用的 N。“0.33 ns/op”表示每一个操作耗费多少时间（纳秒）
//go test -v -bench=. benchmark_test.go
//go test -v -count=2 -bench=Add$ benchmark_test.go  仅测试Add结尾的方法
// b.N 由基准测试框架提供。
func Benchmark_Add(b *testing.B) {
	var n int
	for i := 0; i < b.N; i++ {
		n++
	}
	fmt.Println("n-->:", n)
}

//对一段代码可能存在的内存分配进行统计
//-benchmem参数以显示内存分配情况
//-bench后添加了 Alloc，指定只测试 Benchmark_Alloc() 函数
//“16 B/op”表示每一次调用需要分配 16 个字节，“2 allocs/op”表示每一次调用有两次分配。
//go test -v -bench=Alloc -benchmem benchmark_test.go
func Benchmark_Alloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%d", i)
	}
}

//有些测试需要一定的启动和初始化时间，如果从 Benchmark() 函数开始计时会很大程度上影响测试结果的精准性。
// testing.B 提供了一系列的方法可以方便地控制计时器，从而让计时器只在需要的区间进行测试
//从 Benchmark_xx() 函数开始，Timer 就开始计数。StopTimer() 可以停止这个计数过程，做一些耗时的操作，通过 StartTimer() 重新开始计时。ResetTimer() 可以重置计数器的数据
//计数器内部不仅包含耗时数据，还包括内存分配的数据。
//go test -v -bench=Add_TimerControl -benchmem benchmark_test.go
func Benchmark_Add_TimerControl(b *testing.B) {
	time.Sleep(2 * time.Second)
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	//b.StopTimer()
	// 开始计时器
	//b.StartTimer()
	var n int
	for i := 0; i < b.N; i++ {
		n++
	}
}
