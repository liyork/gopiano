package main

import (
	"fmt"
	"log"
	"net/http"
)

// go run server1.go &
// ./fetch http://localhost:8000
// ./fetch http://localhost:8000/help
func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}
