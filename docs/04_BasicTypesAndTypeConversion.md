<div align="center">
  <h1>Basic Types and Type Conversion</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 16, 2026</sub>
</div>

Go is a statically typed programming language. Every variable in Go has an associated type.

Data types classify a related set of data. They define how the data is stored in memory, what are the possible values that a variable of a particular data type can hold, and the operations that can be done on them.

Go has several built-in data types for representing common values like numbers, booleans, strings etc. In this chapter, We'll look at all these basic data types one by one and understand how they work.

## Numeric Types

Numeric types are used to represent numbers. They can be classified into Integers and Floating point types.

### Integers

Integers are used to store whole numbers. Go has several built-in integer types of varying size for storing signed and unsigned integers.

**Signed Integers**

|  Type   |        Size        |           Range           |
| :-----: | :----------------: | :-----------------------: |
| `int8`  |       8 bits       |        -128 to 127        |
| `int16` |      16 bits       | $-2^{15}$ to $2^{15} - 1$ |
| `int32` |      32 bits       | $-2^{31}$ to $2^{31} - 1$ |
| `int64` |      64 bits       | $-2^{63}$ to $2^{63} - 1$ |
|  `int`  | Platform dependent |    Platform dependent     |

The size of the generic `int` type is platform dependent. It is 32-bits wide on a 32-bit system and 64-bits wide on a 64-bit system.

**Unsigned Integers**

|   Type   |        Size        |       Range        |
| :------: | :----------------: | :----------------: |
| `uint8`  |       8 bits       |      0 to 255      |
| `uint16` |      16 bits       | 0 to $2^{16} - 1$  |
| `uint32` |      32 bits       | 0 to $2^{32} - 1$  |
| `uint64` |      64 bits       | 0 to $2^{64} - 1$  |
|  `uint`  | Platform dependent | Platform dependent |

The size of `uint` type is platform dependent. It is 32-bits wide on a 32-bit system and 64-bits wide on a 64-bit system.

**Tip:** When you are working with integer values, you should always use the `int` data type unless you have a good reason to use the sized or unsigned integer types.

In Go, you can declare octal numbers using prefix `0` and hexadecimal numbers using the prefix `0x` or `0X`. Following is a complete example of integer types:

```go
package main

import "fmt"

func main() {
    var myInt8 int8 = 97

    /*
      When you don't declare any type explicitly, the type inferred is 'int'
      (The default type for integers)
    */
    var myInt = 1200

    var myUint uint = 500

    var myHexNumber = 0xFF  // Use prefix '0x' or '0X' for declaring hexadecimal numbers
    var myOctalNumber = 034 // Use prefix '0' for declaring octal numbers

    fmt.Printf("%d, %d, %d, %#x, %#o\n", myInt8, myInt, myUint, myHexNumber, myOctalNumber)
}
```

**Integer Type Aliases**

Go has two additional integer types called `byte` and `rune` that are aliases for `uint8` and `int32` data types respectively.

In Go, the `byte` and `rune` data types are used to distinguish characters from integer values.

Go doesn't have a `char` data type. It uses `byte` and `rune` to represent character values. The `byte` data type represents ASCII characters and the `rune` data type represents a more broader set of Unicode characters that are encoded in UTF-8 format.

Characters are expressed in Go by enclosing them in single quotes like this `'A'`.

The default type for character values is `rune`. That means, if you don’t declare a type explicitly when declaring a variable with a character value, then Go will infer the type as `rune`.

```go
var firstLetter = 'A' // Type inferred as 'rune' (Default type for character values)
```

You can create a `byte` variable by explicitly specifying the type:

```go
var lastLetter byte = 'Z'
```

Both `byte` and `rune` data types are essentially integers. For example, a `byte` variable with value `'a'` is converted to the integer `97`.

Similarly, a `rune` variable with a unicode value `'♥'` is converted to the corresponding unicode codepoint `U+2665`, where `U+` means unicode and the numbers are hexadecimal, which is essentially an integer.

```go
package main

import "fmt"

func main() {
    var myByte byte = 'a'
    var myRune rune = '♥'

    fmt.Printf("%c = %d and %c = %U", myByte, myByte, myRune, myRune)
}
```

In the above example, I've printed the variable `myByte` in character and decimal format, and the variable `myRune` in character and Unicode format.

### Floating Point Types

Floating point types are used to store numbers with a decimal component (ex - 1.24, 4.50000). Go has two floating point types - `float32` and `float64`.

- `float32` occupies 32 bits in memory and stores values in single-precision floating point format.
- `float64` occupies 64 bits in memory and stores values in double-precision floating point format.

The default type for floating point values is `float64`. So when you initialize a floating point variable with an initial value without specifying a type explicitly, the compiler will infer the type as `float64`.

```go
var a = 9715.635 // Type inferred as 'float64' (Default type for floating-point numbers)
```

### Operations on Numeric Types

Go provides several operators for performing operations on numeric types -

- **Arithmetic Operators:** `+`, `-`, `*`, `/`, `%`
- **Comparison Operators:** `==`, `!=`, `<`, `>`, `<=`, `>=`
- **Bitwise Operators:** `&`, `|`, `^`, `<<`, `>>`
- **Increment and Decrement Operators:** `++`, `--`
- **Assignment Operators:** `+=`, `-=`, `*=`, `/=`, `%=`, `<<=`, `>>=`, `&=`, `|=`, `^=`

Here is an example demonstrating some of the above operators.

