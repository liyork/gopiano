package test

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"testing"
)

type SomeError struct {
	Reason string
}

// 编写一个错误对象很简单，写一个结构体，然后挂载 Error() 方法就可以了。
func (s SomeError) Error() string {
	return s.Reason
}

func TestErrBase(t *testing.T) {
	var err error = SomeError{"something happened"}
	fmt.Println(err)
}

// 对于上面代码中错误对象的形式很常用，所以Go语言内置了一个通用错误类型，在 errors 包里。
func TestNewErr(t *testing.T) {
	var err = errors.New("something happened")
	fmt.Println("err:", err)

	var thing = "something1"
	err = fmt.Errorf("%s happened", thing)
	fmt.Println("err:", err)
}

// 在 Go 语言里面它不会抛异常，而是以返回值的形式来通知上层逻辑来处理错误
// 通过读文件来尝试一下 Go 语言的错误处理。
func TestErrByFile(t *testing.T) {
	// 打开文件
	var f, err = os.Open("main.go") // 返回了两个值,不过 Go 语言普遍没有使用多返回值的习惯，仅仅是在需要返回错误的时候才会需要两个返回值。
	//除了错误之外，还有一个地方需要两个返回值，那就是字典，通过第二个返回值来告知读取的结果是零值还是根本就不存在。
	//如：var score, ok := scores["apple"]
	if err != nil {
		// 文件不存在、权限等原因
		fmt.Println("open file failed reason:" + err.Error())
		return
	}
	// 将文件的关闭调用推迟到当前函数的尾部执行，即使后面的代码抛出了异常，文件关闭也会确保被执行，相当于 Java 语言的 finally 语句块。
	defer f.Close()
	// 存储文件内容
	var content = []byte{}
	// 临时的缓冲，按块读取，一次最多读取 100 字节
	var buf = make([]byte, 100)
	for {
		// 读文件，将读到的内容填充到缓冲,将文件的内容往切片里填充，填充的量不会超过切片的长度(注意不是容量)
		n, err := f.Read(buf) // 如果遇到文件尾了，切片就不会填满。所以需要通过返回值 n 来明确到底读了多少字节
		if n > 0 {
			// 将读到的内容聚合起来,...用于一次性追加多个元素,作用是将切片参数的所有元素展开后传递给 append 函数,展开只是形式上的展开，在实现上其实并没有展开，传递过去的参数本质上还是切片。
			content = append(content, buf[:n]...)
		}
		if err != nil {
			// 遇到流结束或者其它错误
			break
		}
	}
	// 输出文件内容
	fmt.Println(string(content))
}

// 体验 Redis 的错误处理
// 实现一个小功能，获取 Redis 中两个整数值，然后相乘，再存入 Redis 中
// 因为 Go 语言中不轻易使用异常语句，所以对于任何可能出错的地方都需要判断返回值的错误信息
func TestErrByRedis(t *testing.T) {
	// 定义客户端对象，内部包含一个连接池
	var client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 定义三个重要的整数变量值，默认都是零
	var val1, val2, val3 int

	// 获取第一个值
	valstr1, err := client.Get("value1").Result()
	if err == nil {
		val1, err = strconv.Atoi(valstr1)
		if err != nil {
			fmt.Println("value1 not a valid integer")
			return
		}
		// 字符串的零值是空串而不是 nil，不好从字符串内容本身判断出 Redis 是否存在这个 key 还是对应 key 的 value 为空串，
		// 需要通过返回值的错误信息来判断。代码中的 redis.Nil 就是客户端专门为 key 不存在这种情况而定义的错误对象。
	} else if err != redis.Nil {
		fmt.Println("redis access error reason:" + err.Error())
		return
	}

	// 获取第二个值
	valstr2, err := client.Get("value2").Result()
	if err == nil {
		val2, err = strconv.Atoi(valstr2)
		if err != nil {
			fmt.Println("value1 not a valid integer")
			return
		}
	} else if err != redis.Nil {
		fmt.Println("redis access error reason:" + err.Error())
		return
	}

	// 保存第三个值
	val3 = val1 * val2
	ok, err := client.Set("value3", val3, 0).Result()
	if err != nil {
		fmt.Println("set value error reason:" + err.Error())
		return
	}
	fmt.Println(ok)
}
