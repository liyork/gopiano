package base

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadFileByLine(t *testing.T) {

	file, err := os.Open("/export/App/pilot/ifacefile.txt")
	if err != nil {
		fmt.Println("open file err:", err.Error())
		return
	}
	defer file.Close()

	r := bufio.NewReader(file)

	for {
		data, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err", err.Error())
			break
		}

		fmt.Println(string(data))
	}
}

// 一次性读取适合读取内容比较小的文件,大文件读取建议使用第一种方式
func TestReadFile(t *testing.T) {
	// 一次性读取,不需要关闭文件
	data, err := ioutil.ReadFile("/export/App/pilot/ifacefile.txt")
	if err != nil {
		fmt.Println("read file err:", err.Error())
		return
	}

	// dkljflsjj\nsdfsd\nxcvxc\nvxc\nvwe31\n432\n
	fmt.Println(string(data))
}

func TestReadFileAll(t *testing.T) {
	fi, err := os.Open("/export/App/pilot/ifacefile.txt")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	fmt.Println(string(fd))
}
