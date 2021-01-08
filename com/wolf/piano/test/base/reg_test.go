package base

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegBase(t *testing.T) {
	match, _ := regexp.MatchString("\\d{16}", "880218end") //六位连续的数字
	fmt.Println(match)                                     //输出true

	// 任意位置进行匹配1.*.*.*
	match, _ = regexp.MatchString("1\\..*\\..*\\..*", "10.1.1.1")
	fmt.Println("match:", match)

	match, _ = regexp.MatchString("1\\..*\\..*\\..*", "11.10.10.10")
	fmt.Println("match:", match)

	// 以1开头，头尾都要进行匹配。
	match, _ = regexp.MatchString("^1\\..*\\..*\\..*$", "10.1.1.1")
	fmt.Println("match:", match)

	match, _ = regexp.MatchString("^1\\..*\\..*\\..*$", "11.1.1.1")
	fmt.Println("match:", match)
}
