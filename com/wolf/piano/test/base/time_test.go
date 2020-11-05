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

type timestr struct {
	lasttime time.Time
}

func TestTimeDiff(t *testing.T) {
	t1 := time.Now()
	fmt.Println(time.Now().Sub(t1) < 2*time.Second)
	time.Sleep(2 * time.Second)
	fmt.Println(time.Now().Sub(t1) < 2*time.Second)

	i := timestr{} // 默认有值0001-01-01 00:00:00 +0000
	fmt.Println(time.Now().Sub(i.lasttime) < 2*time.Second)

	t2 := time.Now()
	//50s前
	t3 := time.Now().Add(time.Second * 50)
	fmt.Println("t2与t1相差：", t3.Sub(t2))
}
