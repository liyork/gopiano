package main

import (
	"fmt"
	"net"
	"os"
)

// go run MXrecords.go golang.com
// go run MXrecords.go www.mtsoukalos.eu
// host -t mx golang.com
// host -t mx www.mtsoukalos.eu
func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Need a domain name!")
		return
	}

	domain := arguments[1]
	MXs, err := net.LookupMX(domain)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, MX := range MXs {
		fmt.Println(MX.Host)
	}
}
