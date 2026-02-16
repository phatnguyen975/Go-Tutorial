<div align="center">
  <h1>Working with Constants</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 16, 2026</sub>
</div>

## Constants

In Go, we use the term constant to represent fixed (unchanging) values such as `5`, `1.34`, `true`, `"Hello"`, etc.

### Literals are Constants

All the literals in Go, be it integer literals like `5`, `1000`, or floating-point literals like `4.76`, `1.89`, or boolean literals like `true`, `false`, or string literals like `"Hello"`, `"John"` are constants.

### Declaring a Constant

Literals are constants without a name. To declare a constant and give it a name, you can use the `const` keyword like so:

```go
const myFavLanguage = "Python"
const sunRisesInTheEast = true
```

You can also specify a type in the declaration like this:

```go
const a int = 1234
const b string = "Hi"
```

Multiple declarations in a single statement is also possible:

```go
const country, code = "India", 91

const (
    employeeId string = "E101"
    salary float64 = 50000.0
)
```

Constants, as you would expect, cannot be changed. That is, you cannot re-assign a constant to a different value after it is initialized.

```go
const a = 123
a = 321 // Compiler error (Cannot assign to constant)
```

## Typed and Untyped Constants

Constants in Go are special. They work differently from how they work in other languages. To understand why they are special and how they exactly work, we need some background on Go's type system. So let's jump right into it.

### Background

Go is a statically typed programming language. Which means that the type of every variable is known or inferred by the compiler at compile time.

But it goes a step further with its type system and doesn’t even allow you to perform operations that mix numeric types. For example, you cannot add a `float64` variable to an `int`, or even an `int64` variable to an `int`.

```go
var myFloat float64 = 21.54
var myInt int = 562
var myInt64 int64 = 120

var res1 = myFloat + myInt // Not allowed (Compiler error)
var res2 = myInt + myInt64 // Not allowed (Compiler error)
```

For the above operations to work, you'll need to explicitly cast the variables so that all of them are of the same type.

```go
var res1 = myFloat + float64(myInt) // Works
var res2 = myInt + int(myInt64)     // Works
```

If you've worked with other statically typed languages like `C`, `C++` or `Java`, then you must be aware that they automatically convert smaller types to larger types whenever you mix them in any operation. For example, `int` can be automatically converted to `long`, `float` or `double`.

So the obvious question is that - why doesn't Go do the same? why doesn't it perform implicit type conversions like `C`, `C++` or `Java`?

