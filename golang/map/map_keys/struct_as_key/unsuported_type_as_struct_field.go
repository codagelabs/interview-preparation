package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Approach 1: Using hashing to handle non-comparable fields
// ======================================================

type PersonWithUnsupportedFields struct {
	Name      string
	Hobbies   []string              // Slice - normally not comparable
	Addresses map[string]string     // Map - normally not comparable
	Callback  func(string) string   // Function - normally not comparable
}

// Hash generates a unique hash for the struct including non-comparable fields
func (p PersonWithUnsupportedFields) Hash() string {
	h := sha256.New()
	h.Write([]byte(p.Name))
	hobbiesJSON, _ := json.Marshal(p.Hobbies)
	h.Write(hobbiesJSON)
	addressJSON, _ := json.Marshal(p.Addresses)
	h.Write(addressJSON)
	funcPtr := reflect.ValueOf(p.Callback).Pointer()
	funcName := runtime.FuncForPC(funcPtr).Name()
	h.Write([]byte(funcName))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Approach 2: Transform non-comparable fields into comparable types
// =============================================================

type ComparablePerson struct {
	Name           string
	HobbiesString  string    // Transform slice to string
	AddressString  string    // Transform map to string
	CallbackString string    // Transform function to string
}

// Convert transforms PersonWithUnsupportedFields to ComparablePerson
func (p PersonWithUnsupportedFields) Convert() ComparablePerson {
	return ComparablePerson{
		Name:           p.Name,
		HobbiesString:  strings.Join(p.Hobbies, ","),
		AddressString:  mapToString(p.Addresses),
		CallbackString: getFunctionName(p.Callback),
	}
}

// Helper functions for conversion
func mapToString(m map[string]string) string {
	pairs := make([]string, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, k+"="+v)
	}
	return strings.Join(pairs, ";")
}

func getFunctionName(f func(string) string) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func main() {
	// Create test data
	person1 := PersonWithUnsupportedFields{
		Name:    "John",
		Hobbies: []string{"reading", "gaming"},
		Addresses: map[string]string{
			"home": "123 Main St",
			"work": "456 Office Blvd",
		},
		Callback: func(s string) string { return "Hello " + s },
	}

	person2 := PersonWithUnsupportedFields{
		Name:    "Alice",
		Hobbies: []string{"painting", "traveling"},
		Addresses: map[string]string{
			"home":     "789 Park Ave",
			"vacation": "321 Beach Rd",
		},
		Callback: func(s string) string { return "Hi " + s },
	}

	fmt.Println("Approach 1: Using Hashing")
	fmt.Println("========================")
	
	// Using hash approach
	hashedMap := make(map[string]PersonWithUnsupportedFields)
	hashedMap[person1.Hash()] = person1
	hashedMap[person2.Hash()] = person2

	// Demonstrate hash-based retrieval
	fmt.Printf("Person1 hash: %s\n", person1.Hash()[:8])
	retrievedPerson := hashedMap[person1.Hash()]
	fmt.Printf("Retrieved: %s with hobbies: %v\n\n", retrievedPerson.Name, retrievedPerson.Hobbies)

	fmt.Println("Approach 2: Using Transformation")
	fmt.Println("==============================")
	
	// Using transformation approach
	comparableMap := make(map[ComparablePerson]string)
	comparableMap[person1.Convert()] = "Active"
	comparableMap[person2.Convert()] = "Inactive"

	// Demonstrate transformation-based retrieval
	converted1 := person1.Convert()
	fmt.Printf("Original hobbies: %v\n", person1.Hobbies)
	fmt.Printf("Converted hobbies string: %s\n", converted1.HobbiesString)
	fmt.Printf("Status: %s\n", comparableMap[converted1])

	// Demonstrate that both approaches maintain uniqueness
	fmt.Println("\nUniqueness Demonstration")
	fmt.Println("=======================")
	fmt.Printf("Person1 and Person2 have different hashes: %v\n", 
		person1.Hash() != person2.Hash())
	fmt.Printf("Person1 and Person2 have different comparable forms: %v\n",
		person1.Convert() != person2.Convert())
}