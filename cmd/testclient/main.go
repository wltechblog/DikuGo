package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create reader and writer
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Start a goroutine to read from the server
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from server:", err)
				os.Exit(1)
			}
			fmt.Print(line)
		}
	}()

	// Read input from the user and send to the server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if strings.ToLower(input) == "quit" {
			break
		}
		writer.WriteString(input + "\r\n")
		writer.Flush()
	}
}
