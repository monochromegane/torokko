package main

import (
	"fmt"
)

func main() {
	err := Run()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
