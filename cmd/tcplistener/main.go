package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		buffer := make([]byte, 8)
		var parts string

		for {
			n, err := f.Read(buffer)

			// Handle end of file
			if err == io.EOF {
				// Send any remaining content as the last line
				if parts != "" {
					// Remove trailing \r if present
					line := strings.TrimRight(parts, "\r")
					if line != "" {
						out <- line
					}
				}
				break
			}

			if err != nil {
				log.Printf("Read error: %v", err)
				break
			}

			// Store the bytes read
			chunk := string(buffer[:n])
			parts += chunk

			// Process complete lines
			for {
				// Look for \r\n (HTTP standard) or just \n (fallback)
				crlfIndex := strings.Index(parts, "\r\n")
				lfIndex := strings.Index(parts, "\n")

				var lineEndIndex int
				var lineEndLength int

				if crlfIndex != -1 && (lfIndex == -1 || crlfIndex <= lfIndex) {
					// Found \r\n
					lineEndIndex = crlfIndex
					lineEndLength = 2
				} else if lfIndex != -1 {
					// Found \n only
					lineEndIndex = lfIndex
					lineEndLength = 1
				} else {
					// No line ending found
					break
				}

				currentLine := parts[:lineEndIndex]
				if currentLine != "" { // Only send non-empty lines
					out <- currentLine
				}

				// Remove processed line from buffer
				parts = parts[lineEndIndex+lineEndLength:]
			}
		}
	}()

	return out
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error while creating listener, err :: %s", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		// Handle each connection in a goroutine for better concurrency
		go func(c net.Conn) {
			defer c.Close()

			lines := getLinesChannel(c)
			for line := range lines {
				fmt.Printf("read: %s\n", line)
			}
		}(conn)
	}
}
