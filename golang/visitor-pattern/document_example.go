package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// VISITOR PATTERN - DOCUMENT PROCESSING EXAMPLE
// ============================================================================
// This example shows how to use the Visitor pattern to export a document
// structure to different formats (HTML, Markdown, Plain Text).
// ============================================================================

// DocumentVisitor defines the visitor interface for document elements
type DocumentVisitor interface {
	VisitParagraph(p *Paragraph)
	VisitHeading(h *Heading)
	VisitImage(i *Image)
	VisitTable(t *Table)
	VisitCodeBlock(c *CodeBlock)
}

// DocumentElement is the element interface
type DocumentElement interface {
	Accept(v DocumentVisitor)
}

// ============================================================================
// CONCRETE ELEMENTS - Different Document Parts
// ============================================================================

// Paragraph represents a text paragraph
type Paragraph struct {
	Text string
}

func (p *Paragraph) Accept(v DocumentVisitor) {
	v.VisitParagraph(p)
}

// Heading represents a section heading
type Heading struct {
	Text  string
	Level int // 1-6 for H1-H6
}

func (h *Heading) Accept(v DocumentVisitor) {
	v.VisitHeading(h)
}

// Image represents an embedded image
type Image struct {
	URL     string
	AltText string
	Caption string
}

func (i *Image) Accept(v DocumentVisitor) {
	v.VisitImage(i)
}

// Table represents a data table
type Table struct {
	Headers []string
	Rows    [][]string
}

func (t *Table) Accept(v DocumentVisitor) {
	v.VisitTable(t)
}

// CodeBlock represents a code snippet
type CodeBlock struct {
	Language string
	Code     string
}

func (c *CodeBlock) Accept(v DocumentVisitor) {
	v.VisitCodeBlock(c)
}

// ============================================================================
// CONCRETE VISITORS - Different Export Formats
// ============================================================================

// HTMLExporter exports document to HTML
type HTMLExporter struct {
	output strings.Builder
}

func (h *HTMLExporter) VisitParagraph(p *Paragraph) {
	h.output.WriteString(fmt.Sprintf("<p>%s</p>\n", p.Text))
}

func (h *HTMLExporter) VisitHeading(hd *Heading) {
	h.output.WriteString(fmt.Sprintf("<h%d>%s</h%d>\n", hd.Level, hd.Text, hd.Level))
}

func (h *HTMLExporter) VisitImage(i *Image) {
	h.output.WriteString(fmt.Sprintf("<figure>\n"))
	h.output.WriteString(fmt.Sprintf("  <img src=\"%s\" alt=\"%s\">\n", i.URL, i.AltText))
	if i.Caption != "" {
		h.output.WriteString(fmt.Sprintf("  <figcaption>%s</figcaption>\n", i.Caption))
	}
	h.output.WriteString("</figure>\n")
}

func (h *HTMLExporter) VisitTable(t *Table) {
	h.output.WriteString("<table>\n")
	h.output.WriteString("  <thead>\n    <tr>\n")
	for _, header := range t.Headers {
		h.output.WriteString(fmt.Sprintf("      <th>%s</th>\n", header))
	}
	h.output.WriteString("    </tr>\n  </thead>\n")
	h.output.WriteString("  <tbody>\n")
	for _, row := range t.Rows {
		h.output.WriteString("    <tr>\n")
		for _, cell := range row {
			h.output.WriteString(fmt.Sprintf("      <td>%s</td>\n", cell))
		}
		h.output.WriteString("    </tr>\n")
	}
	h.output.WriteString("  </tbody>\n</table>\n")
}

func (h *HTMLExporter) VisitCodeBlock(c *CodeBlock) {
	h.output.WriteString(fmt.Sprintf("<pre><code class=\"language-%s\">\n%s\n</code></pre>\n", c.Language, c.Code))
}

func (h *HTMLExporter) GetOutput() string {
	return h.output.String()
}

// MarkdownExporter exports document to Markdown
type MarkdownExporter struct {
	output strings.Builder
}

func (m *MarkdownExporter) VisitParagraph(p *Paragraph) {
	m.output.WriteString(fmt.Sprintf("%s\n\n", p.Text))
}

func (m *MarkdownExporter) VisitHeading(h *Heading) {
	prefix := strings.Repeat("#", h.Level)
	m.output.WriteString(fmt.Sprintf("%s %s\n\n", prefix, h.Text))
}

func (m *MarkdownExporter) VisitImage(i *Image) {
	m.output.WriteString(fmt.Sprintf("![%s](%s)\n", i.AltText, i.URL))
	if i.Caption != "" {
		m.output.WriteString(fmt.Sprintf("*%s*\n", i.Caption))
	}
	m.output.WriteString("\n")
}

func (m *MarkdownExporter) VisitTable(t *Table) {
	// Headers
	m.output.WriteString("| ")
	for _, header := range t.Headers {
		m.output.WriteString(fmt.Sprintf("%s | ", header))
	}
	m.output.WriteString("\n")

	// Separator
	m.output.WriteString("|")
	for range t.Headers {
		m.output.WriteString("---|")
	}
	m.output.WriteString("\n")

	// Rows
	for _, row := range t.Rows {
		m.output.WriteString("| ")
		for _, cell := range row {
			m.output.WriteString(fmt.Sprintf("%s | ", cell))
		}
		m.output.WriteString("\n")
	}
	m.output.WriteString("\n")
}