```go
package main

import (
    "fmt"
    "math"
)

func main() {
    var a, b = 4, 5
    var res1 = (a + b) * (a + b)/2 // Arithmetic operations

    a++ // Increment a by 1

    b += 10 // Increment b by 10

    var res2 = a ^ b // Bitwise XOR

    var r = 3.5
    var res3 = math.Pi * r * r // Operations on floating-point type

    fmt.Printf("res1 : %v, res2 : %v, res3 : %v", res1, res2, res3)
}
```

## Booleans

Go provides a data type called `bool` to store boolean values. It can have two possible values - `true` and `false`.

```go
var myBoolean = true
var anotherBoolean bool = false
```

### Operations on Boolean Types

You can use the following operators on boolean types.

- **Logical Operators:**
  - `&&`: (logical conjunction, "and")
  - `||`: (logical disjunction, "or")
  - `!`: (logical negation)
- **Equality and Inequality:** `==`, `!=`

The operators `&&` and `||` follow short-circuiting rules. That means, in the expression `E1 && E2`, if `E1` evaluates to `false` then `E2` won't be evaluated. Similarly, in the expression `E1 || E2`, if `E1` evaluates to `true` then `E2` won't be evaluated.

Here is an example of Boolean types.

```go
package main

import "fmt"

func main() {
    var truth = 3 <= 5
    var falsehood = 10 != 10

    // Short Circuiting
    var res1 = 10 > 20 && 5 == 5 // Second operand is not evaluated since first evaluates to false
    var res2 = 2*2 == 4 || 10%3 == 0 // Second operand is not evaluated since first evaluates to true

    fmt.Println(truth, falsehood, res1, res2)
}
```

## Complex Numbers

Complex numbers are one of the basic types in Go. Go has two complex types of different sizes:

- `complex64`: Both real and imaginary parts are of `float32` type.
- `complex128`: Both real and imaginary parts are of `float64` type.

The default type for a complex number in Go is `complex128`. You can create a complex number like this.

```go
var x = 5 + 7i // Type inferred as 'complex128'
```

Go also provides a built-in function named `complex` for creating complex numbers. If you're creating a complex number with variables instead of literals, then you'll need to use the `complex` function.

```go
var a = 3.57
var b = 6.23

// var c = a + bi won't work. Create the complex number like this.
var c = complex(a, b)
```

**Note:** Both real and imaginary parts of the complex number must be of the same floating point type. If you try to create a complex number with different real and imaginary part types, then the compiler will throw an error.

```go
var a float32 = 4.92
var b float64 = 7.38

/*
    The following statement won't compile.
    (Both real and imaginary parts must be of the same floating-point type)
*/
var c = complex(a, b) // Compiler error
```

### Operations on Complex Numbers

You can perform arithmetic operations like addition, subtraction, multiplication, and division on complex numbers.

```go
package main

import "fmt"

func main() {
    var a = 3 + 5i
    var b = 2 + 4i

    var res1 = a + b
    var res2 = a - b
    var res3 = a * b
    var res4 = a / b

    fmt.Println(res1, res2, res3, res4)
}
```

## Strings

In Go, a string is a sequence of bytes.

- Is **immutable** (cannot be changed after creation).
- Is a **read-only slice of bytes**.
- Commonly contains **UTF-8 encoded text**.

Strings in Go are declared either using double quotes as in `"Hello World"` or back ticks as in `` `Hello World` ``.

```go
// Normal String (Can not contain newlines, and can have escape characters like '\n', '\t', etc.)
var name = "Steve Jobs"

// Raw String (Can span multiple lines and escape characters are not interpreted)
var bio = `Steve Jobs was an American entrepreneur and inventor.
           He was the CEO and co-founder of Apple Inc.`
```

Double-quoted strings cannot contain newlines and they can have escape characters like `\n`, `\t`, etc. In double-quoted strings, a `\n` character is replaced with a newline, and a `\t` character is replaced with a tab space, and so on.

Strings enclosed within back ticks are raw strings. They can span multiple lines. Moreover, escape characters don't have any special meaning in raw strings.

```go
package main

import "fmt"

func main() {
    var website = "\thttps://www.callicoder.com\t\n"

    var siteDescription = `\t\tCalliCoder is a programming blog where you can find
                            practical guides and tutorials on programming languages,
                            web development, and desktop app development.\t\n`

    fmt.Println(website, siteDescription)
}
```

## Type Conversion

Go has a strong type system. It doesn’t allow you to mix numeric types in an expression. For example, You cannot add an `int` variable to a `float64` variable or even an `int` variable to an `int64` variable. You cannot even perform an assignment between mixed types.

```go
var a int64 = 4
var b int = a // Compiler error (Cannot use a (type in64) as type int in assignment)

var c int = 500

var result = a + c // Compiler error (Invalid Operation: mismatched types int64 and int)
```

Unlike other statically typed languages like `C`, `C++`, and `Java`, Go doesn't provide any implicit type conversion. To learn why Go is designed this way, check out the next chapter - Working with Constants in Go.

All right! So we cannot add, subtract, compare or perform any kind of operation on two different types even if they are numeric. But what to do if we need to perform such operations?

Well, you'll need to explicitly cast the variables to the target type.

```go
var a int64 = 4
var b int = int(a) // Explicit Type Conversion

var c float64 = 6.5

// Explicit Type Conversion
var result = float64(b) + c // Works
```

The general syntax for converting a value `v` to a type `T` is `T(v)`. Here are few more examples.

```go
var myInt int = 65
var myUint uint = uint(myInt)
var myFloat float64 = float64(myInt)
```
