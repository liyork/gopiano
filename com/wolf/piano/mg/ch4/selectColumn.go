package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// go run selectColumn.go 3 swtag.log
func main() {
	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Printf("usage: use selectColumn column [file1] [file2] [...]\n")
		os.Exit(1)
	}

	temp, err := strconv.Atoi(arguments[1])
	if err != nil {
		fmt.Printf("column value is not an integer", temp)
		return
	}

	column := temp
	if column < 0 {
		fmt.Println("Invalid column value")
		os.Exit(1)
	}

	for _, fileName := range arguments[2:] {
		fmt.Println("\t\t", fileName)
		f, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("error opening file %s\n", err)
			continue
		}
		defer f.Close()

		r := bufio.NewReader(f)
		for {
			// 读取到参数中的字符就停止,swtag.log中即使一行也要有一个回车，表示\n不然就是EOF了，即读到末尾也没发现\n
			line, err := r.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Printf("error reading file %s\n", err)
			}

			data := strings.Fields(line)
			if len(data) > column {
				fmt.Println(data[column-1])
			}
		}
	}
}
