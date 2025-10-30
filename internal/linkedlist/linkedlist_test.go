package linkedlist

import (
	"testing"
)

func TestNew(t *testing.T) {
	ll := New()
	
	if ll.Head != nil {
		t.Error("expected Head to be nil")
	}
	if ll.Tail != nil {
		t.Error("expected Tail to be nil")
	}
	if ll.Size() != 0 {
		t.Errorf("expected size to be 0, got %d", ll.Size())
	}
}

func TestLinkedList_Append(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{
			name:   "append to empty list",
			values: []int{1},
			want:   []int{1},
		},
		{
			name:   "append multiple values",
			values: []int{1, 2, 3, 4, 5},
			want:   []int{1, 2, 3, 4, 5},
		},
		{
			name:   "append same values",
			values: []int{7, 7, 7},
			want:   []int{7, 7, 7},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := New()
			for _, v := range tt.values {
				ll.Append(v)
			}
			
			got := ll.ToSlice()
			if !slicesEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
			
			if ll.Size() != len(tt.want) {
				t.Errorf("size = %d, want %d", ll.Size(), len(tt.want))
			}
		})
	}
}

func TestLinkedList_Prepend(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{
			name:   "prepend to empty list",
			values: []int{1},
			want:   []int{1},
		},
		{
			name:   "prepend multiple values",
			values: []int{1, 2, 3},
			want:   []int{3, 2, 1},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := New()
			for _, v := range tt.values {
				ll.Prepend(v)
			}
			
			got := ll.ToSlice()
			if !slicesEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkedList_InsertAt(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		index     int
		value     int
		want      []int
		wantError bool
	}{
		{
			name:      "insert at beginning",
			initial:   []int{1, 2, 3},
			index:     0,
			value:     0,
			want:      []int{0, 1, 2, 3},
			wantError: false,
		},
		{
			name:      "insert in middle",
			initial:   []int{1, 2, 4},
			index:     2,
			value:     3,
			want:      []int{1, 2, 3, 4},
			wantError: false,
		},
		{
			name:      "insert at end",
			initial:   []int{1, 2, 3},
			index:     3,
			value:     4,
			want:      []int{1, 2, 3, 4},
			wantError: false,
		},
		{
			name:      "insert out of range",
			initial:   []int{1, 2, 3},
			index:     10,
			value:     4,
			want:      []int{1, 2, 3},
			wantError: true,
		},
		{
			name:      "insert negative index",
			initial:   []int{1, 2, 3},
			index:     -1,
			value:     4,
			want:      []int{1, 2, 3},
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := createList(tt.initial)
			err := ll.InsertAt(tt.index, tt.value)
			
			if (err != nil) != tt.wantError {
				t.Errorf("InsertAt() error = %v, wantError %v", err, tt.wantError)
				return
			}
			
			got := ll.ToSlice()
			if !slicesEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkedList_Delete(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		delete    int
		want      []int
		wantError bool
	}{
		{
			name:      "delete from empty list",
			initial:   []int{},
			delete:    1,
			want:      []int{},
			wantError: true,
		},
		{
			name:      "delete first element",
			initial:   []int{1, 2, 3},
			delete:    1,
			want:      []int{2, 3},
			wantError: false,
		},
		{
			name:      "delete middle element",
			initial:   []int{1, 2, 3},
			delete:    2,
			want:      []int{1, 3},
			wantError: false,
		},
		{
			name:      "delete last element",
			initial:   []int{1, 2, 3},
			delete:    3,
			want:      []int{1, 2},
			wantError: false,
		},
		{
			name:      "delete non-existent element",
			initial:   []int{1, 2, 3},
			delete:    99,
			want:      []int{1, 2, 3},
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := createList(tt.initial)
			err := ll.Delete(tt.delete)
			
			if (err != nil) != tt.wantError {
				t.Errorf("Delete() error = %v, wantError %v", err, tt.wantError)
				return
			}
			
			got := ll.ToSlice()
			if !slicesEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkedList_DeleteAt(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		index     int
		want      []int
		wantError bool
	}{
		{
			name:      "delete at beginning",
			initial:   []int{1, 2, 3},
			index:     0,
			want:      []int{2, 3},
			wantError: false,
		},
		{
			name:      "delete in middle",
			initial:   []int{1, 2, 3},
			index:     1,
			want:      []int{1, 3},
			wantError: false,
		},
		{
			name:      "delete at end",
			initial:   []int{1, 2, 3},
			index:     2,
			want:      []int{1, 2},
			wantError: false,
		},
		{
			name:      "delete out of range",
			initial:   []int{1, 2, 3},
			index:     10,
			want:      []int{1, 2, 3},
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := createList(tt.initial)
			err := ll.DeleteAt(tt.index)
			
			if (err != nil) != tt.wantError {
				t.Errorf("DeleteAt() error = %v, wantError %v", err, tt.wantError)
				return
			}
			
			got := ll.ToSlice()
			if !slicesEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkedList_Find(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		find      int
		wantFound bool
	}{
		{
			name:      "find in empty list",
			initial:   []int{},
			find:      1,
			wantFound: false,
		},
		{
			name:      "find existing element",
			initial:   []int{1, 2, 3, 4, 5},
			find:      3,
			wantFound: true,
		},
		{
			name:      "find non-existing element",
			initial:   []int{1, 2, 3},
			find:      99,
			wantFound: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := createList(tt.initial)
			node, found := ll.Find(tt.find)
			
			if found != tt.wantFound {
				t.Errorf("Find() found = %v, want %v", found, tt.wantFound)
			}
			
			if found && node.Value != tt.find {
				t.Errorf("Find() node.Value = %v, want %v", node.Value, tt.find)
			}
		})
	}
}

func TestLinkedList_GetAt(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		index     int
		want      int
		wantError bool
	}{
		{
			name:      "get first element",
			initial:   []int{1, 2, 3},
			index:     0,
			want:      1,
			wantError: false,
		},
		{
			name:      "get middle element",
			initial:   []int{1, 2, 3},
			index:     1,
			want:      2,
			wantError: false,
		},
		{
			name:      "get last element",
			initial:   []int{1, 2, 3},
			index:     2,
			want:      3,
			wantError: false,
		},
		{
			name:      "get out of range",
			initial:   []int{1, 2, 3},
			index:     10,
			want:      0,
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := createList(tt.initial)
			got, err := ll.GetAt(tt.index)
			
			if (err != nil) != tt.wantError {
				t.Errorf("GetAt() error = %v, wantError %v", err, tt.wantError)
				return
			}
			
			if !tt.wantError && got != tt.want {
				t.Errorf("GetAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkedList_Reverse(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		want    []int
	}{
		{
			name:    "reverse empty list",
			initial: []int{},
			want:    []int{},
		},
		{
			name:    "reverse single element",
			initial: []int{1},
			want:    []int{1},
		},
		{
			name:    "reverse multiple elements",
			initial: []int{1, 2, 3, 4, 5},
			want:    []int{5, 4, 3, 2, 1},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := createList(tt.initial)
			ll.Reverse()
			
			got := ll.ToSlice()
			if !slicesEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkedList_IsEmpty(t *testing.T) {
	ll := New()
	if !ll.IsEmpty() {
		t.Error("new list should be empty")
	}
	
	ll.Append(1)
	if ll.IsEmpty() {
		t.Error("list with element should not be empty")
	}
	
	ll.Clear()
	if !ll.IsEmpty() {
		t.Error("cleared list should be empty")
	}
}

func TestLinkedList_Clear(t *testing.T) {
	ll := createList([]int{1, 2, 3, 4, 5})
	
	ll.Clear()
	
	if ll.Size() != 0 {
		t.Errorf("size after Clear() = %d, want 0", ll.Size())
	}
	if ll.Head != nil {
		t.Error("Head should be nil after Clear()")
	}
	if ll.Tail != nil {
		t.Error("Tail should be nil after Clear()")
	}
}

// Helper functions

func createList(values []int) *LinkedList {
	ll := New()
	for _, v := range values {
		ll.Append(v)
	}
	return ll
}

func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

