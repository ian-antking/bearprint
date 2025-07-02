package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Received: %s\n", line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading stdin: %v", err)
	}
}
