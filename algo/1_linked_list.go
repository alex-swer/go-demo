package main

import "fmt"

// Node represents a node of a linked list
type Node struct {
	Value int
	Next  *Node
}

// LinkedList represents a linked list
type LinkedList struct {
	Head *Node
}

// Add adds a new node to the end of the list
func (ll *LinkedList) Add(value int) {
	newNode := &Node{Value: value}
	if ll.Head == nil {
		ll.Head = newNode
		return
	}
	current := ll.Head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
}

// Print outputs all elements of the list
func (ll *LinkedList) Print() {
	current := ll.Head
	for current != nil {
		fmt.Print(current.Value, " ")
		current = current.Next
	}
	fmt.Println()
}

func main() {
	ll := &LinkedList{}
	ll.Add(1)
	ll.Add(2)
	ll.Add(3)
	ll.Print() // Output: 1 2 3
}
