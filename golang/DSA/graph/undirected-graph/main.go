package main

import "fmt"

//Undirected UnDirectedGraph implementation using adjacency list
type UnDirectedGraph struct {
	AdjacencyList  map[int][]int
}

//NewUnDirectedGraph creates a new UnDirectedGraph
func NewUnDirectedGraph() *UnDirectedGraph {
	return &UnDirectedGraph{
		AdjacencyList: make(map[int][]int),
	}
}

//AddVertex adds a vertex to the UnDirectedGraph	
func (g *UnDirectedGraph) AddVertex(vertex int) {
	if _, exists := g.AdjacencyList[vertex]; !exists {
		g.AdjacencyList[vertex] = []int{}
		return
	}
	fmt.Printf("Vertex %d already exists \n", vertex)
}

//AddEdge adds an edge to the UnDirectedGraph
func (g *UnDirectedGraph) AddEdge(v1, v2 int) {
	g.AddVertex(v1)
	g.AddVertex(v2)
	// we are doing undirected UnDirectedGraph so we need to add edge in both directions
	g.AdjacencyList[v1] = append(g.AdjacencyList[v1], v2)
	g.AdjacencyList[v2] = append(g.AdjacencyList[v2], v1)
}

//PrintUnDirectedGraph prints the UnDirectedGraph
func (g *UnDirectedGraph) PrintUnDirectedGraph() {
	for vertex, neighbors := range g.AdjacencyList {
		fmt.Printf("%d: %v\n", vertex, neighbors)
	}
}
//	RemoveEdge removes an edge from the UnDirectedGraph
func (g *UnDirectedGraph) RemoveEdge(v1, v2 int) {
	g.AdjacencyList[v1] = removeFromSlice(g.AdjacencyList[v1], v2)
	g.AdjacencyList[v2] = removeFromSlice(g.AdjacencyList[v2], v1)
}

func removeFromSlice(slice []int, value int) []int {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

//RemoveVertex removes a vertex from the UnDirectedGraph
func (g *UnDirectedGraph) RemoveVertex(vertex int) {
	for _, neighbor := range g.AdjacencyList[vertex] {
		g.AdjacencyList[neighbor] = removeFromSlice(g.AdjacencyList[neighbor], vertex)
	}
	delete(g.AdjacencyList, vertex)
}


func (g *UnDirectedGraph) BFS(start int) {
	visited := make(map[int]bool)
	queue := []int{start}
	visited[start] = true

	fmt.Print("BFS: ")
	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		fmt.Printf("%d ", vertex)

		for _, neighbor := range g.AdjacencyList[vertex] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	fmt.Println()
}

func main() {
	UnDirectedGraph := NewUnDirectedGraph()
	UnDirectedGraph.AddEdge(1, 2)
	UnDirectedGraph.AddEdge(1, 3)
	UnDirectedGraph.AddEdge(2, 4)
	UnDirectedGraph.AddEdge(3, 4)
	UnDirectedGraph.AddEdge(4, 5)
	// 
   	//
	UnDirectedGraph.PrintUnDirectedGraph()
	UnDirectedGraph.BFS(1)
	UnDirectedGraph.BFS(5)
}
