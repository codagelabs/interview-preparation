package main

import "fmt"

// ============================================================================
// VISITOR PATTERN - E-COMMERCE EXAMPLE
// ============================================================================
// This example demonstrates the Visitor pattern with an e-commerce system
// where we need to perform different calculations (tax, shipping, discount)
// on different product types (Electronics, Clothing, Books).
// ============================================================================

// Visitor interface defines methods for visiting each concrete element
type Visitor interface {
	VisitElectronics(e *Electronics)
	VisitClothing(c *Clothing)
	VisitBook(b *Book)
}

// Element interface represents a product that can be visited
type Element interface {
	Accept(v Visitor)
	GetName() string
	GetPrice() float64
}

// ============================================================================
// CONCRETE ELEMENTS - Different Product Types
// ============================================================================

// Electronics represents electronic products
type Electronics struct {
	Name     string
	Price    float64
	Warranty int // warranty period in months
}

func (e *Electronics) Accept(v Visitor) {
	v.VisitElectronics(e)
}

func (e *Electronics) GetName() string {
	return e.Name
}

func (e *Electronics) GetPrice() float64 {
	return e.Price
}

// Clothing represents clothing products
type Clothing struct {
	Name     string
	Price    float64
	Size     string
	Material string
}

func (c *Clothing) Accept(v Visitor) {
	v.VisitClothing(c)
}

func (c *Clothing) GetName() string {
	return c.Name
}

func (c *Clothing) GetPrice() float64 {
	return c.Price
}

// Book represents book products
type Book struct {
	Name      string
	Price     float64
	Pages     int
	Author    string
	Hardcover bool
}

func (b *Book) Accept(v Visitor) {
	v.VisitBook(b)
}

func (b *Book) GetName() string {
	return b.Name
}

func (b *Book) GetPrice() float64 {
	return b.Price
}

// ============================================================================
// CONCRETE VISITORS - Different Operations
// ============================================================================

// TaxCalculator calculates tax for different product types
type TaxCalculator struct {
	TotalTax float64
}

func (tc *TaxCalculator) VisitElectronics(e *Electronics) {
	tax := e.Price * 0.15 // 15% tax on electronics
	tc.TotalTax += tax
	fmt.Printf("  🔌 %s: $%.2f (Tax: $%.2f @ 15%%)\n", e.Name, e.Price, tax)
}

func (tc *TaxCalculator) VisitClothing(c *Clothing) {
	tax := c.Price * 0.08 // 8% tax on clothing
	tc.TotalTax += tax
	fmt.Printf("  👕 %s: $%.2f (Tax: $%.2f @ 8%%)\n", c.Name, c.Price, tax)
}

func (tc *TaxCalculator) VisitBook(b *Book) {
	tax := b.Price * 0.05 // 5% tax on books
	tc.TotalTax += tax
	fmt.Printf("  📚 %s: $%.2f (Tax: $%.2f @ 5%%)\n", b.Name, b.Price, tax)
}

// ShippingCalculator calculates shipping costs for different product types
type ShippingCalculator struct {
	TotalShipping float64
}

func (sc *ShippingCalculator) VisitElectronics(e *Electronics) {
	shipping := 15.0 // Flat $15 for electronics (fragile)
	sc.TotalShipping += shipping
	fmt.Printf("  🔌 %s: $%.2f shipping (fragile item)\n", e.Name, shipping)
}

func (sc *ShippingCalculator) VisitClothing(c *Clothing) {
	shipping := 5.0 // Flat $5 for clothing (lightweight)
	sc.TotalShipping += shipping
	fmt.Printf("  👕 %s: $%.2f shipping (lightweight)\n", c.Name, shipping)
}

func (sc *ShippingCalculator) VisitBook(b *Book) {
	// Books: $3 base + $0.01 per page
	shipping := 3.0 + (float64(b.Pages) * 0.01)
	sc.TotalShipping += shipping
	fmt.Printf("  📚 %s: $%.2f shipping (%d pages)\n", b.Name, shipping, b.Pages)
}

// DiscountCalculator calculates available discounts
type DiscountCalculator struct {
	TotalDiscount float64
}

func (dc *DiscountCalculator) VisitElectronics(e *Electronics) {
	// 10% discount if warranty > 24 months
	var discount float64
	if e.Warranty > 24 {
		discount = e.Price * 0.10
		dc.TotalDiscount += discount
		fmt.Printf("  🔌 %s: -$%.2f discount (extended warranty)\n", e.Name, discount)
	} else {
		fmt.Printf("  🔌 %s: No discount available\n", e.Name)
	}
}

func (dc *DiscountCalculator) VisitClothing(c *Clothing) {
	// 15% discount on cotton material
	var discount float64
	if c.Material == "Cotton" {
		discount = c.Price * 0.15
		dc.TotalDiscount += discount
		fmt.Printf("  👕 %s: -$%.2f discount (cotton material)\n", c.Name, discount)
	} else {
		fmt.Printf("  👕 %s: No discount available\n", c.Name)
	}
}

func (dc *DiscountCalculator) VisitBook(b *Book) {
	// 20% discount on hardcover books
	var discount float64
	if b.Hardcover {
		discount = b.Price * 0.20
		dc.TotalDiscount += discount
		fmt.Printf("  📚 %s: -$%.2f discount (hardcover)\n", b.Name, discount)
	} else {
		fmt.Printf("  📚 %s: No discount available\n", b.Name)
	}
}

// InfoPrinter prints detailed information about products
type InfoPrinter struct{}

