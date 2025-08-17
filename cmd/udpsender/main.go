package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// Resolve UDP address for localhost:42069
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	// Dial UDP connection
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatalf("Failed to dial UDP connection: %v", err)
	}
	defer conn.Close()

	// Create bufio.Reader for reading from stdin
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("UDP sender started. Type messages and press Enter to send.")

	// Infinite loop for user input
	for {
		// Print prompt
		fmt.Print("> ")

		// Read line from stdin
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from stdin: %v", err)
			continue
		}

		// Write line to UDP connection
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("Error writing to UDP connection: %v", err)
			continue
		}
	}
}
