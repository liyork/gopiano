package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"
)

// 实例2： 两个flag操作一个变量
// ./flagdemo --gopher_type=23
// ./flagdemo --g=23
var gopherType string

func init() {
	const (
		defaultGopher = "pocket"
		usage         = "the variety of gopher"
	)
	flag.StringVar(&gopherType, "gopher_type", defaultGopher, usage)
	flag.StringVar(&gopherType, "g", defaultGopher, usage+" (shorthand)")
}

// 实例3：用户定义flag类型(一个时间段的切片)
type interval []time.Duration

// String是一个用来格式化flag值(flag.Value接口的一部分)的方法
// String方法的输出将被用于调试
func (i *interval) String() string {
	return fmt.Sprint(*i)
}

// Set是一个用来设置flag值(flag.Value接口的一部分)的方法
// Set的参数是String类型，用于设置为flag
// 这是一个以逗号为分隔符的数组，我们需要分离它
func (i *interval) Set(value string) error {
	// 如果flag能被设置为多时间，加速度值，如果有如此声明，我们将会删除这些
	// 这些将会允许很多组合，例如"-deltaT 10s -deltaT 15s"
	if len(*i) > 0 {
		return errors.New("interval flag already set")
	}
	for _, dt := range strings.Split(value, ",") {
		duration, err := time.ParseDuration(dt)
		if err != nil {
			return err
		}
		*i = append(*i, duration)
	}
	return nil
}

// 将一个flag定义为堆积期间。因为它还有个特殊类型，我们需要使用Var函数，从而在初始化中创建flag
var intervalFlag interval
var inputFlagvar int
var inputName = flag.String("name", "CHENJIAN", "Input Your Name.")

// 返回的是指针
var inputAge = flag.Int("age", 27, "Input Your Age")
var inputGender = flag.String("gender", "female", "Input Your Gender")

func init() {
	// 将命令行flag与intervalFlag绑定，并设置使用信息
	flag.Var(&intervalFlag, "deltaT", "comma-separated list of intervals to use between events")
	flag.IntVar(&inputFlagvar, "flagname", 1234, "Help")
}

// go build -v flagdemo.go
// ./flagdemo -h
// ./flagdemo youname 33
// ./flagdemo -name balbalba -age 1111 -flagname=12333 dfdf xccccc eette

func main() {
	flag.Parse()
	// Args returns the non-flag command-line arguments.
	// NArg is the number of arguments remaining after flags have been processed.
	fmt.Printf("args=%s, num=%d\n", flag.Args(), flag.NArg())
	for i := 0; i != flag.NArg(); i++ {
		fmt.Printf("arg[%d]=%s\n", i, flag.Arg(i))
	}
	fmt.Println("name=", *inputName)
	fmt.Println("age=", *inputAge)
	fmt.Println("gender=", *inputGender)
	fmt.Println("flagname=", inputFlagvar)
	fmt.Println("gopherType=", gopherType)
}

func example() {
	//flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")

	// 自定义类型需要实现value接口
	//flag.Var(&flagVal, "name", "help message for flagname")

	// 解析函数将会在碰到第一个非flag命令行参数时停止
	//flag.Parse()

	//命令行参数的格式
	//-flag xxx （使用空格，一个 - 符号）
	//--flag xxx （使用空格，两个 - 符号）
	//-flag=xxx （使用等号，一个 - 符号）
	//–-flag=xxx （使用等号，两个 - 符号）
}
