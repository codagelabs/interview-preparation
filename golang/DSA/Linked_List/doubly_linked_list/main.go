package main

import "fmt"

type node struct {
	next *node
	data string
	prev *node
}
type DoublyLinkedList struct {
	head *node
}

func (d *DoublyLinkedList) InsertFromFront(data string) {
	newNode := &node{data: data}
	if d.head == nil {
		d.head = newNode
		return
	}
	currentNode := d.head
	newNode.next = currentNode
	currentNode.prev = newNode
	d.head = newNode

}

func (d *DoublyLinkedList) InsertAtEnd(data string) {
	newNode := &node{data: data}
	if d.head == nil {
		d.head = newNode
		return
	}
	currentNode := d.head
	for currentNode.next != nil {
		currentNode = currentNode.next
	}
	newNode.prev = currentNode
	currentNode.next = newNode

}

func (d *DoublyLinkedList) InsertAfterNodeValue(searchValue, data string) {
	newNode := &node{data: data}
	if d.head == nil {
		d.head = newNode
		return
	}
	currentNode := d.head
	for currentNode.next != nil {
		currentNode = currentNode.next
	}
	newNode.prev = currentNode
	currentNode.next = newNode

}

func (d *DoublyLinkedList) TraverseList() {
	currentNode := d.head
	for currentNode != nil {
		fmt.Println(currentNode.data)
		currentNode = currentNode.next
	}

}

func main() {
	dl := DoublyLinkedList{}
	dl.InsertFromFront("rahul")
	dl.InsertFromFront("shinde")
	dl.InsertAtEnd("bhausaheb")
	dl.InsertAtEnd("mahadu")
	dl.InsertFromFront("shinde1")
	dl.TraverseList()
}
