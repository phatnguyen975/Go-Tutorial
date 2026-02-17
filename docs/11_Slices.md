<div align="center">
  <h1>Slices</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 18, 2026</sub>
</div>

A Slice is a segment of an array. Slices build on arrays and provide more power, flexibility, and convenience compared to arrays.

Just like arrays, slices are indexable and have a length. But unlike arrays, they can be resized.

Internally, a slice is just a reference to an underlying array. In this chapter, we'll learn how to create and use slices, and also understand how they work under the hood.

## Declaring a Slice

A slice of type `T` is declared using `[]T`. For example, here is how you can declare a slice of type `int`.

```go
var s []int
```

The slice is declared just like an array except that we do not specify any size in the brackets `[]`.

## Creating and Initializing a Slice

### Creating a Slice Using Literal

You can create a slice using a slice literal like this:

```go
var s = []int{3, 5, 7, 9, 11, 13, 17}
```

The expression on the right-hand side of the above statement is a slice literal. The slice literal is declared just like an array literal, except that you do not specify any size in the square brackets `[]`.

When you create a slice using a slice literal, it first creates an array and then returns a slice reference to it.

Let's see a complete example:

```go
package main

import "fmt"

func main() {
    // Creating a slice using a slice literal
    var s = []int{3, 5, 7, 9, 11, 13, 17}

    // Short hand declaration
    t := []int{2, 4, 8, 16, 32, 64}

    fmt.Println("s = ", s)
    fmt.Println("t = ", t)
}
```

### Creating a Slice from an Array

Since a slice is a segment of an array, we can create a slice from an array.

To create a slice from an array `a`, we specify two indices `low` (lower bound) and `high` (upper bound) separated by a colon:

```go
a[low:high]
```

The above expression selects a slice from the array `a`. The resulting slice includes all the elements starting from index `low` to `high`, but excluding the element at index `high`.

Let's see an example to make things more clear:

```go
package main

import "fmt"

func main() {
    var a = [5]string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon"}

    // Creating a slice from the array
    var s []string = a[1:4]

    fmt.Println("Array a = ", a)
    fmt.Println("Slice s = ", s)
}
```

The `low` and `high` indices in the slice expression are optional. The default value for `low` is `0`, and `high` is the `length` of the slice.

```go
package main

import "fmt"

func main() {
    a := [5]string{"C", "C++", "Java", "Python", "Go"}

    slice1 := a[1:4]
    slice2 := a[:3]
    slice3 := a[2:]
    slice4 := a[:]

    fmt.Println("Array a = ", a)
    fmt.Println("slice1 = ", slice1)
    fmt.Println("slice2 = ", slice2)
    fmt.Println("slice3 = ", slice3)
    fmt.Println("slice4 = ", slice4)
}
```

### Creating a Slice from Another Slice

A slice can also be created by slicing an existing slice.

```go
package main

import "fmt"

func main() {
    cities := []string{"New York", "London", "Chicago", "Beijing", "Delhi", "Mumbai", "Bangalore", "Hyderabad", "Hong Kong"}

    asianCities := cities[3:]
    indianCities := asianCities[1:5]

    fmt.Println("cities = ", cities)
    fmt.Println("asianCities = ", asianCities)
    fmt.Println("indianCities = ", indianCities)
}
```

## Modifying a slice

Slices are reference types. They refer to an underlying array. Modifying the elements of a slice will modify the corresponding elements in the referenced array. Other slices that refer the same array will also see those modifications.

```go
package main

import "fmt"

func main() {
    a := [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

    slice1 := a[1:]
    slice2 := a[3:]

    fmt.Println("------- Before Modifications -------")
    fmt.Println("a  = ", a)
    fmt.Println("slice1 = ", slice1)
    fmt.Println("slice2 = ", slice2)

    slice1[0] = "TUE"
    slice1[1] = "WED"
    slice1[2] = "THU"

    slice2[1] = "FRIDAY"

    fmt.Println("\n-------- After Modifications --------")
    fmt.Println("a  = ", a)
    fmt.Println("slice1 = ", slice1)
    fmt.Println("slice2 = ", slice2)
}
```

## Length and Capacity of a Slice

A slice consists of three things:

- A **pointer** (reference) to an underlying array.
- The **length** of the segment of the array that the slice contains.
- The **capacity** (the maximum size up to which the segment can grow).

<p align="center">
  <img src="https://www.callicoder.com/static/4712a0660a1c95ca833df9526221a282/a76f4/golang-slices-illustration.png" style="width:70%;" alt="Packer">
