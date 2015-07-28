package main

import (
	"fmt"

	"github.com/monochromegane/cargo"
)

func main() {
	err := cargo.Run()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