And here is what Go designers have to say about this (Quoting from [Go's official doc](https://golang.org/doc/faq#conversions)).

> "The convenience of automatic conversion between numeric types in C is outweighed by the confusion it causes. When is an expression unsigned? How big is the value? Does it overflow? Is the result portable, independent of the machine on which it executes? It also complicates the compiler; "the usual arithmetic conversions" are not easy to implement and inconsistent across architectures. For reasons of portability, we decided to make things clear and straightforward at the cost of some explicit conversions in the code." ([Excerpt from Go's official doc](https://golang.org/doc/faq#conversions))

All right! So Go doesn't provide implicit type conversions and it requires us to do explicit type casting whenever we mix variables of multiple types in an operation.

But how does Go's type system work with constants? Given that all of the following statements are valid in Go.

```go
var myInt32 int32 = 10
var myInt int = 10
var myFloat64 float64 = 10
var myComplex complex64 = 10
```

What is the type of the constant value `10` in the above examples? Moreover, if there are no implicit type conversions in Go, then wouldn't we need to write the above statements like:

```go
var myInt32 int32 = int32(10)
var myFloat64 float64 = float64(10)
// etc...
```

Well, the answers to all theses questions lay in the way constants are handled in Go. So let's find out how they are handled.

### Untyped Constants

Any constant in Go, named or unnamed, is untyped unless given a type explicitly. For example, all of the following constants are untyped:

```go
1       // untyped integer constant
4.5     // untyped floating-point constant
true    // untyped boolean constant
"Hello" // untyped string constant
```

They are untyped even after you give them a name:

```go
const a = 1
const f = 4.5
const b = true
const s = "Hello"
```

Now, you might be wondering that I'm using terms like `integer` constant, `string` constant, and I'm also saying that they are untyped.

Well yes, the value `1` is an integer, `4.5` is a float, and `"Hello"` is a string. But they are just values. They are not given a fixed type yet, like `int32` or `float64` or `string`, that would force them to obey Go's strict type rules.

The fact that the value `1` is untyped allows us to assign it to any variable whose type is compatible with integers.

```go
var myInt int = 1
var myFloat float64 = 1
var myComplex complex64 = 1
```

**Note:** Although the value `1` is untyped, it is an untyped integer. So it can only be used where an integer is allowed. You cannot assign it to a `string` or a `boolean` variable for example. Similarly, an untyped floating-point constant like `4.5` can be used anywhere a floating-point value is allowed

```go
var myFloat32 float32 = 4.5
var myComplex64 complex64 = 4.5
```

**Let's now see an example of an untyped string constant**

In Go, you can create a type alias using the `type` keyword like so:

```go
type RichString string // Type alias of 'string'
```

Given the strongly typed nature of Go, you can't assign a `string` variable to a `RichString` variable.

```go
var myString string = "Hello"
var myRichString RichString = myString // Won't work
```

But, you can assign an untyped string constant to a `RichString` variable because it is compatible with strings.

```go
const myUntypedString = "Hello"
var myRichString RichString = myUntypedString // Works
```

### Constants and Type Inference: Default Type

Go supports type inference. That is, it can infer the type of a variable from the value that is used to initialize it. So you can declare a variable with an initial value, but without any type information, and Go will automatically determine the type.

```go
var a = 5 // Go compiler automatically infers the type of the variable 'a'
```

But how does it work? Given that constants in Go are untyped, what will be the type of the variable a in the above example? Will it be `int8` or `int16` or `int32` or `int64` or `int`?

Well, it turns out that every untyped constant in Go has a **default type**. The default type is used when we assign the constant to a variable that doesn't have any explicit type available.

Following are the default types for various constants in Go:

|         Constants          | Default Type |
| :------------------------: | :----------: |
|   integers (`10`, `76`)    |     int      |
|  floats (`3.14`, `7.92`)   |   float64    |
| complex numbers (`3 + 5i`) |  complex128  |
| characters (`'a'`, `'♠'`)  |     rune     |
| booleans (`true`, `false`) |     bool     |
|    strings (`"Hello"`)     |    string    |

So, in the statement `var a = 5`, since no explicit type information is available, the default type for integer constants is used to determine the type of `a`, which is `int`.

### Typed Constants

In Go, constants are typed when you explicitly specify the type in the declaration like this:

```go
const typedInt int = 1 // Typed constant
```

Just like variables, all the rules of Go's type system applies to typed constant. For example, you cannot assign a typed integer constant to a float variable:

```go
var myFloat64 float64 = typedInt // Compiler error
```

With typed constants, you lose all the flexibility that comes with untyped constants like assigning them to any variable of compatible type or mixing them in mathematical operations. So you should declare a type for a constant only if it's absolutely necessary. Otherwise, just declare constants without a type.

## Constant Expressions

The fact that constants are untyped (unless given a type explicitly) allows you to mix them in any expression freely.

So you can have a contant expression containing a mix of various untyped constants as long as those untyped constants are compatible with each other.

```go
const a = 5 + 7.5 // Valid
const b = 12 / 5  // Valid
const c = 'z' + 1 // Valid

const d = "Hey" + true // Invalid (untyped string constant and untyped boolean constant are not compatible with each other)
```

The evaluation of constant expressions and their result follows certain rules. Let's look at those rules:

- A comparison operation between two untyped constants always outputs an untyped boolean constant.

```go
const a = 7.5 > 5       // true (untyped boolean constant)
const b = "xyz" < "uvw" // false (untyped boolean constant)
```

- For any other operation (except shift):
  - If both the operands are of the same type (e.g, both are untyped integer constants), the result is also of the same type. For example, the expression `25 / 2` yields `12` not `12.5`. Since both the operands are untyped integers, the result is truncated to an integer.
  - If the operands are of different type, the result is of the operand's type that is broader as per the rule: `integer` < `rune` < `floating-point` < `complex`.

```go
const a = 25 / 2       // 12 (untyped integer constant)
const b = (6 + 8i) / 2 // (3 + 4i) (untyped complex constant)
```

- Shift operation rules are a bit complex. First of all, there are some requirements
  - The right operand of a shift expression must either have an unsigned integer type or be an untyped constant that can represent a value of type `uint`.
  - The left operand must either have an integer type or be an untyped constant that can represent a value of type `int`. **The rule** - If the left operand of a shift expression is an untyped constant, the result is an untyped integer constant; otherwise the result is of the same type as the left operand.

```go
const a = 1 << 5        // 32 (untyped integer constant)
const b = int32(1) << 4 // 16 (int32)
const c = 16.0 >> 2     // 4 (untyped integer constant) - 16.0 can represent a value of type 'int'
const d = 32 >> 3.0     // 4 (untyped integer constant) - 3.0 can represent a value of type 'uint'

const e = 10.50 << 2    // ILLEGAL (10.50 can't represent a value of type 'int')
const f = 64 >> -2      // ILLEGAL (The right operand must be an unsigned int or an untyped constant compatible with 'uint')
```
