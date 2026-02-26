<div align="center">
  <h1>Interfaces</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 27, 2026</sub>
</div>

## Go Interface

An interface in Go is a type defined using a set of method signatures. The interface defines the behavior for similar type of objects.

For example, here is an interface that defines the behavior for Geometrical shapes:

```go
// Go Interface - 'Shape'
type Shape interface {
	Area() float64
	Perimeter() float64
}
```

An interface is declared using the type keyword, followed by the name of the interface and the keyword `interface`. Then, we specify a set of method signatures inside curly braces.

## Implementing an Interface

To implement an interface, you just need to implement all the methods declared in the interface.

**Go Interfaces are implemented implicitly**

Unlike other languages like Java, you don't need to explicitly specify that a type implements an interface using something like an `implements` keyword. You just implement all the methods declared in the interface and you're done.

Here are two struct types that implement the `Shape` interface:

```go
// Struct type 'Rectangle' - implements the 'Shape' interface by implementing all its methods
type Rectangle struct {
	Length, Width float64
}

func (r Rectangle) Area() float64 {
	return r.Length * r.Width
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Length + r.Width)
}
```

```go
// Struct type 'Circle' - implements the 'Shape' interface by implementing all its methods
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) Diameter() float64 {
	return 2 * c.Radius
}
```

## Using an Interface Type with Concrete Values

An interface in itself is not that useful unless we use it with a concrete type that implements all its methods.

Let's see how an interface can be used with concrete values.

- **An interface type can hold any value that implements all its methods**

```go
package main

import "fmt"

func main() {
	var s Shape = Circle{5.0}
	fmt.Printf("Shape Type = %T, Shape Value = %v\n", s, s)
	fmt.Printf("Area = %f, Perimeter = %f\n\n", s.Area(), s.Perimeter())

	s = Rectangle{4.0, 6.0}
	fmt.Printf("Shape Type = %T, Shape Value = %v\n", s, s)
	fmt.Printf("Area = %f, Perimeter = %f\n", s.Area(), s.Perimeter())
}
```

- **Using interface types as arguments to functions**

```go
package main

import "fmt"

// Generic function to calculate the total area of multiple shapes of different types
func CalculateTotalArea(shapes ...Shape) float64 {
	totalArea := 0.0
	for _, s := range shapes {
		totalArea += s.Area()
	}
	return totalArea
}

func main() {
	totalArea := CalculateTotalArea(Circle{2}, Rectangle{4, 5}, Circle{10})
	fmt.Println("Total area = ", totalArea)
}
```

- **Using interface types as fields**

```go
package main

import "fmt"

// Interface types can also be used as fields
type MyDrawing struct {
	shapes  []Shape
	bgColor string
	fgColor string
}

func (drawing MyDrawing) Area() float64 {
	totalArea := 0.0
	for _, s := range drawing.shapes {
		totalArea += s.Area()
	}
	return totalArea
}

func main() {
	drawing := MyDrawing{
		shapes: []Shape{
			Circle{2},
			Rectangle{3, 5},
			Rectangle{4, 7},
		},
		bgColor: "red",
		fgColor: "white",
	}

	fmt.Println("Drawing", drawing)
	fmt.Println("Drawing Area = ", drawing.Area())
}
```

## Interface Values: How Does an Interface Type Work with Concrete Values?

Under the hood, an interface value can be thought of as a tuple consisting of a value and a concrete type:

```go
// interface
(value, type)
```

Let's see an example to understand more:

```go
package main

import "fmt"

func main() {
	var s Shape

	s = Circle{5}
	fmt.Printf("(%v, %T)\n", s, s)
	fmt.Printf("Shape area = %v\n", s.Area())

	s = Rectangle{4, 7}
	fmt.Printf("(%v, %T)\n", s, s)
	fmt.Printf("Shape area = %v\n", s.Area())
}
```

Checkout the output of the above program and notice how the variable `s` has information about the value as well as the type of the `Shape` that is assigned to it.

When we call a method on an interface value, a method of the same name on its underlying type is executed.

For example, in the above program, when we call the method `Area()` on the variable `s`, it executes the `Area()` method of its underlying type.

## Footnotes

- [How to use interfaces in Go](http://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go)
- [Go Data Structures: Interfaces](https://research.swtch.com/interfaces)
