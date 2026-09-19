<div align="center">
  <h1>Iterations</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 18, 2026</sub>
</div>

## Table of Contents

1. [Go Has Only One Looping Keyword: `for`](#1-go-has-only-one-looping-keyword-for)
2. [The Classic Three-Component `for` Loop](#2-the-classic-three-component-for-loop)
3. [The "While" Form: Condition-Only](#3-the-while-form-condition-only)
4. [The Infinite Loop](#4-the-infinite-loop)
5. [The `range` Clause: An Overview](#5-the-range-clause-an-overview)
6. [Ranging Over a Slice or Array](#6-ranging-over-a-slice-or-array)
7. [Ranging Over a String](#7-ranging-over-a-string)
8. [Ranging Over a Map](#8-ranging-over-a-map)
9. [Ranging Over a Channel](#9-ranging-over-a-channel)
10. [Ranging Over an Integer (Go 1.22+)](#10-ranging-over-an-integer-go-122)
11. [Ranging Over a Function — Range-Over-Func Iterators (Go 1.23+)](#11-ranging-over-a-function--range-over-func-iterators-go-123)
12. [`break` and `continue`](#12-break-and-continue)
13. [Labeled `break` and `continue`](#13-labeled-break-and-continue)
14. [The Go 1.22 Loop Variable Scoping Change](#14-the-go-122-loop-variable-scoping-change)
15. [Modifying a Slice/Map While Iterating](#15-modifying-a-slicemap-while-iterating)
16. [`goto`: A Rarely Used Alternative Control Flow](#16-goto-a-rarely-used-alternative-control-flow)
17. [Nested Loops](#17-nested-loops)
18. [Common Mistakes and Pitfalls](#18-common-mistakes-and-pitfalls)
19. [Best Practices Summary](#19-best-practices-summary)
20. [Use Case Summary Table](#20-use-case-summary-table)
21. [References](#21-references)

## 1. Go Has Only One Looping Keyword: `for`

Unlike most C-family languages, Go does not have separate `while`, `do-while`, or `foreach` keywords — **`for` is the only looping construct in the language**, and it covers every looping need through several different forms of its syntax. This is a deliberate simplification consistent with Go's overall design philosophy of having a small, easy-to-remember core language.

The forms `for` can take are:

1. The classic three-component form (`for init; condition; post { }`)
2. The condition-only form, equivalent to a `while` loop
3. The infinite form, with no clause at all
4. The `range` form, for iterating over a collection or sequence

## 2. The Classic Three-Component `for` Loop

This is the most traditional form, familiar from C, Java, and similar languages:

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
// 0 1 2 3 4
```

The three components, separated by semicolons, are:

- **Init statement**: runs once, before the loop begins (`i := 0`). Its scope is limited to the loop.
- **Condition**: evaluated before every iteration; the loop continues as long as it's `true` (`i < 5`).
- **Post statement**: runs after every iteration, before the condition is checked again (`i++`).

Unlike C, Go does **not** require (or even allow) parentheses around the clause, but the opening brace `{` of the loop body **must** appear on the same line as the `for` clause — Go's automatic semicolon insertion rules make an opening brace on its own line a syntax error here.

```go
// for (i := 0; i < 5; i++) { ... } // COMPILE ERROR: parentheses not allowed
```

Any of the three components can be omitted:

```go
i := 0
for ; i < 5; i++ { // init omitted (already declared above)
    fmt.Println(i)
}

for i := 0; i < 5; { // post omitted (incremented manually inside the body)
    fmt.Println(i)
    i++
}
```

## 3. The "While" Form: Condition-Only

When only the condition is present (with both semicolons dropped), `for` behaves exactly like a `while` loop in other languages:

```go
n := 1
for n < 100 {
    n *= 2
}
fmt.Println(n) // 128
```

This is the idiomatic way to write a "loop while some condition holds" in Go — there is no separate `while` keyword, so this condition-only `for` form fills that role entirely.

## 4. The Infinite Loop

Omitting the condition entirely (and both semicolons) produces an infinite loop, which continues forever until explicitly stopped with `break`, `return`, or a program-terminating call:

```go
for {
    fmt.Println("looping forever...")
    break // without this, the loop would never end
}
```

This form is commonly used for server-style loops (accepting connections, processing a queue) that are only expected to terminate under specific, checked conditions inside the loop body, or possibly never terminate at all during the program's lifetime (e.g., a long-running worker goroutine's main loop).

```go
for {
    task, ok := getNextTask()
    if !ok {
        break // exit when there's no more work
    }
    process(task)
}
```

## 5. The `range` Clause: An Overview

The `range` clause of a `for` statement iterates over elements of various built-in and (since Go 1.23) custom iterable types. Its general form:

```go
for key, value := range collection {
    // ...
}
```

Both `key` and `value` are optional to actually use — you can capture just the first result, or neither:

```go
for key := range collection { }     // only the first result
for key, value := range collection { } // both results
for range collection { }              // neither — iterate purely for side effects / count
```

What exactly `range` produces as its first and second result depends entirely on the type being ranged over — covered type by type in the sections below.

## 6. Ranging Over a Slice or Array

For a slice or array, `range` produces the **index** as the first value and a **copy of the element** at that index as the second value:

```go
fruits := []string{"apple", "banana", "cherry"}

for i, fruit := range fruits {
    fmt.Println(i, fruit)
}
// 0 apple
// 1 banana
// 2 cherry

for i := range fruits { // indices only
    fmt.Println(i)
}

for _, fruit := range fruits { // values only
    fmt.Println(fruit)
}
```

**Important:** the second value is a **copy** of the element, not a reference to it — mutating it inside the loop body does not affect the original slice/array:

```go
nums := []int{1, 2, 3}
for _, n := range nums {
    n *= 2 // modifies only the loop-local copy
}
fmt.Println(nums) // [1 2 3] — unchanged
```

To mutate elements in place, index back into the slice explicitly:

```go
for i := range nums {
    nums[i] *= 2 // this DOES modify the original slice
}
fmt.Println(nums) // [2 4 6]
```

`range` over a `nil` slice is safe and simply performs zero iterations, exactly like ranging over an empty slice.

## 7. Ranging Over a String

Ranging over a `string` is a special case worth understanding carefully: Go strings are sequences of **bytes**, but string literals and most real-world text are UTF-8 encoded, where a single logical character ("rune") can span multiple bytes. `range` over a string decodes the UTF-8 encoding for you automatically, yielding the **byte index** where each rune starts, and the **rune** (an `int32` representing a Unicode code point) itself:

```go
s := "héllo" // é is a 2-byte UTF-8 sequence

for i, r := range s {
    fmt.Printf("%d: %c (%d)\n", i, r, r)
}
// 0: h (104)
// 1: é (233)   <- note: index jumps from 1 to 3, skipping index 2
// 3: l (108)
// 4: l (108)
// 5: o (111)
```

Notice the index jumps from `1` to `3` — this is because `é` occupies **two bytes** in the underlying UTF-8 encoding, and the reported index is always the **byte offset** where that rune begins within the string, not a simple rune count. If you index into the string directly with `s[i]` instead of using `range`, you get individual **bytes**, not runes — for genuinely correct handling of multi-byte characters, `range` (or explicit use of the `unicode/utf8` package) is the correct approach, not raw byte indexing.

```go
fmt.Println(len(s))    // 6 — byte length, not the number of visible characters (5)
fmt.Println(s[1])       // a single byte of é's 2-byte encoding, NOT the full character
```

## 8. Ranging Over a Map

For a map, `range` produces the **key** as the first value and the corresponding **value** as the second:

```go
ages := map[string]int{"Alice": 30, "Bob": 25}

for name, age := range ages {
    fmt.Println(name, age)
}

for name := range ages { // keys only
    fmt.Println(name)
}
```

**Iteration order over a map is intentionally randomized** by the Go runtime — running the same loop multiple times can (and typically does) visit entries in a different order each time. If a deterministic, reproducible order is needed, collect the keys into a slice, sort that slice, and iterate through the sorted keys instead.

It is safe to `delete` an entry from a map during a `range` loop over that same map (including the entry currently being visited), but inserting a **new** key during iteration has explicitly unspecified behavior regarding whether that new key will be visited later in the same loop — avoid inserting new keys mid-iteration if predictable behavior matters.

## 9. Ranging Over a Channel

Ranging over a channel receives values from it repeatedly, until the channel is both **closed** and fully **drained** of any buffered values, at which point the loop exits automatically:

```go
func producer(ch chan<- int) {
    defer close(ch) // signals "no more values" so the range loop below can terminate
    for i := 0; i < 5; i++ {
        ch <- i
    }
}

func main() {
    ch := make(chan int)
    go producer(ch)

    for v := range ch {
        fmt.Println(v)
    }
    fmt.Println("channel closed and drained")
}
```

`range` over a channel produces only a **single** value per iteration (the received element) — there is no second "key" component, unlike slices and maps. If the channel is never closed and no more values ever arrive, a `range` loop over it blocks forever, which is a common cause of goroutine leaks if not paired with a proper cancellation mechanism.

## 10. Ranging Over an Integer (Go 1.22+)

Go 1.22 added the ability to `range` directly over an integer value, providing a more concise way to express a simple counted loop:

```go
for i := range 5 {
    fmt.Println(i)
}
// 0 1 2 3 4
```

This is equivalent to the traditional three-component form `for i := 0; i < 5; i++`, but reads more concisely for the common case of "do something N times, optionally using the index." The value component can also be omitted entirely when only the repetition count matters, not the index itself:

```go
for range 3 {
    fmt.Println("hello")
}
// hello
// hello
// hello
```

Ranging over an integer only ever counts upward from `0`; it doesn't support a custom start, step, or descending direction — for anything beyond the simple `0` to `n-1` case, the classic three-component form is still the right tool.

## 11. Ranging Over a Function — Range-Over-Func Iterators (Go 1.23+)

Go 1.23 extended `range` to work directly over specially-shaped **functions**, called **iterators**, enabling custom types (trees, linked lists, generators, database cursors, and so on) to support idiomatic `for range` iteration without needing to first materialize a full slice or expose their internals.

### 11.1 The Iterator Function Shapes

An iterator is a function accepting a single `yield` function as its argument, matching one of three shapes (defined generically in the standard library's `iter` package as `iter.Seq[V]` and `iter.Seq2[K, V]`):

```go
type Seq[V any]     func(yield func(V) bool)
type Seq2[K, V any]  func(yield func(K, V) bool)
```

The iterator calls `yield` once per element it wants to produce; `yield` returns `true` if the loop should continue, or `false` if the loop was stopped early (via `break`, an early `return` in the loop body, or an error) — a well-behaved iterator must stop producing further values as soon as `yield` returns `false`.

```go
func Backward[E any](s []E) func(func(int, E) bool) {
    return func(yield func(int, E) bool) {
        for i := len(s) - 1; i >= 0; i-- {
            if !yield(i, s[i]) {
                return // the consuming loop stopped early — stop producing values
            }
        }
    }
}
```

### 11.2 Using an Iterator with `range`

```go
nums := []string{"a", "b", "c"}

for i, v := range Backward(nums) {
    fmt.Println(i, v)
}
// 2 c
// 1 b
// 0 a
```

From the caller's perspective, this reads exactly like ranging over a slice or map — the `for range` syntax itself doesn't change; what's new is that arbitrary functions matching the `Seq`/`Seq2` shape can now be the thing being ranged over, not just built-in collection types.

### 11.3 `break`, `continue`, and `return` Inside a Range-Over-Func Loop

The compiler translates the loop body's control-flow statements into the appropriate return values from `yield`: a `break` becomes `yield` returning `false` (stopping further iteration); an implicit or explicit `continue` at the end of the loop body becomes `yield` returning `true` (letting the iterator proceed to the next value).

```go
for i, v := range Backward(nums) {
    if v == "b" {
        break // translates to yield returning false — Backward's loop stops immediately
    }
    fmt.Println(i, v)
}
```

`defer` and `panic` behave consistently with how they work in ordinary Go code: a `defer` inside the loop body runs when the _surrounding function containing the loop_ returns, while a `defer` inside the iterator function itself runs when the _iterator function_ returns — these are two genuinely separate function scopes, even though the syntax makes the iteration feel seamless.

### 11.4 A Real Motivating Use Case

Range-over-func is particularly valuable for iterating over custom containers (trees, linked lists, database result cursors) or generating values on the fly, without needing to first collect everything into an intermediate slice held entirely in memory — this matters for very large or even conceptually infinite sequences, where materializing a full slice upfront would be wasteful or simply impossible.

```go
func Count(start, end int) func(func(int) bool) {
    return func(yield func(int) bool) {
        for i := start; i < end; i++ {
            if !yield(i) {
                return
            }
        }
    }
}

for n := range Count(10, 15) {
    fmt.Println(n) // 10 11 12 13 14
}
```

**Practical takeaway:** in everyday application code, you're more likely to _consume_ iterators exposed by a library (an `All()` method on a custom collection type, for instance) than to write your own from scratch — but understanding the `yield`-based shape helps make sense of what's actually happening syntactically when you see `for range someFunction`.

## 12. `break` and `continue`

- **`break`** immediately exits the innermost enclosing `for` loop (or `switch`/`select`), skipping any remaining iterations entirely.
- **`continue`** skips the rest of the current iteration's body and jumps directly to the next iteration (evaluating the post statement and condition again, for a three-component loop).

```go
for i := 0; i < 10; i++ {
    if i == 5 {
        break // stop the loop entirely once i reaches 5
    }
    fmt.Println(i)
}
// 0 1 2 3 4

for i := 0; i < 10; i++ {
    if i%2 != 0 {
        continue // skip printing odd numbers
    }
    fmt.Println(i)
}
// 0 2 4 6 8
```

## 13. Labeled `break` and `continue`

By default, `break` and `continue` only affect the **innermost** enclosing loop. To break or continue an **outer** loop from within a nested loop, Go supports labeling a loop and referencing that label explicitly:

```go
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 1 {
            continue outer // skips the rest of the INNER loop AND the rest of THIS OUTER iteration
        }
        fmt.Println(i, j)
    }
}
// 0 0
// 1 0
// 2 0
```

```go
search:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if i == 1 && j == 1 {
            break search // exits BOTH loops immediately, not just the inner one
        }
        fmt.Println(i, j)
    }
}
// 0 0
// 0 1
// 0 2
// 1 0
```

A label immediately precedes the statement it labels (typically a `for` loop), and is referenced by name in the corresponding `break`/`continue` statement. Labeled `break`/`continue` are also valid with `switch` and `select` statements when nested inside a loop, letting you break out of an outer loop from within an inner `switch`/`select` — something a plain, unlabeled `break` cannot do, since it would otherwise only exit the innermost `switch`/`select`.

## 14. The Go 1.22 Loop Variable Scoping Change

This is an important, relatively recent change to be aware of, since it affects the behavior of closures created inside loops.

**Before Go 1.22:** a `for` loop's variables (declared in the init statement, or the single variable in a `range` clause) were a **single set of variables, reused across every iteration** — the loop body ran repeatedly, reassigning the same variable(s) each time, rather than creating fresh ones per iteration.

```go
// Pre-Go 1.22 behavior
var funcs []func()
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() { fmt.Println(i) })
}
for _, f := range funcs {
    f() // prints 3, 3, 3 — every closure captured the SAME, shared i
}
```

**Since Go 1.22:** each iteration of a `for` loop gets its **own, fresh copy** of the loop variable(s), scoped to just that single iteration — closures created inside the loop body now capture a distinct variable per iteration, matching what many developers intuitively expect:

```go
// Go 1.22+ behavior (same source code as above)
var funcs []func()
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() { fmt.Println(i) })
}
for _, f := range funcs {
    f() // prints 0, 1, 2 — each closure captured its OWN iteration's i
}
```

**Which behavior applies depends on the Go language version a module declares** in its `go.mod` file's `go` directive, not merely which Go toolchain version is installed — a module that specifies `go 1.21` or earlier retains the old semantics even when built with a newer Go toolchain, while a module specifying `go 1.22` or later gets the new, per-iteration variable semantics. This was a deliberate, carefully staged language change specifically designed to avoid silently breaking existing code that (rarely, but occasionally) depended on the old shared-variable behavior.

**Practical implication:** if your module targets Go 1.22 or later, the classic "closure captures the wrong loop variable" bug largely disappears for ordinary loop-with-closure code — but being explicit about what a closure captures (`i := i` inside the loop body) remains a reasonable defensive habit in code that must remain compatible with older language versions, or in any codebase where the target Go version isn't immediately obvious to every reader.

## 15. Modifying a Slice/Map While Iterating

### 15.1 Slices

Appending to a slice while ranging over it does **not** extend the range of the loop — `range` evaluates the length of the slice once, at the start of the loop, so elements appended during iteration are never visited in that same loop:

```go
nums := []int{1, 2, 3}
for _, n := range nums {
    if n == 2 {
        nums = append(nums, 99) // does NOT cause 99 to be visited in this loop
    }
    fmt.Println(n)
}
// 1 2 3 — the appended 99 is never printed
```

### 15.2 Maps

As noted in [Section 8](#8-ranging-over-a-map), deleting the current or a not-yet-visited key during a map `range` loop is explicitly safe and well-defined; inserting a new key during that same iteration has unspecified visitation behavior for that particular loop run.

## 16. `goto`: A Rarely Used Alternative Control Flow

Go includes a `goto` statement, which transfers control directly to a labeled statement within the same function. It exists mainly for rare, specific situations (some generated code, certain error-handling patterns predating widespread `defer` usage, or breaking out of deeply nested structures in ways `break`/`continue` labels don't quite cover) rather than as a routine looping tool:

```go
i := 0
loop:
if i < 5 {
    fmt.Println(i)
    i++
    goto loop
}
```

`goto` is subject to strict rules — it cannot jump into the scope of a variable declaration it would skip over, and cannot jump into a block from outside it. In virtually all idiomatic Go code, a structured `for` loop (with `break`/`continue`, including labeled forms) is preferred over `goto` for looping constructs; `goto` sees genuine, occasional use mostly outside of loops entirely (for instance, jumping to a single shared cleanup label near the end of a function in code written before `defer` covered that need as cleanly).

## 17. Nested Loops

Loops can be nested arbitrarily deep, exactly as in most languages, with each level of nesting able to use its own loop variables independently:

```go
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        fmt.Println(i, j)
    }
}
```

As shown in [Section 13](#13-labeled-break-and-continue), labeled `break`/`continue` are specifically how Go lets an inner loop affect an outer loop's control flow directly, since unlabeled `break`/`continue` always apply to only the innermost enclosing loop by default.

## 18. Common Mistakes and Pitfalls

### 18.1 Expecting `range` Over a Slice to Give a Mutable Reference

```go
for _, v := range items {
    v.Field = newValue // modifies only the loop-local copy of the struct, not the original
}
```

Index back into the slice explicitly (`items[i].Field = newValue`) to genuinely mutate elements during iteration — see [Section 6](#6-ranging-over-a-slice-or-array).

### 18.2 Indexing a String Byte-by-Byte and Expecting Characters

```go
s := "héllo"
for i := 0; i < len(s); i++ {
    fmt.Printf("%c", s[i]) // corrupts multi-byte characters like é
}
```

Use `range` over the string (or the `unicode/utf8` package) to correctly decode multi-byte UTF-8 runes — see [Section 7](#7-ranging-over-a-string).

### 18.3 Relying on Map Iteration Order

```go
for k, v := range myMap {
    // assuming a specific, repeatable order — WRONG
}
```

Map iteration order is intentionally randomized; sort keys explicitly if deterministic output is needed — see [Section 8](#8-ranging-over-a-map).

### 18.4 Forgetting a Range-Over-Channel Loop Blocks Forever on an Unclosed Channel

```go
for v := range ch { // blocks forever if ch is never closed and no more values arrive
    fmt.Println(v)
}
```

Ensure the sending side closes the channel once it's done, or pair the loop with a `select`-based cancellation mechanism if indefinite blocking isn't acceptable — see [Section 9](#9-ranging-over-a-channel).

### 18.5 Assuming Pre-1.22 Closure-Capture Semantics on a Go 1.22+ Module (or Vice Versa)

```go
for i := 0; i < 3; i++ {
    go func() { fmt.Println(i) }()
}
```

The actual behavior here depends on the Go language version declared in `go.mod`, not just the installed toolchain — be aware of which semantics apply to your specific module, especially when reading or maintaining code written before this change, or when a codebase mixes files targeting different language versions.

### 18.6 Appending to a Slice Mid-Range and Expecting the New Elements to Be Visited

```go
for _, v := range nums {
    nums = append(nums, v*2) // appended elements are never visited in THIS loop run
}
```

`range` captures the slice's length once at the start; elements appended during iteration aren't visited in that same loop — see [Section 15.1](#151-slices).

### 18.7 Writing a Range-Over-Func Iterator That Ignores `yield`'s Return Value

```go
func BadIterator(yield func(int) bool) {
    for i := 0; i < 1000000; i++ {
        yield(i) // ignoring the return value — keeps producing values even after the consumer stopped early!
    }
}
```

A well-behaved iterator must check `yield`'s return value and stop producing further values as soon as it returns `false` — ignoring it defeats the purpose of early termination via `break` and can cause wasted work or even bugs in the consuming loop.

## 19. Best Practices Summary

1. **Use the condition-only `for` form for "while"-style loops** — there's no separate `while` keyword, and this is the idiomatic equivalent.
2. **Use `range` whenever iterating over a slice, array, map, string, or channel**, rather than manually managing an index, for both clarity and correctness (especially for strings, where manual byte indexing can corrupt multi-byte characters).
3. **Index back into a slice explicitly** when you need to mutate elements during iteration — the `range`-provided value is always a copy.
4. **Never assume a particular map iteration order**; sort keys explicitly when deterministic output matters.
5. **Use labeled `break`/`continue`** when a nested loop needs to affect an outer loop's control flow, rather than resorting to extra boolean flags or `goto`.
6. **Be aware of the Go 1.22 loop-variable scoping change** and which semantics apply to code you're reading or writing, based on the module's declared `go` version.
7. **Use `for range n` (Go 1.22+) for simple counted repetition** when the classic three-component form's extra ceremony isn't needed.
8. **Ensure a channel is eventually closed** (or paired with a cancellation mechanism) before relying on `for range` over it, to avoid blocking forever.
9. **When writing a range-over-func iterator, always check and respect `yield`'s return value**, stopping promptly once it returns `false`.
10. **Reserve `goto` for the narrow, specific situations it's actually suited to** — structured `for` loops with `break`/`continue` (including labeled forms) cover the overwhelming majority of everyday looping needs.

## 20. Use Case Summary Table

| Technique                     | When to Use                                                         | Example Scenario                                                             |
| ----------------------------- | ------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| Three-component `for`         | Classic counted loop with explicit init/condition/post              | `for i := 0; i < n; i++ { ... }`                                             |
| Condition-only `for`          | A "while"-style loop                                                | `for balance > 0 { ... }`                                                    |
| Infinite `for`                | A loop whose termination is conditional and checked inside the body | A server accept loop, a worker's main processing loop                        |
| `range` over slice/array      | Iterating elements with index and/or value                          | Processing every item in a collected list of results                         |
| `range` over string           | Correctly decoding Unicode text character by character              | Text processing that must handle multi-byte characters correctly             |
| `range` over map              | Iterating key-value pairs                                           | Printing or processing all entries in a configuration map                    |
| `range` over channel          | Consuming a stream of values until the sender signals completion    | A worker consuming jobs from a channel until it's closed                     |
| `range` over int (1.22+)      | Simple counted repetition, optionally using the index               | `for i := range 10 { ... }`, `for range 3 { ... }`                           |
| `range` over function (1.23+) | Custom, lazy, or large/infinite iteration sources                   | Iterating a custom tree/list type, or a generator producing values on demand |
| Labeled `break`/`continue`    | Controlling an outer loop from within a nested loop                 | Searching a 2D grid and stopping all loops once a match is found             |
| `goto`                        | Rare, specific control-flow needs outside typical looping           | Legacy-style shared cleanup jump, certain generated code                     |

## 21. References

1. Go Team — _A Tour of Go: For_, and the following pages on `for` variants. https://go.dev/tour/flowcontrol/1
2. Go Team — _Effective Go_ (For, Control structures). https://go.dev/doc/effective_go
3. Go Language Specification — _For statements_. https://go.dev/ref/spec#For_statements
4. Go Team — _Go 1.22 Release Notes_ (range-over-integer, for-loop variable scoping change). https://go.dev/doc/go1.22
5. Go Team — _Go 1.23 Release Notes_ (range-over-function iterators). https://go.dev/doc/go1.23
6. Go Team — _Range Over Function Types_, The Go Blog. https://go.dev/blog/range-functions
7. Go Team — _Go Wiki: Rangefunc Experiment (FAQ)_. https://go.dev/wiki/RangefuncExperiment
8. Go standard library documentation — package `iter`. https://pkg.go.dev/iter
9. DoltHub Blog — _Go range iterators demystified_. https://www.dolthub.com/blog/2024-07-12-golang-range-iters-demystified/
10. Sandeep — _Go 1.23: Range-over-Function Iterators Now Part of the Language_, Medium. https://medium.com/@ksandeeptech07/go-1-23-range-over-function-iterators-now-part-of-the-language-f99ad7072713
11. Pairs Engineering — _A look at iterators in Go (Golang)_. https://medium.com/eureka-engineering/a-look-at-iterators-in-go-f8e86062937c
12. Ardan Labs — _Range-Over Functions in Go_. https://www.ardanlabs.com/blog/2024/04/range-over-functions-in-go.html
