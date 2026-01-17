package main

import "fmt"

func main() {

	//The zero value of a slice is nil. The len and cap functions will both return 0 for a nil slice.
	var nilSlice []int
	fmt.Println(len(nilSlice), cap(nilSlice)) // output: 0 0

	var array [5]int
	fmt.Println(len(array), cap(array)) // output: 5 5

	letters := []string{"a", "b", "c", "d"}
	fmt.Println(len(letters), cap(letters)) // output: 4 4
	//A slice can be created with the built-in function called make, which has the signature,
	//func make([]T, len, cap) []T
	s := make([]int, 1, 5)
	fmt.Println(s)              // output: [0]
	fmt.Println(len(s), cap(s)) //output: 1 5

	//When the capacity argument is omitted, it defaults to the specified length. Here’s a more succinct version of the same code:
	s = make([]int, 5)
	fmt.Println(s)              // output: [0 0 0 0 0]
	fmt.Println(len(s), cap(s)) //output: 5 5

	//s = make([]int, 5, 4)
	// invalid argument: length and capacity swapped
	// capacity cannot be less than len

	//Create slice by slicing it
	alphabates := []string{"a", "b", "c", "d", "e", "f"}
	ef := alphabates[4:]
	fmt.Println(ef)               //output: ["e", "f"]
	fmt.Println(len(ef), cap(ef)) //output: 2 3
	bcd := alphabates[1:4]
	fmt.Println(bcd)                //output: ["b", "c", "d"]
	fmt.Println(len(bcd), cap(bcd)) //output: 3 5
	abcd := alphabates[:4]
	fmt.Println(abcd)                 //output: ["a","b", "c", "d"]
	fmt.Println(len(abcd), cap(abcd)) //output: 4 6

	bcd[0] = "New"
	fmt.Println(bcd)
	fmt.Println(abcd)
	fmt.Println(alphabates)

	bcd = append(alphabates, "g", "n", "m", "l")

	bcd[0] = "New"
	fmt.Println(bcd)
	fmt.Println(abcd)
	fmt.Println(alphabates)

	// In Go, this demonstrates slice growth behavior:
	// 1. Initial slice 's' starts with len=5, cap=5 (from previous make() call)
	// 2. append() behavior in Go:
	//    - When current cap is full, Go allocates new array
	//    - For slices < 1024 elements: new_cap = old_cap * 2
	//    - For slices >= 1024 elements: new_cap = old_cap * 1.25
	// 3. Memory efficiency:
	//    - Go uses this growth strategy to amortize the cost of reallocations
	//    - Prevents too frequent reallocations while managing memory usage
	// 4. Performance note:
	//    - If final size is known, better to pre-allocate with make()
	//    - Helps avoid multiple reallocations during growth
	for i := 0; i < 2000; i++ {
		s = append(s, i)
		fmt.Printf("Len: %d, Cap: %d\n", len(s), cap(s))
	}	

}