func (m *MarkdownExporter) VisitCodeBlock(c *CodeBlock) {
	m.output.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", c.Language, c.Code))
}

func (m *MarkdownExporter) GetOutput() string {
	return m.output.String()
}

// PlainTextExporter exports document to plain text
type PlainTextExporter struct {
	output strings.Builder
}

func (p *PlainTextExporter) VisitParagraph(par *Paragraph) {
	p.output.WriteString(fmt.Sprintf("%s\n\n", par.Text))
}

func (p *PlainTextExporter) VisitHeading(h *Heading) {
	p.output.WriteString(fmt.Sprintf("%s\n", strings.ToUpper(h.Text)))
	p.output.WriteString(strings.Repeat("=", len(h.Text)))
	p.output.WriteString("\n\n")
}

func (p *PlainTextExporter) VisitImage(i *Image) {
	p.output.WriteString(fmt.Sprintf("[IMAGE: %s - %s]\n", i.AltText, i.URL))
	if i.Caption != "" {
		p.output.WriteString(fmt.Sprintf("Caption: %s\n", i.Caption))
	}
	p.output.WriteString("\n")
}

func (p *PlainTextExporter) VisitTable(t *Table) {
	// Calculate column widths
	colWidths := make([]int, len(t.Headers))
	for i, header := range t.Headers {
		colWidths[i] = len(header)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Print headers
	for i, header := range t.Headers {
		p.output.WriteString(fmt.Sprintf("%-*s  ", colWidths[i], header))
	}
	p.output.WriteString("\n")

	// Print separator
	for _, width := range colWidths {
		p.output.WriteString(strings.Repeat("-", width) + "  ")
	}
	p.output.WriteString("\n")

	// Print rows
	for _, row := range t.Rows {
		for i, cell := range row {
			p.output.WriteString(fmt.Sprintf("%-*s  ", colWidths[i], cell))
		}
		p.output.WriteString("\n")
	}
	p.output.WriteString("\n")
}

func (p *PlainTextExporter) VisitCodeBlock(c *CodeBlock) {
	p.output.WriteString(fmt.Sprintf("Code (%s):\n", c.Language))
	p.output.WriteString("----------------------------------------\n")
	p.output.WriteString(c.Code)
	p.output.WriteString("\n----------------------------------------\n\n")
}

func (p *PlainTextExporter) GetOutput() string {
	return p.output.String()
}

// ============================================================================
// DOCUMENT - Client Code
// ============================================================================

// Document holds a collection of document elements
type Document struct {
	Title    string
	elements []DocumentElement
}

func (d *Document) AddElement(element DocumentElement) {
	d.elements = append(d.elements, element)
}

func (d *Document) Export(visitor DocumentVisitor) {
	for _, element := range d.elements {
		element.Accept(visitor)
	}
}

// ============================================================================
// MAIN - Demonstration
// ============================================================================

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║      VISITOR PATTERN - DOCUMENT EXPORT EXAMPLE           ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Create a document
	doc := &Document{Title: "Visitor Pattern Guide"}

	// Add elements to document
	doc.AddElement(&Heading{
		Text:  "Introduction to Visitor Pattern",
		Level: 1,
	})

	doc.AddElement(&Paragraph{
		Text: "The Visitor pattern is a behavioral design pattern that lets you separate algorithms from the objects on which they operate. It's particularly useful when you need to perform various operations across a set of objects with different types.",
	})

	doc.AddElement(&Heading{
		Text:  "Key Benefits",
		Level: 2,
	})

	doc.AddElement(&Table{
		Headers: []string{"Benefit", "Description"},
		Rows: [][]string{
			{"Open/Closed", "Add new operations without modifying classes"},
			{"Single Responsibility", "Separate algorithms from objects"},
			{"Type Safety", "Compile-time checking"},
		},
	})

	doc.AddElement(&Heading{
		Text:  "Example Code",
		Level: 2,
	})

	doc.AddElement(&CodeBlock{
		Language: "go",
		Code: `type Visitor interface {
    VisitElementA(a *ElementA)
    VisitElementB(b *ElementB)
}

type Element interface {
    Accept(v Visitor)
}`,
	})

	doc.AddElement(&Image{
		URL:     "https://example.com/visitor-pattern.png",
		AltText: "Visitor Pattern Diagram",
		Caption: "Structure of the Visitor Pattern",
	})

	// Export to HTML
	fmt.Println("📄 HTML OUTPUT:")
	fmt.Println("═══════════════════════════════════════════════════════════")
	htmlExporter := &HTMLExporter{}
	doc.Export(htmlExporter)
	fmt.Println(htmlExporter.GetOutput())

	// Export to Markdown
	fmt.Println("📝 MARKDOWN OUTPUT:")
	fmt.Println("═══════════════════════════════════════════════════════════")
	mdExporter := &MarkdownExporter{}
	doc.Export(mdExporter)
	fmt.Println(mdExporter.GetOutput())

	// Export to Plain Text
	fmt.Println("📃 PLAIN TEXT OUTPUT:")
	fmt.Println("═══════════════════════════════════════════════════════════")
	txtExporter := &PlainTextExporter{}
	doc.Export(txtExporter)
	fmt.Println(txtExporter.GetOutput())

	fmt.Println("✨ Key Takeaway:")
	fmt.Println("   We exported the same document to 3 different formats")
	fmt.Println("   without modifying any of the document element classes!")
	fmt.Println("   Each exporter (visitor) encapsulates a different export algorithm. 🚀")
}