</p>

Let's consider the following array and the slice obtained from it as an example.

```go
var a = [6]int{10, 20, 30, 40, 50, 60}
var s = [1:4]
```

Here is how the slice `s` in the above example is represented.

<p align="center">
  <img src="https://www.callicoder.com/static/bd787440081265e5d0867add4f1b4b17/d8f62/golang-slices-length-capacity.png" style="width:70%;" alt="Packer">
</p>

The length of the slice is the number of elements in the slice, which is `3` in the above example.

The capacity is the number of elements in the underlying array starting from the first element in the slice. It is `5` in the above example.

You can find the length and capacity of a slice using the built-in functions `len()` and `cap()`.

```go
package main

import "fmt"

func main() {
    a := [6]int{10, 20, 30, 40, 50, 60}
    s := a[1:4]

    fmt.Printf("s = %v, len = %d, cap = %d", s, len(s), cap(s))
}
```

A slice's length can be extended up to its capacity by re-slicing it. Any attempt to extend its length beyond the available capacity will result in a runtime error.

Check out the following example to understand how re-slicing a given slice changes its length and capacity.

```go
package main

import "fmt"

func main() {
    s := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
    fmt.Println("Original Slice")
    fmt.Printf("s = %v, len = %d, cap = %d\n", s, len(s), cap(s))

    s = s[1:5]
    fmt.Println("\nAfter slicing from index 1 to 5")
    fmt.Printf("s = %v, len = %d, cap = %d\n", s, len(s), cap(s))

    s = s[:8]
    fmt.Println("\nAfter extending the length")
    fmt.Printf("s = %v, len = %d, cap = %d\n", s, len(s), cap(s))

    s = s[2:]
    fmt.Println("\nAfter dropping the first two elements")
    fmt.Printf("s = %v, len = %d, cap = %d\n", s, len(s), cap(s))
}
```

## Creating a Slice using The Built-in `make()` Function

Now that we know about the length and capacity of a slice. Let’s look at another way to create a slice.

Go provides a library function called `make()` for creating slices. Following is the signature of `make()` function.

```go
func make([]T, len, cap) []T
```

The `make()` function takes a type, a length, and an optional capacity. It allocates an underlying array with size equal to the given capacity, and returns a slice that refers to that array.

```go
package main

import "fmt"

func main() {
    // Creates an array of size 10, slices it till index 5, and returns the slice reference
    s := make([]int, 5, 10)
    fmt.Printf("s = %v, len = %d, cap = %d", s, len(s), cap(s))
}
```

The `capacity` parameter in the `make()` function is optional. When omitted, it defaults to the specified length.

```go
package main

import "fmt"

func main() {
    // Creates an array of size 5, and returns a slice reference to it
    s := make([]int, 5)
    fmt.Printf("s = %v, len = %d, cap = %d", s, len(s), cap(s))
}
```

## Zero Value of Slices

The zero value of a slice is `nil`. A `nil` slice doesn't have any underlying array, and has a length and capacity of `0`.

```go
package main

import "fmt"

func main() {
    var s []int
    fmt.Printf("s = %v, len = %d, cap = %d\n", s, len(s), cap(s))

    if s == nil {
        fmt.Println("s is nil")
    }
}
```

## Slice Functions

### The `copy()` Function: Copying a Slice

The `copy()` function copies elements from one slice to another. Its signature looks like this:

```go
func copy(dst, src []T) int
```

It takes two slices - a destination slice, and a source slice. It then copies elements from the source to the destination and returns the number of elements that are copied.

The number of elements copied will be the minimum of `len(src)` and `len(dst)`.

```go
package main

import "fmt"

func main() {
    src := []string{"Sublime", "VSCode", "IntelliJ", "Eclipse"}
    dest := make([]string, 2)

    numElementsCopied := copy(dest, src)

    fmt.Println("src = ", src)
    fmt.Println("dest = ", dest)
    fmt.Println("Number of elements copied from src to dest = ", numElementsCopied)
}
```

### The `append()` Function: Appending to a Slice

The `append()` function appends new elements at the end of a given slice. Following is the signature of append function.

```go
func append(s []T, x ...T) []T
```

It takes a slice and a variable number of arguments `x ...T`. It then returns a new slice containing all the elements from the given slice as well as the new elements.

If the given slice doesn't have sufficient capacity to accommodate new elements then a new underlying array is allocated with bigger capacity. All the elements from the underlying array of the existing slice are copied to this new array, and then the new elements are appended.

