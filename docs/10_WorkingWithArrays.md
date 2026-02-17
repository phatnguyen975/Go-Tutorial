<div align="center">
  <h1>Working with Arrays</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 18, 2026</sub>
</div>

An array is a fixed-size collection of elements of the same type. The elements of the array are stored sequentially and can be accessed using their index.

<p align="center">
  <img src="https://www.callicoder.com/static/7f0d6552c97069c1351fc510608f4018/78a22/Golang-Array.png" style="width:80%;" alt="Packer">
</p>

## Declaring an Array

You can declare an array of length `n` and type `T` like so:

```go
var a[n]T
```

Now let's see a complete example:

```go
package main

import "fmt"

func main() {
	var x [5]int // An array of 5 integers
	fmt.Println(x)

	var y [8]string // An array of 8 strings
	fmt.Println(y)

	var z [3]complex128 // An array of 3 complex numbers
	fmt.Println(z)
}
```

By default, all the array elements are initialized with the zero value of the corresponding array type.

For example, if we declare an integer array, all the elements will be initialized with `0`. If we declare a string array, all the elements will be initialized with an empty string `""`, and so on.

## Accessing Array Elements

The elements of an array are stored sequentially and can be accessed by their index. The index starts from `zero` and ends at `length - 1`.

```go
package main

import "fmt"

func main() {
	var x [5]int // An array of 5 integers

	x[0] = 100
	x[1] = 101
	x[3] = 103
	x[4] = 105

	fmt.Printf("x[0] = %d, x[1] = %d, x[2] = %d\n", x[0], x[1], x[2])
	fmt.Println("x = ", x)
}
```

In the above example, since we didn't assign any value to `x[2]`, it has the value 0 (The zero value for integers).

## Initializing an Array Using Literal

You can declare and initialize an array at the same time like this:

```go
var a = [5]int{2, 4, 6, 8, 10}
```

The expression on the right-hand side of the above statement is called an array literal.

**Note:** We do not need to specify the type of the variable `a` as in `var a [5]int`, because the compiler can automatically infer the type from the expression on the right hand side.

You can also use Go's short variable declaration for declaring and initializing an array. The above array declaration can also be written as below inside any function.

```go
a := [5]int{2, 4, 6, 8, 10}
```

Here is a complete example:

```go
package main

import "fmt"

func main() {
	// Declaring and initializing an array at the same time
	var a = [5]int{2, 4, 6, 8, 10}
	fmt.Println(a)

	// Short declaration for declaring and initializing an array
	b := [5]int{2, 4, 6, 8, 10}
	fmt.Println(b)

	// You don't need to initialize all the elements of the array
	// The un-initialized elements will be assigned the zero value of the corresponding array type
	c := [5]int{2}
	fmt.Println(c)
}
```

**Let Go compiler infer the length of the array**

You can also omit the size declaration from the initialization expression of the array, and let the compiler count the number of elements for you.

```go
package main

import "fmt"

func main() {
    // Letting Go compiler infer the length of the array
	a := [...]int{3, 5, 7, 9, 11, 13, 17}
	fmt.Println(a)

    // If you specify the index with :, the elements in between will be zeroed
    b = [...]int{100, 3: 400, 500}
    fmt.Println(b)
}
```

## Exploring More about Arrays

1. **Array's length is part of its type**

The length of an array is part of its type. So the array `a[5]int` and `a[10]int` are completely distinct types, and you cannot assign one to the other.

This also means that you cannot resize an array, because resizing an array would mean changing its type, and you cannot change the type of a variable in Go.

```go
package main

func main() {
	var a = [5]int{3, 5, 7, 9, 11}
	var b [10]int = a // Error, a and b are distinct types
}
```

2. **Arrays in Go are value types**

Arrays in Go are value types unlike other languages like `C`, `C++`, and `Java` where arrays are reference types.

This means that when you assign an array to a new variable or pass an array to a function, the entire array is copied. So if you make any changes to this copied array, the original array won't be affected and will remain unchanged.

```go
package main

import "fmt"

func main() {
	a1 := [5]string{"English", "Japanese", "Spanish", "French", "Hindi"}
	a2 := a1 // A copy of the array `a1` is assigned to 'a2'

	a2[1] = "German"

	fmt.Println("a1 = ", a1) // The array 'a1' remains unchanged
	fmt.Println("a2 = ", a2)
}
```

## Iterating over An Array

You can use the `for` loop to iterate over an array like so:

```go
package main

import "fmt"

func main() {
	names := [3]string{"Mark Zuckerberg", "Bill Gates", "Larry Page"}

	for i := 0; i < len(names); i++ {
		fmt.Println(names[i])
	}
}
```

The `len()` function in the above example is used to find the length of the array.

Let's see another example. In the example below, we find the sum of all the elements of the array by iterating over the array, and adding the elements one by one to the variable `sum`.

```go
package main

import "fmt"

func main() {
	a := [4]float64{3.5, 7.2, 4.8, 9.5}
	sum := float64(0)

	for i := 0; i < len(a); i++ {
		sum = sum + a[i]
	}

	fmt.Printf("Sum of all the elements in array  %v = %f", a, sum)
}
```

### Iterating over An Array Using `range`

Go provides a more powerful form of `for` loop using the `range` operator. Here is how you can use the `range` operator with `for` loop to iterate over an array:

```go
package main

import "fmt"

func main() {
	daysOfWeek := [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

	for index, value := range daysOfWeek {
		fmt.Printf("Day %d of week = %s\n", index, value)
	}
}
```

Let's now write the same `sum` example that we wrote with normal `for` loop using `range` form of the `for` loop:

```go
package main

import "fmt"

func main() {
	a := [4]float64{3.5, 7.2, 4.8, 9.5}
	sum := float64(0)

	for index, value := range a {
		sum = sum + value
	}

	fmt.Printf("Sum of all the elements in array %v = %f", a, sum)
}
```

When you run the above program, it’ll generate an error like this:

```bash
./array_iteration_range.go:9:13: index declared and not used
```

Go compiler doesn't allow creating variables that are never used. You can fix this by using an `_` (underscore) in place of `index`.

```go
package main

import "fmt"

func main() {
	a := [4]float64{3.5, 7.2, 4.8, 9.5}
	sum := float64(0)

	for _, value := range a {
		sum += value
	}

	fmt.Printf("Sum of all the elements in array %v = %f", a, sum)
}
```

The underscore `_` is used to tell the compiler that we don't need this variable.

## Multidimensional Arrays

All the arrays that we created so far in this post are one dimensional. You can also create multi-dimensional arrays in Go.

The following example demonstrates how to create multidimensional:

```go
package main

import "fmt"

func main() {
	a := [2][2]int{
		{3, 5},
		{7, 9},	// This trailing comma is mandatory
	}
	fmt.Println(a)

	// Just like 1D arrays, you don't need to initialize all the elements in a multi-dimensional array
	// Un-initialized array elements will be assigned the zero value of the array type
	b := [3][4]float64{
		{1, 3},
		{4.5, -3, 7.4, 2},
		{6, 2, 11},
	}
	fmt.Println(b)
}
```
