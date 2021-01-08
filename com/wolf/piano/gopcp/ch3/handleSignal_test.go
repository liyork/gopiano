package ch3

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestReceiveSignal(t *testing.T) {
	sigRecv1 := make(chan os.Signal, 1)
	sigs1 := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
	fmt.Printf("Set notification for %s... [sigRecv1]\n", sigs1)
	signal.Notify(sigRecv1, sigs1...)

	sigRecv2 := make(chan os.Signal, 1)
	sigs2 := []os.Signal{syscall.SIGQUIT}
	fmt.Printf("Set notification for %s... [sigRecv2]\n", sigs2)
	signal.Notify(sigRecv2, sigs2...)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for sig := range sigRecv1 {
			fmt.Printf("Received a signal from sigRecv1: %s\n", sig)
		}
		fmt.Printf("End. [sigRecv1]\n")
		wg.Done()
	}()
	go func() {
		for sig := range sigRecv2 {
			fmt.Printf("Received a signal from sigRecv2: %s\n", sig)
		}
		fmt.Printf("End. [sigRecv2]\n")
		wg.Done()
	}()

	fmt.Printf("Wait for 2 seconds...")
	time.Sleep(2 * time.Second)
	fmt.Printf("Stop notification...")
	signal.Stop(sigRecv1)
	close(sigRecv1)
	fmt.Printf("done. [sigRecv1]\n")

	wg.Wait()
}

func TestSendSignal(t *testing.T) {
	cmds := []*exec.Cmd{
		exec.Command("ps", "aux"),
		exec.Command("grep", "Signal"),
		exec.Command("grep", "-v", "grep"),
		exec.Command("grep", "-v", "go run"),
		exec.Command("awk", "{printf $2}"),
	}

	output, err := runCmds(cmds)
	if err != nil {
		fmt.Printf("Command Execution Error: %s\n", err)
		return
	}

	fmt.Println("output,", output)

	for _, pidStr := range output {
		pid, _ := strconv.ParseInt(pidStr, 10, 64)
		proc, _ := os.FindProcess(int(pid))
		err = proc.Signal(syscall.SIGINT)
	}
}

func runCmds(cmds []*exec.Cmd) ([]string, error) {

	var tmpBuf bytes.Buffer
	for i, cmd := range cmds {
		if i != 0 {
			cmd.Stdin = &tmpBuf
		}

		var outputBuf2 bytes.Buffer
		cmd.Stdout = &outputBuf2
		if err := cmd.Start(); err != nil {
			fmt.Printf("Error: The second command can not be startup %s\n", err)
			return nil, err
		}
		if err := cmd.Wait(); err != nil {
			fmt.Printf("Error: Countn't wait for the first command: %s\n", err)
			return nil, err
		}
		tmpBuf.Reset()
		tmpBuf.Write(outputBuf2.Bytes())
	}

	return convert2Array(tmpBuf)
}

func convert2Array(buffer bytes.Buffer) ([]string, error) {
	return nil, nil
}
