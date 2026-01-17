# Basic Golang Questions and Answers

## Table of Contents
1. [Basic Concepts](#basic-concepts)
2. [Concurrency](#concurrency)
3. [Error Handling](#error-handling)
4. [Interfaces](#interfaces)
5. [Packages and Modules](#packages-and-modules)

## Basic Concepts

### Q: What is Go?
A: Go (or Golang) is an open-source programming language developed by Google. It's designed to be simple, efficient, and suitable for modern software development.

### Q: What are the main features of Go?
A: 
- Static typing
- Garbage collection
- Built-in concurrency support
- Fast compilation
- Simple and clean syntax
- Strong standard library

### Q: What is a goroutine?
A: A goroutine is a lightweight thread managed by the Go runtime. It's a function that runs concurrently with other functions.

### Q: What is the difference between `var` and `:=`?
A: 
- `var` is used for variable declaration
- `:=` is used for short variable declaration and initialization
- `:=` can only be used inside functions

## Concurrency

### Q: What is a channel in Go?
A: A channel is a typed conduit through which you can send and receive values with the channel operator `<-`. Channels are used for communication between goroutines.

### Q: What is the difference between buffered and unbuffered channels?
A: 
- Unbuffered channels: Sends and receives are synchronized
- Buffered channels: Can hold a specified number of values before blocking

### Q: What is `select` in Go?
A: The `select` statement lets a goroutine wait on multiple communication operations. It blocks until one of its cases can run.

## Error Handling

### Q: How does error handling work in Go?
A: Go uses explicit error handling through return values. Functions that can fail return an error as their last return value.

### Q: What is the difference between `panic` and `error`?
A: 
- `error` is for expected error conditions
- `panic` is for unexpected, unrecoverable errors

## Interfaces

### Q: What is an interface in Go?
A: An interface is a set of method signatures. A type implements an interface by implementing all its methods.

### Q: What is the empty interface?
A: The empty interface `interface{}` can hold values of any type. It's used when you need to work with values of unknown type.

## Packages and Modules

### Q: What is a package in Go?
A: A package is a way to group related code. Every Go file belongs to a package.

### Q: What is a module in Go?
A: A module is a collection of related Go packages that are versioned together. It's defined by a `go.mod` file.

### Q: What is `go mod tidy`?
A: `go mod tidy` ensures that the `go.mod` file matches the source code in the module. It adds any missing module requirements and removes unused ones.

## Common Code Examples

### Basic Hello World
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Goroutine Example
```go
package main

import (
    "fmt"
    "time"
)

func say(s string) {
    for i := 0; i < 5; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Println(s)
    }
}

func main() {
    go say("world")
    say("hello")
}
```

### Channel Example
```go
package main

import "fmt"

func main() {
    ch := make(chan int)
    go func() {
        ch <- 42
    }()
    value := <-ch
    fmt.Println(value)
}
```

## Best Practices

1. **Error Handling**
   - Always check errors
   - Use meaningful error messages
   - Consider using custom error types

2. **Concurrency**
   - Use goroutines for concurrent operations
   - Use channels for communication
   - Avoid shared memory

3. **Code Organization**
   - Keep packages focused and small
   - Use meaningful names
   - Follow Go's standard formatting

4. **Testing**
   - Write tests for your code
   - Use table-driven tests
   - Test edge cases

## Resources

- [Official Go Documentation](https://golang.org/doc/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Tour](https://tour.golang.org/) 