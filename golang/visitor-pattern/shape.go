package main

import "fmt"

// Visitor interface
type Visitor interface {
	visitCircle(circle *Circle)
	visitRectangle(rectangle *Rectangle)
	visitTriangle(triangle *Triangle)
}

// Shape interface
type Shape interface {
	accept(visitor Visitor)
}

// Circle struct
type Circle struct {
	radius float64
}

func (c *Circle) accept(visitor Visitor) {
	visitor.visitCircle(c)
}

// Rectangle struct
type Rectangle struct {
	width  float64
	height float64
}

func (r *Rectangle) accept(visitor Visitor) {
	visitor.visitRectangle(r)
}

// Triangle struct
type Triangle struct {
	base   float64
	height float64
}

func (t *Triangle) accept(visitor Visitor) {
	visitor.visitTriangle(t)
}

// AreaCalculator struct
type AreaCalculator struct {
	TotalArea float64
}

func (a *AreaCalculator) visitCircle(circle *Circle) {
	a.TotalArea += 3.14 * circle.radius * circle.radius
	fmt.Println("Area of circle is", a.TotalArea)
}

func (a *AreaCalculator) visitRectangle(rectangle *Rectangle) {
	a.TotalArea += rectangle.width * rectangle.height
	fmt.Println("Area of rectangle is", a.TotalArea)
}

func (a *AreaCalculator) visitTriangle(triangle *Triangle) {
	a.TotalArea += 0.5 * triangle.base * triangle.height
	fmt.Println("Area of triangle is", a.TotalArea)
}

func main() {
	circle := &Circle{radius: 10}
	rectangle := &Rectangle{width: 10, height: 20}
	triangle := &Triangle{base: 10, height: 20}
	areaCalculator := &AreaCalculator{}
	areaCalculator.visitCircle(circle)
	areaCalculator.visitRectangle(rectangle)
	areaCalculator.visitTriangle(triangle)
	fmt.Println("Total area is", areaCalculator.TotalArea)
}
