package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
)

func findIp(input string) string {
	partIp := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	grammer := partIp + "\\." + partIp + "\\." + partIp + "\\." + partIp
	matchMe := regexp.MustCompile(grammer)
	return matchMe.FindString(input)
}

// go run findIPv4.go b.log
// `sort -rn`将`findIPv4.go`的输出作为输入，排序后倒序输出，`uniq -c`计算重复ip的出现次数，最后`sort -rn`按照重复ip的出现次数倒叙输出。
// go run findIPv4.go b.log | sort -rn | uniq -c | sort -rn
func main() {
	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Printf("usage: %s logfile\n", filepath.Base(arguments[0]))
		os.Exit(1)
	}
	for _, filename := range arguments[1:] {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Printf("error openning file %s\n", err)
			os.Exit(-1)
		}
		defer f.Close()

		r := bufio.NewReader(f)
		for {
			line, err := r.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Printf("error openning file %s\n", err)
				break
			}

			ip := findIp(line)
			trail := net.ParseIP(ip)
			if trail.To4() == nil {
				continue
			} else {
				fmt.Println(ip)
			}
		}
	}
}
