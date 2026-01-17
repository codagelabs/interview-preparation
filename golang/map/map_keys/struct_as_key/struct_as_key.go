package main

import "fmt"


type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	Salary int
	dept struct {
		deptName string
		deptCode string
	}
}


func (p Person) String() string {
	return fmt.Sprintf("Name: %s, Age: %d", p.Name, p.Age)
}


func main() {

	// map with struct as key
	personMap := make(map[Person]string)
	personMap[Person{Name: "John", Age: 20}] = "New York"
	personMap[Person{Name: "Jane", Age: 30}] = "Los Angeles"

	fmt.Println(personMap)


	// Create a map with Employee struct as key, which embeds Person struct and has additional fields
	// The Employee struct contains:
	// - An embedded Person struct (with Name and Age)
	// - A Salary field
	// - A dept field which is an anonymous struct containing department details

	employeeMap := make(map[Employee]string)
	employeeMap[Employee{Person: Person{Name: "John", Age: 20}, Salary: 100000, dept: struct {
		deptName string
		deptCode string
	}{deptName: "IT", deptCode: "IT001"}}] = "New York"	
	fmt.Println(employeeMap)










}

