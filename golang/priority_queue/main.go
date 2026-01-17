package main

import (
	"container/heap"
	"fmt"
)

type Item struct {
	priority float64
	name     string
}

type PriorityItemQueue []Item

func (p PriorityItemQueue) Len() int { return len(p) }
func (p PriorityItemQueue) Less(i, j int) bool {
	return p[i].priority < p[j].priority // Max-heap
}
func (p PriorityItemQueue) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p *PriorityItemQueue) Push(x interface{}) {
	*p = append(*p, x.(Item))
}

func (p *PriorityItemQueue) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	*p = old[:n-1]
	return item
}

func main() {
	priorityItemQueue := &PriorityItemQueue{}
	heap.Init(priorityItemQueue)

	// MUST use heap.Push to maintain ordering
	heap.Push(priorityItemQueue, Item{1.0, "priority 1"})
	heap.Push(priorityItemQueue, Item{1.0, "priority 1"})
	heap.Push(priorityItemQueue, Item{1.1, "priority 1.1"})
	heap.Push(priorityItemQueue, Item{1.2, "priority 1.2"})
	heap.Push(priorityItemQueue, Item{2.0, "priority 2"})
	heap.Push(priorityItemQueue, Item{3.0, "priority 3"})
	heap.Push(priorityItemQueue, Item{5.0, "priority 5"})
	heap.Push(priorityItemQueue, Item{6.0, "priority 6"})

	for priorityItemQueue.Len() > 0 {
		fmt.Println(heap.Pop(priorityItemQueue).(Item))
	}
}