However, if the slice has enough capacity to accommodate new elements, then the `append()` function re-uses its underlying array and appends new elements to the same array.

Let's see an example to understand things better.

```go
package main

import "fmt"

func main() {
    slice1 := []string{"C", "C++", "Java"}
    slice2 := append(slice1, "Python", "Ruby", "Go")

    fmt.Printf("slice1 = %v, len = %d, cap = %d\n", slice1, len(slice1), cap(slice1))
    fmt.Printf("slice2 = %v, len = %d, cap = %d\n", slice2, len(slice2), cap(slice2))

    slice1[0] = "C#"
    fmt.Println("\nslice1 = ", slice1)
    fmt.Println("slice2 = ", slice2)
}
```

In the above example, since `slice1` has capacity `3`, it can't accommodate more elements. So a new underlying array is allocated with bigger capacity when we append more elements to it.

So if you modify `slice1`, `slice2` won't see those changes because it refers to a different array.

**But what if `slice1` had enough capacity to accommodate new elements?** Well, in that case, no new array would be allocated, and the elements would be added to the same underlying array.

Also, in that case, changes to `slice1` would affect `slice2` as well because both would refer to the same underlying array. This is demonstrated in the following example.

```go
package main

import "fmt"

func main() {
    slice1 := make([]string, 3, 10)
    copy(slice1, []string{"C", "C++", "Java"})

    slice2 := append(slice1, "Python", "Ruby", "Go")

    fmt.Printf("slice1 = %v, len = %d, cap = %d\n", slice1, len(slice1), cap(slice1))
    fmt.Printf("slice2 = %v, len = %d, cap = %d\n", slice2, len(slice2), cap(slice2))

    slice1[0] = "C#"
    fmt.Println("\nslice1 = ", slice1)
    fmt.Println("slice2 = ", slice2)
}
```

**Appending to a nil slice**

When you append values to a `nil` slice, it allocates a new slice and returns the reference of the new slice.

```go
package main

import "fmt"

func main() {
    var s []string

    // Appending to a nil slice
    s = append(s, "Cat", "Dog", "Lion", "Tiger")

    fmt.Printf("s = %v, len = %d, cap = %d", s, len(s), cap(s))
}
```

**Appending one slice to another**

You can directly append one slice to another using the `...` operator. This operator expands the slice to a list of arguments. The following example demonstrates its usage.

```go
package main

import "fmt"

func main() {
    slice1 := []string{"Jack", "John", "Peter"}
    slice2 := []string{"Bill", "Mark", "Steve"}

    slice3 := append(slice1, slice2...)

    fmt.Println("slice1 = ", slice1)
    fmt.Println("slice2 = ", slice2)
    fmt.Println("After appending slice1 & slice2 = ", slice3)
}
```

## Slice of Slices

Slices can be of any type. They can also contain other slices. The example below creates a slice of slices.

```go
package main

import "fmt"

func main() {
    s := [][]string{
      {"India", "China"},
      {"USA", "Canada"},
      {"Switzerland", "Germany"},
    }

    fmt.Println("Slice s = ", s)
    fmt.Println("length = ", len(s))
    fmt.Println("capacity = ", cap(s))
}
```

## Iterating over a Slice

You can iterate over a slice in the same way you iterate over an array.

### Iterating over a Slice Using `for` Loop

```go
package main

import "fmt"

func main() {
    countries := []string{"India", "America", "Russia", "England"}

    for i := 0; i < len(countries); i++ {
        fmt.Println(countries[i])
    }
}
```

### Iterating over a Slice Using the `range` form of `for` Loop

```go
package main

import "fmt"

func main() {
    primeNumbers := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}

    for index, number := range primeNumbers {
        fmt.Printf("PrimeNumber(%d) = %d\n", index + 1, number)
    }
}
```

**Ignoring the `index` from the `range` form of `for` loop using blank identifier**

The `range` form of `for` loop gives you the `index` and the `value` at that index in each iteration. If you don't want to use the `index`, then you can discard it by using an underscore `_`.

The underscore `_` is called the blank identifier. It is used to tell Go compiler that we don’t need this value.

```go
package main

import "fmt"

func main() {
    numbers := []float64{3.5, 7.4, 9.2, 5.4}

    sum := 0.0
    for _, number := range numbers {
        sum += number
    }

    fmt.Printf("Total Sum = %.2f", sum)
}
```
