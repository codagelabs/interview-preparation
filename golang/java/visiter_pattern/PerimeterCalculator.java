// PerimeterCalculator - calculates perimeter of shapes
public class PerimeterCalculator implements ShapeVisitor {
    private double totalPerimeter = 0.0;

    @Override
    public void visitCircle(Circle circle) {
        double perimeter = 2 * Math.PI * circle.getRadius();
        System.out.printf("  Circle perimeter: %.2f%n", perimeter);
        totalPerimeter += perimeter;
    }

    @Override
    public void visitRectangle(Rectangle rectangle) {
        double perimeter = 2 * (rectangle.getWidth() + rectangle.getHeight());
        System.out.printf("  Rectangle perimeter: %.2f%n", perimeter);
        totalPerimeter += perimeter;
    }

    @Override
    public void visitTriangle(Triangle triangle) {
        // Assuming equilateral triangle for simplicity
        double perimeter = 3 * triangle.getBase();
        System.out.printf("  Triangle perimeter: %.2f%n", perimeter);
        totalPerimeter += perimeter;
    }

    public double getTotalPerimeter() {
        return totalPerimeter;
    }
}

