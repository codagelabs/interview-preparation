
package main

import "fmt"

// Example of using interfaces as map keys in Go.
// This demonstrates how interfaces can be used as map keys when they contain comparable types.

func main() {
	// Create a map with interface{} as key type
	m := make(map[interface{}]string)

	// Add values using different comparable types as keys
	m[42] = "Integer key"                // int as key
	m["hello"] = "String key"            // string as key
	m[3.14] = "Float key"               // float as key
	m[true] = "Boolean key"             // bool as key
	m[struct{ x int }{1}] = "Struct key" // struct as key

	// Print all key-value pairs
	fmt.Println("Map contents:")
	for k, v := range m {
		fmt.Printf("Key type: %T, Key value: %v, Value: %s\n", k, k, v)
	}

	// Demonstrate key lookup
	if val, exists := m[42]; exists {
		fmt.Printf("\nFound value for key 42: %s\n", val)
	}

	// Demonstrate that different types are treated as different keys
	// even if they have the same "value"
	m[int64(42)] = "Int64 key"
	fmt.Printf("\nValue for int key 42: %s\n", m[42])
	fmt.Printf("Value for int64 key 42: %s\n", m[int64(42)])

	// Note: The following would cause a panic if uncommented
	// because slices are not comparable:
	// m[[]int{1,2,3}] = "Slice key" // panic!

	// Similarly, maps are not comparable:
	// m[map[string]int{"a": 1}] = "Map key" // panic!
}