func (ip *InfoPrinter) VisitElectronics(e *Electronics) {
	fmt.Printf("  🔌 Electronics: %s\n", e.Name)
	fmt.Printf("     Price: $%.2f\n", e.Price)
	fmt.Printf("     Warranty: %d months\n", e.Warranty)
}

func (ip *InfoPrinter) VisitClothing(c *Clothing) {
	fmt.Printf("  👕 Clothing: %s\n", c.Name)
	fmt.Printf("     Price: $%.2f\n", c.Price)
	fmt.Printf("     Size: %s\n", c.Size)
	fmt.Printf("     Material: %s\n", c.Material)
}

func (ip *InfoPrinter) VisitBook(b *Book) {
	fmt.Printf("  📚 Book: %s\n", b.Name)
	fmt.Printf("     Price: $%.2f\n", b.Price)
	fmt.Printf("     Author: %s\n", b.Author)
	fmt.Printf("     Pages: %d\n", b.Pages)
	fmt.Printf("     Type: ")
	if b.Hardcover {
		fmt.Println("Hardcover")
	} else {
		fmt.Println("Paperback")
	}
}

// ============================================================================
// SHOPPING CART - Client Code
// ============================================================================

type ShoppingCart struct {
	items []Element
}

func (cart *ShoppingCart) AddItem(item Element) {
	cart.items = append(cart.items, item)
}

func (cart *ShoppingCart) ApplyVisitor(v Visitor) {
	for _, item := range cart.items {
		item.Accept(v)
	}
}

func (cart *ShoppingCart) GetTotalPrice() float64 {
	total := 0.0
	for _, item := range cart.items {
		total += item.GetPrice()
	}
	return total
}

// ============================================================================
// MAIN - Demonstration
// ============================================================================

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║       VISITOR PATTERN - E-COMMERCE EXAMPLE               ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Create shopping cart
	cart := &ShoppingCart{}

	// Add products to cart
	cart.AddItem(&Electronics{
		Name:     "Laptop",
		Price:    999.99,
		Warranty: 36,
	})

	cart.AddItem(&Electronics{
		Name:     "Smartphone",
		Price:    699.99,
		Warranty: 12,
	})

	cart.AddItem(&Clothing{
		Name:     "T-Shirt",
		Price:    29.99,
		Size:     "M",
		Material: "Cotton",
	})

	cart.AddItem(&Clothing{
		Name:     "Jeans",
		Price:    79.99,
		Size:     "32",
		Material: "Denim",
	})

	cart.AddItem(&Book{
		Name:      "Clean Code",
		Price:     45.99,
		Pages:     464,
		Author:    "Robert C. Martin",
		Hardcover: true,
	})

	cart.AddItem(&Book{
		Name:      "The Go Programming Language",
		Price:     39.99,
		Pages:     380,
		Author:    "Alan Donovan",
		Hardcover: false,
	})

	// Display products
	fmt.Println("📦 SHOPPING CART ITEMS:")
	fmt.Println("─────────────────────────────────────────────────────────")
	infoPrinter := &InfoPrinter{}
	cart.ApplyVisitor(infoPrinter)

	// Calculate subtotal
	subtotal := cart.GetTotalPrice()
	fmt.Println()
	fmt.Printf("💰 Subtotal: $%.2f\n", subtotal)
	fmt.Println()

	// Calculate tax
	fmt.Println("🧾 TAX CALCULATION:")
	fmt.Println("─────────────────────────────────────────────────────────")
	taxCalc := &TaxCalculator{}
	cart.ApplyVisitor(taxCalc)
	fmt.Printf("\n💳 Total Tax: $%.2f\n", taxCalc.TotalTax)
	fmt.Println()

	// Calculate shipping
	fmt.Println("📮 SHIPPING CALCULATION:")
	fmt.Println("─────────────────────────────────────────────────────────")
	shippingCalc := &ShippingCalculator{}
	cart.ApplyVisitor(shippingCalc)
	fmt.Printf("\n🚚 Total Shipping: $%.2f\n", shippingCalc.TotalShipping)
	fmt.Println()

	// Calculate discounts
	fmt.Println("🎁 DISCOUNT CALCULATION:")
	fmt.Println("─────────────────────────────────────────────────────────")
	discountCalc := &DiscountCalculator{}
	cart.ApplyVisitor(discountCalc)
	fmt.Printf("\n💝 Total Discount: $%.2f\n", discountCalc.TotalDiscount)
	fmt.Println()

	// Calculate final total
	finalTotal := subtotal + taxCalc.TotalTax + shippingCalc.TotalShipping - discountCalc.TotalDiscount
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("💰 Subtotal:        $%8.2f\n", subtotal)
	fmt.Printf("🧾 Tax:             $%8.2f\n", taxCalc.TotalTax)
	fmt.Printf("📮 Shipping:        $%8.2f\n", shippingCalc.TotalShipping)
	fmt.Printf("🎁 Discount:       -$%8.2f\n", discountCalc.TotalDiscount)
	fmt.Println("───────────────────────────────────────────────────────────")
	fmt.Printf("💳 FINAL TOTAL:     $%8.2f\n", finalTotal)
	fmt.Println("═══════════════════════════════════════════════════════════")

	fmt.Println()
	fmt.Println("✨ Key Takeaway:")
	fmt.Println("   We added 4 different operations (Info, Tax, Shipping, Discount)")
	fmt.Println("   without modifying the product classes (Electronics, Clothing, Book)!")
	fmt.Println("   This is the power of the Visitor Pattern! 🚀")
}
