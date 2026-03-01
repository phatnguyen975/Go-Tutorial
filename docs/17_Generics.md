<div align="center">
  <h1>Generics</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>March 02, 2026</sub>
</div>

Generics is a powerful feature introduced in [Go 1.18](https://go.dev/doc/go1.18#generics) that allows you to write reusable and type-safe code. With generics, you can define functions and data structures that work with multiple types without the need for runtime type assertions or type casting.

## What are Generics?

**Generics**, also known as parametric polymorphism, enable you to write code that operates on multiple types without explicitly specifying the types upfront. This leads to more concise, reusable, and type-safe code.

Prior to the introduction of generics, you had to use interfaces and type assertions for achieving similar functionality. However, this approach had its drawbacks, such as a lack of type safety and increased boilerplate code.

**Key benefits:**

- **Code reuse:** One function or type works for many types.
- **Type safety:** Catch type errors at compile time.
- **Cleaner code:** No repetitive code or risky type casting.

For example, a function to find the maximum of two values would need separate versions for `int`, `float64`, etc., without generics. Now, you can write it once using a custom constraint for ordered types.

```go
package main

import "fmt"

type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64 | ~string
}

func Max[T Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

func main() {
    fmt.Println(Max(10, 20))
    fmt.Println(Max(3.14, 2.71))
    fmt.Println(Max("go", "golang"))
}
```

Here, `Ordered` is a custom constraint that includes types supporting comparison operators (`>`, `<`, etc.). The `~` allows custom types with these underlying types (e.g., `type MyInt int`). This function works with `int`, `float64`, `string`, or other ordered types.

## Type Parameters and Type Constraints

Type parameters are the heart of generics. They're placeholders for types, defined in square brackets `[T]` after a function or type name. Type constraints restrict the types that can be used as type arguments.

```go
func FunctionName[T Constraint](param T) T {
    // Body
}
```

- `T` is the type parameter.
- `Constraint` specifies allowed types (e.g., `comparable`, `any`, or a custom interface).

Let's define a simple generic function to understand the syntax:

```go
package main

import "fmt"

func PrintSlice[T any](s []T) {
    for _, v := range s {
        fmt.Println(v)
    }
}

func main() {
    intSlice := []int{1, 2, 3}
    stringSlice := []string{"hello", "world"}

    PrintSlice[int](intSlice)
    PrintSlice[string](stringSlice)
}
```

In the `PrintSlice` function, `T` is a type parameter, and `any` is a type constraint. The `any` constraint is built into the language and allows any type to be used. The function can be called with any slice type, such as `[]int` or `[]string`.

**Common Constraints**

|    Constraint    |             Description             |            Example             |
| :--------------: | :---------------------------------: | :----------------------------: |
|      `any`       |           Any type at all           | `int`, `string`, structs, etc. |
|   `comparable`   |  Types that support `==` and `!=`   |    `int`, `string`, `bool`     |
| Custom Interface | User-defined interface with methods | Custom structs implementing it |

**Tip:** Use `any` for flexibility, but tighter constraints like `comparable` for specific operations.

## Writing Generic Functions

Generic functions let you write one function that works with multiple types. The key is choosing the right constraint to support the operations you need.

Example, a generic `Swap` function):

```go
package main

import "fmt"

func Swap[T any](a, b T) (T, T) {
    return b, a
}

func main() {
    x, y := Swap(10, 20)
    fmt.Println(x, y)

    s1, s2 := Swap("hi", "world")
    fmt.Println(s1, s2)
}
```

**Real-world Use Case**

A generic `Filter` function to keep elements matching a condition:

```go
package main

import "fmt"

func Filter[T any](slice []T, predicate func(T) bool) []T {
    result := []T{}
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

func main() {
    nums := []int{1, 2, 3, 4, 5}
    evens := Filter(nums, func(n int) bool { return n%2 == 0 })
    fmt.Println(evens)

    strs := []string{"cat", "dog", "bird"}
    short := Filter(strs, func(s string) bool { return len(s) == 3 })
    fmt.Println(short)
}
```

**Note:** The `any` constraint works here because we only store and return `T`. For operations like `>`, use `comparable` or a custom constraint.

## Building Generic Types

Generics also work with structs, interfaces, and other types. This is great for data structures like lists, stacks, or trees that need to handle any type.

Example, a generic `Stack` struct:

