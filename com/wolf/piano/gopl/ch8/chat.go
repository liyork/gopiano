package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// go run chat.go &
// cd ..
// 开俩窗口
// go run ./ch8/netcat3.go ./ch8/common.go
// go run ./ch8/netcat3.go ./ch8/common.go
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn4(conn)
	}
}

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages: //消息发送给每个client
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering: // 来存
			clients[cli] = true
		case cli := <-leaving: // 走删
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn4(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "Yor are " + who
	//广播
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
