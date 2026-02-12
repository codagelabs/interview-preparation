package main

import "container/list"

type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	doubly   *list.List
}

type Node struct {
	key   int
	value int
}

func Constructor(capacity int) LRUCache {

	return LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		doubly:   list.New(),
	}		
}

func (this *LRUCache) Get(key int) int {
	
	if _, ok := this.cache[key]; !ok {
		

		return -1
	



			