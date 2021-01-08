package main

import "fmt"

type Flags uint

const (
	FlagUp Flags = 1 << iota // iota从0开始
	FlagBroadcast
	FlagLoopback
	FlagPointToPoint
	FlagMulticast
)

func IsUp(v Flags) bool {
	return v&FlagUp == FlagUp
}

// x&^y 只有x中是1且y的那位不是1才能保留
func TurnDown(v *Flags) {
	*v &^= FlagUp
}

func SetBroadcast(v *Flags) {
	*v |= FlagBroadcast
}

func IsCase(v Flags) bool {
	return v&(FlagBroadcast|FlagMulticast) != 0
}

func main() {
	fmt.Println(FlagUp, FlagBroadcast, FlagLoopback)
	var v = FlagMulticast | FlagUp
	fmt.Printf("%b %t\n", v, IsUp(v))
	TurnDown(&v)
	fmt.Printf("%b %t\n", v, IsUp(v))
	SetBroadcast(&v)
	fmt.Printf("%b %t\n", v, IsUp(v))
	fmt.Printf("%b %t\n", v, IsCase(v))
}
