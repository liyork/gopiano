package base

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestTimeBase(t *testing.T) {
	// 返回本地时间Time
	fmt.Println(time.Now())
	// 返回UTC时间
	fmt.Println(time.Now().UTC())

	// 得到一个UTC时间
	lastLoginTime, err := time.Parse("2006-01-02 15:04:05", "2020-11-30 09:45:54")
	fmt.Println(lastLoginTime, err)

	// 得到本地时区
	loc, _ := time.LoadLocation("Local")
	// 将时间字符串转化为本地时间
	lastLoginTime, err = time.ParseInLocation("2006-01-02 15:04:05", "2020-11-30 09:45:54", loc)
	fmt.Println(lastLoginTime, err)

}

// 10位数的时间戳是以 秒 为单位；
// 13位数的时间戳是以 毫秒 为单位；
// 19位数的时间戳是以 纳秒 为单位；
func TestTimeBase2(t *testing.T) {
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

func TestCompose(t *testing.T) {
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

func TestTimeFormat(t *testing.T) {
	format := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(format)
}

func TestTimeStamp2Time(t *testing.T) {
	i := time.Now().Unix()
	unix := time.Unix(i, 0)
	fmt.Println(unix)
}

//UTC标准、北京、美国洛杉矶在同一时刻的转换：
func TestLocation(t *testing.T) {
	formate := "2006-01-02 15:04:05 Mon"
	now := time.Now()
	local1, err1 := time.LoadLocation("UTC") //输入参数"UTC"，等同于""
	if err1 != nil {
		fmt.Println(err1)
	}
	local2, err2 := time.LoadLocation("Local")
	if err2 != nil {
		fmt.Println(err2)
	}
	local3, err3 := time.LoadLocation("America/Los_Angeles")
	if err3 != nil {
		fmt.Println(err3)
	}

	fmt.Println(now.In(local1).Format(formate))
	fmt.Println(now.In(local2).Format(formate))
	fmt.Println(now.In(local3).Format(formate))
}

func TestTimeInit(t *testing.T) {
	local, _ := time.LoadLocation("America/Los_Angeles")
	timeFormat := "2006-01-02 15:04:05"
	//通过unix标准时间的秒，纳秒设置时间
	time1 := time.Unix(1480390585, 0)
	// 字符串设定时间
	time2, _ := time.ParseInLocation(timeFormat, "2016-11-28 19:36:25", local)
	fmt.Println(time1.In(local).Format(timeFormat))
	fmt.Println(time2.In(local).Format(timeFormat))
	chinaLocal, _ := time.LoadLocation("Local") //运行时，该服务器必须设置为中国时区，否则最好是采用"Asia/Chongqing"之类具体的参数。
	fmt.Println(time2.In(chinaLocal).Format(timeFormat))
}