```go
package main

import "fmt"

type Stack[T any] struct {
    data []T
}

func (s *Stack[T]) Push(v T) {
    s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.data) == 0 {
        var zero T
        return zero, false
    }
    lastIndex := len(s.data) - 1
    value := s.data[lastIndex]
    s.data = s.data[:lastIndex]
    return value, true
}

func (s *Stack[T]) Size() int {
    return len(s.data)
}

func main() {
    intStack := Stack[int]{}
    intStack.Push(1)
    intStack.Push(2)
    intStack.Push(3)

    fmt.Println(intStack.Pop())
    fmt.Println(intStack.Size())

    stringStack := Stack[string]{}
    stringStack.Push("hello")
    stringStack.Push("world")

    fmt.Println(stringStack.Pop())
    fmt.Println(stringStack.Size())
}
```

In this example, we define a generic `Stack` type with a type parameter `T` and a constraint `any`. We then implement `Push`, `Pop`, and `Size` methods on this generic type. The main function demonstrates how to create and use stacks of different types: `int` and `string`.

## Using Constraints to Limit Types

Constraints control which types a generic function or type can use. Go offers `any` and `comparable`, but you can define custom constraints with interfaces.

Example, a custom constraint for numeric types:

```go
package main

import "fmt"

type Number interface {
    ~int | ~float64 | ~float32
}

func Sum[T Number](numbers []T) T {
    var total T
    for _, n := range numbers {
        total += n
    }
    return total
}

func main() {
    ints := []int{1, 2, 3}
    floats := []float64{1.5, 2.5, 3.5}
    fmt.Println(Sum(ints))
    fmt.Println(Sum(floats))
}
```

- The `~` allows types with `int`, `float64`, etc., as their underlying type (e.g., `type MyInt int`).
- Use `|` to combine multiple types.

**Constraint Examples**

|       Constraint Example       |           Allowed Types           |
| :----------------------------: | :-------------------------------: |
|             `~int`             |            `~float64`             |
| `interface{ String() string }` | Any type with a `String()` method |
|          `comparable`          |  Types supporting `==` and `!=`   |

**Tip:** Choose the tightest constraint to make your code safe and clear.

## Type Inference: Write Less, Do More

Go's type inference simplifies generics. You often don't need to specify the type when calling a generic function - Go infers it from the arguments.

```go
package main

import "fmt"

func Print[T any](value T) {
    fmt.Println(value)
}

func main() {
    Print(42)
    Print("hello")
}
```

Sometimes, you need explicit types, like with generic types or interfaces:

```go
package main

import "fmt"

func Process[T comparable](a, b T) bool {
    return a == b
}

func main() {
    var x interface{} = 42
    var y interface{} = 42
    result := Process[int](x.(int), y.(int)) // Explicit type needed
    fmt.Println(result)
}
```

**When inference fails:**

- With slices of interfaces.
- When using generic types directly (e.g., `Stack[T]`).

**Tip:** Lean on inference for simple cases, but specify types for complex scenarios.

## Avoiding Common Generics Pitfalls

Generics are powerful, but they have traps. Here are common issues and fixes.

**1. Overusing generics:**

- Don't use generics when a simple function or interface works. A function for just `string` doesn't need generics.
- **Fix:** Only use generics for code reused across types.

**2. Wrong constraints:**

- Using `any` when you need `+` or `>` causes compile errors.
- **Fix:** Use constraints like `Number` or `comparable`.

**3. Nil and zero values:**

- Generic types may return zero values (e.g., `0` for int, `""` for string). Handle them carefully.
- **Example:**

```go
package main

import "fmt"

func GetFirst[T any](slice []T) T {
    if len(slice) == 0 {
        var zero T
        return zero
    }
    return slice[0]
}

func main() {
    nums := []int{}
    fmt.Println(GetFirst(nums))

    strs := []string{}
    fmt.Println(GetFirst(strs))
}
```

**4. Performance:**

- Generics generate specialized code for each type, increasing binary size.
- **Fix:** Test performance for critical code.

## Real-World Example: A Generic Map Function

Let’s wrap up with a practical `Map` function that transforms a slice of one type into another.

```go
package main

import (
    "fmt"
    "strings"
)

func Map[T, U any](slice []T, transform func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = transform(v)
    }
    return result
}

func main() {
    nums := []int{1, 2, 3}
    doubles := Map(nums, func(n int) int { return n * 2 })
    fmt.Println(doubles)

    strs := []string{"a", "b", "c"}
    uppers := Map(strs, strings.ToUpper)
    fmt.Println(uppers)
}
```

**Why this rocks:**

- Works with any input/output type.
- Reusable for countless transformations.
- Type-safe and concise.

**Tip:** Use Map for data transformations like type conversions or calculations.

**Other Use Cases**

- Writing generic algorithms, like sorting, searching, or filtering
- Creating generic data structures, like linked lists, trees, or queues
- Developing reusable utility functions for error handling, logging, or caching
