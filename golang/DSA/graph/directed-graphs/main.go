package main

import (
	"fmt"
)
// DirectedGraph represents a simple directed graph using an adjacency list
type DirectedGraph struct {
	adjacencyList map[int][]int
}
// NewDirectedGraph initializes and returns a new DirectedGraph
func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		adjacencyList: make(map[int][]int),
	}
}
// AddVertex adds a new vertex to the directed graph
func (g *DirectedGraph) AddVertex(vertex int) {
	if _, exists := g.adjacencyList[vertex]; !exists {
		g.adjacencyList[vertex] = []int{}
		return
	}
	fmt.Printf("Vertex %d already exists \n", vertex)
}

// AddEdge adds a directed edge from v1 to v2
func (g *DirectedGraph) AddEdge(v1, v2 int) {
	g.AddVertex(v1)
	g.AddVertex(v2)
	g.adjacencyList[v1] = append(g.adjacencyList[v1], v2)

}

// PrintGraph prints the adjacency list of the directed graph
func (g *DirectedGraph) PrintGraph() {
	for vertex, neighbors := range g.adjacencyList {
		fmt.Printf("%d: %v\n", vertex, neighbors)
	}
}


func (g *DirectedGraph) RemoveEdge(v1, v2 int) {
	g.adjacencyList[v1] = removeFromSlice(g.adjacencyList[v1], v2)
}
func removeFromSlice(slice []int, value int) []int {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}


func (g *DirectedGraph) RemoveVertex(vertex int) {
	delete(g.adjacencyList, vertex)
	for v, neighbors := range g.adjacencyList {
		g.adjacencyList[v] = removeFromSlice(neighbors, vertex)
	}
}


func (	g *DirectedGraph) BFS(start int) {
	visited := make(map[int]bool)
	queue := []int{start}
	visited[start] = true

	fmt.Print("BFS: ")
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		fmt.Print(current, " ")

		for _, neighbor := range g.adjacencyList[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	fmt.Println()
}

func (g *DirectedGraph) DFS(start int) {
	visited := make(map[int]bool)
	stack := []int{start}

	fmt.Print("DFS: ")
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if !visited[current] {
			visited[current] = true
			fmt.Print(current, " ")

			for i := len(g.adjacencyList[current]) - 1; i >= 0; i-- {
				neighbor := g.adjacencyList[current][i]
				if !visited[neighbor] {
					stack = append(stack, neighbor)
				}
			}
		}
	}
	fmt.Println()
}

func (g *DirectedGraph) dfsHelper(vertex int, visited map[int]bool) {
	visited[vertex] = true
	fmt.Printf("%d ", vertex)
	for _, neighbor := range g.adjacencyList[vertex] {
		if !visited[neighbor] {
			g.dfsHelper(neighbor, visited)
		}
	}
}


func (g *DirectedGraph) DFS_recursion(start int) {
	visited := make(map[int]bool)
	fmt.Print("DFS (recursion): ")	
	g.dfsHelper(start, visited)
	fmt.Println()



}
func main() {
	graph := NewDirectedGraph()
	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(2, 4)
	graph.AddEdge(3, 4)
	graph.AddEdge(4, 5)

	graph.PrintGraph()
	graph.BFS(1)
		graph.BFS(5)
		graph.DFS_recursion(1)
}
