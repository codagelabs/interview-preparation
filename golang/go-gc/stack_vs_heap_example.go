package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

// assignment  Stack vs Heap Allocation Examples
// stack and heap
// memory management and allocation on stack 
// memory usage and garbage collection


// memory usage

// Example 1: Stack Allocation - Simple values
func stackAllocationExample() {
	// These values are allocated on the stack
	var (
		integer int = 42
		float   float64 = 3.14
		boolean bool = true
		str     string = "hello"
		array   [5]int = [5]int{1, 2, 3, 4, 5}
		point   struct {
			x, y int
			name string
		} = struct {
			x, y int
			name string
		}{10, 20, "point"}
	)

	fmt.Printf("Stack allocated values:\n")
	fmt.Printf("  Integer: %d (size: %d bytes)\n", integer, unsafe.Sizeof(integer))
	fmt.Printf("  Float: %f (size: %d bytes)\n", float, unsafe.Sizeof(float))
	fmt.Printf("  Boolean: %t (size: %d bytes)\n", boolean, unsafe.Sizeof(boolean))
	fmt.Printf("  String: %s (size: %d bytes)\n", str, unsafe.Sizeof(str))
	fmt.Printf("  Array: %v (size: %d bytes)\n", array, unsafe.Sizeof(array))
	fmt.Printf("  Point: %+v (size: %d bytes)\n", point, unsafe.Sizeof(point))
}

// Example 2: Stack Allocation - Small slices with known size
func stackAllocationSmallSlices() {
	// Small slices with known size at compile time can be stack allocated
	smallSlice := make([]int, 10)
	for i := range smallSlice {
		smallSlice[i] = i
	}
	fmt.Printf("Small slice (stack): %v (size: %d bytes)\n", smallSlice, unsafe.Sizeof(smallSlice))
}

// Example 3: Heap Allocation - Values that escape
func heapAllocationExample() *int {
	// This value escapes to heap because it's returned
	value := 42
	fmt.Printf("Heap allocated value: %d (returned pointer)\n", value)
	return &value // This causes the value to escape to heap
}

// Example 4: Heap Allocation - Dynamic slices
func heapAllocationDynamicSlices(size int) {
	// This slice escapes to heap because size is determined at runtime
	dynamicSlice := make([]int, size)
	for i := range dynamicSlice {
		dynamicSlice[i] = i
	}
	fmt.Printf("Dynamic slice (heap): size=%d, len=%d, cap=%d\n", 
		size, len(dynamicSlice), cap(dynamicSlice))
}

// Example 5: Heap Allocation - Large values
func heapAllocationLargeValues() {
	// Large arrays escape to heap
	largeArray := [1000]int{}
	for i := range largeArray {
		largeArray[i] = i
	}
	fmt.Printf("Large array (heap): size=%d bytes\n", unsafe.Sizeof(largeArray))
}

// Example 6: Heap Allocation - Values referenced by pointers
func heapAllocationPointerReferences() {
	// Values referenced by pointers that escape
	type Person struct {
		Name string
		Age  int
	}

	// This struct escapes to heap because we return a pointer to it
	person := Person{Name: "Alice", Age: 30}
	personPtr := &person

	fmt.Printf("Person struct (heap): %+v (size: %d bytes)\n", 
		*personPtr, unsafe.Sizeof(*personPtr))
}

// Example 7: Mixed allocation - Some stack, some heap
func mixedAllocationExample() {
	// Stack allocated
	stackValue := 100
	stackArray := [5]int{1, 2, 3, 4, 5}

	// Heap allocated (because we take their addresses)
	heapValue := 200
	heapArray := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Taking addresses forces heap allocation
	stackPtr := &stackValue
	stackArrayPtr := &stackArray
	heapPtr := &heapValue
	heapArrayPtr := &heapArray

	fmt.Printf("Mixed allocation:\n")
	fmt.Printf("  Stack value: %d (ptr: %p)\n", *stackPtr, stackPtr)
	fmt.Printf("  Stack array: %v (ptr: %p)\n", *stackArrayPtr, stackArrayPtr)
	fmt.Printf("  Heap value: %d (ptr: %p)\n", *heapPtr, heapPtr)
	fmt.Printf("  Heap array: %v (ptr: %p)\n", *heapArrayPtr, heapArrayPtr)
}

