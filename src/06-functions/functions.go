package main

import (
	"errors"
	"fmt"
)

// 1. Basic Function
func sayHello() {
	fmt.Println("Hello, Go!")
}

// 2. Function With Parameters
func greet(name string) {
	fmt.Println("Hello,", name)
}

// 3. Function With Return Value
func add(a int, b int) int {
	return a + b
}

// 4. Multiple Return Values
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("Cannot divide by zero")
	}
	return a / b, nil
}

// 5. Named Return Values
func rectangle(width, height float64) (area float64, perimeter float64) {
	area = width * height
	perimeter = 2 * (width + height)
	return
}

// 6. Variadic Function
func sum(numbers ...int) int {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

// 7. Anonymous Function
func main() {
	sayHello()

	greet("John")

	fmt.Println("Add:", add(5, 3))

	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Divide:", result)
	}

	area, perimeter := rectangle(5, 3)
	fmt.Println("Area:", area, "Perimeter:", perimeter)

	fmt.Println("Sum:", sum(1, 2, 3, 4, 5))

	// Anonymous function
	func(message string) {
		fmt.Println("Anonymous:", message)
	}("Hello from anonymous function")

	// 8. Function As Value (First-class citizen)
	multiply := func(a, b int) int {
		return a * b
	}
	fmt.Println("Multiply:", multiply(4, 5))

	// 9. Function As Parameter
	operate := func(a, b int, op func(int, int) int) int {
		return op(a, b)
	}

	fmt.Println("Operate Add:", operate(10, 5, add))
	fmt.Println("Operate Multiply:", operate(10, 5, multiply))

	// 10. Closure
	counter := func() func() int {
		count := 0
		return func() int {
			count++
			return count
		}
	}

	c := counter()
	fmt.Println("Counter:", c())
	fmt.Println("Counter:", c())
	fmt.Println("Counter:", c())

	// 11. Defer Function
	defer fmt.Println("This runs at the end of main")

	fmt.Println("End of main")
}
