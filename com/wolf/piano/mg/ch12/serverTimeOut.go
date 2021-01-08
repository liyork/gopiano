package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func myHandler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Fprintf(w, "Served: %s\n", r.Host)
}

func timeHandler2(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format(time.RFC1123)
	Body := "The current time is:"
	fmt.Fprintf(w, "<h1 align=\"center\">%s</h1>", Body)
	fmt.Fprintf(w, "<h2 align=\"center\">%s</h2>n", t)
	fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Printf("Served time for :%s\n", r.Host)
}

// go run serverTimeOut.go
// time nc localhost 8001
func main() {
	PORT := ":8001"
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Printf("Listening on http://0.0.0.0%s\n", PORT)
	} else {
		PORT = ":" + arguments[1]
		fmt.Printf("Listening on http://0.0.0.0%s\n", PORT)
	}

	m := http.NewServeMux()
	srv := &http.Server{
		Addr:         PORT,
		Handler:      m,
		ReadTimeout:  3 * time.Second, // 读取整个请求，包括消息体的最大持续时间
		WriteTimeout: 3 * time.Second, // 从请求头读取结束到响应写入结束的时间。
	}

	m.HandleFunc("/time", timeHandler2)
	m.HandleFunc("/", myHandler2)

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