// Example 8: Function parameters and return values
func functionParameterExample(value int) int {
	// Parameter 'value' is stack allocated
	// Return value is also stack allocated
	return value * 2
}

func functionParameterHeapExample(value *int) *int {
	// Parameter 'value' points to heap-allocated memory
	// Return value points to heap-allocated memory
	result := *value * 2
	return &result // This escapes to heap
}

// Example 9: Interface values
func interfaceAllocationExample() {
	// Interface values are heap allocated
	var interfaceValue interface{} = "hello"
	var interfaceSlice []interface{} = []interface{}{1, "world", 3.14}

	fmt.Printf("Interface values (heap):\n")
	fmt.Printf("  Interface: %v (size: %d bytes)\n", interfaceValue, unsafe.Sizeof(interfaceValue))
	fmt.Printf("  Interface slice: %v (size: %d bytes)\n", interfaceSlice, unsafe.Sizeof(interfaceSlice))
}

// Example 10: Channel and map allocation
func channelMapAllocationExample() {
	// Channels and maps are always heap allocated
	ch := make(chan int, 10)
	m := make(map[string]int)
	m["key"] = 42

	fmt.Printf("Channel and map (heap):\n")
	fmt.Printf("  Channel: %v (size: %d bytes)\n", ch, unsafe.Sizeof(ch))
	fmt.Printf("  Map: %v (size: %d bytes)\n", m, unsafe.Sizeof(m))
}

// Example 11: String allocation
func stringAllocationExample() {
	// Small strings can be stack allocated
	smallString := "hello"
	
	// Large strings or concatenated strings may be heap allocated
	largeString := "this is a very long string that might be allocated on the heap"
	concatenatedString := smallString + " " + largeString

	fmt.Printf("String allocation:\n")
	fmt.Printf("  Small string: %s (size: %d bytes)\n", smallString, unsafe.Sizeof(smallString))
	fmt.Printf("  Large string: %s (size: %d bytes)\n", largeString, unsafe.Sizeof(largeString))
	fmt.Printf("  Concatenated string: %s (size: %d bytes)\n", concatenatedString, unsafe.Sizeof(concatenatedString))
}

// Example 12: Slice allocation patterns
func sliceAllocationPatterns() {
	// Stack allocated (small, known size)
	smallSlice := make([]int, 5)
	
	// Heap allocated (large size)
	largeSlice := make([]int, 1000)
	
	// Heap allocated (dynamic size)
	dynamicSlice := make([]int, runtime.NumCPU())
	
	// Heap allocated (capacity > size)
	capacitySlice := make([]int, 5, 100)

	fmt.Printf("Slice allocation patterns:\n")
	fmt.Printf("  Small slice: len=%d, cap=%d (likely stack)\n", len(smallSlice), cap(smallSlice))
	fmt.Printf("  Large slice: len=%d, cap=%d (heap)\n", len(largeSlice), cap(largeSlice))
	fmt.Printf("  Dynamic slice: len=%d, cap=%d (heap)\n", len(dynamicSlice), cap(dynamicSlice))
	fmt.Printf("  Capacity slice: len=%d, cap=%d (heap)\n", len(capacitySlice), cap(capacitySlice))
}

// Example 13: Struct allocation
func structAllocationExample() {
	// Small struct - stack allocated
	type SmallStruct struct {
		x, y int
		name string
	}
	smallStruct := SmallStruct{x: 10, y: 20, name: "point"}

	// Large struct - heap allocated
	type LargeStruct struct {
		data [1000]int
		name string
	}
	largeStruct := LargeStruct{name: "large"}

	// Struct with pointers - heap allocated
	type PointerStruct struct {
		ptr *int
		name string
	}
	value := 42
	pointerStruct := PointerStruct{ptr: &value, name: "pointer"}

	fmt.Printf("Struct allocation:\n")
	fmt.Printf("  Small struct: %+v (size: %d bytes)\n", smallStruct, unsafe.Sizeof(smallStruct))
	fmt.Printf("  Large struct: %+v (size: %d bytes)\n", largeStruct, unsafe.Sizeof(largeStruct))
	fmt.Printf("  Pointer struct: %+v (size: %d bytes)\n", pointerStruct, unsafe.Sizeof(pointerStruct))
}

