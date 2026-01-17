# Visitor Pattern in Go

## What is the Visitor Pattern?

The Visitor pattern is a behavioral design pattern that lets you separate algorithms from the objects on which they operate. It allows you to add new operations to existing object structures without modifying those structures.

## When to Use

- When you need to perform operations on all elements of a complex object structure (like a composite tree)
- When you want to add new operations without changing the classes of the elements
- When many distinct and unrelated operations need to be performed on objects in an object structure
- When the object structure rarely changes, but you often need to define new operations over it

## Benefits

✅ **Open/Closed Principle**: You can introduce new behaviors without changing existing code  
✅ **Single Responsibility Principle**: Multiple versions of the same behavior can be moved into the same class  
✅ **Clean code**: Business logic is separated from the data structure  
✅ **Type safety**: Compile-time type checking for all operations

## Drawbacks

❌ Must update all visitors when adding/removing element classes  
❌ Elements might not have access to private fields/methods when needed by visitors  
❌ Can make code more complex for simple use cases

## Structure

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       │ uses
       │
       ▼
┌─────────────┐           ┌──────────────┐
│  Visitor    │◄─────────│  Element     │
│ (Interface) │           │ (Interface)  │
└──────┬──────┘           └──────┬───────┘
       │                          │
       │                          │
   ┌───┴────┬────────┐       ┌───┴────┬─────────┐
   │        │        │       │        │         │
┌──▼───┐ ┌─▼────┐ ┌─▼───┐ ┌─▼────┐ ┌─▼──────┐ ...
│Visitor││Visitor││Visitor││Element││Element │
│  A    ││  B    ││  C    ││   A   ││   B    │
└───────┘ └──────┘ └─────┘ └──────┘ └────────┘
```

## Key Components

1. **Visitor Interface**: Declares visit methods for each concrete element type
2. **Concrete Visitors**: Implement the visitor interface with specific operations
3. **Element Interface**: Declares an Accept method that takes a visitor
4. **Concrete Elements**: Implement the Accept method by calling the visitor's visit method

## Real-World Examples

### Example 1: E-commerce System
- **Elements**: Different product types (Electronics, Clothing, Books)
- **Visitors**: Tax calculator, Shipping cost calculator, Discount calculator

### Example 2: Document Processing
- **Elements**: Paragraph, Image, Table, Heading
- **Visitors**: HTML exporter, PDF exporter, Plain text exporter

### Example 3: Company Structure
- **Elements**: Engineers, Managers, Executives
- **Visitors**: Salary calculator, Performance evaluator, Vacation calculator

## Code Examples

This repository contains three examples:

1. **`main.go`** - E-commerce system with product types and calculations
2. **`document_example.go`** - Document structure with different exporters
3. **`shape_example.go`** - Simple geometric shapes with different operations

## How It Works

```go
// 1. Define the Visitor interface
type Visitor interface {
    VisitElementA(e *ElementA)
    VisitElementB(e *ElementB)
}

// 2. Define the Element interface
type Element interface {
    Accept(v Visitor)
}

// 3. Concrete Elements implement Accept
type ElementA struct { /* ... */ }

func (e *ElementA) Accept(v Visitor) {
    v.VisitElementA(e)  // Double dispatch!
}

// 4. Concrete Visitors implement operations
type ConcreteVisitor struct { /* ... */ }

func (cv *ConcreteVisitor) VisitElementA(e *ElementA) {
    // Perform operation on ElementA
}

// 5. Client code
element := &ElementA{}
visitor := &ConcreteVisitor{}
element.Accept(visitor)  // Executes the operation
```

## Double Dispatch

The Visitor pattern uses **double dispatch** - a technique that determines the method to call based on two objects:
1. The type of the element (first dispatch via Accept)
2. The type of the visitor (second dispatch via Visit method)

This allows you to add new operations (visitors) without modifying element classes.

## Visitor vs Other Patterns

| Pattern | Purpose | When to Use |
|---------|---------|-------------|
| **Visitor** | Add operations to objects | Object structure is stable, operations change |
| **Strategy** | Encapsulate algorithms | Switch between algorithms at runtime |
| **Command** | Encapsulate requests | Queue, log, or undo operations |
| **Iterator** | Traverse collections | Access elements sequentially |

## Best Practices

1. **Keep Element interface stable**: Avoid adding/removing element types frequently
2. **Use meaningful names**: Name visitors after their operation (e.g., `TaxCalculator`, `HTMLExporter`)
3. **Consider type safety**: Go's type system ensures all visitors handle all elements
4. **Aggregate results**: Store results in visitor fields for later retrieval
5. **Handle errors gracefully**: Return errors from Visit methods if needed

## Common Pitfalls

❌ **Don't use for frequently changing object structures**: Every new element type requires updating all visitors  
❌ **Don't overuse**: Simple operations might not need a visitor pattern  
❌ **Don't break encapsulation**: Visitors shouldn't access private element data unnecessarily

## Running the Examples

```bash
# Run the e-commerce example
go run main.go

# Run the document example
go run document_example.go

# Run the shape example
go run shape_example.go
```

## Further Reading

- [Design Patterns: Elements of Reusable Object-Oriented Software](https://en.wikipedia.org/wiki/Design_Patterns) (Gang of Four)
- [Refactoring Guru - Visitor Pattern](https://refactoring.guru/design-patterns/visitor)
- [Source Making - Visitor Pattern](https://sourcemaking.com/design_patterns/visitor)

## Summary

The Visitor pattern is powerful when you need to perform many different operations on a stable object structure. It promotes clean separation of concerns and makes adding new operations easy. However, it comes with the trade-off of making it harder to add new element types.

Choose the Visitor pattern when operations change more frequently than the object structure itself.


