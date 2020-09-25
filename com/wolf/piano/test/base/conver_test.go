package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestString2Int(t *testing.T) {
	s := "1"
	// string到int
	i, _ := strconv.Atoi(s)
	fmt.Println("string2int:", i)

	// base：进位制（2进制到36 进制）,如果 base 为 0，则根据字符串的前缀判断进位制（0x:16，0:8，其它:10）
	// bitSize：指定整数类型（0:int、8:int8、16:int16、32:int32、64:int64）
	// string到int64
	int64, _ := strconv.ParseInt(s, 10, 64)
	fmt.Println("string2int64:", int64)
}

func TestInt2String(t *testing.T) {
	// int到string
	i := 1
	var s = strconv.Itoa(i)
	fmt.Println("int2string:", s)

	// int64到string
	var x int64 = 64
	s = strconv.FormatInt(x, 10)
	fmt.Println("int642string10:", s)
	s = strconv.FormatInt(x, 2)
	fmt.Println("int642string2:", s)

	// AppendInt 将 int 型整数 i 转换为字符串形式，并追加到 dst 的尾部
	// i：要转换的字符串
	// base：进位制
	b := make([]byte, 0)
	b = strconv.AppendInt(b, -2048, 16)
	fmt.Printf("%s", b) // -800
}

// int/float/bool/字符-->str，用fmt.sprintf进行转换，%q展示双引号括起来的字符串面值
func TestConvert2String(t *testing.T) {
	var (
		num4 = 1
		char = 'w'
		bl   = true
		str1 string
	)
	//int-->>string
	str1 = fmt.Sprintf("%d", num4)
	fmt.Printf("str=%q\n", str1)

	//str-->>sting
	str1 = fmt.Sprintf("%c", char)
	fmt.Printf("str=%q\n", str1)
	//bool-->>string
	str1 = fmt.Sprintf("%t", bl)
	fmt.Printf("str=%q\n", str1)
}

// fmt：格式标记（b、e、E、f、g、G）
// prec：精度（数字部分的长度，不包括指数部分）
// bitSize：指定浮点类型（32:float32、64:float64）
//
// 转出格式标记：
// 'b' (-ddddp±ddd，二进制指数)
// 'e' (-d.dddde±dd，十进制指数)
// 'E' (-d.ddddE±dd，十进制指数)
// 'f' (-ddd.dddd，没有指数)
// 'g' ('e':大指数，'f':其它情况)
// 'G' ('E':大指数，'f':其它情况)
//
// 如果格式标记为 'e'，'E'和'f'，则 prec 表示小数点后的数字位数
// 如果格式标记为 'g'，'G'，则 prec 表示总的数字位数（整数部分+小数部分）
//func FormatFloat(f float64, fmt byte, prec, bitSize int) string
func TestFloat2String(t *testing.T) {
	f := 100.12345678901234567890123456789
	fmt.Println(strconv.FormatFloat(f, 'b', 5, 32))
	// 13123382p-17
	fmt.Println(strconv.FormatFloat(f, 'e', 5, 32))
	// 1.00123e+02
	fmt.Println(strconv.FormatFloat(f, 'E', 5, 32))

	// AppendFloat 将浮点数 f 转换为字符串值，并将转换结果追加到 dst 的尾部
	b := make([]byte, 0)
	b = strconv.AppendFloat(b, f, 'f', 5, 32)
	b = append(b, " "...)
	b = strconv.AppendFloat(b, f, 'e', 5, 32)
	fmt.Printf("%s", b) // 100.12346 1.00123e+0
}

// str转int/float/bool，用strconv包的函数进行转换
func TestConvertFromString(t *testing.T) {
	var (
		str2 = "true"
		b1   bool
		str3 = "122333"
		n1   int64
	)
	//string转bool，%v是展示值的默认格式
	b1, _ = strconv.ParseBool(str2)
	fmt.Printf("b=%v\n", b1)

	//string转int
	n1, _ = strconv.ParseInt(str3, 10, 64)
	fmt.Printf("n1=%v\n", n1)

}

func TestInt2Byte(t *testing.T) {
	a := 1
	x := int32(a)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, x)
	fmt.Println("byte:", buffer.Bytes())

}