// Example 14: Escape analysis demonstration
func escapeAnalysisDemo() {
	fmt.Printf("\n=== Escape Analysis Demo ===\n")
	
	// Case 1: No escape - value stays on stack
	noEscape := 42
	fmt.Printf("No escape: %d (stack allocated)\n", noEscape)

	// Case 2: Escape - value moves to heap
	escape := 42
	escapePtr := &escape
	fmt.Printf("Escape: %d (heap allocated, ptr: %p)\n", *escapePtr, escapePtr)

	// Case 3: Interface escape
	var interfaceValue interface{} = 42
	fmt.Printf("Interface escape: %v (heap allocated)\n", interfaceValue)

	// Case 4: Slice escape
	slice := make([]int, 10)
	slicePtr := &slice
	fmt.Printf("Slice escape: %v (heap allocated, ptr: %p)\n", *slicePtr, slicePtr)
}

// Example 15: Memory usage demonstration
func memoryUsageDemo() {
	fmt.Printf("\n=== Memory Usage Demo ===\n")
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	fmt.Printf("Initial heap stats:\n")
	fmt.Printf("  Heap Alloc: %d bytes\n", m.HeapAlloc)
	fmt.Printf("  Heap Sys: %d bytes\n", m.HeapSys)
	fmt.Printf("  Heap Objects: %d\n", m.HeapObjects)

	// Allocate some heap memory
	heapSlice := make([]int, 10000)
	for i := range heapSlice {
		heapSlice[i] = i
	}

	runtime.ReadMemStats(&m)
	fmt.Printf("\nAfter heap allocation:\n")
	fmt.Printf("  Heap Alloc: %d bytes\n", m.HeapAlloc)
	fmt.Printf("  Heap Sys: %d bytes\n", m.HeapSys)
	fmt.Printf("  Heap Objects: %d\n", m.HeapObjects)

	// Force garbage collection
	runtime.GC()

	runtime.ReadMemStats(&m)
	fmt.Printf("\nAfter garbage collection:\n")
	fmt.Printf("  Heap Alloc: %d bytes\n", m.HeapAlloc)
	fmt.Printf("  Heap Sys: %d bytes\n", m.HeapSys)
	fmt.Printf("  Heap Objects: %d\n", m.HeapObjects)
}

func main() {
	fmt.Printf("=== Go Stack vs Heap Allocation Examples ===\n\n")

	// Run all examples
	stackAllocationExample()
	fmt.Println()

	stackAllocationSmallSlices()
	fmt.Println()

	heapValue := heapAllocationExample()
	fmt.Printf("Returned heap value: %d\n", *heapValue)
	fmt.Println()

	heapAllocationDynamicSlices(100)
	fmt.Println()

	heapAllocationLargeValues()
	fmt.Println()

	heapAllocationPointerReferences()
	fmt.Println()

	mixedAllocationExample()
	fmt.Println()

	result := functionParameterExample(21)
	fmt.Printf("Function parameter result: %d\n", result)
	fmt.Println()

	value := 10
	heapResult := functionParameterHeapExample(&value)
	fmt.Printf("Function parameter heap result: %d\n", *heapResult)
	fmt.Println()

	interfaceAllocationExample()
	fmt.Println()

	channelMapAllocationExample()
	fmt.Println()

	stringAllocationExample()
	fmt.Println()

	sliceAllocationPatterns()
	fmt.Println()

	structAllocationExample()
	fmt.Println()

	escapeAnalysisDemo()
	fmt.Println()

	memoryUsageDemo()
	fmt.Println()

	fmt.Printf("=== Examples Complete ===\n")
	fmt.Printf("\nKey Points:\n")
	fmt.Printf("1. Stack allocation: Small, simple values with known lifetime\n")
	fmt.Printf("2. Heap allocation: Large values, dynamic sizes, escaped values\n")
	fmt.Printf("3. Escape analysis determines allocation location\n")
	fmt.Printf("4. Pointers, interfaces, channels, maps always heap allocated\n")
	fmt.Printf("5. Use profiling tools to understand actual allocation patterns\n")
} 