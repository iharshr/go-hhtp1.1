package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("error while reading messages.txt , err :: %s", err)
	}
	defer file.Close()
	buffer := make([]byte, 8)

	for {
		chunk, err := file.Read(buffer)
		if err != nil {
			log.Fatal("error reading file")
			break
		}
		fmt.Printf("read: %s\n", string(buffer[:chunk]))

	}
}
