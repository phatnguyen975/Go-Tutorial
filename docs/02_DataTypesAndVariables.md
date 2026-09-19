<div align="center">
  <h1>Data Types and Variables</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 19, 2026</sub>
</div>

## Table of Contents

1. [Go Is Statically and Strongly Typed](#1-go-is-statically-and-strongly-typed)
2. [Basic Data Types: Overview](#2-basic-data-types-overview)
3. [Boolean Type](#3-boolean-type)
4. [Integer Types](#4-integer-types)
5. [Floating-Point Types](#5-floating-point-types)
6. [Complex Number Types](#6-complex-number-types)
7. [The `string` Type](#7-the-string-type)
8. [`byte` and `rune`: Aliases with Meaning](#8-byte-and-rune-aliases-with-meaning)
9. [Declaring Variables: `var`](#9-declaring-variables-var)
10. [Declaring Variables: Short Declaration (`:=`)](#10-declaring-variables-short-declaration-)
11. [Zero Values](#11-zero-values)
12. [Multiple Variable Declaration](#12-multiple-variable-declaration)
13. [Variable Scope and Shadowing](#13-variable-scope-and-shadowing)
14. [The Blank Identifier](#14-the-blank-identifier)
15. [Constants: the `const` Keyword](#15-constants-the-const-keyword)
16. [Typed vs. Untyped Constants](#16-typed-vs-untyped-constants)
17. [`iota`: Auto-Incrementing Constants](#17-iota-auto-incrementing-constants)
18. [Constant Expressions and Their Limits](#18-constant-expressions-and-their-limits)
19. [Type Conversion: Go's Strict, Explicit Rules](#19-type-conversion-gos-strict-explicit-rules)
20. [Numeric Type Conversion](#20-numeric-type-conversion)
21. [String Conversions](#21-string-conversions)
22. [Converting Between Named Types](#22-converting-between-named-types)
23. [Type Conversion vs. Type Assertion](#23-type-conversion-vs-type-assertion)
24. [Common Mistakes and Pitfalls](#24-common-mistakes-and-pitfalls)
25. [Best Practices Summary](#25-best-practices-summary)
26. [Use Case Summary Table](#26-use-case-summary-table)
27. [References](#27-references)

## 1. Go Is Statically and Strongly Typed

Every variable in Go has a fixed type, determined either explicitly or by inference, **at compile time** — this is what "statically typed" means. Go is also **strongly typed**: unlike languages such as C, JavaScript, or Python, Go does **not** perform implicit conversions between different types, even between closely related ones (like `int` and `float64`, or `int32` and `int64`). Every conversion between distinct types must be written explicitly.

```go
var i int = 42
var f float64 = i // COMPILE ERROR: cannot use i (type int) as type float64
```

This strictness is a deliberate design choice that eliminates an entire class of subtle bugs common in more permissive languages, at the cost of a small amount of extra verbosity at conversion points.

## 2. Basic Data Types: Overview

Go's basic (predeclared) types fall into four broad families:

| Family                   | Types                                                                                              |
| ------------------------ | -------------------------------------------------------------------------------------------------- |
| Boolean                  | `bool`                                                                                             |
| Numeric (integer)        | `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `uintptr` |
| Numeric (floating-point) | `float32`, `float64`                                                                               |
| Numeric (complex)        | `complex64`, `complex128`                                                                          |
| Text                     | `string`                                                                                           |
| Aliases                  | `byte` (alias for `uint8`), `rune` (alias for `int32`)                                             |

## 3. Boolean Type

`bool` represents a truth value, either `true` or `false`. Its zero value is `false`.

```go
var isActive bool = true
var hasPermission bool // zero value: false
```

**No implicit conversion from other types.** Unlike C or JavaScript, no other type (integers included) can be used in a boolean context — `if 1 { ... }` simply does not compile in Go.

## 4. Integer Types

Go provides both **sized** integer types (with an explicit, fixed bit width) and two **platform-dependent** integer types:

| Type      | Size                                                                          | Range                                                                                          |
| --------- | ----------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| `int8`    | 8 bits                                                                        | -128 to 127                                                                                    |
| `int16`   | 16 bits                                                                       | -32,768 to 32,767                                                                              |
| `int32`   | 32 bits                                                                       | -2,147,483,648 to 2,147,483,647                                                                |
| `int64`   | 64 bits                                                                       | -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807                                        |
| `uint8`   | 8 bits                                                                        | 0 to 255                                                                                       |
| `uint16`  | 16 bits                                                                       | 0 to 65,535                                                                                    |
| `uint32`  | 32 bits                                                                       | 0 to 4,294,967,295                                                                             |
| `uint64`  | 64 bits                                                                       | 0 to 18,446,744,073,709,551,615                                                                |
| `int`     | platform-dependent (32 or 64 bits, effectively always 64 on modern platforms) | at least -2^31 to 2^31-1                                                                       |
| `uint`    | platform-dependent, matching `int`'s width                                    | at least 0 to 2^32-1                                                                           |
| `uintptr` | large enough to hold the bit pattern of any pointer                           | implementation-specific; used for low-level pointer arithmetic, rarely needed in ordinary code |

**`int` is the default, idiomatic choice** for ordinary integer values in Go, unless a specific bit width is genuinely required (for binary file formats, network protocols, or matching a specific external API's exact type). Its zero value, like every integer type, is `0`.

```go
var age int = 30
var smallNumber int8 = 100
var largeCounter uint64 = 18000000000
```

**Overflow wraps around silently** for fixed-width integer arithmetic — Go does not panic or error on integer overflow by default:

```go
var b int8 = 127
b++ // wraps around to -128, silently — no panic, no error
```

## 5. Floating-Point Types

Go provides two IEEE-754 floating-point types:

| Type      | Size    | Precision                          |
| --------- | ------- | ---------------------------------- |
| `float32` | 32 bits | ~7 decimal digits of precision     |
| `float64` | 64 bits | ~15-17 decimal digits of precision |

`float64` is the default, idiomatic choice for floating-point values unless memory constraints or interoperability with a specific `float32`-based API require otherwise.

```go
var price float64 = 19.99
var ratio float32 = 0.5
```

**Floating-point comparison caution:** due to how binary floating-point representation works, direct equality comparison (`==`) between floating-point values that result from computation can be unreliable — two mathematically equal values computed via different paths may differ in their least-significant bits. Comparing against a small tolerance (an "epsilon") is the common workaround when equality is genuinely needed after computation, rather than relying on exact `==`.

## 6. Complex Number Types

Go has built-in support for complex numbers, a feature relatively rare among mainstream general-purpose languages:

```go
var c1 complex64 = complex(2, 3)   // 2+3i, using two float32 components
var c2 complex128 = complex(1.5, 2.5) // 1.5+2.5i, using two float64 components

fmt.Println(real(c2)) // 1.5 — the real part
fmt.Println(imag(c2)) // 2.5 — the imaginary part
```

`complex64` is built from two `float32` values (real and imaginary parts); `complex128` from two `float64` values. These types see use primarily in scientific/numeric computing contexts (signal processing, certain mathematical algorithms) and are rarely encountered in typical application-level Go code.

## 7. The `string` Type

A `string` in Go is an **immutable sequence of bytes**, conventionally (but not enforced by the type system) holding UTF-8 encoded text. Its zero value is the empty string `""`.

```go
var name string = "Gopher"
greeting := "Hello, " + name // string concatenation with +
```

**Strings are immutable** — once created, a string's contents can never be changed in place; any "modification" (like uppercasing or replacing a substring) actually produces a brand-new string value, leaving the original unchanged:

```go
s := "hello"
// s[0] = 'H' // COMPILE ERROR: cannot assign to s[0] (strings are immutable)

upper := strings.ToUpper(s) // creates a NEW string; s itself is untouched
```

Indexing a string (`s[i]`) yields a single **byte**, not a character — for text that may contain multi-byte UTF-8 characters, iterating with `for range` (which decodes UTF-8 automatically) or using the `unicode/utf8` package is necessary for correct rune-by-rune handling, since a raw byte index can land in the middle of a multi-byte character.

Multi-line/raw strings use backticks instead of double quotes, and interpret no escape sequences at all (everything between the backticks, including literal newlines, is taken verbatim):

```go
raw := `Line one
Line two
No \n escape processing here`
```

## 8. `byte` and `rune`: Aliases with Meaning

Go defines two predeclared **aliases** for existing numeric types, specifically to make code that works with bytes and characters more self-documenting:

- **`byte`** is an alias for `uint8` — used when a value genuinely represents a single raw byte (e.g., a byte read from a file, a byte in a network buffer).
- **`rune`** is an alias for `int32` — used when a value represents a single Unicode code point (a "character," loosely speaking).

```go
var b byte = 'A'   // 65 — a single byte
var r rune = '世'    // 19990 — a Unicode code point, needs more than one byte in UTF-8

fmt.Printf("%T %T\n", b, r) // uint8 int32 — the underlying types, confirming they're aliases
```

Since these are true aliases (not distinct named types), `byte` and `uint8` (and `rune` and `int32`) are **completely interchangeable** everywhere — the compiler treats them as exactly the same type, and the alias is purely a readability aid for the programmer, signaling intent rather than creating a genuinely distinct type.

## 9. Declaring Variables: `var`

The `var` keyword declares a variable, optionally with an explicit type and/or an initial value:

```go
var age int = 30        // explicit type and initial value
var name = "Alice"       // type inferred from the value (string)
var count int             // no initial value — gets int's zero value, 0
```

`var` can be used both at the package level (outside any function) and inside a function body — it's the **only** form usable at the package level, since the short declaration form (`:=`, covered next) is restricted to inside function bodies.

```go
var appVersion = "1.0.0" // package-level variable — must use var, not :=

func main() {
    var localVar = 42 // var also works fine inside a function
}
```

### 9.1 Grouped `var` Declarations

Multiple `var` declarations can be grouped in a single block, which is commonly used for related package-level variables:

```go
var (
    host     = "localhost"
    port     = 8080
    timeout  = 30 * time.Second
)
```

## 10. Declaring Variables: Short Declaration (`:=`)

Inside a function body, the short variable declaration form combines declaration and initialization, with the type always inferred from the assigned value:

```go
func main() {
    age := 30           // int, inferred
    name := "Alice"       // string, inferred
    price := 19.99        // float64, inferred
}
```

**Restrictions on `:=`:**

- It can **only** be used inside a function body (including at the top level of `if`/`for`/`switch` init clauses) — never at the package level.
- **At least one** variable on the left-hand side must be genuinely new — `:=` can mix a new variable with an already-declared one in the same statement, as long as at least one is new (commonly seen with multiple-return-value calls):

```go
value, err := doSomething() // both new
result, err := doSomethingElse() // result is new; err is REUSED (assigned to), not redeclared
```

If every variable on the left is already declared in the _same scope_, `:=` produces a compile error (`no new variables on left side of :=`); if at least one is genuinely new, the already-declared ones are simply assigned to, not shadowed, as long as they're in the same scope block.

### 10.1 `var` vs. `:=`: When to Use Which

- `:=` is by far the most common, idiomatic choice inside function bodies — it's shorter and lets the type be inferred naturally from the initializer.
- `var` is necessary at the package level, and is also preferred when you want to declare a variable **without** an initial value (relying on its zero value), or when you want to be explicit about a type that differs from what would be inferred (e.g., declaring a `float64` from an integer literal, or an interface-typed variable holding a specific concrete value).

```go
var total float64 = 0 // explicit, even though the literal 0 would normally infer as int
```

## 11. Zero Values

Every variable declared without an explicit initial value is automatically set to its type's **zero value** — Go has no concept of an "uninitialized," garbage-containing variable the way some lower-level languages do.

| Type                                              | Zero value                                         |
| ------------------------------------------------- | -------------------------------------------------- |
| Numeric types (`int`, `float64`, etc.)            | `0` (or `0.0` for floats)                          |
| `bool`                                            | `false`                                            |
| `string`                                          | `""` (empty string)                                |
| Pointer, slice, map, channel, function, interface | `nil`                                              |
| Struct                                            | every field set to its own zero value, recursively |
| Array                                             | every element set to its type's zero value         |

```go
var i int       // 0
var f float64   // 0
var s string    // ""
var b bool      // false
var p *int      // nil
var sl []int    // nil
```

This design — every variable is automatically, predictably initialized — eliminates an entire category of "uninitialized variable" bugs common in languages that leave freshly-declared variables holding arbitrary memory contents.

## 12. Multiple Variable Declaration

Both `var` and `:=` support declaring several variables in a single statement:

```go
var x, y, z int = 1, 2, 3       // same type, shared declaration
var name, age = "Alice", 30     // different types, inferred individually

a, b := 10, 20                  // short form, multiple variables
a, b = b, a                     // idiomatic swap, no temporary variable needed
```

The swap idiom (`a, b = b, a`) works because Go evaluates the entire right-hand side **before** performing any of the assignments on the left, making a manual temporary variable for swapping unnecessary — this pattern is very commonly seen in sorting and array-manipulation code.

## 13. Variable Scope and Shadowing

A variable's scope is generally the block (delimited by `{ }`) it's declared in, extending from the point of declaration to the end of that block. An inner block can declare a variable with the **same name** as one in an outer scope — this is called **shadowing**, and it's legal, though it can be a source of subtle bugs if unintentional:

```go
x := 10
if true {
    x := 20 // shadows the outer x — this is a NEW, separate variable
    fmt.Println(x) // 20
}
fmt.Println(x) // 10 — the outer x was never touched
```

Shadowing is especially easy to trigger by accident with the `:=` operator inside a nested block (like an `if`'s init clause, or a new block scope) when the intent was actually to assign to the existing outer variable — see [Section 24.4](#244-accidental-shadowing-with-) for the common pitfall this causes.

## 14. The Blank Identifier

The underscore `_` is a special, write-only identifier — assigning to it discards a value entirely, without declaring an actual variable, and without triggering an "unused variable" compile error:

```go
_, err := someFunction() // discard the first return value, keep only err

for _, value := range mySlice { // discard the index, keep only the value
    fmt.Println(value)
}
```

The blank identifier can be used anywhere a value needs to be explicitly and deliberately ignored — multiple return values, `range` loop variables, or even a top-level `var _ SomeInterface = (*MyType)(nil)` compile-time interface-satisfaction check.

## 15. Constants: the `const` Keyword

A **constant** is a named value whose value is fixed at **compile time** and can never change during the program's execution. Constants are declared with the `const` keyword:

```go
const Pi = 3.14159
const AppName = "MyApp"
const MaxRetries = 3
```

Unlike variables, a constant has **no memory address** at runtime — the compiler inlines its value directly at every point it's used, rather than allocating storage for it the way a variable's value is stored. Because of this, you cannot take the address of a constant (`&Pi` is a compile error), and a constant's value must be something the compiler can fully evaluate during compilation.

### 15.1 Grouped Constants

Like `var`, constants can be grouped in a block:

```go
const (
    StatusPending  = "pending"
    StatusActive   = "active"
    StatusInactive = "inactive"
)
```

## 16. Typed vs. Untyped Constants

This is one of the more distinctive and genuinely useful aspects of Go's type system, and it has no close equivalent in most other statically typed languages.

### 16.1 Untyped Constants

A constant declared without an explicit type is **untyped** — the compiler assigns it a "default kind" (boolean, rune, integer, floating-point, complex, or string) based on its literal form, but the constant itself doesn't commit to one specific concrete type until it's actually used in a context that requires one:

```go
const MaxRetries = 3 // untyped integer constant

var a int32 = MaxRetries   // fine — 3 fits comfortably as an int32
var b int64 = MaxRetries   // also fine — the SAME untyped constant adapts to int64 here
var c float64 = MaxRetries // also fine — even converts to float64 without an explicit cast
```

The same untyped constant `MaxRetries` was used as an `int32`, an `int64`, and even a `float64` above, with **no explicit conversion needed anywhere** — this is only possible because it's untyped; the compiler implicitly converts it to whatever concrete type the context actually requires, as long as the value can be represented in that type without loss.

### 16.2 Typed Constants

A constant declared with an explicit type is **typed**, and behaves much more like a regular value of that type — it cannot be implicitly used as a different type without an explicit conversion, exactly like a regular typed variable:

```go
const TypedPi float64 = 3.14159

var x float32 = TypedPi // COMPILE ERROR: cannot use TypedPi (type float64) as type float32
var y float32 = float32(TypedPi) // OK — explicit conversion required
```

### 16.3 Default Types

When an untyped constant is used in a context with **no** explicit type requirement (like a `:=` short declaration with no other type information), it takes on a **default type**, following fixed rules: `bool` for boolean constants, `rune` for character-literal constants, `int` for integer constants, `float64` for floating-point constants, `complex128` for complex constants, and `string` for string constants.

```go
x := 3.14 // untyped float constant → defaults to float64 here, since := has no other type info
```

### 16.4 Precision

Untyped constants are, by specification, evaluated with significantly **higher precision** than any of Go's concrete numeric types provide at runtime — the specification requires at least 256 bits of precision for constant arithmetic. This means an expression like `1.0 / 3.0`, computed entirely as untyped constants, retains far more precision than a `float64` division performed at runtime, right up until the moment the constant is actually converted to a concrete type for use.

```go
const Precise = 1.0 / 3.0 // computed with extended constant-arithmetic precision
fmt.Printf("%.50f\n", Precise) // prints far more accurate digits than a runtime float64 division would
```

**Practical guidance:** untyped constants are generally the preferred, more flexible choice for most constant declarations, specifically _because_ of this ability to adapt to whatever type context they're used in — reach for an explicitly typed constant only when you specifically need to pin the value to one exact type (for instance, to satisfy a particular interface, or to prevent accidental use in an unintended numeric context).

## 17. `iota`: Auto-Incrementing Constants

`iota` is a predeclared identifier that provides a simple, elegant way to generate sequences of related constants, most commonly used to implement enumeration-like patterns (which Go doesn't have as a dedicated language feature the way some other languages do).

### 17.1 Basic `iota`

Inside a `const` block, `iota` starts at `0` and increments by `1` for each subsequent constant specification (each line) within that same block, resetting to `0` again in any new `const` block:

```go
const (
    Sunday = iota // 0
    Monday        // 1 — iota implicitly increments; the expression "= iota" is implicitly repeated
    Tuesday       // 2
    Wednesday     // 3
    Thursday      // 4
    Friday        // 5
    Saturday      // 6
)
```

Notice that only the first line explicitly writes `= iota` — Go implicitly repeats the previous line's expression for every subsequent constant in the same block that has no explicit value of its own, which is exactly what makes this sequence auto-increment.

### 17.2 Skipping Values with `_`

The blank identifier can be used within an `iota` block to deliberately skip a value in the sequence:

```go
const (
    _  = iota // 0, discarded
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                    // 1 << 20 = 1,048,576
    GB                    // 1 << 30 = 1,073,741,824
)
```

### 17.3 Expressions Involving `iota`

`iota` can be used within arbitrary constant expressions, not just as a bare value — a common pattern is bit-shifting to produce power-of-two flag values, or arithmetic to produce a custom sequence:

```go
const (
    Ldate = 1 << iota // 1  (binary: ...0001)
    Ltime              // 2  (binary: ...0010)
    Lmicroseconds       // 4  (binary: ...0100)
    Llongfile            // 8  (binary: ...1000)
)

const (
    SNo1 = iota + 1 // 1
    SNo2             // 2
    SNo3             // 3
)
```

### 17.4 `iota` Resets Per Block

Each new `const (...)` block resets `iota` back to `0` — it does not carry a running count across separate `const` blocks:

```go
const (
    A = iota // 0
    B        // 1
)

const (
    C = iota // 0 again — a fresh block
    D        // 1
)
```

**Note:** `iota`'s count still increments for every line in a `const` block, even lines that don't actually reference `iota` at all — if a block mixes `iota`-based and explicitly-valued constants, the position (line count) still advances `iota`'s underlying counter for subsequent lines that do use it.

## 18. Constant Expressions and Their Limits

Constants must be expressions the compiler can fully evaluate at **compile time** — this restricts what can appear in a constant declaration considerably compared to an ordinary variable:

```go
const x = 5 + 3              // fine — a compile-time-evaluable expression
const y = "hello" + " world" // fine — string concatenation of constants
// const z = time.Now()      // COMPILE ERROR: time.Now() is a runtime function call, not a constant expression
```

Only booleans, runes, strings, integers, floating-point, and complex-number values (and expressions built purely from other constants of those kinds) can be constants — structs, slices, maps, and results of most function calls cannot be, since none of those can be fully determined at compile time in the general case (with a small number of specific built-in exceptions, like `len()` applied to a string literal or an array type, which the compiler _can_ evaluate at compile time).

## 19. Type Conversion: Go's Strict, Explicit Rules

Go requires an **explicit conversion** to change a value's type — there is no implicit numeric promotion, no implicit widening or narrowing, and no automatic conversion between related types the way many other languages provide. The general syntax is `T(value)`, where `T` is the target type:

```go
var i int = 42
var f float64 = float64(i) // explicit conversion required
var u uint = uint(f)       // explicit conversion required
```

This strictness applies even between numeric types that seem "obviously" compatible, like `int` and `int64`:

```go
var a int32 = 10
var b int64 = 20
// c := a + b // COMPILE ERROR: mismatched types int32 and int64
c := int64(a) + b // must explicitly convert first
```

**Why Go is this strict:** implicit numeric conversions in other languages are a well-documented source of subtle bugs — silent precision loss, unexpected sign changes, and surprising overflow behavior that's easy to overlook when a conversion happens invisibly. Go's requirement to write every conversion explicitly makes every place where a value's representation actually changes visible directly in the source code.

## 20. Numeric Type Conversion

Converting between numeric types can lose information, and Go does **not** protect against this at compile time or runtime — it's the programmer's responsibility to convert safely when precision or range matters:

```go
var f float64 = 3.99
var i int = int(f) // truncates toward zero: i is 3, NOT rounded to 4

var big int64 = 300
var small int8 = int8(big) // overflow: wraps around silently, producing an unexpected value (44, due to wraparound)
```

**Float-to-integer conversion truncates** (discards the fractional part), it does **not** round — `int(3.99)` is `3`, and `int(-3.99)` is `-3` (truncation toward zero, not toward negative infinity). If rounding is actually needed, use `math.Round` (or `math.Floor`/`math.Ceil`) before converting to an integer type.

**Converting to a smaller integer type can silently overflow/wrap**, with no error or panic — always ensure a value genuinely fits within the target type's range before converting, if correctness depends on it, since Go performs no runtime bounds checking on this kind of conversion.

## 21. String Conversions

### 21.1 Numeric Type to String: A Common Trap

Converting an **integer** directly to a `string` using the `T(value)` syntax does **not** produce the decimal digits of that number as text — it instead interprets the integer as a single **Unicode code point** and produces the corresponding character:

```go
n := 65
s := string(n) // s is "A" — the character for code point 65, NOT "65"!
```

This is a well-known and frequently surprising trap for newcomers. To get the actual decimal text representation of a number, use the `strconv` package instead:

```go
import "strconv"

n := 65
s := strconv.Itoa(n) // s is "65" — correct decimal string representation
```

### 21.2 String to Numeric Types via `strconv`

Similarly, converting text like `"65"` into a numeric value requires `strconv`, not a direct `T(value)` conversion (which doesn't even compile for a string-to-int conversion in most forms):

```go
i, err := strconv.Atoi("65")   // parses "65" into the int 65
f, err := strconv.ParseFloat("3.14", 64) // parses "3.14" into the float64 3.14
```

Both return an error if the string isn't a validly formatted number for the target type — always check this error, since a malformed input string (not a valid number at all) is a very real, common failure mode.

### 21.3 String to `[]byte` and `[]rune`, and Back

```go
s := "hello"
b := []byte(s) // converts to a slice of the string's raw bytes
r := []rune(s)  // converts to a slice of the string's decoded Unicode code points

s2 := string(b) // converts back to a string
s3 := string(r) // also converts back to a string
```

These conversions **copy** the underlying data — a `[]byte` obtained this way is independent of the original string (necessary, since strings are immutable but byte slices are mutable) and modifying it has no effect on the original string.

## 22. Converting Between Named Types

Go allows explicit conversion between any two types that share the same **underlying type** — this includes converting between a user-defined named type and its underlying basic type, or between two different named types built on the same underlying type:

```go
type Celsius float64
type Fahrenheit float64

c := Celsius(100)
f := Fahrenheit(c) // valid: both have the same underlying type, float64

var plainFloat float64 = float64(c) // converting a named type back to its underlying basic type
```

This is a common, idiomatic pattern for giving numeric or string values additional type safety and meaning (preventing, say, accidentally passing a `Celsius` value where a `Fahrenheit` value was expected, since they're distinct types even though both are ultimately `float64`), while still allowing explicit, deliberate conversion when actually needed.

## 23. Type Conversion vs. Type Assertion

These two concepts are easy to confuse by name, but serve entirely different purposes:

|                  | Type Conversion                                                                                                                                 | Type Assertion                                                                                                            |
| ---------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| Syntax           | `T(value)`                                                                                                                                      | `value.(T)`                                                                                                               |
| Applies to       | Any two types with a compatible underlying representation (numeric types, string/byte/rune conversions, named types sharing an underlying type) | Specifically extracting the concrete type stored inside an **interface** value                                            |
| What it does     | Genuinely reinterprets/converts a value from one concrete type to another                                                                       | Checks (and retrieves) what concrete type an interface value is actually holding — it does not "convert" anything         |
| Failure behavior | A compile-time error if the types are fundamentally incompatible                                                                                | A runtime panic (or a `false` "ok" result, with the safe two-result form) if the interface doesn't hold the asserted type |

```go
// Type conversion — genuinely changes representation
var i int = 42
var f float64 = float64(i)

// Type assertion — extracts a concrete type from an interface
var any1 any = "hello"
s, ok := any1.(string) // checks what concrete type any1 actually holds
```

## 24. Common Mistakes and Pitfalls

### 24.1 Expecting Implicit Numeric Conversion

```go
var a int32 = 5
var b int64 = 10
// c := a + b // COMPILE ERROR: mismatched types
```

Every conversion between distinct numeric types must be explicit — see [Section 19](#19-type-conversion-gos-strict-explicit-rules).

### 24.2 Converting an Integer to a String and Getting a Character, Not Digits

```go
n := 65
s := string(n) // "A", not "65" — a very common beginner surprise
```

Use `strconv.Itoa` for the decimal text representation of a number — see [Section 21.1](#211-numeric-type-to-string-a-common-trap).

### 24.3 Assuming Float-to-Int Conversion Rounds

```go
f := 3.99
i := int(f) // 3, not 4 — truncation, not rounding
```

Use `math.Round(f)` before converting to an integer type if rounding (not truncation) is what's actually needed.

### 24.4 Accidental Shadowing with `:=`

```go
err := doSomething()
if err != nil {
    // handle
}

if value, err := doSomethingElse(); err != nil { // this err SHADOWS the outer err, in this if's scope
    // handle — this looks like it's reusing the outer err, but it's actually a new one
}
// the outer err is completely untouched by the if block above
```

Be careful with `:=` inside a nested block when the intent is to assign to an already-declared outer variable — inside a new block scope, `:=` always creates fresh variables for every name on its left side that isn't declared _in that exact same block_, even if a same-named variable exists in an outer scope.

### 24.5 Silent Integer Overflow on Narrowing Conversion

```go
var big int64 = 1000
var small int8 = int8(big) // silently wraps to an unexpected value — no error, no panic
```

Go performs no runtime bounds checking on numeric conversions — verify a value fits the target type's range before converting, when correctness depends on it.

### 24.6 Forgetting Constants Have No Memory Address

```go
const Pi = 3.14159
// p := &Pi // COMPILE ERROR: cannot take the address of a constant
```

A constant is inlined at compile time and has no runtime storage location — if you need an addressable value, use a `var` instead.

### 24.7 Assuming Constants Must Be Explicitly Typed to Be Useful

```go
const MaxItems int = 100 // works, but unnecessarily rigid

var smallCount int32 = MaxItems // COMPILE ERROR if MaxItems is explicitly typed int, not int32
```

Leaving a constant untyped (`const MaxItems = 100`) is often more flexible, letting the same constant adapt to whatever numeric type context it's used in — see [Section 16](#16-typed-vs-untyped-constants).

## 25. Best Practices Summary

1. **Prefer `:=` inside function bodies** for its brevity and automatic type inference; reserve `var` for package-level declarations, zero-value-only declarations, or when you need an explicit type differing from what would be inferred.
2. **Default to `int` for integers and `float64` for floating-point values** unless a specific bit width is genuinely required for interoperability or memory constraints.
3. **Leave constants untyped unless you have a specific reason to pin their type** — untyped constants are more flexible and adapt naturally to whatever numeric context they're used in.
4. **Use `iota` for sequential, related constant groups** (status codes, day-of-week enumerations, bit flags), rather than manually numbering each one.
5. **Use `strconv`, not a direct type conversion, for converting between numbers and their decimal text representations** — a direct `string(n)` conversion produces a character, not digits.
6. **Always check the error returned by `strconv` parsing functions** — malformed input is a real, common failure mode, not an edge case to ignore.
7. **Verify a value fits the target type's range before a narrowing numeric conversion**, since Go performs no runtime overflow checking on these conversions.
8. **Use `math.Round`/`Floor`/`Ceil` before converting a float to an integer** if rounding behavior (rather than truncation) is actually needed.
9. **Watch for accidental shadowing** when using `:=` inside a nested block with a variable name that also exists in an outer scope.
10. **Give numeric or string values additional type safety through named types** (`type Celsius float64`) when mixing conceptually distinct values of the same underlying representation could otherwise lead to bugs.

## 26. Use Case Summary Table

| Technique                       | When to Use                                                                    | Example Scenario                                                                               |
| ------------------------------- | ------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------- |
| `var x T`                       | Package-level declarations, or declaring with only a zero value                | `var count int`, `var appName = "MyApp"` at package scope                                      |
| `x := value`                    | Local variable declarations inside a function body                             | The overwhelming majority of everyday local variable declarations                              |
| Untyped `const`                 | A constant value that should flexibly adapt to different numeric contexts      | `const MaxRetries = 3` used as both `int32` and `int64`                                        |
| Typed `const`                   | A constant that must be pinned to one specific type                            | Satisfying a specific interface, or a value that must never be misused as another numeric type |
| `iota`                          | Sequential, related constant groups                                            | Enumerated status codes, days of the week, bit-flag values                                     |
| `strconv.Itoa` / `strconv.Atoi` | Converting between an integer and its decimal text form                        | Formatting a number for display, or parsing a numeric CLI argument                             |
| `T(value)` numeric conversion   | Explicitly converting between compatible numeric types                         | Converting an `int` result into a `float64` for a division                                     |
| `[]byte(s)` / `[]rune(s)`       | Getting mutable access to a string's raw bytes or decoded characters           | Text processing that needs to manipulate individual bytes/runes                                |
| Named type conversion           | Adding type safety to a basic type while keeping explicit conversion available | `type UserID int`, `type Celsius float64`                                                      |
| Blank identifier `_`            | Deliberately discarding a value without declaring a variable                   | `_, err := f()`, `for _, v := range s`                                                         |

## 27. References

1. Go Team — _A Tour of Go: Basic types, Variables, Constants_. https://go.dev/tour/basics/11
2. Go Team — _Effective Go_ (Constants, Variables sections). https://go.dev/doc/effective_go
3. Go Language Specification — _Types, Constants, Variables, Conversions, Iota_. https://go.dev/ref/spec
4. Go standard library documentation — package `strconv`. https://pkg.go.dev/strconv
5. Go standard library documentation — package `math` (`Round`, `Floor`, `Ceil`). https://pkg.go.dev/math
6. Go101 — _Constants and Variables_. https://go101.org/article/constants-and-variables.html
7. OneUptime Engineering Blog — _How to Understand Constant Types and Untyped Constants in Go_. https://oneuptime.com/blog/post/2026-01-23-go-constant-types/view
8. env.dev — _Constants in Go: const, iota & Type-Safe Enumerations_. https://env.dev/guides/go-constants
9. Application Architect — _Go Constants: const and iota Explained_. https://www.application-architect.com/posts/go-constants-const-and-iota-explained/
10. Devpriya Shivani — _Go Constants_, Medium. https://medium.com/@dpshiv12/go-constants-fa29270bd170
11. Dataplexa — _807-Go – Lesson 7: Constants & iota_. https://dataplexa.com/807-go-lesson-7-constants-and-iota/
