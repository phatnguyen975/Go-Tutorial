<div align="center">
  <h1>Methods</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 27, 2026</sub>
</div>

Technically, Go is not an object-oriented programming language. It doesn't have classes, objects, and inheritance.

However, Go has types. And, you can define methods on types. This allows for an object-oriented style of programming in Go.

## Go Methods

A method is nothing but a function with a special **receiver** argument.

The **receiver** argument has a name and a type. It appears between the `func` keyword and the method name.

```go
func (receiver Type) MethodName(parameterList) (returnTypes) {
    // Body
}
```

The receiver can be either a struct type or a non-struct type.

Let's see an example to understand how to define a method on a `type` and how to invoke such a method.

```go
package main

import "fmt"

// Struct type - 'Point'
type Point struct {
	X, Y float64
}

// Method with receiver 'Point'
func (p Point) IsAbove(y float64) bool {
	return p.Y > y
}

func main() {
	p := Point{2.0, 4.0}

	fmt.Println("Point: ", p)

	fmt.Println("Is Point p located above the line y = 1.0 ? : ", p.IsAbove(1))
}
```

Notice the way we called the method `IsAbove()` on the `Point` instance `p`. It's exactly like the way you call methods in an object-oriented programming language.

## Methods are Functions

Since a method is just a function with a receiver argument. We can write the above example using a regular function as well.

```go
package main

import "fmt"

// Struct type - 'Point'
type Point struct {
	X, Y float64
}

func IsAboveFunc(p Point, y float64) bool {
	return p.Y > y
}

/*
    Compare the above function with the corresponding method:
    func (p Point) IsAbove(y float64) bool {
        return p.Y > y
    }
*/

func main() {
	p := Point{2.5, -3.0}

	fmt.Println("Point : ", p)

	fmt.Println("Is Point p located above the line y = 1.0 ? : ", IsAboveFunc(p, 1))
}
```

## Why Methods instead of Functions?

Methods help you write object-oriented style code in Go. Method calls are much easier to read and understand than function calls.

Moreover, Methods help you avoid naming conflicts. Since a method is tied to a particular receiver, you can have the same method names on different receiver types.

Let's see an example.

```go
package main

import (
	"fmt"
	"math"
)

type ArithmeticProgression struct {
	A, D float64
}

// Method with receiver 'ArithmeticProgression'
func (ap ArithmeticProgression) NthTerm(n int) float64 {
	return ap.A + float64(n-1)*ap.D
}

type GeometricProgression struct {
	A, R float64
}

// Method with receiver 'GeometricProgression'
func (gp GeometricProgression) NthTerm(n int) float64 {
	return gp.A * math.Pow(gp.R, float64(n-1))
}

func main() {
	ap := ArithmeticProgression{1, 2} // AP: 1 3 5 7 9 ...
	gp := GeometricProgression{1, 2}  // GP: 1 2 4 8 16 ...

	fmt.Println("5th Term of the Arithmetic series = ", ap.NthTerm(5))
	fmt.Println("5th Term of the Geometric series = ", gp.NthTerm(5))
}
```

## Methods with Pointer Receivers

All the examples that we saw in the previous sections had a value receiver.

With a value receiver, the method operates on a copy of the value passed to it. Therefore, any modifications done to the receiver inside the method is not visible to the caller.

**Go also allows you to define a method with a pointer receiver.**

```go
func (receiver *Type) MethodName(parameterList) (returnTypes) {
    // Body
}
```

Methods with pointer receivers can modify the value to which the receiver points. Such modifications are visible to the caller of the method as well.

```go
package main

import "fmt"

type Point struct {
	X, Y float64
}

/*
    Translates the current Point, at location (X,Y), by dx along the
    x axis and dy along the y axis so that it now represents the
    point (X+dx, Y+dy).
*/
func (p *Point) Translate(dx, dy float64) {
	p.X = p.X + dx
	p.Y = p.Y + dy
}

func main() {
	p := Point{3, 4}
	fmt.Println("Point p = ", p)

	p.Translate(7, 9)
	fmt.Println("After Translate, p = ", p)
}
```

## Methods with Pointer Receivers as Functions

Just like methods with value receivers, we can also write methods with pointer receivers as functions. The following example shows how to rewrite the previous example with functions.

