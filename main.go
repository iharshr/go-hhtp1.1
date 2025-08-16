package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("error while reading messages.txt , err :: %s", err)
	}
	defer file.Close()
	buffer := make([]byte, 8)
	var parts string

	for {
		f, err := file.Read(buffer)
		// handeling end of file here
		if err == io.EOF {
			if parts != "" {
				fmt.Printf("read: %s\n", string(buffer[:f]))
			}
			break
		}
		// store the 8 bites here
		chunk := string(buffer[:f])
		// add to line buffer
		// this line buffer is of length till \n
		parts += chunk

		for {
			new_line_index := strings.Index(parts, "\n")
			if new_line_index == -1 {
				// this means no new line
				// we are still reading
				// so break from this loop
				break
			}

			current_line := parts[:new_line_index]
			fmt.Printf("read: %s\n", string(current_line))

			// clear buffer for new line
			// so we can read next line
			parts = parts[new_line_index+1:]
		}

	}
}
