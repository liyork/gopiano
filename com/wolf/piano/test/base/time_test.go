package base

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

// 10位数的时间戳是以 秒 为单位；
// 13位数的时间戳是以 毫秒 为单位；
// 19位数的时间戳是以 纳秒 为单位；
func TestTimeBase(t *testing.T) {
	fmt.Printf("时间戳（秒）：%v;\n", time.Now().Unix())
	fmt.Printf("时间戳（纳秒）：%v;\n", time.Now().UnixNano())
	fmt.Printf("时间戳（毫秒）：%v;\n", time.Now().UnixNano()/1e6)
	fmt.Printf("时间戳（纳秒转换为秒）：%v;\n", time.Now().UnixNano()/1e9)
}

// 年月日时分秒转time
func TestCreateTime(t *testing.T) {
	theTime := time.Date(2014, 1, 7, 5, 50, 4, 0, time.Local)
	unixTime := theTime.Unix()
	fmt.Println(unixTime)
}

func TestGetInsKey(t *testing.T) {
	now := time.Now().UnixNano() / 1e6
	formatInt := strconv.FormatInt(now, 10)
	i := len(formatInt) - 5
	fmt.Println(i)
	s := string([]byte(formatInt)[i:])
	fmt.Println(fmt.Sprintf("%s_%d_%s", "11111", os.Getpid(), s))
}
