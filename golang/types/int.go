package main

import "fmt"

func main() {
	var a int8 = 127 // Minimum value for int8
	a++
	fmt.Println(a) // Output: 127 (wraps around to the maximum value)

}

func checkType(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Println("Type is int:", v)
	case int8:
		fmt.Println("Type is int8:", v)
	case int16:
		fmt.Println("Type is int16:", v)
	case int32:
		fmt.Println("Type is int32:", v)
	case int64:
		fmt.Println("Type is int64:", v)
	case string:
		fmt.Println("Type is string:", v)
	case float64:
		fmt.Println("Type is float64:", v)
	default:
		fmt.Println("Unknown type")
	}
}

//128
//256
//512
//1024
//2048
//4098
//8196
