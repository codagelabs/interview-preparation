// DrawVisitor - "draws" shapes (prints ASCII representation)
public class DrawVisitor implements ShapeVisitor {

    @Override
    public void visitCircle(Circle circle) {
        System.out.println("  Drawing Circle:");
        System.out.println("     ___");
        System.out.println("    /   \\");
        System.out.println("   |  O  |");
        System.out.println("    \\___/");
    }

    @Override
    public void visitRectangle(Rectangle rectangle) {
        System.out.println("  Drawing Rectangle:");
        System.out.println("   +-----+");
        System.out.println("   |     |");
        System.out.println("   |     |");
        System.out.println("   +-----+");
    }

    @Override
    public void visitTriangle(Triangle triangle) {
        System.out.println("  Drawing Triangle:");
        System.out.println("      /\\");
        System.out.println("     /  \\");
        System.out.println("    /____\\");
    }
}

