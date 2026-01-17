package main

import "fmt"


type Graph struct {
	vertex []*vertex

}

type vertex struct {
	Key int
	adjacent []*vertex
}

func NewGraph() *Graph {	
	return &Graph{	
		vertex: []*vertex{},

	}

}

func (g *Graph)AddVertex(key int) {
	newVertex := &vertex{
		Key: key,
	}
	if g.Contains(key) {
		fmt.Printf("Vertex %d already exists \n", key)
		return
	}
	g.vertex = append(g.vertex, newVertex)
}

func (g *Graph)AddEdge(from , to  int) {
	var fromVertex, toVertex *vertex
	for _, v := range g.vertex {
		if v.Key == from {
			fromVertex = v
		}
		if v.Key == to {
			toVertex = v
		}
	}
	if fromVertex == nil || toVertex == nil {
		fmt.Printf("invalid edge %d ---> %d \n", from, to)
		return
	}
	fromVertex.adjacent = append(fromVertex.adjacent, toVertex)
	
}

func (g Graph)Contains(key int) bool {
	for _, v := range g.vertex {
		if v.Key == key {
			return true
		}
	}
	return false
}

func (g *Graph)PrintGraph() {
	for _, v := range g.vertex {
		print("Vertex ", v.Key, " : ")
		for _, a := range v.adjacent {
			print(a.Key, " ")
		}
		println()
	}
}


func main(){
	
	graph := NewGraph()
	graph.AddVertex(1)
	graph.AddVertex(2)
	graph.AddVertex(3)
	graph.AddVertex(1)
	graph.AddEdge(1,2)
	graph.AddEdge(4,3)	
	
graph.PrintGraph()	
}