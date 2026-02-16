<div align="center">
  <h1>Variables and Data Types</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 16, 2026</sub>
</div>

## Introduction to Variables and Data Types

Every program needs to store some data/information in memory. The data is stored in memory at a particular memory location.

A variable is just a convenient name given to a memory location where the data is stored. Apart from a name, every variable also has an associated type.

<p align="center">
  <img src="https://www.callicoder.com/static/7a5355bf9eaecbbf401f7637654f0508/a76f4/golang-variables-data-types-illustration.png" style="width:90%;" alt="Variables">
</p>

Data Types or simply Types, categorize related set of data, define the way they are stored, the range of values they can hold, and the operations that can be done on them.

For example, Golang has a data type called `int8`. It represents 8-bit integers whose values can range from `-128` to `127`. It also defines the operations that can be done on `int8` data type such as addition, subtraction, multiplication, division etc.

We also have an `int` data type in Golang whose size is machine dependent. It is 32 bits wide on a 32-bit system and 64 bits wide on a 64-bit system.

Other examples of data types in Golang are `bool`, `string`, `float32`, `float64`, etc. You’ll learn more about these data types in the next section. I gave a brief idea of data types here because it is necessary to understand them before we dive deep into Golang variables.

## Golang Variables in Depth

### Declaring Variables

In Golang, We use the `var` keyword to declare variables:

```go
var firstName string
var lastName string
var age int
```

You can also declare multiple variables at once like so:

```go
var (
	firstName string
	lastName  string
	age       int
)
```

You can even combine multiple variable declarations of the same type with comma:

```go
var (
	firstName, lastName string
	age                 int
)
```

### Zero Values

Any variable declared without an initial value will have a zero-value depending on the type of the variable:

|             Type             | Zero Value |
| :--------------------------: | :--------: |
|            `bool`            |  `false`   |
|           `string`           |    `""`    |
| `int`, `int8`, `int16`, etc. |    `0`     |
|     `float32`, `float64`     |   `0.0`    |

The example below demonstrates the concept of zero values:

```go
package main

import "fmt"

func main() {
	var (
		firstName, lastName string
		age                 int
		salary              float64
		isConfirmed         bool
	)

	fmt.Printf("firstName: %s, lastName: %s, age: %d, salary: %f, isConfirmed: %t\n",
		firstName, lastName, age, salary, isConfirmed)
}
```

### Declaring Variables with Initial Value

Here is how you can initialize variables during declaration:

```go
var firstName string = "Satoshi"
var lastName string = "Nakamoto"
var age int = 35
```

You can also use multiple declarations like this:

```go
var (
	firstName string = "Satoshi"
	lastName  string = "Nakamoto"
	age       int    = 35
)
```

Or even combine multiple variable declarations of the same type with comma and initialize them like so:

```go
var (
	firstName, lastName string = "Satoshi", "Nakamoto"
	age int = 35
)
```

### Type Inference

Although Go is a statically typed language, it doesn't require you to explicitly specify the type of every variable you declare.

When you declare a variable with an initial value, Go automatically infers the type of the variable from the value on the right-hand side. So you need not specify the type when you're initializing the variable at the time of declaration:

```go
package main

import "fmt"

func main() {
    var name = "Rajeev Singh"

    fmt.Printf("Variable 'name' is of type %T", name)
}
```

In the above example, Golang automatically infers the type of the variable as `string` from the value on the right-hand side. If you try to reassign the variable to a value of some other type, then the compiler will throw an error:

```go
var name = "Rajeev Singh" // Type inferred as 'string'
name = 1234 // Compiler error
```

Type inference allows us to declare and initialize multiple variables of different data types in a single line like so:

```go
package main

import "fmt"

func main() {
    var firstName, lastName, age, salary = "John", "Maxwell", 28, 50000.0

    fmt.Printf("firstName: %T, lastName: %T, age: %T, salary: %T",
        firstName, lastName, age, salary)
}
```

### Short Declaration

Go provides a short variable declaration syntax using `:=` operator. It is a shorthand for declaring and initializing a variable (with inferred type).

For example, the shorthand for `var name = "Rajeev"` is `name := "Rajeev"`. Here is a complete example:

```go
package main

import "fmt"

func main() {
	name := "Rajeev Singh"
	age, salary, isProgrammer := 35, 50000.0, true

	fmt.Println(name, age, salary, isProgrammer)
}
```

**Note:** A short variable declaration can **only be used inside a function**. Outside a function, every statement needs to begin with a keyword like `var`, `func`, etc., and therefore, `:=` operator is not available.
