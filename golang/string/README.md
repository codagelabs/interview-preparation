# Go Strings Package

> **Source**: [Go Strings Package Documentation](https://pkg.go.dev/strings)

Package `strings` implements simple functions to manipulate UTF-8 encoded strings. This package provides a comprehensive set of utilities for string manipulation, searching, and transformation.

## Table of Contents
1. [Overview](#overview)
2. [String Functions](#string-functions)
3. [Types](#types)
4. [Usage Examples](#usage-examples)
5. [Best Practices](#best-practices)

## Overview

The `strings` package is part of Go's standard library and provides efficient string manipulation functions. All functions in this package work with UTF-8 encoded strings, making them suitable for international text processing.

**Key Features:**
- UTF-8 aware string operations
- Efficient memory usage
- Thread-safe operations
- Comprehensive function coverage

## String Functions

### String Comparison Functions

#### `Compare(a, b string) int`
Compares two strings lexicographically.
- Returns -1 if `a < b`
- Returns 0 if `a == b`
- Returns +1 if `a > b`

```go
import "strings"

result := strings.Compare("apple", "banana") // Returns -1
result = strings.Compare("banana", "banana") // Returns 0
result = strings.Compare("zebra", "apple")  // Returns +1
```

#### `EqualFold(s, t string) bool`
Case-insensitive string comparison.

```go
result := strings.EqualFold("Hello", "hello") // Returns true
result = strings.EqualFold("Go", "golang")   // Returns false
```

### String Search Functions

#### `Contains(s, substr string) bool`
Checks if a string contains a substring.

```go
result := strings.Contains("Hello, World!", "World") // Returns true
result = strings.Contains("Hello, World!", "Python") // Returns false
```

#### `ContainsAny(s, chars string) bool`
Checks if a string contains any character from the given set.

```go
result := strings.ContainsAny("Hello", "aeiou") // Returns true (contains 'e', 'o')
result = strings.ContainsAny("Hello", "xyz")   // Returns false
```

#### `ContainsRune(s string, r rune) bool`
Checks if a string contains a specific Unicode rune.

```go
result := strings.ContainsRune("Hello", 'l') // Returns true
result = strings.ContainsRune("Hello", 'x') // Returns false
```

#### `ContainsFunc(s string, f func(rune) bool) bool`
Checks if a string contains any rune that satisfies the given function.

```go
result := strings.ContainsFunc("Hello123", unicode.IsDigit) // Returns true
result = strings.ContainsFunc("Hello", unicode.IsUpper)    // Returns true
```

#### `Index(s, substr string) int`
Returns the index of the first occurrence of a substring.

```go
index := strings.Index("Hello, World!", "World") // Returns 7
index = strings.Index("Hello, World!", "Python") // Returns -1
```

#### `LastIndex(s, substr string) int`
Returns the index of the last occurrence of a substring.

```go
index := strings.LastIndex("Hello, Hello!", "Hello") // Returns 7
```

#### `IndexAny(s, chars string) int`
Returns the index of the first occurrence of any character from the given set.

```go
index := strings.IndexAny("Hello", "aeiou") // Returns 1 (first vowel 'e')
```

#### `IndexByte(s string, c byte) int`
Returns the index of the first occurrence of a byte.

```go
index := strings.IndexByte("Hello", 'l') // Returns 2
```

#### `IndexRune(s string, r rune) int`
Returns the index of the first occurrence of a Unicode rune.

```go
index := strings.IndexRune("Hello", 'l') // Returns 2
```

#### `IndexFunc(s string, f func(rune) bool) int`
Returns the index of the first rune that satisfies the given function.

```go
index := strings.IndexFunc("Hello123", unicode.IsDigit) // Returns 5
```

### String Manipulation Functions

#### `Clone(s string) string`
Creates a copy of the string. Useful for ensuring string immutability.

```go
original := "Hello"
cloned := strings.Clone(original)
// cloned is a separate copy of original
```

#### `Repeat(s string, count int) string`
Repeats a string the specified number of times.

```go
result := strings.Repeat("Go ", 3) // Returns "Go Go Go "
```

#### `Replace(s, old, new string, n int) string`
Replaces occurrences of `old` with `new` in string `s`. If `n < 0`, all occurrences are replaced.

```go
result := strings.Replace("Hello Hello Hello", "Hello", "Hi", 2)
// Returns "Hi Hi Hello"
```

#### `ReplaceAll(s, old, new string) string`
Replaces all occurrences of `old` with `new` in string `s`.

```go
result := strings.ReplaceAll("Hello Hello Hello", "Hello", "Hi")
// Returns "Hi Hi Hi"
```

#### `Map(mapping func(rune) rune, s string) string`
Applies a mapping function to each rune in the string.

```go
result := strings.Map(func(r rune) rune {
    if r >= 'a' && r <= 'z' {
        return r - 32 // Convert to uppercase
    }
    return r
}, "hello world")
// Returns "HELLO WORLD"
```

### String Splitting Functions

#### `Split(s, sep string) []string`
Splits a string by the given separator.

```go
parts := strings.Split("apple,banana,cherry", ",")
// Returns ["apple", "banana", "cherry"]
```

#### `SplitN(s, sep string, n int) []string`
Splits a string by the given separator, limiting the number of parts.

```go
parts := strings.SplitN("apple,banana,cherry", ",", 2)
// Returns ["apple", "banana,cherry"]
```

#### `SplitAfter(s, sep string) []string`
Splits a string after the given separator (keeps the separator).

```go
parts := strings.SplitAfter("apple,banana,cherry", ",")
// Returns ["apple,", "banana,", "cherry"]
```

#### `SplitAfterN(s, sep string, n int) []string`
Splits a string after the given separator, limiting the number of parts.

```go
parts := strings.SplitAfterN("apple,banana,cherry", ",", 2)
// Returns ["apple,", "banana,cherry"]
```

#### `Fields(s string) []string`
Splits a string by whitespace.

```go
parts := strings.Fields("  hello   world  ")
// Returns ["hello", "world"]
```

#### `FieldsFunc(s string, f func(rune) bool) []string`
Splits a string using a custom function to determine separators.

```go
parts := strings.FieldsFunc("hello,world;golang", func(r rune) bool {
    return r == ',' || r == ';'
})
// Returns ["hello", "world", "golang"]
```

### String Joining Functions

#### `Join(elems []string, sep string) string`
Joins a slice of strings with the given separator.

```go
result := strings.Join([]string{"apple", "banana", "cherry"}, ", ")
// Returns "apple, banana, cherry"
```

### String Trimming Functions

#### `Trim(s, cutset string) string`
Removes characters from both ends of the string.

```go
result := strings.Trim("!!!Hello!!!", "!")
// Returns "Hello"
```

#### `TrimLeft(s, cutset string) string`
Removes characters from the left end of the string.

```go
result := strings.TrimLeft("!!!Hello!!!", "!")
// Returns "Hello!!!"
```

#### `TrimRight(s, cutset string) string`
Removes characters from the right end of the string.

```go
result := strings.TrimRight("!!!Hello!!!", "!")
// Returns "!!!Hello"
```

#### `TrimSpace(s string) string`
Removes leading and trailing whitespace.

```go
result := strings.TrimSpace("  \t\n Hello, World! \n\t\r\n  ")
// Returns "Hello, World!"
```

#### `TrimPrefix(s, prefix string) string`
Removes the specified prefix if present.

```go
result := strings.TrimPrefix("Hello, World!", "Hello, ")
// Returns "World!"
```

#### `TrimSuffix(s, suffix string) string`
Removes the specified suffix if present.

```go
result := strings.TrimSuffix("Hello, World!", ", World!")
// Returns "Hello"
```

#### `TrimFunc(s string, f func(rune) bool) string`
Removes characters from both ends based on a function.

```go
result := strings.TrimFunc("!!!Hello!!!", func(r rune) bool {
    return r == '!'
})
// Returns "Hello"
```

#### `TrimLeftFunc(s string, f func(rune) bool) string`
Removes characters from the left end based on a function.

```go
result := strings.TrimLeftFunc("!!!Hello!!!", func(r rune) bool {
    return r == '!'
})
// Returns "Hello!!!"
```

#### `TrimRightFunc(s string, f func(rune) bool) string`
Removes characters from the right end based on a function.

```go
result := strings.TrimRightFunc("!!!Hello!!!", func(r rune) bool {
    return r == '!'
})
// Returns "!!!Hello"
```

### String Case Functions

#### `ToUpper(s string) string`
Converts string to uppercase.

```go
result := strings.ToUpper("hello world")
// Returns "HELLO WORLD"
```

#### `ToLower(s string) string`
Converts string to lowercase.

```go
result := strings.ToLower("HELLO WORLD")
// Returns "hello world"
```

#### `ToTitle(s string) string`
Converts string to title case.

```go
result := strings.ToTitle("hello world")
// Returns "HELLO WORLD"
```

#### `Title(s string) string`
Converts string to title case (first letter of each word capitalized).

```go
result := strings.Title("hello world")
// Returns "Hello World"
```

### String Cutting Functions

#### `Cut(s, sep string) (before, after string, found bool)`
Cuts a string at the first occurrence of the separator.

```go
before, after, found := strings.Cut("hello,world", ",")
// Returns "hello", "world", true
```

#### `CutPrefix(s, prefix string) (after string, found bool)`
Cuts the prefix from the string if present.

```go
after, found := strings.CutPrefix("hello,world", "hello,")
// Returns "world", true
```

#### `CutSuffix(s, suffix string) (before string, found bool)`
Cuts the suffix from the string if present.

```go
before, found := strings.CutSuffix("hello,world", ",world")
// Returns "hello", true
```

### String Counting Functions

#### `Count(s, substr string) int`
Counts non-overlapping occurrences of a substring.

```go
count := strings.Count("hello hello hello", "hello")
// Returns 3
```

### String Validation Functions

#### `ToValidUTF8(s, replacement string) string`
Returns a copy of the string with invalid UTF-8 sequences replaced.

```go
result := strings.ToValidUTF8("Hello\xffWorld", "?")
// Returns "Hello?World"
```

### String Sequence Functions (Go 1.23+)

#### `Lines(s string) iter.Seq[string]`
Returns an iterator over the lines in the string.

```go
for line := range strings.Lines("line1\nline2\nline3") {
    fmt.Println(line)
}
```

#### `FieldsSeq(s string) iter.Seq[string]`
Returns an iterator over the fields in the string.

```go
for field := range strings.FieldsSeq("  hello   world  ") {
    fmt.Println(field)
}
```

#### `SplitSeq(s, sep string) iter.Seq[string]`
Returns an iterator over the parts when splitting by separator.

```go
for part := range strings.SplitSeq("a,b,c", ",") {
    fmt.Println(part)
}
```

## Types

### Builder

`Builder` is used to efficiently build strings using write methods. It minimizes memory copying.

```go
var b strings.Builder
b.WriteString("Hello")
b.WriteString(" ")
b.WriteString("World")
result := b.String() // "Hello World"
```

**Builder Methods:**
- `Write(p []byte) (int, error)` - Writes bytes
- `WriteString(s string) (int, error)` - Writes string
- `WriteByte(c byte) error` - Writes single byte
- `WriteRune(r rune) (int, error)` - Writes rune
- `Grow(n int)` - Pre-allocates space
- `Reset()` - Clears the builder
- `Cap() int` - Returns capacity
- `Len() int` - Returns current length

### Reader

`Reader` implements various I/O interfaces for reading from a string.

```go
reader := strings.NewReader("Hello, World!")
data := make([]byte, 5)
n, err := reader.Read(data)
// data contains "Hello"
```

**Reader Methods:**
- `Read(b []byte) (n int, err error)` - Implements io.Reader
- `ReadAt(b []byte, off int64) (n int, err error)` - Implements io.ReaderAt
- `ReadByte() (byte, error)` - Implements io.ByteReader
- `ReadRune() (ch rune, size int, err error)` - Implements io.RuneReader
- `Seek(offset int64, whence int) (int64, error)` - Implements io.Seeker
- `Reset(s string)` - Resets to read from new string
- `Len() int` - Returns unread bytes
- `Size() int64` - Returns total string size

### Replacer

`Replacer` replaces multiple strings with replacements. It's safe for concurrent use.

```go
r := strings.NewReplacer("<", "&lt;", ">", "&gt;")
result := r.Replace("This is <b>HTML</b>!")
// Returns "This is &lt;b&gt;HTML&lt;/b&gt;!"
```

**Replacer Methods:**
- `Replace(s string) string` - Returns string with replacements
- `WriteString(w io.Writer, s string) (n int, err error)` - Writes to writer with replacements

## Usage Examples

### Basic String Operations

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    // String concatenation with Builder
    var b strings.Builder
    b.WriteString("Hello")
    b.WriteString(" ")
    b.WriteString("World")
    fmt.Println(b.String()) // "Hello World"

    // String splitting and joining
    text := "apple,banana,cherry"
    parts := strings.Split(text, ",")
    result := strings.Join(parts, " | ")
    fmt.Println(result) // "apple | banana | cherry"

    // String trimming
    dirty := "  \t\n  Hello, World!  \n\t\r\n  "
    clean := strings.TrimSpace(dirty)
    fmt.Printf("'%s'\n", clean) // 'Hello, World!'

    // String replacement
    html := "<b>Hello</b> <i>World</i>"
    r := strings.NewReplacer("<b>", "<strong>", "</b>", "</strong>",
                             "<i>", "<em>", "</i>", "</em>")
    result = r.Replace(html)
    fmt.Println(result) // "<strong>Hello</strong> <em>World</em>"
}
```

### Advanced String Processing

```go
package main

import (
    "fmt"
    "strings"
    "unicode"
)

func main() {
    // Custom string processing
    text := "Hello123World456"
    
    // Extract only letters
    letters := strings.Map(func(r rune) rune {
        if unicode.IsLetter(r) {
            return r
        }
        return -1 // Remove non-letters
    }, text)
    fmt.Println(letters) // "HelloWorld"

    // Extract only digits
    digits := strings.Map(func(r rune) rune {
        if unicode.IsDigit(r) {
            return r
        }
        return -1 // Remove non-digits
    }, text)
    fmt.Println(digits) // "123456"

    // Case-insensitive search
    if strings.Contains(strings.ToLower(text), "hello") {
        fmt.Println("Found 'hello' (case-insensitive)")
    }

    // Multiple replacements
    r := strings.NewReplacer(
        "Hello", "Hi",
        "World", "Universe",
        "123", "ABC",
        "456", "XYZ",
    )
    result := r.Replace(text)
    fmt.Println(result) // "HiABCUniverseXYZ"
}
```

### String Validation and Cleaning

```go
package main

import (
    "fmt"
    "strings"
    "unicode"
)

func main() {
    // Clean user input
    userInput := "  \t\n  John Doe  \n\t\r\n  "
    cleanName := strings.TrimSpace(userInput)
    fmt.Printf("Clean name: '%s'\n", cleanName)

    // Validate email format (basic)
    email := "user@example.com"
    if strings.Contains(email, "@") && strings.Contains(email, ".") {
        fmt.Println("Email format looks valid")
    }

    // Remove special characters
    dirtyText := "Hello!!! World??? How are you???"
    cleanText := strings.TrimFunc(dirtyText, func(r rune) bool {
        return r == '!' || r == '?'
    })
    fmt.Println(cleanText) // "Hello!!! World??? How are you"

    // Count words
    text := "The quick brown fox jumps over the lazy dog"
    words := strings.Fields(text)
    fmt.Printf("Word count: %d\n", len(words)) // 9
}
```

### Performance Optimization Examples

```go
package main

import (
    "fmt"
    "strings"
    "time"
)

func main() {
    // Efficient string building
    start := time.Now()
    
    var b strings.Builder
    b.Grow(1000) // Pre-allocate space
    
    for i := 0; i < 1000; i++ {
        b.WriteString("word")
        if i < 999 {
            b.WriteString(" ")
        }
    }
    
    result := b.String()
    elapsed := time.Since(start)
    fmt.Printf("Built string in %v\n", elapsed)
    fmt.Printf("Length: %d\n", len(result))

    // Efficient string reading
    reader := strings.NewReader("Large text content here...")
    
    // Read in chunks
    buffer := make([]byte, 10)
    for {
        n, err := reader.Read(buffer)
        if n > 0 {
            fmt.Printf("Read %d bytes: %s\n", n, buffer[:n])
        }
        if err != nil {
            break
        }
    }
}
```

## Best Practices

### 1. Use Builder for String Concatenation

```go
// Good: Use Builder for multiple concatenations
var b strings.Builder
for i := 0; i < 1000; i++ {
    b.WriteString("word")
    b.WriteString(" ")
}
result := b.String()

// Bad: Multiple string concatenations
result := ""
for i := 0; i < 1000; i++ {
    result += "word" + " " // Creates new strings each time
}
```

### 2. Pre-allocate Builder Capacity

```go
var b strings.Builder
b.Grow(estimatedSize) // Pre-allocate to avoid reallocations
```

### 3. Use Appropriate Functions

```go
// For simple contains check
if strings.Contains(text, "search") { ... }

// For case-insensitive check
if strings.EqualFold(text, "search") { ... }

// For custom logic
if strings.ContainsFunc(text, unicode.IsDigit) { ... }
```

### 4. Handle UTF-8 Properly

```go
// Use rune-aware functions for international text
index := strings.IndexRune(text, 'ñ')
count := strings.Count(text, "café")
```

### 5. Use Replacer for Multiple Replacements

```go
// Good: Use Replacer for multiple replacements
r := strings.NewReplacer("old1", "new1", "old2", "new2")
result := r.Replace(text)

// Bad: Multiple Replace calls
result := strings.Replace(text, "old1", "new1", -1)
result = strings.Replace(result, "old2", "new2", -1)
```

## Performance Considerations

- **Builder**: Most efficient for building strings incrementally
- **Join**: More efficient than manual concatenation for slices
- **Trim functions**: More efficient than manual character removal
- **Replacer**: More efficient than multiple Replace calls
- **Reader**: Efficient for reading large strings in chunks

## Common Use Cases

1. **Text Processing**: Cleaning, formatting, and transforming text
2. **Data Parsing**: Splitting CSV, parsing configuration files
3. **Input Validation**: Cleaning user input, checking formats
4. **String Building**: Constructing dynamic strings efficiently
5. **Text Search**: Finding patterns, counting occurrences
6. **Case Conversion**: Normalizing text case for comparisons

## References

- [Go Strings Package Documentation](https://pkg.go.dev/strings)
- [Go Blog: Strings, bytes, runes and characters in Go](https://blog.golang.org/strings)
- [Go by Example: Strings and Runes](https://gobyexample.com/strings-and-runes)
