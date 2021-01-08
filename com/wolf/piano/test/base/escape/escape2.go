package main

func Slice() {
	//s := make([]int, 1000, 1000)
	s := make([]int, 10000, 10000) // escapes

	for index, _ := range s {
		s[index] = index
	}
}

// go run -gcflags=-m escape2.go
// 栈空间不足逃逸,是否逃逸取决于栈空间是否足够大
// 实际上当栈空间不足以存放当前对象时或无法判断当前切片长度时会将对象分配到堆中
func main() {
	Slice()
}
