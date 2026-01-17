package main

import (
	"fmt"
	"strings"
)

type Node struct {
	Data string
	Next *Node
}
type LinkedList struct {
	Head *Node
}

func NewList() *LinkedList {
	return &LinkedList{}
}

func (l *LinkedList) InsertAfterNodeValue(searchData string, dataValue string) {
	currentNode := l.Head
	for currentNode != nil {
		if strings.EqualFold(searchData, currentNode.Data) {
			currentNode.Next = &Node{Data: dataValue, Next: currentNode.Next}
			return
		}
		currentNode = currentNode.Next

	}
}

func (l *LinkedList) InsertBeforeNodeValue(searchData string, dataValue string) {
	currentNode := l.Head
	if strings.EqualFold(searchData, currentNode.Data) {
		l.Head = &Node{Data: dataValue, Next: currentNode}
		return
	}
	prev := currentNode
	currentNode = l.Head

	for currentNode != nil {
		if strings.EqualFold(searchData, currentNode.Data) {
			prev.Next = &Node{Data: dataValue, Next: currentNode}
			return
		}
		prev = currentNode
		currentNode = currentNode.Next
	}
}
func (l *LinkedList) AddNodeAtEnd(dataValue string) {
	newNode := &Node{Data: dataValue, Next: nil}
	if l.Head == nil {
		l.Head = newNode
		return
	}
	currentNode := l.Head
	for currentNode.Next != nil {
		currentNode = currentNode.Next
	}
	currentNode.Next = newNode

}

func (l *LinkedList) AddNodeAtTheFront(dataValue string) {
	newNode := &Node{Data: dataValue, Next: nil}
	newNode.Next = l.Head
	l.Head = newNode
	return
}

func (l *LinkedList) DeleteFirstNode() {
	nodeToRemove := l.Head
	l.Head = nodeToRemove.Next

}
func (l *LinkedList) DeleteLastNode() {
	currentNode := l.Head
	if l.Head == nil {
		return
	}
	if l.Head.Next == nil {
		l.Head = nil
		return
	}
	prev := l.Head
	for currentNode.Next != nil {
		prev = currentNode
		currentNode = currentNode.Next
	}
	prev.Next = nil

}

func (l *LinkedList) ListValues() {
	currentNode := l.Head
	for currentNode != nil {
		fmt.Println(currentNode.Data)
		currentNode = currentNode.Next

	}
}

func main() {
	ll := NewList()
	ll.AddNodeAtEnd("Rahul")
	ll.AddNodeAtEnd("Bhausaheb")
	ll.ListValues()
	println("------")
	ll.DeleteLastNode()
	ll.ListValues()

}