func TestByte2Int(t *testing.T) {
	x := int32(1000)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, x)

	b := buffer.Bytes()
	bin_buf := bytes.NewBuffer(b)
	var y int32
	binary.Read(bin_buf, binary.BigEndian, &y)
	fmt.Println(y)

	//fmt.Println(strings.Repeat("~", 20))
}

//struct String
//{
//        byte*   str;
//        intgo   len;
//};
// 一块连续的内存，首地址+长度
// 如果p里边有\0，他不会做处理这个时候，如果再对这个string做其他处理就可能出问题了，比如strconv.Atoi转成int就有错误
func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

func TestString2Byte(t *testing.T) {
	str2 := "hello"

	data2 := []byte(str2)
	fmt.Println(data2)
	str2 = string(data2)
	fmt.Println(str2)
}

func TestInt2Bytes(t *testing.T) {
	x := int32(256)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	fmt.Println(bytesBuffer.Bytes())
}

func TestBytesToInt(t *testing.T) {
	y := int32(256)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, y)

	bytesBuffer = bytes.NewBuffer(bytesBuffer.Bytes())

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	fmt.Println(int(x))
}

func TestString2Bool(t *testing.T) {
	// 它接受真值：1, t, T, TRUE, true, True
	// 它接受假值：0, f, F, FALSE, false, False.
	// 其它任何值都返回一个错误
	b, e := strconv.ParseBool("1")
	fmt.Println(b, e)

}

// FormatUint 将 uint 型整数 i 转换为字符串形式
// base：进位制（2 进制到 36 进制）
// 大于 10 进制的数，返回值使用小写字母 'a' 到 'z'
func TestUnit2String(t *testing.T) {
	i := uint64(2048)
	fmt.Println(strconv.FormatUint(i, 2))  // 100000000000
	fmt.Println(strconv.FormatUint(i, 10)) // 2048

	// AppendUint 将 uint 型整数 i 转换为字符串形式，并追加到 dst 的尾部
	var b = strconv.AppendUint([]byte("xxx"), 2048, 10)
	fmt.Printf("%s", b) // 800
}

func TestString2Unit(t *testing.T) {
	fmt.Println(strconv.ParseUint("FF", 16, 8))
}

// bitSize：指定浮点类型（32:float32、64:float64）
// 如果 s 是合法的格式，而且接近一个浮点值，
// 则返回浮点数的四舍五入值（依据 IEEE754 的四舍五入标准）
// 如果 s 不是合法的格式，则返回“语法错误”
// 如果转换结果超出 bitSize 范围，则返回“超出范围”
func TestString2Float(t *testing.T) {
	s := "0.12345678901234567890"
	f, err := strconv.ParseFloat(s, 32)
	fmt.Println(f, err) // 0.12345679104328156
}

//layout的时间必须是"2006-01-02 15:04:05"这个时间，不管格式如何，如换一个时间解析出来的时间就不对了，要特别注意这一点。
func TestString2Time(t *testing.T) {
	layout := "2006-01-02 15:04:05.999999"
	ti, _ := time.Parse(layout, "2013-10-05 18:30:50")
	fmt.Println(ti.Year())
	fmt.Println("ti:", ti)
	// 当前时区时间格式化转换
	strTime := "2018-03-24T20:01:00+08:00"
	tim, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", strTime, time.Local)
	fmt.Println("tim:", tim)
}

func TestInt2Time(t *testing.T) {
	tt := time.Unix(0, 1515049539324129700) //将纳秒转换为 time 类型
	fmt.Println(tt.String())

	// 1000毫秒
	tt = time.Unix(0, 1000*1e6) //将纳秒转换为 time 类型
	fmt.Println(tt.String())

	tt = time.Unix(1136214245, 0) //将秒转换为 time 类型
	fmt.Println(tt.String())
}

func TestBool2String(t *testing.T) {
	fmt.Println(strconv.FormatBool(0 < 1)) // true
	fmt.Println(strconv.FormatBool(0 > 1)) // false

	// AppendBool 将布尔值 b 转换为字符串 "true" 或 "false"
	rst := make([]byte, 0)
	rst = strconv.AppendBool(rst, 0 < 1)
	fmt.Printf("%s\n", rst) // true
	rst = strconv.AppendBool(rst, 0 > 1)
	fmt.Printf("%s\n", rst) // truefalse
}

func TestTime2String(t *testing.T) {
	format := time.Now().Format("2006-01-02 15:04")
	fmt.Println("format:", format)
}
