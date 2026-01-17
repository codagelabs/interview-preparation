// AreaCalculator - calculates area of shapes
public class AreaCalculator implements ShapeVisitor {
    private double totalArea = 0.0;

    @Override
    public void visitCircle(Circle circle) {
        double area = Math.PI * circle.getRadius() * circle.getRadius();
        System.out.printf("  Circle area: %.2f%n", area);
        totalArea += area;
    }

    @Override
    public void visitRectangle(Rectangle rectangle) {
        double area = rectangle.getWidth() * rectangle.getHeight();
        System.out.printf("  Rectangle area: %.2f%n", area);
        totalArea += area;
    }

    @Override
    public void visitTriangle(Triangle triangle) {
        double area = 0.5 * triangle.getBase() * triangle.getHeight();
        System.out.printf("  Triangle area: %.2f%n", area);
        totalArea += area;
    }

    public double getTotalArea() {
        return totalArea;
    }
}

