package main

import "fmt"

func main() {

	//Golang Array can be used as keys for map
	var a = [...]int{10, 20, 20}
	var b = [3]int{10, 20, 20}
	sliceOfValues := make(map[[3]int]string)
	sliceOfValues[a] = "dsdsds"
	sliceOfValues[b] = "wqwq"
	fmt.Println(sliceOfValues)

	//Golang Map with simple struct as key
	type Student struct {
		Name   string
		RollNo int
	}
	var yash Student = Student{
		"Yash",
		20,
	}
	var shiva = Student{
		"shiva",
		30,
	}

	om := Student{
		"om",
		30,
	}

	structOfValues := make(map[Student]string)
	structOfValues[yash] = "some text"
	structOfValues[shiva] = "some text"
	structOfValues[om] = "OM is good student"

	fmt.Println(structOfValues)

	unitMap := map[uint8]string{}
	unitMap[120] = "some value"
	fmt.Println(unitMap)

}
