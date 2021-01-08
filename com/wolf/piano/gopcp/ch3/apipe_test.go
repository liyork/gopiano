package ch3

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestBasePipe(t *testing.T) {
	cmd0 := exec.Command("echo", "-n", "My first command comes from golang.")
	stdout0, err := cmd0.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:Couldn't obtain the stdout pipe for command No.0:%s\n", err)
		return
	}

	if err := cmd0.Start(); err != nil {
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
		return
	}

	var outputBuf0 bytes.Buffer
	for {
		tempOutput := make([]byte, 5)
		n, err := stdout0.Read(tempOutput)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Error: Coundn't read data from the pipe: %s\n", err)
				return
			}
		}
		if n > 0 {
			outputBuf0.Write(tempOutput[:n])
		}
	}
	fmt.Printf("%s\n", outputBuf0.String())
}

func TestUsingBuf(t *testing.T) {
	cmd0 := exec.Command("echo", "-n", "My first command comes from golang.")
	stdout0, err := cmd0.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:Couldn't obtain the stdout pipe for command No.0:%s\n", err)
		return
	}

	if err := cmd0.Start(); err != nil {
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
		return
	}

	// 返回缓冲读取器，默认4096
	outputBuf0 := bufio.NewReader(stdout0)

	output0, isReadFinish, err := outputBuf0.ReadLine()
	if err != nil {
		fmt.Printf("Error: Countn't read data from the pipe: %s\n", err)
		return
	}
	fmt.Printf("%s, %v\n", string(output0), isReadFinish)
}

// ps aux | grep go
func TestPipeTwoCommand(t *testing.T) {
	cmd1 := exec.Command("ps", "aux")

	// 媒介
	var outputBuf1 bytes.Buffer
	cmd1.Stdout = &outputBuf1
	if err := cmd1.Start(); err != nil {
		fmt.Printf("Error: The first command can not be startup %s\n", err)
		return
	}
	if err := cmd1.Wait(); err != nil {
		fmt.Printf("Error: Countn't wait for the first command: %s\n", err)
		return
	}

	//fmt.Printf("%s\n", outputBuf1.Bytes())

	// grep一个不存在则cmd2.Wait()提示错误ExitError,!state.Success()，应该是grep进程没查到数据导致status非0
	//cmd2 := exec.Command("grep", "go")
	cmd2 := exec.Command("grep", "xxxxx")

	cmd2.Stdin = &outputBuf1
	var outputBuf2 bytes.Buffer
	cmd2.Stdout = &outputBuf2
	if err := cmd2.Start(); err != nil {
		fmt.Printf("Error: The second command can not be startup %s\n", err)
		return
	}
	if err := cmd2.Wait(); err != nil {
		fmt.Printf("Error: Countn't wait for the second command: %s\n", err)
		return
	}
	fmt.Printf("%s\n", outputBuf2.Bytes())
}

func TestNamedPipe(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Printf("create pipe error, err:%v\n", err)
		return
	}

	go func() {
		n, err := writer.Write([]byte("inputxxxxx"))
		if err != nil {
			fmt.Printf("Error: Couldn't write data to the named pipe: %s\n", err)
			return
		}
		fmt.Printf("Written %d byte(s). [file-based pipe]\n", n)
	}()

	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Couldn't read data from the named pipe: %s\n", err)
			return
		}
		fmt.Printf("Read %d byte(s). [file-based pipe], content:%v\n", n, string(output[:n]))
	}()

	time.Sleep(2 * time.Second)
}

// go在标准库中提供了一个基于内存的有原子性操作保证的管道
//reader, writer := io.Pipe()
//内部使用sync代码包中api，以此保证操作的原子性，可以放心地并发写入和读取数据。
//这种管道不是基于文件系统的，并没有作为中介的缓冲区，所以通过它传递的数据只会被复制一次。提高了数据传递的效率
func TestMemoryPipe(t *testing.T) {
	reader, writer := io.Pipe()

	go func() {
		n, err := writer.Write([]byte("inputxxxxx"))
		if err != nil {
			fmt.Printf("Error: Couldn't write data to the named pipe: %s\n", err)
			return
		}
		fmt.Printf("Written %d byte(s). [file-based pipe]\n", n)
	}()

	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Couldn't read data from the named pipe: %s\n", err)
			return
		}
		fmt.Printf("Read %d byte(s). [file-based pipe], content:%v\n", n, string(output[:n]))
	}()

	time.Sleep(2 * time.Second)
}
