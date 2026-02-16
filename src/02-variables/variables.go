package main

import "fmt"

func zeroValue() {
	var (
		firstName, lastName string
		age                 int
		salary              float64
		isConfirmed         bool
	)

	fmt.Printf("firstName: %s, lastName: %s, age: %d, salary: %f, isConfirmed: %t",
		firstName, lastName, age, salary, isConfirmed)
}

func initialValue() {
	var firstName, lastName, age, salary = "John", "Maxwell", 28, 50000.0

	fmt.Printf("firstName: %T, lastName: %T, age: %T, salary: %T",
		firstName, lastName, age, salary)
}

func shorthand() {
	name := "Rajeev Singh"
	age, salary, isProgrammer := 35, 50000.0, true

	fmt.Println(name, age, salary, isProgrammer)
}

func main() {
	zeroValue()
	initialValue()
	shorthand()
}
