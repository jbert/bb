package main

import (
	"fmt"

	"github.com/jbert/bb/ff"
)

func main() {
	for i := range 16 {
		for j := range 16 {
			p := ff.Poly(i*16 + j)
			fmt.Printf("%02x, ", byte(p.Mul(2)))
		}
		fmt.Printf("\n")
	}
}
