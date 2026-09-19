<div align="center">
  <h1>Arrays</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 18, 2026</sub>
</div>

## Table of Contents

1. [What Is an Array?](#1-what-is-an-array)
2. [Declaring and Initializing Arrays](#2-declaring-and-initializing-arrays)
3. [The Length Is Part of the Type](#3-the-length-is-part-of-the-type)
4. [Accessing and Modifying Elements](#4-accessing-and-modifying-elements)
5. [Arrays Are Value Types](#5-arrays-are-value-types)
6. [Passing Arrays to Functions](#6-passing-arrays-to-functions)
7. [Array Comparability](#7-array-comparability)
8. [Iterating Over an Array](#8-iterating-over-an-array)
9. [Multidimensional Arrays](#9-multidimensional-arrays)
10. [Pointers to Arrays](#10-pointers-to-arrays)
11. [Array Literals with Indexed Elements](#11-array-literals-with-indexed-elements)
12. [Converting Between Arrays and a View of Their Elements](#12-converting-between-arrays-and-a-view-of-their-elements)
13. [Why Arrays Are Rarely Used Directly in Everyday Go Code](#13-why-arrays-are-rarely-used-directly-in-everyday-go-code)
14. [When Arrays Are Actually the Right Choice](#14-when-arrays-are-actually-the-right-choice)
15. [Common Mistakes and Pitfalls](#15-common-mistakes-and-pitfalls)
16. [Best Practices Summary](#16-best-practices-summary)
17. [Use Case Summary Table](#17-use-case-summary-table)
18. [References](#18-references)

## 1. What Is an Array?

An **array** in Go is a fixed-length, ordered sequence of elements, all of the same type. Once declared, an array's length can never change — it is baked into the array's type itself.

```go
var numbers [5]int // an array of exactly 5 ints
```

The type `[5]int` and the type `[10]int` are two entirely different, incompatible types in Go's type system, even though both hold `int` elements — this is one of the most distinctive things about Go arrays compared to arrays in many other languages, and it's the root cause of why arrays are used far less often in everyday Go code than the more flexible slice type (a separate, related type built on top of arrays).

## 2. Declaring and Initializing Arrays

### 2.1 Zero-Valued Declaration

```go
var scores [3]int
fmt.Println(scores) // [0 0 0] — every element gets the element type's zero value
```

### 2.2 Array Literal

```go
primes := [5]int{2, 3, 5, 7, 11}
```

### 2.3 Letting the Compiler Count the Elements (`...`)

```go
primes := [...]int{2, 3, 5, 7, 11} // the compiler infers the length (5) from the literal
fmt.Println(len(primes)) // 5
```

The `[...]` form is convenient when you want the array literal's own length to define the array's size, without needing to count and write the number yourself (and without risking a mismatch if you add or remove elements later without updating a hardcoded count).

## 3. The Length Is Part of the Type

This is the single most important fact to understand about Go arrays: **the array's length is part of its type**, not just a runtime property of a given value.

```go
var a [3]int
var b [5]int

// a = b // COMPILE ERROR: cannot use b (type [5]int) as type [3]int
```

`[3]int` and `[5]int` are as distinct from each other, type-wise, as `int` and `string` are — you cannot assign one to the other, and a function declared to accept a `[3]int` parameter cannot be called with a `[5]int` argument, even though both are "arrays of int." This rigidity is precisely why Go's more flexible **slice** type exists as a separate, dynamically-sized construct layered on top of arrays — slices decouple "how many elements right now" from the type itself.

## 4. Accessing and Modifying Elements

Elements are accessed with zero-based indexing, exactly as in most C-family languages:

```go
primes := [5]int{2, 3, 5, 7, 11}

fmt.Println(primes[0]) // 2
primes[0] = 100
fmt.Println(primes)    // [100 3 5 7 11]
```

Accessing an index outside the array's valid range (`0` to `len(array)-1`) causes an immediate runtime panic:

```go
var a [3]int
_ = a[5] // panic: runtime error: index out of range [5] with length 3
```

When the index is a compile-time constant and clearly out of bounds for a fixed-size array, the Go compiler can even catch this as a **compile-time error** rather than waiting for a runtime panic, since the array's size is statically known.

## 5. Arrays Are Value Types

Unlike slices, maps, and channels, an array in Go is a genuine **value type** — assigning an array, or passing it as a function argument, copies **every single element**. There is no shared underlying storage between the original and the copy.

```go
a := [3]int{1, 2, 3}
b := a       // b is a completely independent COPY of a
b[0] = 100

fmt.Println(a) // [1 2 3] — unaffected
fmt.Println(b) // [100 2 3]
```

This full-copy behavior is consistent with how Go treats structs (another value type) and is fundamentally different from arrays/lists in many other languages, where assignment often just copies a reference to the same underlying data.

## 6. Passing Arrays to Functions

Because arrays are value types, passing an array to a function copies the **entire array**, regardless of its size — a function receiving an array parameter gets its own independent copy, and any modification inside the function has no effect on the caller's original array:

```go
func zeroOut(arr [5]int) {
    for i := range arr {
        arr[i] = 0 // modifies only the local copy
    }
}

func main() {
    nums := [5]int{1, 2, 3, 4, 5}
    zeroOut(nums)
    fmt.Println(nums) // [1 2 3 4 5] — unchanged
}
```

If a function needs to mutate the caller's actual array, it must take a **pointer to the array** instead (see [Section 10](#10-pointers-to-arrays)):

```go
func zeroOutInPlace(arr *[5]int) {
    for i := range arr {
        arr[i] = 0 // Go automatically dereferences arr for indexing
    }
}

nums := [5]int{1, 2, 3, 4, 5}
zeroOutInPlace(&nums)
fmt.Println(nums) // [0 0 0 0 0]
```

**Performance implication:** for large arrays, passing by value copies a potentially significant amount of data on every function call — this copying cost is one of the practical reasons arrays are used sparingly for anything beyond small, fixed-size data in performance-conscious code.

## 7. Array Comparability

Unlike slices and maps, arrays **can** be compared directly with `==`/`!=`, as long as their element type is itself comparable. Two arrays are equal if they have the same length (which, since length is part of the type, is only possible if they're the same array type to begin with) and every corresponding element is equal:

```go
a := [3]int{1, 2, 3}
b := [3]int{1, 2, 3}
c := [3]int{1, 2, 4}

fmt.Println(a == b) // true — same type, same elements
fmt.Println(a == c) // false — differs at index 2
```

This makes arrays usable as map keys (since map keys must be comparable) in situations where a fixed-size collection of values needs to serve as a composite key — something a slice, which is never comparable with `==`, cannot do directly.

```go
type Coord3D [3]int // a fixed 3-element array type

visited := make(map[Coord3D]bool)
visited[Coord3D{1, 2, 3}] = true
```

## 8. Iterating Over an Array

`for range` works on arrays exactly as it does on slices:

```go
primes := [5]int{2, 3, 5, 7, 11}

for i, v := range primes {
    fmt.Println(i, v)
}

for _, v := range primes { // values only
    fmt.Println(v)
}
```

**Important:** because `range` over an array (not a pointer to one) iterates over a **copy** of the array by default in older, pre-1.22 semantics of the loop variable's own scope, mutating the value `v` inside the loop body never affects the original array — you'd need to index back into the array explicitly (`primes[i] = ...`) to actually change it:

```go
for _, v := range primes {
    v = v * 2 // has no effect on primes itself — v is just a loop-local copy of each element
}
fmt.Println(primes) // unchanged

for i := range primes {
    primes[i] = primes[i] * 2 // this DOES modify the original array
}
```

## 9. Multidimensional Arrays

Go supports arrays of arrays, which is how multidimensional, grid-like fixed structures are represented:

```go
var grid [3][3]int // a 3x3 grid of ints, all zero-valued initially

grid[1][1] = 5 // set the center element

for _, row := range grid {
    fmt.Println(row)
}
```

A multidimensional array literal nests accordingly:

```go
board := [2][3]int{
    {1, 2, 3},
    {4, 5, 6},
}
```

Each dimension's size is, as always, part of the overall type — `[3][3]int` and `[2][3]int` are different types from each other, following the same rigid-length rule described in [Section 3](#3-the-length-is-part-of-the-type).

## 10. Pointers to Arrays

A pointer to an array (`*[N]T`) allows a function to operate on the original array without copying it, while still benefiting from the array's fixed-size guarantees:

```go
func sumAndDouble(arr *[5]int) int {
    sum := 0
    for i := range arr { // Go automatically dereferences arr here
        sum += arr[i]
        arr[i] *= 2 // mutates the original array
    }
    return sum
}

nums := [5]int{1, 2, 3, 4, 5}
total := sumAndDouble(&nums)
fmt.Println(total) // 15 (the sum before doubling)
fmt.Println(nums)   // [2 4 6 8 10] — mutated in place
```

Just as with pointers to structs, Go automatically dereferences a pointer to an array for indexing (`arr[i]`) and for `range`, so the syntax reads identically to working with the array value directly — you don't need to manually write `(*arr)[i]`.

## 11. Array Literals with Indexed Elements

Go allows specifying particular indices explicitly in an array literal, leaving all unspecified indices at their zero value — useful for sparse initialization or for making certain positions' significance explicit:

```go
// An array of 100 elements where only index 0 and 99 are set; everything else is 0
var lookup = [100]int{0: -1, 99: -1}

// Combined with the ... length inference: length becomes (highest index + 1)
days := [...]string{0: "Sunday", 6: "Saturday"}
fmt.Println(len(days)) // 7
```

## 12. Converting Between Arrays and a View of Their Elements

An array can be **sliced**, producing a value that shares the array's underlying storage (this operation produces a distinct type — a slice — built on top of the array, rather than another array):

```go
arr := [5]int{1, 2, 3, 4, 5}
view := arr[1:4] // a slice, sharing arr's underlying storage; view == [2 3 4]

view[0] = 100
fmt.Println(arr) // [1 100 3 4 5] — the underlying array was mutated through the slice
```

Since Go 1.20, it's also possible to convert an array **pointer** to a slice of the same element type explicitly, and Go 1.17+ already supported converting a slice to an array pointer (`(*[4]int)(mySlice)`), which panics if the slice is shorter than the target array length. These conversions exist mainly for interoperating with APIs that expect one form or the other, and are less commonly needed in typical application code.

## 13. Why Arrays Are Rarely Used Directly in Everyday Go Code

In practice, the overwhelming majority of "list of things" needs in Go are met with **slices**, not arrays, for a few concrete reasons rooted in everything covered above:

- **Fixed size baked into the type** makes arrays inflexible for data whose length varies at runtime or isn't known in advance — which describes most real-world collections (user input, query results, parsed data).
- **Full-copy value semantics** mean passing a large array around is comparatively expensive, and easy to do by accident (forgetting a `&` when a pointer was intended).
- The standard library, and idiomatic Go generally, is built overwhelmingly around slices for representing sequences — functions like `append`, and most standard-library APIs dealing with sequences of data, are written in terms of slices, not arrays.

Arrays are not a "worse slice," though — they are a distinct, specialized tool that slices are actually built on top of internally (a slice's underlying storage is, conceptually, an array). Arrays earn their place specifically in situations where their unique properties (fixed compile-time size, full value semantics, direct comparability) are exactly what's needed, as covered next.

## 14. When Arrays Are Actually the Right Choice

- **A fixed-size, unchanging collection known at compile time**, such as the days of the week, a fixed set of RGB channels (`[3]byte`), or a cryptographic hash's fixed-size output (`[32]byte` for a SHA-256 digest).
- **A composite map key** that needs value-comparability and a small, fixed number of components — arrays support `==`, letting them serve as map keys, unlike slices (see [Section 7](#7-array-comparability)).
- **Value semantics are genuinely desired**, such as wanting a guaranteed independent copy every time a value is passed around, without any risk of aliasing/shared mutation through an underlying reference.
- **Interfacing with APIs, protocols, or hardware-level data** that inherently deal in a fixed number of bytes/elements (network packet headers, fixed-width binary formats, certain cryptography APIs).
- **Small, stack-friendly, fixed buffers** in performance-sensitive code, where avoiding a heap allocation (which a slice's backing array typically requires once it escapes) by using a small, stack-allocated fixed-size array can be a deliberate optimization.

## 15. Common Mistakes and Pitfalls

### 15.1 Expecting an Array Parameter to Behave Like a Reference

```go
func modify(arr [5]int) {
    arr[0] = 999 // only changes the local copy
}

nums := [5]int{1, 2, 3, 4, 5}
modify(nums)
fmt.Println(nums) // [1 2 3 4 5] — unaffected
```

Fix: pass a pointer (`*[5]int`) if the function needs to mutate the caller's array — see [Section 6](#6-passing-arrays-to-functions).

### 15.2 Assuming Two Different-Length Arrays Are the Same Type

```go
func sum(arr [5]int) int { /* ... */ return 0 }

var nums [10]int
// sum(nums) // COMPILE ERROR: [10]int is not [5]int
```

There is no automatic conversion between different array lengths — if you need a function to work with sequences of varying lengths, use a slice parameter instead, which decouples length from the type entirely.

### 15.3 Mutating a Loop-Local Copy Instead of the Array

```go
for _, v := range arr {
    v *= 2 // modifies only the loop-local copy of each element, not arr itself
}
```

Fix: index back into the array explicitly (`arr[i] = arr[i] * 2`) if you need to actually mutate its elements during a `range` loop — see [Section 8](#8-iterating-over-an-array).

### 15.4 Overusing Arrays Where a Slice Was Actually Needed

```go
var buffer [1024]byte // fine for a genuinely fixed-size buffer...
// ...but awkward and inflexible if the actual amount of data varies and needs to grow
```

If the true size of a collection isn't a fixed, compile-time-known constant, a slice is almost always the more appropriate and more idiomatic choice — reach for an array specifically when its fixed-size, value-type nature is a deliberate, needed property (see [Section 14](#14-when-arrays-are-actually-the-right-choice)), not by default.

### 15.5 Forgetting That Comparing Arrays of Different Lengths Doesn't Even Compile

```go
a := [3]int{1, 2, 3}
b := [4]int{1, 2, 3, 0}
// a == b // COMPILE ERROR: mismatched types [3]int and [4]int
```

Since length is part of the type, this isn't a runtime "false" result — it's a compile-time type error, since `[3]int` and `[4]int` are simply different, incompatible types.

## 16. Best Practices Summary

1. **Reach for a slice by default** for any collection whose size isn't a fixed, compile-time-known constant — arrays are the exception, not the default choice, in everyday Go code.
2. **Use arrays when you specifically need value semantics, comparability with `==`, or a genuinely fixed, compile-time-known size** (hash digests, fixed protocol headers, composite map keys).
3. **Pass a pointer to an array (`*[N]T`) when a function needs to mutate the caller's array** or when the array is large enough that copying it would be wasteful.
4. **Remember that `range` gives you a copy of each element by value** — index back into the array explicitly if you need to mutate it during iteration.
5. **Use the `[...]` literal form** when you want the compiler to infer an array's length from its literal contents, reducing the risk of a manually miscounted size.
6. **Don't try to write functions that are generic over array length without generics** — a function taking `[5]int` cannot accept a `[10]int`; use a slice parameter, or Go's generics with an array-length type parameter if that specific capability is genuinely needed.
7. **Leverage array comparability for composite map keys** when a small, fixed number of comparable components need to act as a single lookup key.

## 17. Use Case Summary Table

| Technique                          | When to Use                                                        | Example Scenario                                                          |
| ---------------------------------- | ------------------------------------------------------------------ | ------------------------------------------------------------------------- |
| `[N]T` value                       | A genuinely fixed-size, compile-time-known collection              | Days of the week, RGB channels (`[3]byte`), a hash digest (`[32]byte`)    |
| `[...]T{...}` literal              | Let the compiler infer the length from the literal's contents      | A lookup table whose size naturally matches its listed entries            |
| Array as a map key                 | A composite, fixed-size, comparable key                            | A 3D grid coordinate (`[3]int`) used to key a visited-set                 |
| `*[N]T` parameter                  | Mutate the caller's array, or avoid copying a large one            | An in-place transformation function operating on a fixed-size buffer      |
| Array comparison (`==`)            | Checking exact equality of two same-typed, fixed-size sequences    | Comparing two hash digests for equality                                   |
| Multidimensional array (`[N][M]T`) | A fixed-size grid or matrix known at compile time                  | A fixed board size for a game (e.g., a 3x3 tic-tac-toe grid)              |
| Fixed-size buffer in hot paths     | Avoiding a heap allocation for small, bounded, stack-friendly data | A small scratch buffer reused across many calls, sized to a known maximum |

## 18. References

1. Go Team — _A Tour of Go: Arrays_. https://go.dev/tour/moretypes/6
2. Go Language Specification — _Array types_. https://go.dev/ref/spec#Array_types
3. Go Team — _Go Slices: usage and internals_, The Go Blog (background on how arrays relate to slices). https://go.dev/blog/slices-intro
4. Go standard library documentation — package `crypto/sha256` (fixed-size `[32]byte` digest as a real-world example of a fixed-size array use case). https://pkg.go.dev/crypto/sha256
5. Go Team — _Go 1.17 Release Notes_ (slice-to-array-pointer conversion). https://go.dev/doc/go1.17
6. Go Team — _Go 1.20 Release Notes_ (slice-to-array conversion additions). https://go.dev/doc/go1.20
