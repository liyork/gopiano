package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	logger         *log.Logger
	signalChan     = make(chan os.Signal, 1)
	connFiles      sync.Map
	serverListener net.Listener
	// 老进程接到信号后，设定true，不再接收建联和已存在连接上的读写，交给新进程处理
	isUpdate = false
)

func init() {
	//file, err := os.OpenFile("./server.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	//if err != nil {
	//	panic(err)
	//}
	logger = log.New(os.Stdout, "", 11)
	go beforeStart()
	go signalHandler()
}

// go run server.go common.go
// kill -USR2 pid
func main() {
	fmt.Println(os.Getpid(), " ", "main start running")
	var err error
	serverListener, err = net.Listen("tcp", ":7000")
	if err != nil {
		panic(err)
	}
	for {
		if isUpdate == true {
			continue
		}
		conn, err := serverListener.Accept()
		if err != nil {
			logger.Println(os.Getpid(), " ", "conn error, err:", err)
			continue
		}
		c := conn.(*net.TCPConn)
		// 每个客户端建联新开协程
		go connectionHandler(c)
	}
}

func connectionHandler(conn *net.TCPConn) {
	file, _ := conn.File()
	connFiles.Store(file, true)
	logger.Printf("conn fd %d\n", file.Fd())
	defer func() {
		connFiles.Delete(file)
		_ = conn.Close()
	}()
	for {
		if isUpdate == true {
			continue
		}
		// 读超时
		err := conn.SetReadDeadline(time.Now().Add(readTimeout))
		if err != nil {
			logger.Println(os.Getpid(), " ", err.Error())
			return
		}
		// 读取4字节,ping
		rBuf := make([]byte, 4)
		_, err = conn.Read(rBuf)
		if err != nil {
			logger.Println(os.Getpid(), " ", err.Error())
			return
		}
		if string(rBuf) != "ping" {
			logger.Println(os.Getpid(), " ", "failed to parse the message "+string(rBuf))
			return
		}
		// 写超时
		err = conn.SetWriteDeadline(time.Now().Add(writeTimeout))
		if err != nil {
			logger.Println(os.Getpid(), " ", err.Error())
			return
		}
		// 写回pong
		_, err = conn.Write([]byte(`pong`))
		if err != nil {
			logger.Println(os.Getpid(), " ", err.Error())
			return
		}
	}
}

func beforeStart() {
	connInterface, err := net.Dial("unix", "./conn.sock")
	if err != nil { // 没有则表明当前不存在老进程
		logger.Println(os.Getpid(), " ", err.Error())
		return
	}

	// 新老传递fd
	defer func() {
		_ = connInterface.Close()
	}()

	unixConn := connInterface.(*net.UnixConn)

	b := make([]byte, 1)
	oob := make([]byte, 32)
	for {
		err = unixConn.SetWriteDeadline(time.Now().Add(time.Minute * 3))
		if err != nil {
			fmt.Println(os.Getpid(), " ", err.Error())
			return
		}
		n, oobn, _, _, err := unixConn.ReadMsgUnix(b, oob)
		if err != nil {
			logger.Println(os.Getpid(), " ", err.Error())
			return
		}
		if n != 1 {
			logger.Printf("recv fd type error: %d\n", n)
			return
		}

		if b[0] != 0 {
			logger.Println(os.Getpid(), " ", "init finish, oobn:", oobn)
			return
		}
		scms, err := unix.ParseSocketControlMessage(oob[0:oobn])
		if err != nil {
			logger.Println(os.Getpid(), " ", err.Error())
			return
		}
		if len(scms) != 1 {
			logger.Printf("recv fd num != 1 : %d\n", len(scms))
			return
		}
		fds, err := unix.ParseUnixRights(&scms[0])
		if err != nil {
			logger.Println(os.Getpid(), " ", err.Error())
			return
		}

		if len(fds) != 1 {
			logger.Printf("recv fd num != 1 : %d\n", len(fds))
			return
		}
		logger.Printf("recv fd %d\n", fds[0])
		// 从fd还原成file
		file := os.NewFile(uintptr(fds[0]), "fd-from-old")
		// 从file还原成conn
		conn, err := net.FileConn(file)
		if err != nil {
			logger.Println(os.Getpid(), " ", err.Error())
			return
		}
		// 继续针对已有conn进行读写
		go connectionHandler(conn.(*net.TCPConn))
	}
}

func signalHandler() {
	// 初始化信号
	signal.Notify(
		signalChan,
		syscall.SIGUSR2,
	)
	for {
		sc := <-signalChan
		switch sc {
		case syscall.SIGUSR2:
			gracefulExit()
		default:
			continue
		}
	}
}

func gracefulExit() {
	var connWait sync.WaitGroup
	// 重新监听conn.sock
	_ = syscall.Unlink("./conn.sock")
	listenerInterface, err := net.Listen("unix", "./conn.sock")
	if err != nil {
		logger.Println(os.Getpid(), " ", err.Error())
		return
	}
	defer func() {
		_ = listenerInterface.Close()
	}()
	unixListener := listenerInterface.(*net.UnixListener)
	connWait.Add(1)
	// 异步接收子进程请求
	go func() {
		defer connWait.Done()
		unixConn, err := unixListener.AcceptUnix()
		if err != nil {
			logger.Println(os.Getpid(), " ", err.Error())
			return
		}
		defer func() {
			_ = unixConn.Close()
		}()
		// 将已存在的conn传递给新进程
		connFiles.Range(func(key, value interface{}) bool {
			if key == nil || value == nil {
				return false
			}
			file := key.(*os.File)
			defer func() {
				_ = file.Close()
			}()
			// Send filedescriptors with a 1 byte message.
			// 写一字节的消息
			buf := make([]byte, 1)
			buf[0] = 0
			// encodes xx into a socket control message
			rights := syscall.UnixRights(int(file.Fd()))
			_, _, err := unixConn.WriteMsgUnix(buf, rights, nil)
			if err != nil {
				logger.Println(os.Getpid(), " ", err.Error())
			}
			logger.Printf("send fd %d, rights:%v\n", file.Fd(), rights)
			return true
		})
		// 写回完成信号
		finish := make([]byte, 1)
		finish[0] = 1
		_, _, err = unixConn.WriteMsgUnix(finish, nil, nil)
		if err != nil {
			logger.Println(os.Getpid(), " ", err.Error())
		}
	}()

	isUpdate = true
	execSpec := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: append([]uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}),
	}

	// fork出子进程，用同样命令和参数，得到新进程pid
	pid, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		logger.Println(os.Getpid(), " ", err.Error())
		return
	}
	logger.Printf("old process %d new process %d\n", os.Getpid(), pid)
	// 成功启动新进程后，这里不再进行接收新建联，交给新进程
	_ = serverListener.Close()

	connWait.Wait()
	os.Exit(0)
}
