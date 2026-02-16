package main

import "fmt"

func main() {
	var (
		firstName, lastName string
		age                 int
		salary              float64
		isConfirmed         bool
	)

	fmt.Printf("firstName: %s, lastName: %s, age: %d, salary: %f, isConfirmed: %t",
		firstName, lastName, age, salary, isConfirmed)
}
