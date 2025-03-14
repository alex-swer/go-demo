package main

import "fmt"

// Node представляет узел связного списка
type Node struct {
	Value int
	Next  *Node
}

// LinkedList представляет связный список
type LinkedList struct {
	Head *Node
}

// Add добавляет новый узел в конец списка
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

// Print выводит все элементы списка
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
	ll.Print() // Вывод: 1 2 3
}
