<div align="center">
  <h1>Pointers</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 18, 2026</sub>
</div>

A pointer is a variable that stores the memory address of another variable. Confused? Let me explain.

Let's first understand what a variable is. Well, whenever we write any program, we need to store some data/information in memory. The data is stored in memory at a particular address. The memory addresses look something like `0xAFFFF` (That's a hexadecimal representation of a memory address).

Now, to access the data, we need to know the address where it is stored. We can keep track of all the memory addresses where the data related to our program is stored. But imagine how hard it would be to remember all those memory addresses and access data using them.

That is why we have the concept of variables. A variable is just a convenient name given to a memory location where the data is stored.

A pointer is also a variable. But it's a special kind of variable because the data that it stores is not just any normal value like a simple integer or a string, it's a memory address of another variable.

<p align="center">
  <img src="https://www.callicoder.com/static/5cd5cace553fc78068e470cdc5312ec4/ace37/golang-pointers-illustration.png" style="width:80%;" alt="Pointer">
</p>

In the above example, the pointer `p` contains the value `0x0001` which is the address of the variable `a`.

## Declaring a Pointer

A pointer of type `T` is declared using the following syntax:

```go
var p *T
```

The type `T` is the type of the variable that the pointer points to. For example, following is a pointer of type `int`.

```go
var p *int
```

The above pointer can only store the memory address of `int` variables.

The zero value of a pointer is `nil`. That means any uninitialized pointer will have the value `nil`. Let's see a complete example.

```go
package main

import "fmt"

func main() {
	var p *int
	fmt.Println("p = ", p)
}
```

## Initializing a Pointer

You can initialize a pointer with the memory address of another variable. The address of a variable can be retrieved using the `&` operator.

```go
var x = 100
var p *int = &x
```

Notice how we use the `&` operator with the variable `x` to get its address, and then assign the address to the pointer `p`.

Just like any other variable in Go, the type of a pointer variable is also inferred by the compiler. So you can omit the type declaration from the pointer `p` in the above example and write it like so:

```go
var p = &x
```

Let's see a complete example to make things more clear:

```go
package main

import "fmt"

func main() {
	var a = 5.67
	var p = &a

	fmt.Println("Value stored in variable a = ", a)
	fmt.Println("Address of variable a = ", &a)
	fmt.Println("Value stored in variable p = ", p)
}
```

## Dereferencing a Pointer

You can use the `*` operator on a pointer to access the value stored in the variable that the pointer points to. This is called **dereferencing** or **indirecting**.

```go
package main

import "fmt"

func main() {
	var a = 100
	var p = &a

	fmt.Println("a = ", a)
	fmt.Println("p = ", p)
	fmt.Println("*p = ", *p)
}
```

You can not only access the value of the pointed variable using `*` operator, but you can change it as well. The following example sets the value stored in the variable `a` through the pointer `p`.

```go
package main

import "fmt"

func main() {
	var a = 1000
	var p = &a

	fmt.Println("a (before) = ", a)

	// Changing the value stored in the pointed variable through the pointer
	*p = 2000

	fmt.Println("a (after) = ", a)
}
```

## Creating a Pointer Using The Built-in `new()` Function

You can also create a pointer using the built-in `new()` function. The `new()` function takes a type as an argument, allocates enough memory to accommodate a value of that type, and returns a pointer to it.

```go
package main

import "fmt"

func main() {
	ptr := new(int) // Pointer to an 'int' type
	*ptr = 100

	fmt.Printf("Ptr = %#x, Ptr value = %d", ptr, *ptr)
}
```

## Pointer to Pointer

A pointer can point to a variable of any type. It can point to another pointer as well. The following example shows how to create a pointer to another pointer:

```go
package main

import "fmt"

func main() {
	var a = 7.98
	var p = &a
	var pp = &p

	fmt.Println("a = ", a)
	fmt.Println("Address of a = ", &a)

	fmt.Println("p = ", p)
	fmt.Println("Address of p = ", &p)

	fmt.Println("pp = ", pp)

	// Dereferencing a pointer to pointer
	fmt.Println("*pp = ", *pp)
	fmt.Println("**pp = ", **pp)
}
```

## No Pointer Arithmetic in Go

If you have worked with C/C++, then you must be aware that these languages support pointer arithmetic. For example, you can increment/decrement a pointer to move to the next/previous memory address. You can add or subtract an integer value to/from a pointer. You can also compare two pointers using relational operators `==`, `<`, `>`, etc.

But Go doesn't support such arithmetic operations on pointers. Any such operation will result in a compile time error.

```go
package main

func main() {
	var x = 67
	var p = &x

	var p1 = p + 1 // Compiler error: Invalid operation
}
```

You can, however, compare two pointers of the same type for equality using `==` operator.

```go
package main

import "fmt"

func main() {
	var a = 75
	var p1 = &a
	var p2 = &a

	if p1 == p2 {
		fmt.Println("Both pointers p1 and p2 point to the same variable")
	}
}
```