```go
package main

import "fmt"

type Point struct {
	X, Y float64
}

/*
    Translates the current Point, at location (X,Y), by dx along the
    x axis and dy along the y axis so that it now represents the
    point (X+dx, Y+dy).
*/
func TranslateFunc(p *Point, dx, dy float64) {
	p.X = p.X + dx
	p.Y = p.Y + dy
}

func main() {
	p := Point{3, 4}
	fmt.Println("Point p = ", p)

	TranslateFunc(&p, 7, 9)
	fmt.Println("After Translate, p = ", p)
}
```

## Methods and Pointer indirection

### Methods with Pointer Ceceivers vs Functions with Pointer Arguments

A method with a pointer receiver can accept both a pointer and a value as the receiver argument. But, a function with a pointer argument can only accept a pointer.

```go
package main

import "fmt"

type Point struct {
	X, Y float64
}

// Method with Pointer receiver
func (p *Point) Translate(dx, dy float64) {
	p.X = p.X + dx
	p.Y = p.Y + dy
}

// Function with Pointer argument
func TranslateFunc(p *Point, dx, dy float64) {
	p.X = p.X + dx
	p.Y = p.Y + dy
}

func main() {
	p := Point{3, 4}
	ptr := &p
	fmt.Println("Point p = ", p)

	// Calling a Method with Pointer receiver
	p.Translate(2, 6)		// Valid
	ptr.Translate(5, 10)	// Valid

	// Calling a Function with a Pointer argument
	TranslateFunc(ptr, 20, 30) // Valid
	TranslateFunc(p, 20, 30)   // Not Valid
}
```

> Notice how both `p.Translate()` and `ptr.Translate()` calls are valid. Since Go knows that the method `Translate()` has a pointer receiver, It interprets the statement `p.Translate()` as `(&p).Translate()`. It's a syntactic sugar provider by Go for convenience.

### Methods with Value Receivers vs Functions with Value Arguments

A method with a value receiver can accept both a value and a pointer as the receiver argument. But, a function with a value argument can only accept a value.

```go
package main

import "fmt"

// Struct type - 'Point'
type Point struct {
	X, Y float64
}

func (p Point) IsAbove(y float64) bool {
	return p.Y > y
}

func IsAboveFunc(p Point, y float64) bool {
	return p.Y > y
}

func main() {
	p := Point{0, 4}
	ptr := &p

	fmt.Println("Point p = ", p)

	// Calling a Method with Value receiver
	p.IsAbove(1)   // Valid
	ptr.IsAbove(1) // Valid

	// Calling a Function with a Value argument
	IsAboveFunc(p, 1)   // Valid
	IsAboveFunc(ptr, 1) // Not Valid
}
```

> In the above example, both `p.IsAbove()` and `ptr.IsAbove()` statements are valid. Go knows that the method `IsAbove()` has a value receiver. So for convenience, It interprets the statement `ptr.IsAbove()` as `(*ptr).IsAbove()`.

## Method Definition Restrictions

Note that, for you to be able to define a method on a receiver, the receiver type must be defined in the same package.

Go doesn't allow you to define a method on a receiver type that is defined in some other package (this includes built-in types such as `int` as well).

In all the previous examples, the structs and the methods were defined in the same package `main`. Therefore, they worked. However, if you try to define a method on a type defined in some other package, the compiler will complain.

Let's see an example. Consider the following package hierarchy:

```text
src/example
	main/
		- main.go
    model/
        - person.go
```

Here are the contents of `person.go`:

```go
package model

type Person struct {
	FirstName string
	LastName string
	Age int
}
```

Let's try to define a method on the struct `Person`:

```go
package main

import "example/model"

// Error: cannot define new methods on non-local types model.Person
func (p model.Person) getFullName() string {
	return p.FirstName + " " + p.LastName
}

func main() {
}
```

## Defining Methods on Non-struct Types

Go allows you to define methods on non-struct types too. In the following example, I’ve defined a method called `reverse()` on the type `MyString`.

```go
package main

import "fmt"

type MyString string

func (myStr MyString) reverse() string {
	s := string(myStr)
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func main() {
	myStr := MyString("OLLEH")
	fmt.Println(myStr.reverse())
}
```
