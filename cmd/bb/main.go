package main

import (
	"fmt"

	"github.com/jbert/bb"
)

func main() {
	for i := range 5 {
		rc := bb.Rcon(byte(i))
		fmt.Printf("rcon(%d) %02X\n", i, rc>>24)
	}
}
