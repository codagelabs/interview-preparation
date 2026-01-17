package main

import (
	"fmt"
	"math"
)

// ============================================================================
// VISITOR PATTERN - GEOMETRIC SHAPES EXAMPLE
// ============================================================================
// This is a simple example showing how to use the Visitor pattern with
// geometric shapes. We can perform different operations (area, perimeter,
// drawing) without modifying the shape classes.
// ============================================================================

// ShapeVisitor defines the visitor interface
type ShapeVisitor interface {
	VisitCircle(c *Circle)
	VisitRectangle(r *Rectangle)
	VisitTriangle(t *Triangle)
}

// Shape is the element interface
type Shape interface {
	Accept(v ShapeVisitor)
}

// ============================================================================
// CONCRETE ELEMENTS - Different Shapes
// ============================================================================

// Circle represents a circle
type Circle struct {
	Radius float64
	X, Y   float64 // center coordinates
}

func (c *Circle) Accept(v ShapeVisitor) {
	v.VisitCircle(c)
}

// Rectangle represents a rectangle
type Rectangle struct {
	Width  float64
	Height float64
	X, Y   float64 // top-left corner coordinates
}

func (r *Rectangle) Accept(v ShapeVisitor) {
	v.VisitRectangle(r)
}

// Triangle represents a triangle
type Triangle struct {
	Base   float64
	Height float64
	X, Y   float64 // base point coordinates
}

func (t *Triangle) Accept(v ShapeVisitor) {
	v.VisitTriangle(t)
}

// ============================================================================
// CONCRETE VISITORS - Different Operations
// ============================================================================

// AreaCalculator calculates the area of shapes
type AreaCalculator struct {
	TotalArea float64
}

func (a *AreaCalculator) VisitCircle(c *Circle) {
	area := math.Pi * c.Radius * c.Radius
	a.TotalArea += area
	fmt.Printf("  ⭕ Circle (radius: %.2f): Area = %.2f\n", c.Radius, area)
}

func (a *AreaCalculator) VisitRectangle(r *Rectangle) {
	area := r.Width * r.Height
	a.TotalArea += area
	fmt.Printf("  ▭ Rectangle (%.2f × %.2f): Area = %.2f\n", r.Width, r.Height, area)
}

func (a *AreaCalculator) VisitTriangle(t *Triangle) {
	area := 0.5 * t.Base * t.Height
	a.TotalArea += area
	fmt.Printf("  △ Triangle (base: %.2f, height: %.2f): Area = %.2f\n", t.Base, t.Height, area)
}

// PerimeterCalculator calculates the perimeter of shapes
type PerimeterCalculator struct {
	TotalPerimeter float64
}

func (p *PerimeterCalculator) VisitCircle(c *Circle) {
	perimeter := 2 * math.Pi * c.Radius
	p.TotalPerimeter += perimeter
	fmt.Printf("  ⭕ Circle (radius: %.2f): Perimeter = %.2f\n", c.Radius, perimeter)
}

func (p *PerimeterCalculator) VisitRectangle(r *Rectangle) {
	perimeter := 2 * (r.Width + r.Height)
	p.TotalPerimeter += perimeter
	fmt.Printf("  ▭ Rectangle (%.2f × %.2f): Perimeter = %.2f\n", r.Width, r.Height, perimeter)
}

func (p *PerimeterCalculator) VisitTriangle(t *Triangle) {
	// Assuming equilateral triangle for simplicity
	// In real scenario, you'd need all three sides
	side := t.Base
	perimeter := 3 * side
	p.TotalPerimeter += perimeter
	fmt.Printf("  △ Triangle (side: %.2f): Perimeter ≈ %.2f\n", side, perimeter)
}

// SVGDrawer generates SVG code for shapes
type SVGDrawer struct {
	svgElements []string
}

func (s *SVGDrawer) VisitCircle(c *Circle) {
	svg := fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="blue" />`, c.X, c.Y, c.Radius)
	s.svgElements = append(s.svgElements, svg)
	fmt.Printf("  ⭕ Circle at (%.2f, %.2f) with radius %.2f\n", c.X, c.Y, c.Radius)
}

func (s *SVGDrawer) VisitRectangle(r *Rectangle) {
	svg := fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="green" />`, r.X, r.Y, r.Width, r.Height)
	s.svgElements = append(s.svgElements, svg)
	fmt.Printf("  ▭ Rectangle at (%.2f, %.2f) with size %.2f × %.2f\n", r.X, r.Y, r.Width, r.Height)
}

func (s *SVGDrawer) VisitTriangle(t *Triangle) {
	x1, y1 := t.X, t.Y
	x2, y2 := t.X+t.Base, t.Y
	x3, y3 := t.X+t.Base/2, t.Y-t.Height
	svg := fmt.Sprintf(`<polygon points="%.2f,%.2f %.2f,%.2f %.2f,%.2f" fill="red" />`, x1, y1, x2, y2, x3, y3)
	s.svgElements = append(s.svgElements, svg)
	fmt.Printf("  △ Triangle at (%.2f, %.2f) with base %.2f and height %.2f\n", t.X, t.Y, t.Base, t.Height)
}

