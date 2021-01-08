package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func lineByLine(file string) error {
	var err error

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		// 指示`bufio.ReadString()`一直读取，直到碰到行分隔符(\n)为止
		line, err := r.ReadString('\n')
		if err == io.EOF { // 文件末尾
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			break
		}
		// line中已经包含换行符
		fmt.Printf(line)
	}
	return nil
}

// go run byLine.go swtag.log | wc
// wc swtag.log
func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Printf("usage: byLine <file1> [<file2> ...]\n")
		return
	}

	for _, file := range flag.Args() {
		err := lineByLine(file)
		if err != nil {
			fmt.Println(err)
		}
	}
}
