package main

import (
	"fmt"  // Import package for formatted I/O
	"sync" // Import package for goroutine synchronization
)

var wg sync.WaitGroup // Create WaitGroup variable to wait for goroutines

// MyFunc function receives a channel for string data
func MyFunc(data chan string) {
	defer wg.Done() // Decrease WaitGroup counter when function completes
	fmt.Println("Waiting for data...") // Message that function is waiting for data
	text := <-data // Receive data from channel
	fmt.Println(text) // Print received data
}

func main() {
	dataChan := make(chan string) // Create channel for string transmission
	fmt.Println("Run MyFunc goroutine") // Message about starting goroutine
	wg.Add(1) // Increase WaitGroup counter by 1
	go MyFunc(dataChan) // Start MyFunc as a goroutine
	dataChan <- "hehey" // Send string "hehey" to channel
	wg.Wait() // Wait for all goroutines to finish
}

