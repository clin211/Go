package main

import (
	"log"
	"tour/cmd"
)

var name string

func main() {
	err := cmd.Execute()

	if err != nil {
		log.Fatalf("cmd.Execute err: %v", err)
	}
}
