# Visitor Pattern in Java - Shape Example

## Overview

This is a simple implementation of the Visitor design pattern using geometric shapes in Java.

## What is the Visitor Pattern?

The Visitor pattern is a behavioral design pattern that lets you separate algorithms from the objects on which they operate. It allows you to add new operations to existing object structures without modifying those structures.

## Structure

```
┌─────────────────┐
│     Shape       │ (Interface)
│   + accept()    │
└────────┬────────┘
         │
    ┌────┴────┬─────────┬──────────┐
    │         │         │          │
┌───▼──┐  ┌──▼──┐  ┌───▼────┐    ...
│Circle│  │Rect │  │Triangle│
└──────┘  └─────┘  └────────┘

┌──────────────────┐
│  ShapeVisitor    │ (Interface)
│ + visitCircle()  │
│ + visitRect()    │
│ + visitTriangle()│
└────────┬─────────┘
         │
    ┌────┴────┬──────────┬─────────┐
    │         │          │         │
┌───▼──────┐ ┌▼─────┐ ┌─▼──────┐ ...
│Area      │ │Peri  │ │Draw    │
│Calculator│ │Calc  │ │Visitor │
└──────────┘ └──────┘ └────────┘
```

## Files

- **Shape.java** - Interface for all shapes
- **ShapeVisitor.java** - Visitor interface
- **Circle.java** - Circle shape implementation
- **Rectangle.java** - Rectangle shape implementation
- **Triangle.java** - Triangle shape implementation
- **AreaCalculator.java** - Visitor that calculates areas
- **PerimeterCalculator.java** - Visitor that calculates perimeters
- **DrawVisitor.java** - Visitor that draws shapes
- **VisitorPatternDemo.java** - Main demo class

## How It Works

1. **Define Visitor Interface**: Declares visit methods for each concrete element type
2. **Define Element Interface**: Declares an accept method that takes a visitor
3. **Concrete Elements**: Implement accept by calling visitor.visit(this)
4. **Concrete Visitors**: Implement specific operations for each element type

## Key Benefits

✅ **Open/Closed Principle**: Add new operations without changing shape classes  
✅ **Single Responsibility**: Separate operations from data structures  
✅ **Type Safety**: Compile-time checking for all operations  
✅ **Clean Code**: Business logic separated from data

## Running the Example

### Compile

```bash
javac *.java
```

### Run

```bash
java VisitorPatternDemo
```

## Expected Output

```
═══════════════════════════════════════════
  VISITOR PATTERN - Simple Shape Example
═══════════════════════════════════════════

📐 CALCULATING AREAS:
───────────────────────────────────────────
  Circle area: 78.54
  Rectangle area: 24.00
  Triangle area: 24.00

✓ Total Area: 126.54

📏 CALCULATING PERIMETERS:
───────────────────────────────────────────
  Circle perimeter: 31.42
  Rectangle perimeter: 20.00
  Triangle perimeter: 18.00

✓ Total Perimeter: 69.42

🎨 DRAWING SHAPES:
───────────────────────────────────────────
  Drawing Circle:
     ___
    /   \
   |  O  |
    \___/

  Drawing Rectangle:
   ┌─────┐
   │     │
   │     │
   └─────┘

  Drawing Triangle:
      /\
     /  \
    /____\

═══════════════════════════════════════════
✨ Key Benefit:
   We performed 3 operations on 3 shape types
   WITHOUT modifying the shape classes!
   Adding new operations is easy - just create
   a new Visitor!
═══════════════════════════════════════════
```

## Adding New Operations

To add a new operation (e.g., JSON export):

1. Create a new class implementing `ShapeVisitor`
2. Implement all visit methods
3. Use it with existing shapes - no changes needed to shape classes!

```java
public class JsonExporter implements ShapeVisitor {
    @Override
    public void visitCircle(Circle circle) {
        System.out.println("{\"type\":\"circle\",\"radius\":" + circle.getRadius() + "}");
    }
    
    @Override
    public void visitRectangle(Rectangle rectangle) {
        System.out.println("{\"type\":\"rectangle\",\"width\":" + rectangle.getWidth() + ",\"height\":" + rectangle.getHeight() + "}");
    }
    
    @Override
    public void visitTriangle(Triangle triangle) {
        System.out.println("{\"type\":\"triangle\",\"base\":" + triangle.getBase() + ",\"height\":" + triangle.getHeight() + "}");
    }
}
```

## When to Use

- Object structure is stable but operations change frequently
- Many distinct operations need to be performed on objects
- You want to keep related operations together
- You want to avoid polluting element classes with many operations

## When NOT to Use

- Object structure changes frequently (you'll need to update all visitors)
- Simple operations that don't justify the pattern's complexity
- When you need to break encapsulation to access private fields

