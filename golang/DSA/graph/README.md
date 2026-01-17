# Graphs

This repository explains **Graph data structures** with definitions, examples, and Go (Golang) representations.

---

## What is a Graph?

A **graph** is a data structure consisting of:
- **Vertices (Nodes)**
- **Edges (Connections between nodes)**

Graphs are used to represent networks, relationships, and dependencies.

---

## Types of Graphs

### 1. Undirected Graph
Edges have **no direction**.

Example:
```
1 —— 2
|    |
3 —— 4
```

Go structure:
```go
type UndirectedGraph struct {
    Adj map[int][]int
}

func NewUndirectedGraph() *UndirectedGraph {
    return &UndirectedGraph{Adj: make(map[int][]int)}
}

func (g *UndirectedGraph) AddEdge(u, v int) {
    g.Adj[u] = append(g.Adj[u], v)
    g.Adj[v] = append(g.Adj[v], u)
}
```

---

### 2. Directed Graph
Edges have a **direction (u → v)**.

Example:
```
1 → 2 → 3
↑
4
```

Go structure:
```go
type DirectedGraph struct {
    Adj map[int][]int
}

func NewDirectedGraph() *DirectedGraph {
    return &DirectedGraph{Adj: make(map[int][]int)}
}

func (g *DirectedGraph) AddEdge(u, v int) {
    g.Adj[u] = append(g.Adj[u], v)
}
```

---

## Cyclic and Acyclic Graphs

### Cyclic Graph
A graph that contains at least one **cycle**.

Example:
```
1 → 2 → 3 → 1
```

### Acyclic Graph
A graph with **no cycles**.

Directed Acyclic Graph (DAG):
```
1 → 2 → 3
    ↓
    4
```

---

## Paths in Graphs

A **path** is a sequence of vertices where each consecutive pair has an edge.

Valid Path:
```
1 → 2 → 3 → 4
```

Not a Path:
```
1 → 3
1 → 2 → 3 → 2
```

---

## Degree of a Graph

### Undirected Graph
Degree = number of edges connected to a vertex.

Total Degree:
```
Total Degree = 2 × Number of edges
```

### Directed Graph
- **In-degree**
- **Out-degree**

---

## Edge Weight

### Weighted Graph
Edges have weights.

Example:
```
1 --5-- 2 --3-- 3
```

### Unit Weight Graph
All edges have weight **1**.

---

## Language
- Go (Golang)
