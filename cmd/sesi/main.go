package main

import (
	"fmt"
	"log"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	fmt.Println("hi")
	return nil
}
