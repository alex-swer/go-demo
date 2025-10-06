package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Hello, World!")
	
	// Wait for Enter key
	fmt.Printf("\nPress Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

