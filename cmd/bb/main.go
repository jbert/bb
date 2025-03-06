package main

import (
	"log"

	"github.com/jbert/bb"

	"fmt"
)

func main() {
	ptxt := "this is one text"
	state, err := bb.PlainBlockToState(ptxt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("state\n%s", state)
}
