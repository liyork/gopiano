package base

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"testing"
)

func TestSystemBase(t *testing.T) {
	println(`系统类型：`, runtime.GOOS)
	println(`系统架构：`, runtime.GOARCH)
	println(`CPU 核数：`, runtime.GOMAXPROCS(0))
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	println(`电脑名称：`, name)
}

func TestGetLocalIp(t *testing.T) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						fmt.Println(ipnet.IP.String())
						break
					}
				}
			}
		}
	}
}

// /var/log/messages、/var/log/syslog
func TestGetPid(t *testing.T) {
	pid := os.Getpid()
	fmt.Printf("进程 PID: %d n", pid)

	prc := exec.Command("ps", "-p", strconv.Itoa(pid), "-v")
	out, err := prc.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}

// 要在goland中设定Environment，单纯机器的env不能用
// 或者export xxxx=3 &  go test -v system_test.go -test.run TestGetEvn
func TestGetEvn(t *testing.T) {
	//for index, key := range os.Environ() {
	//	fmt.Println("index = ", index, " key = ", key)
	//}

	// Getenv 方法有个缺点，即使在未设置环境变量的情况下，它也返回一个空的字符串。
	xxxx := os.Getenv("xxxx")
	fmt.Println("xxxx:", xxxx)

	xxxx2, exist := os.LookupEnv("xxxx2")
	fmt.Printf("xxxx2:%s, err:%v\n", xxxx2, exist)

	os.Setenv("xxxx2", "22222")
	val := os.Getenv("xxxx2")
	fmt.Println("xxxx2值是 :" + val)

	os.Unsetenv("xxxx2")
	val = os.Getenv("xxxx2")
	fmt.Println("xxxx2值是 :" + val)
}
