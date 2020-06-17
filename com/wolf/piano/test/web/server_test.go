package web

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "welcome page")
}

// http://localhost:8080/
func TestServerBase(t *testing.T) {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handle(conn net.Conn) {
	fmt.Fprint(conn, "welcome user", time.Now().String())
	conn.Close()
}

// telnet 127.0.0.1 6666
func TestServerTcpBase(t *testing.T) {
	l, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle(conn)
	}
}
