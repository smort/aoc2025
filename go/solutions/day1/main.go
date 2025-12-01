package main

import (
	"fmt"

	"github.com/smort/aoc2025/util"
)

func main() {
	// part1("input.txt")
	part2("input.txt")
}

func part1(filename string) {
	result := 0
	dial := getDial()
	lines := util.GetLines(filename)
	curr := dial.head
	for {
		if curr.data == 50 {
			break
		}

		curr = curr.next
	}

	for _, line := range lines {
		// println(line)
		if line == "" {
			continue
		}

		direction := string(line[0])
		distance := util.MustConvAtoi(line[1:])
		// println(direction, distance)
		for range distance {
			if direction == "R" {
				curr = curr.next
			} else {
				curr = curr.prev
			}
			// println(curr.data)
		}

		if curr.data == 0 {
			result++
		}
	}

	fmt.Println(result)
}

func part2(filename string) {
	result := 0
	dial := getDial()
	lines := util.GetLines(filename)
	curr := dial.head
	for {
		if curr.data == 50 {
			break
		}

		curr = curr.next
	}

	for _, line := range lines {
		// println(line)
		if line == "" {
			continue
		}

		direction := string(line[0])
		distance := util.MustConvAtoi(line[1:])
		// println(direction, distance)
		for range distance {
			if direction == "R" {
				curr = curr.next
			} else {
				curr = curr.prev
			}
			// println(curr.data)
			if curr.data == 0 {
				result++
			}
		}
	}

	fmt.Println(result)
}

func getDial() CircularDoublyLinkedList {
	dial := CircularDoublyLinkedList{}
	for i := range 100 {
		dial.InsertEnd(&Node{data: i})
	}

	return dial
}

type Node struct {
	data int
	next *Node
	prev *Node
}

type CircularDoublyLinkedList struct {
	head   *Node
	tail   *Node
	length int
}

func (list *CircularDoublyLinkedList) InsertToEmptyLinkedList(n *Node) {
	// Set the head and tail nodes to the new node
	list.head = n
	list.tail = list.head
	// Make the list circular
	list.head.next = list.tail
	list.tail.prev = list.head
	list.tail.next = list.head
	list.head.prev = list.tail
	// Increment the length of the list
	list.length++
}

func (list *CircularDoublyLinkedList) InsertBeginning(n *Node) {
	// Check if the list is empty
	if list.head == nil {
		list.InsertToEmptyLinkedList(n)
	} else {
		// Insert the node at the beginning
		n.next = list.head
		n.prev = list.tail
		list.tail.next = n
		list.head.prev = n
		list.head = n
		// Increment the length of the list
		list.length++
	}
}

func (list *CircularDoublyLinkedList) InsertEnd(n *Node) {
	// Check if the list is empty
	if list.head == nil {
		list.InsertToEmptyLinkedList(n)
	} else {
		// Insert the node at the end
		n.prev = list.tail
		list.tail.next = n
		list.tail = n
		// Make the list circular
		list.tail.next = list.head
		list.head.prev = list.tail
		// Increment the length of the list
		list.length++
	}
}
