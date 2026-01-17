// Shape interface - all shapes must accept a visitor
public interface Shape {
    void accept(ShapeVisitor visitor);
}

