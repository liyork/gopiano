package main

import (
	"fmt"
	"log"
	"net/http"
)

type database map[string]dollars

func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

// go run http1.go &
// curl http://localhost:8000
func main() {
	db := database{"shoes": 50, "socks:": 5}
	log.Fatal(http.ListenAndServe("localhost:8000", db))
}