func (s *SVGDrawer) GetSVG() string {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 500 500">` + "\n"
	for _, element := range s.svgElements {
		svg += "  " + element + "\n"
	}
	svg += "</svg>"
	return svg
}

// JSONExporter exports shape data as JSON
type JSONExporter struct {
	jsonData []string
}

func (j *JSONExporter) VisitCircle(c *Circle) {
	json := fmt.Sprintf(`{"type":"circle","radius":%.2f,"center":{"x":%.2f,"y":%.2f}}`, c.Radius, c.X, c.Y)
	j.jsonData = append(j.jsonData, json)
}

func (j *JSONExporter) VisitRectangle(r *Rectangle) {
	json := fmt.Sprintf(`{"type":"rectangle","width":%.2f,"height":%.2f,"position":{"x":%.2f,"y":%.2f}}`, r.Width, r.Height, r.X, r.Y)
	j.jsonData = append(j.jsonData, json)
}

func (j *JSONExporter) VisitTriangle(t *Triangle) {
	json := fmt.Sprintf(`{"type":"triangle","base":%.2f,"height":%.2f,"position":{"x":%.2f,"y":%.2f}}`, t.Base, t.Height, t.X, t.Y)
	j.jsonData = append(j.jsonData, json)
}

func (j *JSONExporter) GetJSON() string {
	json := "[\n"
	for i, data := range j.jsonData {
		json += "  " + data
		if i < len(j.jsonData)-1 {
			json += ","
		}
		json += "\n"
	}
	json += "]"
	return json
}

// ============================================================================
// DRAWING - Client Code
// ============================================================================

// Drawing holds a collection of shapes
type Drawing struct {
	Name   string
	shapes []Shape
}

func (d *Drawing) AddShape(shape Shape) {
	d.shapes = append(d.shapes, shape)
}

func (d *Drawing) ApplyVisitor(visitor ShapeVisitor) {
	for _, shape := range d.shapes {
		shape.Accept(visitor)
	}
}

// ============================================================================
// MAIN - Demonstration
// ============================================================================

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║      VISITOR PATTERN - GEOMETRIC SHAPES EXAMPLE           ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Create a drawing
	drawing := &Drawing{Name: "My Shapes"}

	// Add shapes to drawing
	drawing.AddShape(&Circle{
		Radius: 50,
		X:      100,
		Y:      100,
	})

	drawing.AddShape(&Rectangle{
		Width:  80,
		Height: 60,
		X:      200,
		Y:      50,
	})

	drawing.AddShape(&Triangle{
		Base:   70,
		Height: 90,
		X:      350,
		Y:      150,
	})

	drawing.AddShape(&Circle{
		Radius: 30,
		X:      250,
		Y:      300,
	})

	drawing.AddShape(&Rectangle{
		Width:  100,
		Height: 40,
		X:      50,
		Y:      250,
	})

	// Calculate areas
	fmt.Println("📐 AREA CALCULATION:")
	fmt.Println("─────────────────────────────────────────────────────────")
	areaCalc := &AreaCalculator{}
	drawing.ApplyVisitor(areaCalc)
	fmt.Printf("\n📊 Total Area: %.2f square units\n", areaCalc.TotalArea)
	fmt.Println()

	// Calculate perimeters
	fmt.Println("📏 PERIMETER CALCULATION:")
	fmt.Println("─────────────────────────────────────────────────────────")
	perimeterCalc := &PerimeterCalculator{}
	drawing.ApplyVisitor(perimeterCalc)
	fmt.Printf("\n📊 Total Perimeter: %.2f units\n", perimeterCalc.TotalPerimeter)
	fmt.Println()

	// Generate SVG
	fmt.Println("🎨 SVG GENERATION:")
	fmt.Println("─────────────────────────────────────────────────────────")
	svgDrawer := &SVGDrawer{}
	drawing.ApplyVisitor(svgDrawer)
	fmt.Println("\n📄 Generated SVG:")
	fmt.Println(svgDrawer.GetSVG())
	fmt.Println()

	// Export to JSON
	fmt.Println("📋 JSON EXPORT:")
	fmt.Println("─────────────────────────────────────────────────────────")
	jsonExporter := &JSONExporter{}
	drawing.ApplyVisitor(jsonExporter)
	fmt.Println(jsonExporter.GetJSON())
	fmt.Println()

	fmt.Println("✨ Key Takeaway:")
	fmt.Println("   We performed 4 different operations (Area, Perimeter, SVG, JSON)")
	fmt.Println("   on 3 shape types without modifying the shape classes!")
	fmt.Println("   Adding a new operation is as simple as creating a new visitor. 🚀")
}
