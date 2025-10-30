package main

import (
	"fmt"
	"go-demo/internal/linkedlist"
)

func linkedlistExample() {
	fmt.Println("=== Linked List Examples ===")
	fmt.Println()
	
	basicOperations()
	fmt.Println()
	
	advancedOperations()
	fmt.Println()
	
	errorHandling()
}

func basicOperations() {
	fmt.Println("1. Basic Operations:")
	
	ll := linkedlist.New()
	
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)
	fmt.Printf("   After appending 1, 2, 3: %v\n", ll.ToSlice())
	
	ll.Prepend(0)
	fmt.Printf("   After prepending 0: %v\n", ll.ToSlice())
	
	fmt.Printf("   Size: %d\n", ll.Size())
}

func advancedOperations() {
	fmt.Println("2. Advanced Operations:")
	
	ll := linkedlist.New()
	for i := 1; i <= 5; i++ {
		ll.Append(i)
	}
	fmt.Printf("   Original list: %v\n", ll.ToSlice())
	
	value, err := ll.GetAt(2)
	if err == nil {
		fmt.Printf("   Value at index 2: %d\n", value)
	}
	
	node, found := ll.Find(3)
	if found {
		fmt.Printf("   Found node with value: %d\n", node.Value)
	}
	
	ll.Reverse()
	fmt.Printf("   After reverse: %v\n", ll.ToSlice())
}

func errorHandling() {
	fmt.Println("3. Error Handling:")
	
	ll := linkedlist.New()
	
	err := ll.Delete(1)
	if err != nil {
		fmt.Printf("   Expected error (empty list): %v\n", err)
	}
	
	ll.Append(1)
	ll.Append(2)
	
	err = ll.InsertAt(10, 99)
	if err != nil {
		fmt.Printf("   Expected error (out of range): %v\n", err)
	}
	
	_, err = ll.GetAt(-1)
	if err != nil {
		fmt.Printf("   Expected error (negative index): %v\n", err)
	}
}

