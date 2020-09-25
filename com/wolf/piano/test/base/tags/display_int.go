// +build int

package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type DisplayName uint64

func Print(name DisplayName) {
	fmt.Printf("%d\n", name)
}

func MakeDisplayName(name string) DisplayName {
	h, err := hex.DecodeString(name)
	if err != nil {
		panic(fmt.Sprintf("decode hex string failed. cause: %v\n", err))
	}
	fmt.Printf("data: %v\n", h)

	value := binary.BigEndian.Uint16(h)
	return DisplayName(value)
}
