package linkedlist

import (
	"errors"
	"fmt"
)

var (
	// ErrEmptyList is returned when an operation cannot be performed on an empty list.
	ErrEmptyList = errors.New("list is empty")
	// ErrIndexOutOfRange is returned when the specified index is invalid.
	ErrIndexOutOfRange = errors.New("index out of range")
)

// Node represents a single node in the linked list.
type Node struct {
	Value int
	Next  *Node
}

// LinkedList represents a singly linked list data structure.
type LinkedList struct {
	Head *Node
	Tail *Node
	size int
}

// New creates and returns a new empty LinkedList.
func New() *LinkedList {
	return &LinkedList{
		Head: nil,
		Tail: nil,
		size: 0,
	}
}

// Append adds a new node with the given value to the end of the list.
// Time complexity: O(1)
func (ll *LinkedList) Append(value int) {
	newNode := &Node{Value: value}
	
	if ll.Head == nil {
		ll.Head = newNode
		ll.Tail = newNode
	} else {
		ll.Tail.Next = newNode
		ll.Tail = newNode
	}
	ll.size++
}

// Prepend adds a new node with the given value to the beginning of the list.
// Time complexity: O(1)
func (ll *LinkedList) Prepend(value int) {
	newNode := &Node{Value: value, Next: ll.Head}
	ll.Head = newNode
	
	if ll.Tail == nil {
		ll.Tail = newNode
	}
	ll.size++
}

// InsertAt inserts a new node with the given value at the specified index.
// Returns ErrIndexOutOfRange if the index is invalid.
// Time complexity: O(n)
func (ll *LinkedList) InsertAt(index, value int) error {
	if index < 0 || index > ll.size {
		return ErrIndexOutOfRange
	}
	
	if index == 0 {
		ll.Prepend(value)
		return nil
	}
	
	if index == ll.size {
		ll.Append(value)
		return nil
	}
	
	newNode := &Node{Value: value}
	current := ll.Head
	
	for i := 0; i < index-1; i++ {
		current = current.Next
	}
	
	newNode.Next = current.Next
	current.Next = newNode
	ll.size++
	
	return nil
}

// Delete removes the first occurrence of the specified value from the list.
// Returns ErrEmptyList if the list is empty, or an error if the value is not found.
// Time complexity: O(n)
func (ll *LinkedList) Delete(value int) error {
	if ll.Head == nil {
		return ErrEmptyList
	}
	
	if ll.Head.Value == value {
		ll.Head = ll.Head.Next
		ll.size--
		if ll.Head == nil {
			ll.Tail = nil
		}
		return nil
	}
	
	current := ll.Head
	for current.Next != nil {
		if current.Next.Value == value {
			if current.Next == ll.Tail {
				ll.Tail = current
			}
			current.Next = current.Next.Next
			ll.size--
			return nil
		}
		current = current.Next
	}
	
	return fmt.Errorf("value %d not found in list", value)
}

// DeleteAt removes the node at the specified index.
// Returns ErrIndexOutOfRange if the index is invalid.
// Time complexity: O(n)
func (ll *LinkedList) DeleteAt(index int) error {
	if index < 0 || index >= ll.size {
		return ErrIndexOutOfRange
	}
	
	if index == 0 {
		ll.Head = ll.Head.Next
		ll.size--
		if ll.Head == nil {
			ll.Tail = nil
		}
		return nil
	}
	
	current := ll.Head
	for i := 0; i < index-1; i++ {
		current = current.Next
	}
	
	if current.Next == ll.Tail {
		ll.Tail = current
	}
	current.Next = current.Next.Next
	ll.size--
	
	return nil
}

// Find searches for the first node with the given value.
// Returns the node and true if found, nil and false otherwise.
// Time complexity: O(n)
func (ll *LinkedList) Find(value int) (*Node, bool) {
	current := ll.Head
	for current != nil {
		if current.Value == value {
			return current, true
		}
		current = current.Next
	}
	return nil, false
}

// GetAt returns the value at the specified index.
// Returns ErrIndexOutOfRange if the index is invalid.
// Time complexity: O(n)
func (ll *LinkedList) GetAt(index int) (int, error) {
	if index < 0 || index >= ll.size {
		return 0, ErrIndexOutOfRange
	}
	
	current := ll.Head
	for i := 0; i < index; i++ {
		current = current.Next
	}
	
	return current.Value, nil
}

// Size returns the number of nodes in the list.
// Time complexity: O(1)
func (ll *LinkedList) Size() int {
	return ll.size
}

// IsEmpty returns true if the list is empty.
// Time complexity: O(1)
func (ll *LinkedList) IsEmpty() bool {
	return ll.size == 0
}

// Clear removes all nodes from the list.
// Time complexity: O(1)
func (ll *LinkedList) Clear() {
	ll.Head = nil
	ll.Tail = nil
	ll.size = 0
}

// ToSlice converts the linked list to a slice of integers.
// Time complexity: O(n)
func (ll *LinkedList) ToSlice() []int {
	if ll.size == 0 {
		return []int{}
	}
	
	result := make([]int, 0, ll.size)
	current := ll.Head
	
	for current != nil {
		result = append(result, current.Value)
		current = current.Next
	}
	
	return result
}

// Reverse reverses the linked list in place.
// Time complexity: O(n)
func (ll *LinkedList) Reverse() {
	if ll.Head == nil || ll.Head.Next == nil {
		return
	}
	
	var prev *Node
	current := ll.Head
	ll.Tail = ll.Head
	
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}
	
	ll.Head = prev
}

