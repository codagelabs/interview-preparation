package main

import "fmt"

func main() {
	type Person struct {
		Name string
	}

	p1 := &Person{Name: "Alice"}
	p2 := &Person{Name: "Bob"}
	p3 := &Person{Name: "Alice"} // Different instance with the same value

	// Map with pointer keys
	m := make(map[*Person]string)
	m[p1] = "Friend"
	m[p2] = "Colleague"
	m[p3] = "Neighbor"

	// Accessing values
	fmt.Println("p1:", m[p1]) // Output: Friend
	fmt.Println("p2:", m[p2]) // Output: Colleague
	fmt.Println("p3:", m[p3]) // Output: Neighbor
	//Since pointers reference memory locations, different objects (even with the same value) will have unique pointers.
	// This makes them suitable for cases where you want to uniquely identify objects.

	p4 := p1 //Different instance with the reference to same value
	m[p4] = "Relative"

	fmt.Println("p1:", m[p1]) // Output: Relative
	fmt.Println("p4:", m[p4]) // Output: Relative

	fmt.Println("Before mutation:", m[p1]) // Output: Friend

	// Mutating the value
	p1.Name = "Neighbor"
	fmt.Println("After mutation:", m[p1]) // Still Output: Friend
	//Since pointers reference memory locations, the map will still see the original pointer value.
	//Be cautious when mutating the objects referenced by pointers. While the pointer as a key won’t change, the value it points to might, leading to logical confusion.

	// Check if p1 and p3 are treated as different keys
	fmt.Printf("p1 == p3: %v\n", p1 == p3) // Output: false
}
