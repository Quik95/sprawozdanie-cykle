package Cykle

import (
	"fmt"
	"strings"
)

// LinkedList holds the elements, where each element points to the next and previous element
type LinkedList struct {
	first *element
	last  *element
	size  int
}

type element struct {
	value int
	prev  *element
	next  *element
	other *element
}

// NewLinkedList instantiates a new list and adds the passed values, if any, to the list
func NewLinkedList(values ...int) *LinkedList {
	list := &LinkedList{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

// Add appends a value (one or more) at the end of the list (same as Append())
func (list *LinkedList) Add(values ...int) {
	for _, value := range values {
		newElement := &element{value: value, prev: list.last}
		if list.size == 0 {
			list.first = newElement
			list.last = newElement
		} else {
			list.last.next = newElement
			list.last = newElement
		}
		list.size++
	}
}

// AddSingle appends a value at the end of the list
func (list *LinkedList) AddSingle(value, currentValue int, otherList *LinkedList) {
	otherList.Add(currentValue)

	newElement := &element{value: value, prev: list.last, other: otherList.last}
	if list.size == 0 {
		list.first = newElement
		list.last = newElement
	} else {
		list.last.next = newElement
		list.last = newElement
	}
	otherList.last.other = newElement
	list.size++
}

// Append appends a value (one or more) at the end of the list (same as Add())
func (list *LinkedList) Append(values ...int) {
	list.Add(values...)
}

// Prepend prepends a values (or more)
func (list *LinkedList) Prepend(values ...int) {
	for v := len(values) - 1; v >= 0; v-- {
		newElement := &element{value: values[v], next: list.first}
		if list.size == 0 {
			list.first = newElement
			list.last = newElement
		} else {
			list.first.prev = newElement
			list.first = newElement
		}
		list.size++
	}
}

// Get returns the element at index.
// Second return parameter is true if index is within bounds of the array and array is not empty, otherwise false.
func (list *LinkedList) Get(index int) (int, bool) {

	if !list.withinRange(index) {
		return -1, false
	}

	// determine traversal direction, last to first or first to last
	if list.size-index < index {
		element := list.last
		for e := list.size - 1; e != index; e, element = e-1, element.prev {
		}
		return element.value, true
	}
	element := list.first
	for e := 0; e != index; e, element = e+1, element.next {
	}
	return element.value, true
}

// Remove removes the element with the given value from the list.
func (list *LinkedList) Remove(value int, otherList *LinkedList) {
	if !list.Contains(value) {
		return
	}

	element := list.first
	for element.next != nil {
		if element.value == value {
			break
		}
		element = element.next
	}

	if element == list.first {
		list.first = element.next
	}
	if element.other == otherList.first {
		otherList.first = element.other.next
	}
	if element == list.last {
		list.last = element.prev
	}
	if element.other == otherList.last {
		otherList.last = element.other.prev
	}
	if element.prev != nil {
		element.prev.next = element.next
	}
	if element.other.prev != nil {
		element.other.prev.next = element.other.next
	}
	if element.next != nil {
		element.next.prev = element.prev
	}
	if element.other.next != nil {
		element.other.next.prev = element.other.prev
	}

	element.other = nil
	element = nil

	list.size--
	otherList.size--
}

// Contains check if values (one or more) are present in the set.
// All values have to be present in the set for the method to return true.
// Performance time complexity of n^2.
// Returns true if no arguments are passed at all, i.e. set is always super-set of empty set.
func (list *LinkedList) Contains(values ...int) bool {
	if len(values) == 0 {
		return true
	}
	if list.size == 0 {
		return false
	}
	for _, value := range values {
		found := false
		for element := list.first; element != nil; element = element.next {
			if element.value == value {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Values returns all elements in the list.
func (list *LinkedList) Values() []int {
	values := make([]int, list.size, list.size)
	for e, element := 0, list.first; element != nil; e, element = e+1, element.next {
		values[e] = element.value
	}
	return values
}

//IndexOf returns index of provided element
func (list *LinkedList) IndexOf(value int) int {
	if list.size == 0 {
		return -1
	}
	for index, element := range list.Values() {
		if element == value {
			return index
		}
	}
	return -1
}

// Empty returns true if list does not contain any elements.
func (list *LinkedList) Empty() bool {
	return list.size == 0
}

// Size returns number of elements within the list.
func (list *LinkedList) Size() int {
	return list.size
}

// Clear removes all elements from the list.
func (list *LinkedList) Clear() {
	list.size = 0
	list.first = nil
	list.last = nil
}

// String returns a string representation of container
func (list *LinkedList) String() string {
	var values []string
	for element := list.first; element != nil; element = element.next {
		values = append(values, fmt.Sprintf("%v", element.value))
	}
	return strings.Join(values, " -> ")
}

// Check that the index is within bounds of the list
func (list *LinkedList) withinRange(index int) bool {
	return index >= 0 && index < list.size
}
