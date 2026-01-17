import java.util.ArrayList;
import java.util.List;

/**
 * ============================================================================
 * VISITOR PATTERN - SHAPES EXAMPLE (Java)
 * ============================================================================
 * This example demonstrates the Visitor pattern with geometric shapes.
 * We can perform different operations (area, perimeter, drawing) without
 * modifying the shape classes.
 * ============================================================================
 */
public class VisitorPatternDemo {

    public static void main(String[] args) {
        System.out.println("===========================================");
        System.out.println("  VISITOR PATTERN - Simple Shape Example");
        System.out.println("===========================================");
        System.out.println();

        // Create different shapes
        List<Shape> shapes = new ArrayList<>();
        shapes.add(new Circle(5.0));
        shapes.add(new Rectangle(4.0, 6.0));
        shapes.add(new Triangle(6.0, 8.0));

        // Operation 1: Calculate Areas
        System.out.println("CALCULATING AREAS:");
        System.out.println("-------------------------------------------");
        AreaCalculator areaCalc = new AreaCalculator();
        for (Shape shape : shapes) {
            shape.accept(areaCalc); // Each shape accepts the visitor
        }
        System.out.printf("%nTotal Area: %.2f%n", areaCalc.getTotalArea());
        System.out.println();

        // Operation 2: Calculate Perimeters
        System.out.println("CALCULATING PERIMETERS:");
        System.out.println("-------------------------------------------");
        PerimeterCalculator perimeterCalc = new PerimeterCalculator();
        for (Shape shape : shapes) {
            shape.accept(perimeterCalc);
        }
        System.out.printf("%nTotal Perimeter: %.2f%n", perimeterCalc.getTotalPerimeter());
        System.out.println();

        // Operation 3: Draw Shapes
        System.out.println("DRAWING SHAPES:");
        System.out.println("-------------------------------------------");
        DrawVisitor drawer = new DrawVisitor();
        for (Shape shape : shapes) {
            shape.accept(drawer);
            System.out.println();
        }

        // Key Benefit
        System.out.println("===========================================");
        System.out.println("Key Benefit:");
        System.out.println("   We performed 3 operations on 3 shape types");
        System.out.println("   WITHOUT modifying the shape classes!");
        System.out.println("   Adding new operations is easy - just create");
        System.out.println("   a new Visitor!");
        System.out.println("===========================================");
    }
}

