<div align="center">
  <h1>Generics</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 17, 2026</sub>
</div>

## Table of Contents

1. [What Are Generics, and Why Did Go Add Them?](#1-what-are-generics-and-why-did-go-add-them)
2. [Life Before Generics](#2-life-before-generics)
3. [Basic Syntax: Type Parameters](#3-basic-syntax-type-parameters)
4. [Type Constraints](#4-type-constraints)
5. [Built-in Constraints: `any` and `comparable`](#5-built-in-constraints-any-and-comparable)
6. [Custom Constraints with Type Sets](#6-custom-constraints-with-type-sets)
7. [The Approximation Element `~`](#7-the-approximation-element-)
8. [Type Inference](#8-type-inference)
9. [Generic Types (Structs, and Beyond)](#9-generic-types-structs-and-beyond)
10. [Generic Methods — What's Allowed and What Isn't](#10-generic-methods--whats-allowed-and-whats-not)
11. [Multiple Type Parameters](#11-multiple-type-parameters)
12. [The Standard Library's Generic Packages](#12-the-standard-librarys-generic-packages)
13. [Common Generic Patterns](#13-common-generic-patterns)
14. [Generics vs. Interfaces: Choosing the Right Tool](#14-generics-vs-interfaces-choosing-the-right-tool)
15. [What the Go Team Deliberately Did _Not_ Generify](#15-what-the-go-team-deliberately-did-not-generify)
16. [Limitations and Things Generics Cannot Do](#16-limitations-and-things-generics-cannot-do)
17. [Performance Considerations](#17-performance-considerations)
18. [Common Mistakes and Pitfalls](#18-common-mistakes-and-pitfalls)
19. [Best Practices Summary](#19-best-practices-summary)
20. [Use Case Summary Table](#20-use-case-summary-table)
21. [References](#21-references)

## 1. What Are Generics, and Why Did Go Add Them?

**Generics** let you write functions and types that operate on a range of types specified by the caller, while still checking type correctness at compile time. In Go, this feature is implemented through **type parameters**: functions and types can declare placeholder types (conventionally named `T`, `U`, `K`, `V`, etc.) that are filled in with a concrete type argument at the call site.

Generics were added in **Go 1.18** (released March 2022) after years of community discussion and several design proposals. The core motivation was to eliminate a recurring pattern in pre-1.18 Go code: writing near-duplicate functions for different types (e.g., `SumInts`, `SumFloats`) or falling back to `interface{}` and giving up compile-time type safety.

Go's generics design has a distinctive character compared to generics/templates in other languages: type parameters are constrained not by a declared subtyping/inheritance relationship, but by **structural constraints** expressed as interfaces — describing exactly which types are legal to substitute in, and what operations the generic code is allowed to perform on them.

## 2. Life Before Generics

To appreciate what generics solve, consider the pre-1.18 alternatives for a function that sums a slice of numbers:

**Option A — duplicate code per type:**

```go
func SumInts(nums []int) int {
    var total int
    for _, n := range nums {
        total += n
    }
    return total
}

func SumFloats(nums []float64) float64 {
    var total float64
    for _, n := range nums {
        total += n
    }
    return total
}
```

This works and is fully type-safe, but the two functions are identical except for the type — every additional numeric type needs another copy-pasted function.

**Option B — `interface{}` (now `any`):**

```go
func Sum(nums []interface{}) interface{} {
    // ... would need type assertions and a type switch inside,
    // loses compile-time type safety, and can't even use `+` directly
    // on the interface{} values without asserting their type first.
}
```

This avoids duplication but sacrifices compile-time type checking, forces callers and implementers to deal with type assertions, and carries the runtime cost of boxing values into an interface (see [Section 17](#17-performance-considerations)).

**Generics solve both problems at once:** one function definition, full compile-time type safety, and no boxing overhead for the common case.

```go
type Number interface {
    int | int64 | float64
}

func SumNumbers[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

func main() {
    ints := []int{1, 2, 3}
    floats := []float64{1.1, 2.2, 3.3}
    fmt.Println(SumNumbers(ints))   // 6
    fmt.Println(SumNumbers(floats)) // 6.6
}
```

## 3. Basic Syntax: Type Parameters

A generic function declares its type parameter(s) in square brackets **before** the regular parameter list, with each type parameter given a **constraint**:

```go
func Print[T any](value T) {
    fmt.Println(value)
}
```

Here, `T` is the type parameter, and `any` is its constraint (meaning "any type at all is allowed"). Calling this function looks almost identical to calling a normal function — the type argument is usually inferred:

```go
Print(42)        // T inferred as int
Print("hello")   // T inferred as string
Print(3.14)      // T inferred as float64
Print[string]("explicit") // type argument given explicitly, when needed
```

**Key rule:** within the generic function's body, you may only use operations that are guaranteed to be valid for **every** type permitted by the constraint. A type parameter constrained by `any`, for example, supports assignment and being passed around, but not arithmetic operators like `+`, since not every type supports `+`.

## 4. Type Constraints

A **type constraint** is a kind of meta-type for a type parameter: it specifies the set of concrete types that calling code is allowed to substitute in, and correspondingly, what operations the generic code can perform on values of that type parameter. Constraints are themselves written as **interfaces**.

While a type parameter's constraint represents a _set_ of permissible types, inside the generic function body at any given call, the type parameter stands for one specific, single type — the type argument actually provided (or inferred) for that call. If the type argument doesn't satisfy the constraint, the code fails to compile.

```go
func Min[T int | float64](a, b T) T {
    if a < b {
        return a
    }
    return b
}
```

Here the constraint `int | float64` (a **union** of type terms) permits exactly two types; `Min` compiles because both `int` and `float64` support the `<` operator used in the body.

## 5. Built-in Constraints: `any` and `comparable`

Go's predeclared identifiers include two constraints usable directly, without importing anything:

### 5.1 `any`

`any` is an alias (added in Go 1.18) for the empty interface `interface{}`. As a constraint, it permits every type, with no guaranteed operations beyond what's true of every Go value (assignment, being passed as an argument, being returned).

```go
func Identity[T any](v T) T {
    return v
}
```

### 5.2 `comparable`

`comparable` is a built-in constraint satisfied by any type whose values can be compared with `==` and `!=` — this includes numeric types, strings, booleans, pointers, channels, interfaces, arrays of comparable types, and structs made entirely of comparable fields. It notably **excludes** slices, maps, and functions, since those are not comparable with `==` in Go.

```go
func Contains[T comparable](slice []T, target T) bool {
    for _, v := range slice {
        if v == target { // requires T to support ==
            return true
        }
    }
    return false
}

fmt.Println(Contains([]int{1, 2, 3}, 2))          // true
fmt.Println(Contains([]string{"a", "b"}, "c"))    // false
// Contains([][]int{{1}}, []int{1}) // COMPILE ERROR: []int does not satisfy comparable
```

`comparable` is exactly what you need whenever generic code relies on `==`/`!=`, or wants to use the type parameter as a map key (since Go map keys must be comparable).

## 6. Custom Constraints with Type Sets

Beyond `any` and `comparable`, you can define your own constraint interfaces that list a **union of specific types** using `|`:

```go
type Number interface {
    int | int8 | int16 | int32 | int64 |
        uint | uint8 | uint16 | uint32 | uint64 |
        float32 | float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}
```

Declaring a reusable constraint like `Number` this way is the idiomatic pattern when several generic functions or types need the same set of permitted types — you define it once as an interface and reference it wherever needed.

A constraint interface can also embed regular methods, exactly like an ordinary interface, in which case a type must both belong to the listed type set (if any) **and** implement those methods:

```go
type Stringer interface {
    String() string
}

func Join[T Stringer](items []T, sep string) string {
    parts := make([]string, len(items))
    for i, item := range items {
        parts[i] = item.String()
    }
    return strings.Join(parts, sep)
}
```

## 7. The Approximation Element `~`

Sometimes you want a constraint to accept not just a specific named type, but also any other type whose **underlying type** matches it — for example, your own `type Celsius float64` should still count as a "float64-like" type for a generic numeric function. The `~` (tilde) prefix on a type term in a constraint enables exactly this:

```go
type Number interface {
    ~int | ~int64 | ~float64
}

type Celsius float64

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

temps := []Celsius{20.5, 21.0, 19.8}
fmt.Println(Sum(temps)) // works, because Celsius's underlying type is float64
```

Without `~float64` (i.e., with just `float64` in the union), `Celsius` — despite being defined as `float64` underneath — would **not** satisfy the constraint, because `float64` in a union without `~` only matches the exact type `float64`, not types defined in terms of it.

**When constraint type inference kicks in:** if a constraint is written using `~` on a type built from another type parameter (for example, a slice-based constraint of the form `~[]E` for some element type `E`), and the compiler already knows the concrete argument for one type parameter, it can often infer the other automatically — this is called **constraint type inference** and is what lets certain generic functions omit explicit type arguments even when the relationship between type parameters is only implicit in the constraint's shape.

## 8. Type Inference

In most calls, you don't need to specify type arguments explicitly — the compiler infers them from the types of the ordinary function arguments you pass in, exactly as shown in every example so far (`SumNumbers(ints)` infers `T = int`, no need to write `SumNumbers[int](ints)`).

Type inference either **succeeds** — in which case calling a generic function looks exactly like calling an ordinary one — or **fails**, in which case the compiler reports a clear error and you must supply the type argument(s) explicitly:

```go
result := Min[float64](3, 4.5) // explicit type argument when inference can't determine it alone
```

Inference can fail in various situations — for example, when a type parameter appears only in the function's return type and nowhere in its regular parameters, since there's nothing for the compiler to infer it _from_.

```go
func Zero[T any]() T {
    var zero T
    return zero
}

x := Zero[int]() // must specify T explicitly — nothing to infer it from
```

The full rules of type inference are intricate (they involve a unification algorithm across both regular function-argument types and constraint type sets), but using it day-to-day is simple: write the call the natural way, and let the compiler tell you if it needs an explicit type argument.

## 9. Generic Types (Structs, and Beyond)

Type parameters aren't limited to functions — you can parameterize a **type declaration** itself, most commonly a struct, to build a generic data structure:

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    last := len(s.items) - 1
    item := s.items[last]
    s.items = s.items[:last]
    return item, true
}

func main() {
    intStack := &Stack[int]{}
    intStack.Push(1)
    intStack.Push(2)
    v, ok := intStack.Pop() // v == 2, ok == true

    strStack := &Stack[string]{}
    strStack.Push("hello")
}
```

Here, `Stack[T]` is a **generic type**; `Stack[int]` and `Stack[string]` are two different **instantiations** of it, each acting like its own concrete type once instantiated. This is the generic equivalent of writing `IntStack`, `StringStack`, etc. by hand — one definition serves every element type, fully type-checked.

A parameterized type can also have multiple type parameters, and generic types other than structs are possible too (e.g., a generic named slice type, a generic map type, or a generic channel type), following the same bracketed type-parameter syntax.

## 10. Generic Methods — What's Allowed and What's Not

A method defined on a generic type automatically has access to that type's type parameter(s) — you don't (and can't) redeclare them on the method itself:

```go
func (s *Stack[T]) Len() int {
    return len(s.items) // T is already in scope from the receiver
}
```

**Important limitation:** a method cannot introduce **new** type parameters of its own beyond the ones the receiver type already has. Go methods, unlike free functions, cannot be independently generic — all of a method's type parameterization must come from its receiver type's own type parameter list. If you need a genuinely separate type parameter for one specific operation (e.g., a `Map` operation that transforms a `Stack[T]` into a `Stack[U]`), you must write it as a standalone generic function instead of a method:

```go
// NOT possible as a method: func (s *Stack[T]) Map[U any](f func(T) U) *Stack[U] { ... }

// Must be a free function instead:
func MapStack[T, U any](s *Stack[T], f func(T) U) *Stack[U] {
    out := &Stack[U]{}
    for _, item := range s.items {
        out.Push(f(item))
    }
    return out
}
```

## 11. Multiple Type Parameters

A generic function or type can declare more than one type parameter, each with its own (possibly different) constraint:

```go
func Map[T, U any](items []T, f func(T) U) []U {
    result := make([]U, len(items))
    for i, item := range items {
        result[i] = f(item)
    }
    return result
}

func main() {
    ints := []int{1, 2, 3}
    strs := Map(ints, func(i int) string {
        return fmt.Sprintf("num:%d", i)
    })
    fmt.Println(strs) // [num:1 num:2 num:3]
}
```

`Map[T, U any]` declares two type parameters, `T` and `U`, both constrained by `any` — `T` is inferred from the input slice's element type, and `U` is inferred from the return type of the function argument `f`.

## 12. The Standard Library's Generic Packages

Go 1.21 (August 2023) added three new standard-library packages built entirely on generics, which cover many everyday needs without any third-party dependency:

### 12.1 `cmp`

Defines the `cmp.Ordered` constraint (types whose values can be compared with `<`, `<=`, `>`, `>=`) and two generic functions:

```go
import "cmp"

cmp.Compare(3, 5)  // -1 (negative, zero, or positive, like strings.Compare)
cmp.Less(3, 5)     // true
```

Go 1.21 also added generic built-in functions `min`, `max`, and `clear` directly to the language (not the `cmp` package), which compute the minimum/maximum of ordered values or clear a map/slice, respectively.

### 12.2 `slices`

Provides common slice operations, generic over the element type: `slices.Contains`, `slices.Index`, `slices.Sort`, `slices.SortFunc`, `slices.Clone`, `slices.Equal`, `slices.Reverse`, and more.

```go
import "slices"

nums := []int{5, 2, 8, 1}
slices.Sort(nums)
fmt.Println(nums) // [1 2 5 8]
fmt.Println(slices.Contains(nums, 8)) // true
```

### 12.3 `maps`

Provides common map operations, generic over key and value type: `maps.Keys`, `maps.Values`, `maps.Clone`, `maps.Equal`, `maps.Copy`, and more.

```go
import "maps"

m := map[string]int{"a": 1, "b": 2}
clone := maps.Clone(m)
```

**Practical takeaway:** before reaching for a third-party generics utility library, check whether `slices`, `maps`, or `cmp` already provide what you need — since Go 1.21, a large share of common generic container operations are covered directly by the standard library, and these standard implementations are typically well-optimized and maintained as part of Go itself.

## 13. Common Generic Patterns

### 13.1 Generic Numeric Helpers

```go
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
```

(Since Go 1.21, prefer the built-in `max`/`min` functions or `cmp.Ordered` from the standard library over hand-rolling this, unless you need something they don't cover.)

### 13.2 Generic Functional Helpers (Map / Filter / Reduce)

```go
func Filter[T any](items []T, predicate func(T) bool) []T {
    var result []T
    for _, item := range items {
        if predicate(item) {
            result = append(result, item)
        }
    }
    return result
}

func Reduce[T, U any](items []T, initial U, f func(U, T) U) U {
    acc := initial
    for _, item := range items {
        acc = f(acc, item)
    }
    return acc
}

nums := []int{1, 2, 3, 4, 5}
evens := Filter(nums, func(n int) bool { return n%2 == 0 })      // [2 4]
total := Reduce(nums, 0, func(acc, n int) int { return acc + n }) // 15
```

### 13.3 Generic Data Structures

Beyond the `Stack[T]` example in [Section 9](#9-generic-types-structs-and-beyond), the same technique applies to queues, linked lists, binary trees, sets, and other containers — write the structure once, parameterized over its element type, instead of hand-writing (or code-generating) a version per concrete type.

### 13.4 Generic Set Type

```go
type Set[T comparable] struct {
    items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
    return &Set[T]{items: make(map[T]struct{})}
}

func (s *Set[T]) Add(item T)      { s.items[item] = struct{}{} }
func (s *Set[T]) Has(item T) bool { _, ok := s.items[item]; return ok }
```

### 13.5 Type-Safe Wrapper Around `sync.Once`-Style Lazy Values

The standard library's own `sync.OnceValue[T]` (Go 1.21+) is itself a generic function, illustrating how generics let previously `interface{}`-based (or type-duplicated) utilities become fully type-safe:

```go
var getConfig = sync.OnceValue(func() *Config {
    return loadConfig() // runs exactly once, lazily, on first call
})
```

## 14. Generics vs. Interfaces: Choosing the Right Tool

Generics and interfaces solve different problems and are frequently used **together** (a constraint is itself an interface), not as competitors:

|                   | Interfaces                                                                                         | Generics                                                                     |
| ----------------- | -------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| Abstracts over    | **Behavior** — what a type can _do_                                                                | **The type itself** — writing one algorithm/structure for many types         |
| Dispatch          | Dynamic, at runtime, via the interface's method set                                                | Resolved at compile time per instantiation                                   |
| Typical fit       | Swappable implementations, dependency injection, heterogeneous collections handled polymorphically | Type-safe containers and algorithms whose logic doesn't change between types |
| Runtime cost      | Some overhead from dynamic dispatch and potential boxing                                           | None beyond what hand-written per-type code would cost                       |
| Canonical example | `io.Writer`, `sort.Interface`, `error`                                                             | `slices.Sort[T]`, a generic `Stack[T]`, `cmp.Ordered`                        |

**Rule of thumb:** if you need to hold _different_ concrete types behind one variable and call the _same_ method on each with _different_ underlying behavior (a `Shape` that could be a `Circle` or `Rectangle`, each computing its own area differently) — that's polymorphism, and it calls for an interface. If you're writing one _identical_ algorithm or structure that should work for `int`, `float64`, or your own numeric type without behavior changing per type — that's generic code, and it calls for a type parameter.

## 15. What the Go Team Deliberately Did _Not_ Generify

Since generics shipped in Go 1.18, the Go team has been notably conservative about retrofitting existing `interface{}`-based standard library APIs into generic ones, which is itself an instructive signal about idiomatic use. For example, `container/list` and `container/heap`, which predate generics and are built around `interface{}`/specific interfaces, were **not** rewritten to use type parameters, even though the opportunity existed. The lesson generally drawn from this restraint is that generics are best reserved for cases with a clear, immediate type-safety or duplication problem to solve, rather than applied reflexively to any container or algorithm just because the language now supports them.

## 16. Limitations and Things Generics Cannot Do

- **No new type parameters on methods.** As covered in [Section 10](#10-generic-methods--whats-allowed-and-whats-not), a method's type parameters must come entirely from its receiver's type parameter list.
- **Operations are limited to what the constraint guarantees.** You cannot use `+` inside a generic function unless the constraint's type set guarantees every permitted type supports `+` (this is why `Number`-style constraints explicitly list arithmetic types, rather than using `any`).
- **No specialization/overloading per concrete type** the way C++ templates allow — a single generic function body must work uniformly for every type satisfying its constraint; you can't write different code paths that get selected purely by which concrete type was substituted (a type switch inside the function is possible but is a very different mechanism from template specialization).
- **`comparable` doesn't mean "has an `Equal` method"** — it specifically means "supports `==`/`!=`," which excludes slices, maps, and functions even if they have a custom `Equal` method.
- **Generic type aliases had restrictions** in earlier Go versions — the ability to parameterize a type alias itself (not just a `type` definition) was limited before Go 1.24, which added broader support for generic type aliases.
- **No compile-time metaprogramming.** Go's generics design deliberately does not support arbitrary compile-time code generation/execution the way C++ templates or Rust's macro system can — the feature is scoped specifically to parametric polymorphism over types.

## 17. Performance Considerations

- **No boxing for the common case.** Unlike `interface{}`-based "generic" code, a properly constrained generic function is compiled (conceptually) as if it were written specifically for each concrete type actually used — there's no heap allocation to box a value into an interface, and no dynamic dispatch overhead, in the typical case.
- **Compilation strategy details can vary** across Go compiler versions and are an internal implementation detail (Go has used different approaches, including a mix of per-type instantiation and shared "GC-shaped" implementations for pointer-like types, to balance compile time, binary size, and runtime speed) — application code generally doesn't need to reason about this directly.
- **Constraints with methods still involve interface-like dispatch** when the constraint requires calling a method on the type parameter, since that method call is resolved through the same mechanism as calling a method through an interface.
- **In practice**, generic code is typically as fast as, or very close to, hand-written per-type code, and meaningfully faster than the equivalent `interface{}`-based code that needed type assertions and boxing — this was a central design goal of Go's generics, not an incidental benefit.

## 18. Common Mistakes and Pitfalls

### 18.1 Using `any` When a Narrower Constraint Is Needed

```go
func Sum[T any](nums []T) T { // COMPILE ERROR inside the body
    var total T
    for _, n := range nums {
        total += n // `any` doesn't guarantee `+` is supported
    }
    return total
}
```

Fix: constrain `T` to a type set that actually supports `+` (e.g., a `Number` constraint as shown earlier).

### 18.2 Forgetting `~` When a Custom Named Type Should Qualify

```go
type Number interface {
    int | float64 // no ~
}
type Meters float64

func Sum[T Number](nums []T) T { /* ... */ }

// Sum([]Meters{1, 2, 3}) // COMPILE ERROR: Meters does not satisfy Number
```

Fix: use `~int | ~float64` if user-defined types with those underlying types should also be accepted.

### 18.3 Expecting `comparable` to Cover Slices/Maps

```go
func Contains[T comparable](s []T, v T) bool { /* uses == */ }

// Contains([][]int{...}, []int{1,2}) // COMPILE ERROR: []int is not comparable
```

There's no way to make `comparable` accept slices/maps directly, since Go itself doesn't define `==` for them; you'd need a different approach (e.g., accepting an explicit equality function parameter) for those types.

### 18.4 Introducing Generics Where a Plain Interface (or No Abstraction) Would Do

Reaching for generics reflexively — parameterizing a type or function "just in case," when there is realistically only ever going to be one concrete type used, or when an interface would express the actual need (different behavior per type) more naturally — adds unnecessary complexity and a less approachable API. Follow the standard library's own restraint (see [Section 15](#15-what-the-go-team-deliberately-did-not-generify)): generify when duplication or lost type-safety is a real, current problem.

### 18.5 Trying to Add Type Parameters to a Method Directly

```go
type Container[T any] struct{ items []T }

// func (c *Container[T]) Convert[U any]() []U { ... } // COMPILE ERROR
```

As covered in [Section 10](#10-generic-methods--whats-allowed-and-whats-not), this must be a standalone function instead of a method.

### 18.6 Not Checking the Standard Library First

Hand-rolling a generic `Contains`, `Sort`, `Keys`, or `Max` helper when Go 1.21's `slices`, `maps`, or `cmp` packages (or the built-in `min`/`max`/`clear` functions) already provide it adds unnecessary code and forfeits the standard library's own optimizations and maintenance.

## 19. Best Practices Summary

1. **Reach for generics to solve a real, current problem** — eliminating genuine code duplication across types, or replacing an `interface{}`-based API that's currently losing type safety — not preemptively.
2. **Write the narrowest constraint that expresses what your generic code actually needs**, rather than defaulting to `any` and hoping the compiler figures it out.
3. **Reuse constraint interfaces** (like a shared `Number` or `Ordered` type) across multiple generic functions/types instead of redefining similar unions repeatedly.
4. **Use `~` in a constraint's type terms** whenever user-defined types with that underlying type should also be allowed to satisfy it.
5. **Prefer the standard library's `slices`, `maps`, and `cmp` packages** (Go 1.21+), and the built-in `min`/`max`/`clear` functions, over hand-writing common generic container operations.
6. **Remember methods can't introduce new type parameters** — if you need a transformation with a different output type parameter than the receiver's, write a standalone generic function instead.
7. **Use `comparable` specifically for `==`/`!=`-based logic**, and understand it does not cover slices, maps, or functions.
8. **Don't confuse generics with interfaces' purpose** — generics for "same logic, many types"; interfaces for "different behavior, one contract."
9. **Keep an eye on the Go team's own restraint** with retrofitting existing APIs — not every container or algorithm benefits from becoming generic, even when it technically could.
10. **Let type inference do the work** in normal calls; only supply explicit type arguments when the compiler reports that inference failed.

## 20. Use Case Summary Table

| Technique                                        | When to Use                                                                                        | Example Scenario                                                    |
| ------------------------------------------------ | -------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------- |
| `func F[T any](...)`                             | Logic that works identically regardless of type, with no operations beyond assignment/pass-through | A generic `Identity`, `Print`, or container element type            |
| `func F[T comparable](...)`                      | Logic relying on `==`/`!=`, or using `T` as a map key                                              | `Contains`, `Unique`, a generic `Set[T]`                            |
| Custom numeric/ordered constraint                | Arithmetic or ordering operations across multiple numeric types                                    | `Sum`, `Max`, `Average` over `int`/`float64`/custom numeric types   |
| `~` in a constraint                              | Custom named types with a matching underlying type should also qualify                             | A `Celsius` or `Meters` type used with a generic numeric helper     |
| Generic struct (`type Stack[T any] struct{...}`) | A container whose structure/logic is identical regardless of element type                          | `Stack[T]`, `Queue[T]`, `Set[T]`, a generic binary tree             |
| Multiple type parameters (`[T, U any]`)          | A transformation between two different types                                                       | `Map[T, U]`, `Reduce[T, U]`                                         |
| `slices` / `maps` / `cmp` (stdlib)               | Common slice/map operations without reinventing them                                               | Sorting, searching, cloning, comparing slices and maps              |
| Plain interface instead of generics              | Different concrete types need genuinely different behavior behind one contract                     | `io.Writer`, `Shape` with `Area()` implemented differently per type |

## 21. References

1. Go Team — _Tutorial: Getting started with generics_. https://go.dev/doc/tutorial/generics
2. Go Team — _An Introduction To Generics_, The Go Blog. https://go.dev/blog/intro-generics
3. Go Team — _Why Generics?_, The Go Blog. https://go.dev/blog/why-generics
4. Ian Lance Taylor & Robert Griesemer — _Type Parameters Proposal_. https://go.dev/design/43651-type-parameters (mirrored at https://tip.golang.org/design/43651-type-parameters)
5. Go Team — _Go 1.21 Release Notes_ (`slices`, `maps`, `cmp` packages; `min`/`max`/`clear` built-ins). https://go.dev/doc/go1.21 (also https://tip.golang.org/doc/go1.21)
6. Go standard library documentation — package `cmp`. https://pkg.go.dev/cmp
7. Go standard library documentation — package `slices`. https://pkg.go.dev/slices
8. Go standard library documentation — package `maps`. https://pkg.go.dev/maps
9. Go standard library documentation — package `sync` (`sync.OnceValue` and related). https://pkg.go.dev/sync
10. Erik Engheim — _Golang Generics For Non-Beginners_, ITNEXT. https://itnext.io/golang-generics-for-non-beginners-6ca7a4716aa9
11. OneUptime Engineering Blog — _How to Use Generics with Type Constraints in Go_. https://oneuptime.com/blog/post/2026-01-23-go-generics-constraints/view
12. Gabriel Anhaia — _Go Generics, 4 Years In: The 3 Cases Where They're the Right Answer_, DEV Community. https://dev.to/gabrielanhaia/go-generics-4-years-in-the-3-cases-where-theyre-the-right-answer-5bip
13. Michael Whatcott — _A generic 'leaderboard' map type in Go_. https://michaelwhatcott.com/go-leaderboard-map-type/
14. Educative — _Creating Constraints_, from _Advanced Techniques in Go Programming_. https://educative.io/courses/advanced-techniques-in-go-programming/creating-constraints
15. Educative — _Type constraints_, from _Go Programming: From Beginner to Professional_, Packt. https://www.packtpub.com/en-us/product/go-programming-from-beginner-to-professional-second-edition-10-9781803243054/chapter/chapter-8-generic-algorithm-superpowers/type-constraints
16. Packt — _Mastering Go: The cmp package_. https://packtpub.com/en-ee/product/mastering-go-9781805127147/chapter/go-generics-4/section/the-cmp-package-ch04lvl1sec38
17. goframe.org — _Go 1.21 (2023-08-08) release summary_. https://goframe.org/en/release/golang/go1.21
18. Leapcell — _Go 1.25 Highlights: How Generics and Performance Define the Future of Go_, Medium. https://leapcell.medium.com/go-1-25-highlights-how-generics-and-performance-define-the-future-of-go-8bd81296b4bf
19. Go Language Specification — _Type parameters, Type constraints, Instantiations_. https://go.dev/ref/spec#Type_parameter_declarations
