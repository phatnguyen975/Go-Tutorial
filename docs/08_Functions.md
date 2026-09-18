<div align="center">
  <h1>Functions</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 18, 2026</sub>
</div>

## Table of Contents

1. [What Is a Function?](#1-what-is-a-function)
2. [Basic Function Declaration](#2-basic-function-declaration)
3. [Parameters](#3-parameters)
4. [Return Values](#4-return-values)
5. [Multiple Return Values](#5-multiple-return-values)
6. [Named Return Values](#6-named-return-values)
7. [The Blank Identifier with Return Values](#7-the-blank-identifier-with-return-values)
8. [Variadic Functions](#8-variadic-functions)
9. [Functions Are First-Class Values](#9-functions-are-first-class-values)
10. [Function Types](#10-function-types)
11. [Anonymous Functions (Function Literals)](#11-anonymous-functions-function-literals)
12. [Closures](#12-closures)
13. [The Classic Loop-Variable-Capture Pitfall](#13-the-classic-loop-variable-capture-pitfall)
14. [Higher-Order Functions](#14-higher-order-functions)
15. [Recursion](#15-recursion)
16. [Functions as Struct Fields](#16-functions-as-struct-fields)
17. [The `error` Return Convention](#17-the-error-return-convention)
18. [Function Naming Conventions](#18-function-naming-conventions)
19. [Common Mistakes and Pitfalls](#19-common-mistakes-and-pitfalls)
20. [Best Practices Summary](#20-best-practices-summary)
21. [Use Case Summary Table](#21-use-case-summary-table)
22. [References](#22-references)

## 1. What Is a Function?

A function is a named, reusable block of code that takes zero or more input parameters and produces zero or more output values. Functions are the fundamental unit of behavior in Go — every program's execution begins at a function (`main`), and virtually all logic is organized into functions and methods (a method being, at its core, a function with an extra receiver parameter).

```go
func add(a int, b int) int {
    return a + b
}
```

Go's functions support several features that go noticeably beyond a typical C-style function signature: multiple return values, named return values, variadic parameters, and full first-class treatment (functions can be stored in variables, passed as arguments, and returned from other functions) — this combined feature set gives Go a genuinely functional-programming-friendly style, despite being a statically typed, imperative-first language.

## 2. Basic Function Declaration

The general form of a function declaration:

```go
func functionName(parameterName ParameterType, ...) (returnType, ...) {
    // function body
}
```

```go
func greet(name string) string {
    return "Hello, " + name
}

func main() {
    message := greet("Alice")
    fmt.Println(message) // Hello, Alice
}
```

A function with no parameters and no return value is written simply as:

```go
func sayHello() {
    fmt.Println("Hello!")
}
```

## 3. Parameters

### 3.1 Basic Parameters

Each parameter is declared with a name followed by its type:

```go
func multiply(x int, y int) int {
    return x * y
}
```

### 3.2 Grouping Parameters of the Same Type

When consecutive parameters share the same type, the type can be written once, after the last parameter in that group:

```go
func multiply(x, y int) int { // equivalent to (x int, y int)
    return x * y
}

func describe(name string, age, height int) { // name is string; age and height are both int
    // ...
}
```

### 3.3 Passing by Value (and Simulating Reference Semantics with Pointers)

Every argument in Go is passed **by value** — the function receives a copy of whatever was passed in. For basic types and structs, this means the function operates on independent data; to let a function observe or mutate the caller's original value, a pointer must be passed explicitly:

```go
func increment(n int) {
    n++ // modifies only the local copy
}

func incrementPtr(n *int) {
    *n++ // modifies the original value through the pointer
}

x := 5
increment(x)
fmt.Println(x) // 5 — unchanged

incrementPtr(&x)
fmt.Println(x) // 6 — changed
```

## 4. Return Values

### 4.1 A Single Return Value

```go
func square(n int) int {
    return n * n
}
```

### 4.2 No Return Value

A function with no return type simply omits it — such a function is called purely for its side effects:

```go
func logMessage(msg string) {
    fmt.Println("[LOG]", msg)
}
```

### 4.3 The Function Body Must Cover Every Path (for Non-Void Functions)

If a function declares a return type, the compiler requires that every possible execution path through the function body actually reaches a `return` statement — a function that could "fall off the end" without returning a value is a compile-time error:

```go
func classify(n int) string {
    if n > 0 {
        return "positive"
    }
    // COMPILE ERROR: missing return at end of function
    // (the compiler cannot prove every path returns, even if you believe n is always > 0)
}
```

Fixing this requires an explicit `else` branch, a trailing `return` after the `if`, or some other structure that satisfies the compiler that every path returns a value.

## 5. Multiple Return Values

Unlike most C-family languages, Go functions can return **more than one value** — this is used pervasively throughout the standard library and idiomatic Go code, most famously for the `(result, error)` pattern:

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result) // 5
}
```

Multiple return values are also commonly used for the "value, ok" idiom, seen with map lookups, type assertions, and channel receives:

```go
func getLanguages() (string, string, int, bool) {
    return "Golang", "TypeScript", 8, true
}

lang1, lang2, count, active := getLanguages()
```

A function returning multiple values can also be used directly as the sole argument to another function whose parameter list matches — for example, passing the two-value result of one function directly into `fmt.Println`, as long as the types line up.

## 6. Named Return Values

Go allows a function's return values to be given **names** directly in the signature, turning them into local variables that are automatically declared (and initially zero-valued) at the start of the function body:

```go
func divide(a, b float64) (result float64, err error) {
    if b == 0 {
        err = errors.New("division by zero")
        return // "naked" return: returns the current values of result and err
    }
    result = a / b
    return
}
```

### 6.1 The "Naked" Return

When return values are named, a bare `return` statement (with no expressions after it) automatically returns whatever the named variables currently hold — this is called a **naked return**. It can reduce repetition, but overusing it in longer functions can hurt readability, since a reader has to scroll back up to the signature to know what a bare `return` is actually sending back.

### 6.2 Documenting Intent

Named return values can also serve a purely documentary purpose, even in functions that always use explicit `return` expressions — clarifying what each returned value represents, which is especially useful when a function returns several values of the same type:

```go
func divmod(a, b int) (quotient, remainder int) {
    quotient = a / b
    remainder = a % b
    return quotient, remainder // explicit, even though naked `return` would also work
}
```

### 6.3 Interaction with `defer`

Named return values have a special, powerful interaction with `defer`: a deferred function can read and modify the named return values before the function actually hands control back to its caller, since the values are set by the `return` statement, but the function hasn't fully returned until all deferred calls have run:

```go
func writeToFile(path string, data []byte) (err error) {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = cerr // attach a close failure to the function's own error result
        }
    }()

    _, err = f.Write(data)
    return err
}
```

This pattern — capturing a cleanup operation's own potential failure into the function's already-set error result — is one of the most common, genuinely useful reasons to reach for named return values deliberately, beyond simple documentation.

## 7. The Blank Identifier with Return Values

When a function returns multiple values but the caller only needs some of them, the blank identifier `_` discards any values that aren't needed, without triggering an "unused variable" compile error:

```go
func f() (int, int) {
    return 10, 20
}

a, _ := f()   // discard the second value
_, b := f()   // discard the first value
fmt.Println(a, b) // 10 20
```

## 8. Variadic Functions

A **variadic** function accepts a variable number of arguments of a specified type, declared by prefixing the parameter's type with `...`. Inside the function body, the variadic parameter behaves like an ordinary slice of that type:

```go
func sum(numbers ...int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

fmt.Println(sum())           // 0 — zero arguments is fine
fmt.Println(sum(1, 2, 3))     // 6
fmt.Println(sum(1, 2, 3, 4, 5)) // 15
```

### 8.1 Spreading an Existing Slice into a Variadic Call

If you already have a slice and want to pass its elements as individual variadic arguments (rather than as a single slice argument), append `...` after the slice at the call site:

```go
nums := []int{10, 20, 30}
total := sum(nums...) // spreads nums' elements as individual arguments
```

### 8.2 A Variadic Parameter Must Be Last

A function can have at most one variadic parameter, and it must be the **final** parameter in the list — any fixed parameters must come before it:

```go
func logWithPrefix(prefix string, values ...int) {
    fmt.Println(prefix, values)
}
```

### 8.3 A Well-Known Real-World Example

`fmt.Println` and its relatives are themselves variadic, over `any`, which is exactly what lets them accept any number of arguments of any type:

```go
func Println(a ...any) (n int, err error)
```

## 9. Functions Are First-Class Values

In Go, functions are values just like integers, strings, or structs — they can be assigned to variables, stored in data structures, passed as arguments, and returned from other functions:

```go
func add(a, b int) int { return a + b }

var operation func(int, int) int = add // storing a function in a variable
result := operation(3, 4)
fmt.Println(result) // 7
```

This first-class treatment is what enables Go's support for a functional-programming style — closures, higher-order functions, and callback-based APIs all rest on this foundation.

## 10. Function Types

A function's **type** is fully described by its parameter types and return types — the function's name plays no part in its type at all. This means two functions with entirely different names but identical signatures share exactly the same type, and can be used interchangeably wherever that function type is expected:

```go
type BinaryOp func(int, int) int // a named function type

func add(a, b int) int      { return a + b }
func multiply(a, b int) int  { return a * b }

var op BinaryOp = add
fmt.Println(op(3, 4)) // 7

op = multiply
fmt.Println(op(3, 4)) // 12
```

Declaring a named function type (like `BinaryOp` above) is a common, idiomatic way to make function-typed parameters and variables more readable and self-documenting, especially when the same function signature is used in multiple places across a package.

## 11. Anonymous Functions (Function Literals)

A function can be declared **inline**, without a name — this is called a function literal, or informally an anonymous function:

```go
square := func(n int) int {
    return n * n
}

fmt.Println(square(5)) // 25
```

An anonymous function can also be invoked **immediately**, right where it's defined, by appending a call directly after its closing brace — a pattern known as an immediately-invoked function expression (IIFE):

```go
result := func(a, b int) int {
    return a + b
}(3, 4)

fmt.Println(result) // 7
```

**Common use cases:** passing a short, one-off callback into another function (e.g., `sort.Slice`'s comparison function), wrapping a block of code to create its own local scope (often combined with `defer`, to limit a deferred call's lifetime to just that block), and, most importantly, forming closures — covered next.

## 12. Closures

A **closure** is a function literal that references variables from outside its own body — it "closes over" those variables, capturing them by reference rather than by value, so the function can read and modify them even after the surrounding code that declared them has finished executing:

```go
func counter() func() int {
    count := 0
    return func() int {
        count++ // this inner function closes over `count`
        return count
    }
}

increment := counter()
fmt.Println(increment()) // 1
fmt.Println(increment()) // 2
fmt.Println(increment()) // 3

anotherIncrement := counter() // a completely separate closure, with its OWN independent count
fmt.Println(anotherIncrement()) // 1 — unaffected by the first counter's state
```

Each call to `counter()` creates a **new**, independent `count` variable and a new closure over it — the returned function retains access to that specific `count` for as long as the closure itself is reachable, even though `counter()`'s own stack frame has long since returned. This is possible because Go's compiler automatically allocates `count` on the heap (via escape analysis) once it detects that a closure captures it and needs to survive past the enclosing function's return.

**Use cases for closures:**

- **Maintaining private state across calls**, without resorting to a package-level global variable — as shown by `counter()` above.
- **Configuring a function's behavior via parameters at creation time**, producing specialized functions from a general template (a pattern sometimes called "function factories" or partial application).
- **Callbacks that need access to surrounding context**, such as an HTTP handler closure that captures a shared database connection or configuration struct.

```go
// A closure-based function factory: produces a validator bound to a specific max value
func maxValidator(max int) func(int) bool {
    return func(n int) bool {
        return n <= max
    }
}

isValidAge := maxValidator(120)
fmt.Println(isValidAge(30))  // true
fmt.Println(isValidAge(200)) // false
```

## 13. The Classic Loop-Variable-Capture Pitfall

Because closures capture variables **by reference**, capturing a loop's iteration variable directly used to be a very common source of bugs before Go 1.22:

```go
// Go 1.21 and earlier: a well-known bug
var funcs []func()
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() {
        fmt.Println(i) // captures the SAME variable i across all iterations
    })
}
for _, f := range funcs {
    f() // prints 3, 3, 3 — not 0, 1, 2, as might be expected!
}
```

Prior to Go 1.22, a `for` loop's index variable (`i` here) was a single variable reused across every iteration — every closure created inside the loop captured that same, single variable, so by the time any of the closures actually ran, `i` had already reached its final value from the loop's last iteration.

**Pre-1.22 fix:** explicitly create a new, iteration-local variable to capture instead:

```go
for i := 0; i < 3; i++ {
    i := i // shadow: a fresh variable, scoped to this single iteration
    funcs = append(funcs, func() {
        fmt.Println(i) // now captures each iteration's own distinct i
    })
}
// prints 0, 1, 2
```

**Go 1.22+:** the language specification was changed so that each iteration of a `for` loop gets its own, fresh copy of the loop variable automatically, eliminating this specific bug for code compiled against the Go 1.22+ language version. Being explicit about what a closure captures remains good defensive practice regardless, particularly in code that must remain compatible with older Go versions.

## 14. Higher-Order Functions

A **higher-order function** is a function that either accepts another function as a parameter, returns a function, or both. Go's first-class function support makes this a natural, commonly used pattern.

```go
// Takes a function as a parameter
func applyTwice(f func(int) int, x int) int {
    return f(f(x))
}

double := func(n int) int { return n * 2 }
fmt.Println(applyTwice(double, 3)) // 12 — double(double(3)) = double(6) = 12
```

```go
// Returns a function
func multiplier(factor int) func(int) int {
    return func(n int) int {
        return n * factor
    }
}

triple := multiplier(3)
fmt.Println(triple(7)) // 21
```

**Real-world examples in the standard library:** `sort.Slice(slice, less func(i, j int) bool)` takes a comparison function; `http.HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))` registers a handler function; middleware patterns in web frameworks commonly take a handler function and return a new, wrapped handler function.

```go
people := []struct {
    Name string
    Age  int
}{
    {"Alice", 30}, {"Bob", 25}, {"Carol", 35},
}

sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age // higher-order: a function passed as an argument
})
```

## 15. Recursion

A function in Go can call itself, either directly or indirectly (through a chain of other function calls) — this is recursion, and Go supports it without any special syntax:

```go
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

fmt.Println(factorial(5)) // 120
```

**Mutual recursion** (two or more functions calling each other) is also fully supported, since Go resolves function calls after the entire package's declarations have been parsed, regardless of the order functions are written in the source file:

```go
func isEven(n int) bool {
    if n == 0 {
        return true
    }
    return isOdd(n - 1)
}

func isOdd(n int) bool {
    if n == 0 {
        return false
    }
    return isEven(n - 1)
}
```

**A practical caution:** Go does not guarantee tail-call optimization (eliminating stack growth for a recursive call in tail position) the way some functional languages do — very deep recursion can exhaust the available stack and crash the program (`fatal error: stack overflow`), so an iterative approach is often preferable for algorithms that could recurse extremely deeply on realistic input sizes.

## 16. Functions as Struct Fields

Because functions are ordinary values, a struct field can itself hold a function, enabling flexible, pluggable behavior configured per-instance:

```go
type Validator struct {
    Name  string
    Check func(value string) bool
}

emailValidator := Validator{
    Name: "email",
    Check: func(value string) bool {
        return strings.Contains(value, "@")
    },
}

fmt.Println(emailValidator.Check("user@example.com")) // true
```

This pattern is common for building configuration-driven behavior, strategy-style designs (different struct instances carrying different "strategies" as function fields), and simple dependency injection without needing a full interface hierarchy for very small, single-operation abstractions.

## 17. The `error` Return Convention

By strong, near-universal convention, a Go function that can fail returns an `error` as its **final** return value — `nil` signaling success, and a non-nil value describing what went wrong:

```go
func parseAge(s string) (int, error) {
    age, err := strconv.Atoi(s)
    if err != nil {
        return 0, fmt.Errorf("invalid age %q: %w", s, err)
    }
    return age, nil
}
```

This convention is so consistently followed throughout the standard library and the broader Go ecosystem that deviating from it (for example, panicking instead of returning an error for an ordinary, expected failure condition) is generally considered surprising and non-idiomatic for most application-level code.

## 18. Function Naming Conventions

- **MixedCaps**, not underscores: `calculateTotal`, not `calculate_total`.
- **Capitalize to export** a function (make it visible outside its package); lowercase to keep it package-private — this is Go's sole visibility mechanism, with no separate access-modifier keywords.
- **Name a function for what it does**, using a verb or verb phrase for functions with side effects (`SaveUser`, `SendEmail`) and a noun phrase for functions that primarily compute and return a value (`Total`, `IsValid`).
- **Avoid redundantly repeating the package name** inside a function's own name, since the package name already prefixes every call site (`strings.ToUpper`, not `strings.StringToUpper`).
- **Constructor-style functions conventionally start with `New`** (`NewClient`, `NewServer`), and functions expected to panic on failure by design are conventionally prefixed with `Must` (`regexp.MustCompile`).

## 19. Common Mistakes and Pitfalls

### 19.1 Forgetting Every Path Must Return a Value

```go
func classify(n int) string {
    if n > 0 {
        return "positive"
    }
    // COMPILE ERROR: missing return
}
```

Fix: ensure every possible execution path explicitly returns, typically with a trailing `return` or an `else` branch — see [Section 4.3](#43-the-function-body-must-cover-every-path-for-non-void-functions).

### 19.2 Overusing Naked Returns in Long Functions

```go
func process() (result int, err error) {
    // ... 40 lines later ...
    return // what does this actually return? the reader has to scroll back up to check
}
```

Naked returns work well in short functions but hurt readability once a function grows — prefer explicit `return result, err` in longer functions, even with named return values declared.

### 19.3 Capturing a Loop Variable Incorrectly (Pre-Go 1.22)

```go
var funcs []func()
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() { fmt.Println(i) }) // all print the same final value on pre-1.22
}
```

Covered in depth in [Section 13](#13-the-classic-loop-variable-capture-pitfall) — shadow the variable (`i := i`) inside the loop body for code that must support pre-1.22 semantics.

### 19.4 Expecting a Variadic Function's Slice Parameter to Be Mutation-Safe by Default

```go
func modify(nums ...int) {
    nums[0] = 999 // mutates the underlying array — which the caller's own slice may share!
}

original := []int{1, 2, 3}
modify(original...)
fmt.Println(original) // [999 2 3] — the caller's slice was mutated
```

When a slice is spread into a variadic call with `...`, the variadic parameter shares the same underlying array as the original slice — mutating elements inside the function is visible to the caller, exactly like passing a slice normally. This is expected slice-aliasing behavior, but easy to overlook specifically in the variadic-spread context.

### 19.5 Assuming Deep Recursion Is Always Safe

```go
func sumTo(n int) int {
    if n == 0 {
        return 0
    }
    return n + sumTo(n-1) // no tail-call optimization — deep n can overflow the stack
}
```

For algorithms that might recurse to very large depths on realistic inputs, prefer an iterative implementation, or restructure with explicit accumulator state managed in a loop.

### 19.6 Ignoring an Error Return Value

```go
result, _ := riskyOperation() // silently discards a potential failure
```

Discarding an error return without at least a deliberate, documented reason is a common source of silent bugs — check and handle (or explicitly and knowingly ignore, with a comment explaining why) every error a function returns.

## 20. Best Practices Summary

1. **Keep functions focused on one clear responsibility** — a function that's hard to name concisely is often a sign it's doing too much.
2. **Follow the `(result, error)` convention** for any function that can fail, with `error` as the final return value.
3. **Use named return values primarily for documentation or when a deferred function needs to modify the result** (e.g., attaching a `Close()` failure to an already-set error), not as a default habit for every function.
4. **Avoid naked returns in anything but short, simple functions**, where the returned values are still obvious at a glance.
5. **Use variadic parameters for genuinely open-ended argument counts**, and remember a spread slice shares its underlying array with the variadic parameter.
6. **Reach for closures to maintain small amounts of private, per-instance state** without resorting to package-level globals.
7. **Be explicit about what a closure captures**, especially in loops, if the code must remain compatible with pre-Go-1.22 semantics.
8. **Prefer iteration over deep recursion** for algorithms whose recursion depth could realistically become large, since Go doesn't guarantee tail-call optimization.
9. **Name functions after what they do**, avoiding redundant restatement of the package name they already live in.
10. **Use named function types** to make function-typed parameters, fields, and variables more self-documenting when the same signature recurs across a codebase.

## 21. Use Case Summary Table

| Technique                             | When to Use                                                          | Example Scenario                                                |
| ------------------------------------- | -------------------------------------------------------------------- | --------------------------------------------------------------- |
| Multiple return values                | A function naturally produces more than one related output           | `(result, error)`, `(value, ok)`                                |
| Named return values                   | Documentation, or a deferred function needs to set the result        | Attaching a cleanup error to a function's own error return      |
| Variadic parameters                   | An unbounded, flexible number of arguments of one type               | `fmt.Println(a ...any)`, a `sum(nums ...int)` helper            |
| Anonymous function (function literal) | A short, one-off callback or a scoped block of code                  | Passed into `sort.Slice`, or an IIFE for local scoping          |
| Closure                               | Small amounts of private state persisting across calls               | A counter, a memoization cache, a configured validator          |
| Higher-order function                 | Behavior needs to be pluggable/parameterized by the caller           | Comparison functions, middleware, event handlers                |
| Named function type                   | The same function signature recurs across a codebase                 | `type Middleware func(http.Handler) http.Handler`               |
| Function as a struct field            | Per-instance, pluggable behavior without a full interface            | A `Validator` struct carrying a `Check func(string) bool` field |
| Recursion                             | The problem is naturally self-similar and recursion depth is bounded | Tree traversal, simple mathematical recursive definitions       |
| `Must`-prefixed function              | The function's contract is "succeed, or the program has a bug"       | `regexp.MustCompile` with a hardcoded, known-valid pattern      |

## 22. References

1. Go Team — _A Tour of Go: Functions_. https://go.dev/tour/basics/4
2. Go Team — _First-Class Functions in Go_ (codewalk). https://go.dev/doc/codewalk/functions
3. Go Team — _Effective Go_ (Functions, Named result parameters). https://go.dev/doc/effective_go
4. Go Language Specification — _Function declarations, Function types, Function literals_. https://go.dev/ref/spec#Function_declarations
5. Go Team — _Go 1.22 Release Notes_ (for-loop variable scoping change). https://go.dev/doc/go1.22
6. Leonardo Souza — _Mastering Functions in Go: A Deep Dive_. https://leogsouza.dev/blog/mastering-functions-in-go-a-deep-dive/
7. S Soumyakanta — _Go Functions: Basics to Advanced Closures_. https://www.s-soumyakanta.com/blog/learn-go-functions-from-basic-concepts-to-advanced-closures
8. KodeKloud — _Return Types: Multiple, Named, Variadic_. https://notes.kodekloud.com/docs/Golang/Using-Functions/Return-Types-Multiple-Named-Variadic/page
9. Caffeine Algorithm — _Functions in Go_. https://caffeinealgorithm.com/blog/functions-in-go
10. Stanza — _Go Syntax Cheatsheet_. https://www.stanza.dev/cheatsheet/go-syntax
11. Abati Babatunde Daniel — _Fundamentals of Functions in Golang_, Medium. https://medium.com/@danielabatibabatunde1/fundamentals-of-functions-in-golang-df4dd0c3072f
12. Go standard library documentation — package `sort` (`sort.Slice`, as a higher-order function example). https://pkg.go.dev/sort
13. Go standard library documentation — package `regexp` (`MustCompile` naming convention). https://pkg.go.dev/regexp
