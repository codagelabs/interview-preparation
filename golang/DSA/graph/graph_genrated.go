package main

import (
	"fmt"
)

// Graph represents a simple undirected graph using an adjacency list
type Graph struct {
	adjacencyList map[int][]int
}

// NewGraph initializes and returns a new Graph
func NewGraph() *Graph {
	return &Graph{
		adjacencyList: make(map[int][]int),
	}
}

// AddVertex adds a new vertex to the graph
func (g *Graph) AddVertex(vertex int) {
	if _, exists := g.adjacencyList[vertex]; !exists {
		g.adjacencyList[vertex] = []int{}
	}
}

// AddEdge adds an undirected edge between two vertices
func (g *Graph) AddEdge(v1, v2 int) {
	g.AddVertex(v1)
	g.AddVertex(v2)
	g.adjacencyList[v1] = append(g.adjacencyList[v1], v2)
	g.adjacencyList[v2] = append(g.adjacencyList[v2], v1)
}

// RemoveEdge removes an undirected edge between two vertices
func (g *Graph) RemoveEdge(v1, v2 int) {
	g.adjacencyList[v1] = removeFromSlice(g.adjacencyList[v1], v2)
	g.adjacencyList[v2] = removeFromSlice(g.adjacencyList[v2], v1)
}

// RemoveVertex removes a vertex and all its edges from the graph
func (g *Graph) RemoveVertex(vertex int) {
	for _, neighbor := range g.adjacencyList[vertex] {
		g.adjacencyList[neighbor] = removeFromSlice(g.adjacencyList[neighbor], vertex)
	}
	delete(g.adjacencyList, vertex)
}

// PrintGraph prints the adjacency list of the graph
func (g *Graph) PrintGraph() {
	for vertex, neighbors := range g.adjacencyList {
		fmt.Printf("%d: %v\n", vertex, neighbors)
	}
}

// BFS performs Breadth-First Search starting from a given vertex
func (g *Graph) BFS(start int) {
	visited := make(map[int]bool)
	queue := []int{start}
	visited[start] = true

	fmt.Print("BFS: ")
	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		fmt.Printf("%d ", vertex)

		for _, neighbor := range g.adjacencyList[vertex] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	fmt.Println()
}

// DFS performs Depth-First Search starting from a given vertex
func (g *Graph) DFS(start int) {
	visited := make(map[int]bool)
	fmt.Print("DFS: ")
	g.dfsHelper(start, visited)
	fmt.Println()
}

func (g *Graph) dfsHelper(vertex int, visited map[int]bool) {
	visited[vertex] = true
	fmt.Printf("%d ", vertex)
	for _, neighbor := range g.adjacencyList[vertex] {
		if !visited[neighbor] {
			g.dfsHelper(neighbor, visited)
		}
	}
}

// Helper function to remove an element from a slice
func removeFromSlice(slice []int, value int) []int {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func main() {
	graph := NewGraph()
	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(2, 4)
	graph.AddEdge(3, 4)
	graph.AddEdge(4, 5)

	fmt.Println("Graph adjacency list:")
	graph.PrintGraph()

	graph.BFS(1)
	graph.DFS(1)

	fmt.Println("Removing edge (1,3) and vertex 5...")
	graph.RemoveEdge(1, 3)
	graph.RemoveVertex(5)
	graph.PrintGraph()
}
