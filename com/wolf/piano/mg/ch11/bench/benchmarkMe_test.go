package main

import "testing"

// -bench=.
// 每个基准测试默认执行至少 1 秒。如果基准函数在少于 1 秒的时间内返回，则 `b.N` 的值增加并且该函数会再次运行。`b.N` 的值第一次是 1，然后变为 2，5，10，20，50 等等。这是因为函数运行的越快，您就需要越多次的运行它来获得准确结果
// 第一列结尾处的 `-8` 表示该函数被执行期间的 goroutines 数，本质上它是 `GOMAXPROCS` 环境变量的值。
// 第二列显示了相关函数的执行次数。较快的函数比较慢的函数被执行了多次。例如，`Benchmark30fibo3()` 函数执行了 500,000 次，而 `Benchmark50fibo2()` 函数仅执行了一次！输出的第三列显示了每个运行的平均值。

//-benchmem多包含了两列。第四列显示了平均分配给每个执行的基准函数的内存数。
// 第五列显示了用于分配第四列的内存值的分配数。

var result int

// go test -bench=. benchmarkMe_test.go benchmarkMe.go
// go test -bench=. -benchmem benchmarkMe_test.go benchmarkMe.go

// //存放 `fibo1(n)` 的结果在一个名为 `r` 的变量中，并在之后使用另一个名为 `result` 的全局变量的原因是很微妙。
// 此技巧用于阻止编译器执行任何优化，这些优化将排除您要测量的函数，因为它的结果从未被使用过！
func benchmarkfibo1(b *testing.B, n int) {
	var r int
	for i := 0; i < b.N; i++ {
		r = fibo1(n)
	}
	result = r
}

func benchmarkfibo2(b *testing.B, n int) {
	var r int
	for i := 0; i < b.N; i++ {
		r = fibo2(n)
	}
	result = r
}

func benchmarkfibo3(b *testing.B, n int) {
	var r int
	for i := 0; i < b.N; i++ {
		r = fibo3(n)
	}
	result = r
}

func Benchmark30fibo1(b *testing.B) {
	benchmarkfibo1(b, 30)
}

func Benchmark30fibo2(b *testing.B) {
	benchmarkfibo2(b, 30)
}

func Benchmark30fibo3(b *testing.B) {
	benchmarkfibo3(b, 30)
}

func Benchmark50fibo1(b *testing.B) {
	benchmarkfibo1(b, 50)
}

func Benchmark50fibo2(b *testing.B) {
	benchmarkfibo2(b, 50)
}

func Benchmark50fibo3(b *testing.B) {
	benchmarkfibo3(b, 50)
}

// 错误
//func BenchmarkFiboI(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		_ = fibo1(i) // 这里用的i
//	}
//}

// 执行 `go test` 命令后不能从它获得任何输出。
// 原因是根据上一节的描述，随着 `b.N` 值的增加，由于 `for`循环这个基准测试函数的运行次数也会增加。这一事实阻止了 `BenchmarkFiboI()` 收敛于一个稳定的数字，从而妨碍了函数完成并返回。

//由于算法原因，下面的基准测试函数也实现错了：
//func BenchmarkfiboII(b *testing.B) {
//    for i := 0 ; i < b.N ; i++ {
//        _ = fibo2(b.N)
//    }
//}

//下面俩个基准测试函数的实现没有错误：
//func BenchmarkFiboIV(b *testing.B) {
//    for i := 0; i < b.N; i++ {
//        _ = fibo3(10)
//    }
//}
//
//func BenchmarkFiboIII(b *testing.B) {
//    _ = fibo3(b.N)
//}
