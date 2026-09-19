<div align="center">
  <h1>Defer</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 18, 2026</sub>
</div>

## Table of Contents

1. [What Is `defer`?](#1-what-is-defer)
2. [Basic Syntax and Behavior](#2-basic-syntax-and-behavior)
3. [Multiple Defers: LIFO Order](#3-multiple-defers-lifo-order)
4. [Arguments Are Evaluated Immediately, at Defer Time](#4-arguments-are-evaluated-immediately-at-defer-time)
5. [The Receiver Is Also Evaluated at Defer Time](#5-the-receiver-is-also-evaluated-at-defer-time)
6. [Capturing a Value at Execution Time Instead](#6-capturing-a-value-at-execution-time-instead)
7. [Defer and Named Return Values](#7-defer-and-named-return-values)
8. [Defer and `panic`/`recover`](#8-defer-and-panicrecover)
9. [Common Use Cases for Defer](#9-common-use-cases-for-defer)
10. [Defer Inside Loops: The Classic Pitfall](#10-defer-inside-loops-the-classic-pitfall)
11. [How Defer Works Internally](#11-how-defer-works-internally)
12. [Performance: Open-Coded, Stack-Allocated, and Heap-Allocated Defers](#12-performance-open-coded-stack-allocated-and-heap-allocated-defers)
13. [Defer with Method Values](#13-defer-with-method-values)
14. [Defer and Multiple Return Paths](#14-defer-and-multiple-return-paths)
15. [Common Mistakes and Pitfalls](#15-common-mistakes-and-pitfalls)
16. [Best Practices Summary](#16-best-practices-summary)
17. [Use Case Summary Table](#17-use-case-summary-table)
18. [References](#18-references)

## 1. What Is `defer`?

`defer` is a statement that schedules a function call to run **just before the surrounding function returns**, regardless of how that function exits — whether it reaches a normal `return`, falls off the end, or is unwinding due to a `panic`. The word "defer" means exactly what it says: postpone this call until later, specifically until the enclosing function is about to give control back to its caller.

```go
func main() {
    defer fmt.Println("This runs last")
    fmt.Println("This runs first")
}
// Output:
// This runs first
// This runs last
```

`defer` is most commonly used for cleanup work that must happen no matter which code path a function takes to exit — closing files, unlocking mutexes, closing database connections, or recovering from a panic.

## 2. Basic Syntax and Behavior

The syntax is simply the `defer` keyword followed by a function or method call:

```go
defer functionName(arguments)
defer object.Method(arguments)
defer func() {
    // arbitrary code in an anonymous function
}()
```

The deferred call is **registered** at the point the `defer` statement executes, but its **body doesn't actually run** until the enclosing function is about to return:

```go
func greet() {
    defer fmt.Println("goodbye")
    fmt.Println("hello")
}
// Output:
// hello
// goodbye
```

## 3. Multiple Defers: LIFO Order

When a function contains several `defer` statements, they execute in **Last-In-First-Out (LIFO)** order — the most recently deferred call runs first, and the first one deferred runs last:

```go
func main() {
    defer fmt.Println("first")
    defer fmt.Println("second")
    defer fmt.Println("third")
}
// Output:
// third
// second
// first
```

A useful way to think about this: reading the deferred calls **bottom-up** in the source gives you the actual order they'll execute in. This LIFO ordering is deliberate and important for resource management — it means resources are released in the reverse order they were acquired, which is exactly the safe unwinding order you want (the last thing opened is the first thing closed):

```go
func process() error {
    db, err := openDatabase()
    if err != nil {
        return err
    }
    defer db.Close() // closed LAST (since it was acquired FIRST)

    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // closed FIRST (since it was acquired LAST, or defensively rolled back if never committed)

    // ... use tx ...
    return tx.Commit()
}
```

## 4. Arguments Are Evaluated Immediately, at Defer Time

This is one of the most important and most frequently misunderstood rules about `defer`: **the arguments to a deferred function call are evaluated the moment the `defer` statement runs, not when the deferred call actually executes later.** Only the _execution_ of the function body is postponed — the _argument values_ are captured right away.

```go
func main() {
    i := 0
    defer fmt.Println("Deferred print:", i) // i is evaluated NOW, capturing 0
    i++
    fmt.Println("Regular print:", i) // 1
}
// Output:
// Regular print: 1
// Deferred print: 0
```

Even though `i` becomes `1` before the function ends, the deferred `fmt.Println` already locked in the value `0` at the moment the `defer` statement executed — it prints `0`, not the value `i` holds later.

## 5. The Receiver Is Also Evaluated at Defer Time

The same "evaluate now, call later" rule applies to a method's **receiver**, not just its explicit arguments — the receiver expression is evaluated and effectively "frozen" the instant the `defer` statement runs:

```go
type Person struct {
    Name string
}

func (p Person) SayHi() {
    fmt.Println("Hi, my name is", p.Name)
}

func main() {
    writer := Person{Name: "Joe"}
    defer writer.SayHi() // the VALUE of writer (a copy, since SayHi has a value receiver) is captured NOW
    writer.Name = "Aiden"
}
// Output:
// Hi, my name is Joe
```

Because `SayHi` has a **value receiver**, `defer writer.SayHi()` captures a _copy_ of `writer` as it existed at the moment of the `defer` statement — the later mutation of `writer.Name` has no effect on what gets printed. If `SayHi` instead had a **pointer receiver** (`func (p *Person) SayHi()`), the deferred call would capture the _pointer_ to `writer`, and would see the mutation, since the pointer still refers to the same, later-mutated struct:

```go
func (p *Person) SayHi() { // pointer receiver
    fmt.Println("Hi, my name is", p.Name)
}

writer := Person{Name: "Joe"}
defer writer.SayHi() // captures &writer now, but *reads* p.Name later, when SayHi actually runs
writer.Name = "Aiden"
// Output: Hi, my name is Aiden
```

## 6. Capturing a Value at Execution Time Instead

If you actually want a deferred call to see a variable's value **as of when it runs**, rather than as of when it was deferred, wrap it in a closure (an anonymous function with no arguments) instead of passing the variable as a direct argument:

```go
func main() {
    i := 0
    defer func() {
        fmt.Println(i) // reads i from the enclosing scope WHEN THE CLOSURE ACTUALLY RUNS
    }()
    i++
}
// Output: 1
```

Here, nothing is "evaluated" at defer time except the closure itself (which captures `i` by reference, as any Go closure does) — the `fmt.Println(i)` inside the closure body only actually runs, and only reads `i`, once the closure executes later. This distinction — passing a value directly as an argument (frozen immediately) versus reading it inside a closure body (read live, later) — is the key mechanism for controlling exactly when a deferred operation observes a variable's value.

## 7. Defer and Named Return Values

A deferred function can read **and modify** a surrounding function's named return values, because a deferred call runs _after_ the `return` statement has set those values, but _before_ the function actually hands control back to its caller. This is one of `defer`'s most powerful — and most easily misused — capabilities.

```go
func c() (i int) {
    defer func() { i++ }()
    return 1
}

fmt.Println(c()) // 2
```

Here's the sequence: `return 1` sets the named return value `i` to `1`; then, before the function truly returns, the deferred closure runs and increments `i` to `2`; only then does the function actually return `2` to its caller.

### 7.1 The Canonical Use Case: Modifying an Error Return Value

This mechanism is commonly used to enrich or override a function's error return based on cleanup that happens in a defer — for example, surfacing a failure from closing a resource that would otherwise be silently discarded:

```go
func writeToFile(path string, data []byte) (err error) {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = cerr // only overwrite err if nothing else already failed
        }
    }()

    _, err = f.Write(data)
    return err
}
```

Without a named return value here, the deferred closure would have no way to influence what the function ultimately returns — this pattern specifically relies on the named return `err` being a real, addressable variable in the function's scope that the defer can read and reassign.

### 7.2 A Caution

While powerful, mutating named return values from a defer can make control flow noticeably harder to follow — a reader has to know this specific Go behavior to correctly predict what a function returns. Use it deliberately and sparingly (primarily for the "wrap/attach a cleanup error" pattern above), and consider a code comment noting that the defer modifies the named return, since it's easy to miss on a quick read.

## 8. Defer and `panic`/`recover`

`defer` is the mechanism that makes `recover()` possible, and understanding their interaction is essential for building robust Go programs.

### 8.1 Deferred Calls Still Run During a Panic

When a function panics, Go doesn't immediately terminate — it stops the function's normal execution, but **still runs all of that function's deferred calls**, in the usual LIFO order, before propagating the panic up to the caller (which then behaves the same way: its own deferreds run, then the panic propagates further, and so on, until either something calls `recover()` or the panic reaches the top of the program and crashes it).

```go
func f() {
    defer fmt.Println("f: deferred cleanup ran")
    panic("something went wrong")
}

func main() {
    defer fmt.Println("main: deferred cleanup ran")
    f()
    fmt.Println("this line never runs")
}
// Output:
// f: deferred cleanup ran
// main: deferred cleanup ran
// panic: something went wrong
// (plus a stack trace, and the program exits)
```

This guarantee — that deferred cleanup always runs, even when a function panics — is precisely why `defer` is the idiomatic tool for resource cleanup in Go, rather than relying on code simply "reaching" a cleanup line at the end of a function, which a panic would otherwise skip entirely.

### 8.2 `recover()` Only Works Inside a Deferred Function

`recover()` stops a panic from propagating further, but it only has an effect when called **directly inside a deferred function** — calling it in ordinary (non-deferred) code does nothing useful, since by the time ordinary code would call it, no panic is actually in progress there yet.

```go
func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
        }
    }()
    result = a / b // panics if b == 0
    return
}

r, err := safeDivide(10, 0)
fmt.Println(r, err) // 0 recovered from panic: runtime error: integer divide by zero
```

This combines two of the ideas already covered: the deferred closure both calls `recover()` to stop the panic, and assigns to the named return value `err` (per [Section 7](#7-defer-and-named-return-values)) so the caller receives a proper error instead of the program crashing.

## 9. Common Use Cases for Defer

### 9.1 Closing Resources

```go
file, err := os.Open("data.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close() // guaranteed to run, however the function exits
```

### 9.2 Unlocking a Mutex

```go
func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock() // idiomatic: the line right after Lock()
    c.count++
}
```

The idiom of putting `defer mu.Unlock()` immediately after `mu.Lock()` is extremely common in Go, since it guarantees the lock is released on every exit path from the function — including any early `return` or an unexpected `panic` — without needing to duplicate the unlock call at every exit point.

### 9.3 Closing an HTTP Response Body

```go
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close() // must always be closed to avoid leaking connections
```

### 9.4 Logging Function Entry/Exit or Timing

```go
func process() {
    defer func(start time.Time) {
        fmt.Println("process took", time.Since(start))
    }(time.Now()) // note: time.Now() is evaluated NOW, at defer time — see Section 4

    // ... do work ...
}
```

### 9.5 Rolling Back a Transaction on Failure

```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback() // if Commit() succeeded already, Rollback() on a committed tx is a safe no-op in most drivers

// ... perform operations ...
return tx.Commit()
```

## 10. Defer Inside Loops: The Classic Pitfall

`defer` schedules a call to run when the **enclosing function** returns — not when the current loop iteration ends, and not when the current block ends. Placing `defer` inside a loop that runs many times, especially one processing an unbounded or large number of items, accumulates deferred calls that all pile up until the whole function finally returns:

```go
// PROBLEMATIC: opens many files, but doesn't close any of them until processFiles() itself returns
func processFiles(paths []string) error {
    for _, path := range paths {
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close() // accumulates — none of these run until processFiles() returns!
        // ... process file ...
    }
    return nil
}
```

If `paths` contains thousands of entries, this function can hold thousands of file handles open simultaneously, potentially exhausting the operating system's per-process file-descriptor limit — a real and common resource-leak bug that Go's compiler does not warn about, since the code is perfectly valid, just poorly suited to a loop with many iterations.

### 10.1 Fix 1: Wrap Each Iteration in Its Own Function

```go
func processFiles(paths []string) error {
    for _, path := range paths {
        if err := processOneFile(path); err != nil {
            return err
        }
    }
    return nil
}

func processOneFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close() // runs at the end of THIS function call, i.e. after each individual file
    // ... process file ...
    return nil
}
```

### 10.2 Fix 2: Use an Explicit Anonymous Function Scope Inside the Loop

```go
func processFiles(paths []string) error {
    for _, path := range paths {
        err := func() error {
            file, err := os.Open(path)
            if err != nil {
                return err
            }
            defer file.Close() // scoped to this anonymous function, runs each iteration
            // ... process file ...
            return nil
        }()
        if err != nil {
            return err
        }
    }
    return nil
}
```

Both fixes work by giving each loop iteration its **own function scope**, so `defer` has a smaller, iteration-sized function to attach to instead of the whole outer loop's enclosing function.

## 11. How Defer Works Internally

Conceptually, each `defer` statement in a function that hasn't been optimized into the "open-coded" fast path (see [Section 12](#12-performance-open-coded-stack-allocated-and-heap-allocated-defers)) creates a small record describing the deferred call — the function to invoke, its already-evaluated arguments, and a pointer to the next record in the chain:

```go
// Simplified conceptual structure — actual runtime internals vary by Go version
type _defer struct {
    fn   func()   // the function to call
    args []uintptr // captured argument values
    next *_defer   // pointer to the next deferred call in the chain
}
```

Each new `defer` statement prepends its record to the front of this per-goroutine chain (associated with the current function's stack frame). When the function is about to return, the runtime walks this chain from the front, invoking each deferred call in turn — since the most recently added record sits at the front, this naturally produces the LIFO order described in [Section 3](#3-multiple-defers-lifo-order).

## 12. Performance: Open-Coded, Stack-Allocated, and Heap-Allocated Defers

Modern Go actually implements `defer` in **three different ways** internally, chosen automatically by the compiler depending on the shape of the code, with meaningfully different performance:

| Implementation                  | Approximate cost                       | When the compiler uses it                                                                                                  |
| ------------------------------- | -------------------------------------- | -------------------------------------------------------------------------------------------------------------------------- |
| **Open-coded defer** (Go 1.14+) | ~6ns (close to a direct function call) | The common, simple case: a small, fixed number of defers (historically up to 8) in a function, not inside a loop           |
| **Stack-allocated defer**       | ~35ns                                  | Cases the compiler can't open-code but can still avoid a heap allocation for                                               |
| **Heap-allocated defer**        | Higher, plus GC pressure               | `defer` inside a loop, or other cases where the number/shape of defers isn't statically simple enough for the faster paths |

**Open-coded defers**, introduced in Go 1.14, work by having the compiler inline the deferred call directly at each of the function's return points, tracking which defers actually fired with a small bitmap, rather than routing every defer through the general-purpose runtime machinery. According to the original proposal, this optimization applied to the vast majority of real-world defer call sites in the Go toolchain's own source, delivering a substantial (roughly 30%) performance improvement over the pre-1.14 approach for those sites, and reduced typical defer overhead from about 35ns down to about 6ns — close enough to a plain function call that the historical advice to "avoid defer in hot paths" is far less relevant on modern Go than it once was.

**Practical takeaway:** in ordinary application code, `defer`'s overhead is negligible and should not discourage you from using it for correctness and clarity, including in moderately hot code paths. The exception that still matters is `defer` **inside a loop** (as covered in [Section 10](#10-defer-inside-loops-the-classic-pitfall)) — that pattern falls off the fast path onto the slower, heap-allocated implementation for every iteration, on top of the resource-lifetime bug it usually causes, making it doubly worth avoiding in genuinely hot, high-iteration-count loops.

## 13. Defer with Method Values

`defer` works naturally with method values (see the discussion of methods for background) — you can defer a method call exactly like a function call, and the same argument/receiver evaluation timing rules from [Sections 4](#4-arguments-are-evaluated-immediately-at-defer-time) and [5](#5-the-receiver-is-also-evaluated-at-defer-time) apply:

```go
type Resource struct {
    name string
}

func (r *Resource) Release() {
    fmt.Println("releasing:", r.name)
}

func use() {
    r := &Resource{name: "connection-1"}
    defer r.Release() // pointer receiver: captures the pointer now, reads r.name when Release() actually runs later
    // ... use r ...
}
```

A very common real-world instance of this is `defer wg.Done()` with `sync.WaitGroup`, or `defer cancel()` with a `context.CancelFunc` — both are simply method/function values being deferred, following the exact same rules as any other deferred call.

## 14. Defer and Multiple Return Paths

One of `defer`'s biggest practical benefits is collapsing cleanup logic that would otherwise need to be duplicated at every one of a function's exit points into a single statement, placed once, right after the resource is acquired:

```go
// WITHOUT defer: cleanup must be repeated at every return point — easy to forget one
func readConfigWithoutDefer(path string) (Config, error) {
    f, err := os.Open(path)
    if err != nil {
        return Config{}, err
    }

    cfg, err := parseConfig(f)
    if err != nil {
        f.Close() // must remember this here...
        return Config{}, err
    }

    if err := validate(cfg); err != nil {
        f.Close() // ...and here...
        return Config{}, err
    }

    f.Close() // ...and here.
    return cfg, nil
}

// WITH defer: cleanup is written once, and runs correctly no matter which return statement executes
func readConfigWithDefer(path string) (Config, error) {
    f, err := os.Open(path)
    if err != nil {
        return Config{}, err
    }
    defer f.Close() // handles ALL exit paths below, including any future ones added later

    cfg, err := parseConfig(f)
    if err != nil {
        return Config{}, err
    }

    if err := validate(cfg); err != nil {
        return Config{}, err
    }

    return cfg, nil
}
```

This is arguably `defer`'s single biggest practical value: it eliminates an entire class of bugs where a new early-return is added later by someone who forgets to also add the corresponding cleanup call at that new exit point.

## 15. Common Mistakes and Pitfalls

### 15.1 Deferring Inside a Loop (Resource Accumulation)

Covered in depth in [Section 10](#10-defer-inside-loops-the-classic-pitfall) — the single most common `defer`-related bug in real-world Go code, especially in functions that process a list of files, connections, or other limited resources.

### 15.2 Expecting a Deferred Argument to Reflect Later Changes

```go
i := 0
defer fmt.Println(i) // captures 0 now
i = 100
// prints 0, not 100 — surprising if you expected "the current value" at return time
```

Fix: wrap in a closure if you need the value as of execution time, not defer time (see [Section 6](#6-capturing-a-value-at-execution-time-instead)).

### 15.3 Deferring `Close()` Before Checking the Open Error

```go
file, err := os.Open("data.txt")
defer file.Close() // PANICS if err != nil, since file would be nil
if err != nil {
    log.Fatal(err)
}
```

Always check the error from an open/create call **before** deferring the corresponding close — never defer a close on a value that might still be `nil`.

### 15.4 Ignoring a Deferred Call's Returned Error

```go
defer file.Close() // if Close() fails, the error is silently discarded
```

For most read-only file operations this is a widely accepted convenience, but when a `Close()` (or `Flush()`) failure genuinely matters — for example, on writes where data integrity is important — check it explicitly, or capture it into a named return value as shown in [Section 7.1](#71-the-canonical-use-case-modifying-an-error-return-value), rather than relying on a bare, unchecked `defer`.

### 15.5 Overusing Named-Return Mutation for Non-Cleanup Logic

```go
func compute() (result int) {
    defer func() { result = result * 2 }() // works, but obscures what compute() actually returns
    result = 21
    return
}
```

While legal and occasionally useful for genuine cleanup-driven adjustments (like attaching a `Close()` error to an already-set `err`), using this mechanism for ordinary business logic transformations makes the function's actual return value much harder to trace by reading top-to-bottom. Prefer straightforward, explicit `return` expressions for anything that isn't specifically post-return cleanup.

### 15.6 Assuming `recover()` Works Anywhere

```go
func mightPanic() {
    if r := recover(); r != nil { // does nothing useful here — no panic is in flight yet at this point
        fmt.Println("recovered:", r)
    }
    panic("boom")
}
```

`recover()` only has an effect when called directly within a deferred function, during an active panic unwind — calling it in ordinary code before a panic happens (or in a nested, non-deferred call from within the deferred function) does not stop anything.

## 16. Best Practices Summary

1. **Defer resource cleanup immediately after successfully acquiring the resource** (`defer f.Close()` right after a successful `os.Open`, `defer mu.Unlock()` right after `mu.Lock()`) — this is the single most reliable habit for avoiding leaks.
2. **Never defer a cleanup call before checking that the corresponding acquisition succeeded** — closing/unlocking a `nil` or invalid value can panic.
3. **Avoid `defer` inside loops with many iterations**; extract the loop body into its own function or wrap it in an immediately-invoked anonymous function so each iteration gets its own scope for the defer to attach to.
4. **Remember that arguments (and receivers) are evaluated at defer time, not execution time** — wrap the deferred call in a closure if you need it to observe a variable's value as of when it actually runs.
5. **Use named return values plus a deferred closure specifically for the "attach a cleanup error to the function's result" pattern** — and comment it, since it's a subtle behavior for readers unfamiliar with it.
6. **Put `recover()` directly inside a deferred function**, and nowhere else, if you need to stop a panic from propagating.
7. **Don't shy away from `defer` in ordinary or moderately hot code for performance reasons** — modern Go's open-coded defers make its overhead close to a direct function call in the common case.
8. **Check a deferred call's error explicitly (via a named return or an inline check) whenever that failure genuinely matters**, rather than silently discarding it through a bare `defer`.
9. **Use `defer` to collapse cleanup logic across multiple return paths into one place**, reducing the risk that a newly added early return forgets to also clean up.
10. **Keep deferred closures focused on cleanup/recovery**, not general business logic, to keep a function's control flow easy to follow.

## 17. Use Case Summary Table

| Technique                                                 | When to Use                                                          | Example Scenario                                                    |
| --------------------------------------------------------- | -------------------------------------------------------------------- | ------------------------------------------------------------------- |
| `defer f.Close()` right after opening                     | Guarantee a resource is released on every exit path                  | Closing a file, database connection, or network socket              |
| `defer mu.Unlock()` right after `mu.Lock()`               | Guarantee a lock is always released, including on early return/panic | Protecting a critical section in a method                           |
| `defer resp.Body.Close()`                                 | Prevent leaking HTTP connections                                     | After a successful `http.Get`/`http.Client.Do` call                 |
| Deferred closure + named return value                     | Attach a cleanup failure to the function's error result              | Surfacing a `Close()`/`Flush()` error alongside a write's own error |
| Deferred closure calling `recover()`                      | Stop a panic from propagating and convert it into a normal error     | Recovering inside an HTTP handler or goroutine entry point          |
| Wrapping a loop body in its own function before deferring | Avoid accumulating deferred calls across many iterations             | Processing a large list of files, each needing its own `Close()`    |
| `defer` for timing/logging                                | Measure or log how long a function took, regardless of exit path     | `defer func(start time.Time) { ... }(time.Now())`                   |
| `defer tx.Rollback()` after `Begin()`                     | Ensure a database transaction isn't left open on any failure path    | Multi-step transactional database operations                        |

## 18. References

1. Go Team — _Defer, Panic, and Recover_, The Go Blog. https://go.dev/blog/defer-panic-and-recover
2. Go Team — _A Tour of Go: Defer_. https://go.dev/tour/flowcontrol/12
3. Go Language Specification — _Defer statements_. https://go.dev/ref/spec#Defer_statements
4. Dan Scales, Keith Randall, Austin Clements — _Proposal: Low-cost defers through inline code, and extra funcdata to manage the panic case_ (open-coded defers design doc). https://go.googlesource.com/proposal/+show/d74d825331d9b16ee286ea77c0e4caeaf0efbe30/design/34481-opencoded-defers.md
5. Go Team — _Go 1.14 Release Notes_ (open-coded defer performance improvement). https://go.dev/doc/go1.14
6. tpaschalis — _What is a defer? And how many can you run?_. https://tpaschalis.me/defer-internals/
7. VictoriaMetrics Blog — _Golang Defer: From Basic To Traps_. https://victoriametrics.com/blog/defer-in-go/
8. Gabriel Anhaia — _Defer Has 3 Performance Cliffs. Here's How to See Them_, DEV Community. https://dev.to/gabrielanhaia/defer-has-3-performance-cliffs-heres-how-to-see-them-l98
9. Gabriel Anhaia — _defer in Loops: The Resource Leak Go Still Lets You Write_, DEV Community. https://dev.to/gabrielanhaia/defer-in-loops-the-resource-leak-go-still-lets-you-write-j9l
10. Shadeshsaha — _Defer in Golang_, Medium. https://medium.com/@shadeshsaha45/defer-in-golang-d5d9d1d07233
11. devtrovert — _Go Secret - Defer: What you know about DEFER in Go is not enough!_. https://blog.devtrovert.com/p/go-secret-defer-what-you-know-about
12. ZetCode — _Using Defer in Go_ and _Understanding the Defer Statement in Golang_. https://zetcode.com/golang/defer/ and https://www.zetcode.com/golang/defer-keyword/
13. saifulire — _Go's Defer: Simple Rules, Deep Runtime Truths with intuitions_, DEV Community. https://dev.to/saifulire/gos-defer-complex-thing-in-simple-manner-with-low-level-intuitions-43gj
14. OneUptime Engineering Blog — _How to Use defer Correctly in Go_. https://oneuptime.com/blog/post/2026-01-23-go-defer/view
15. ocrosby — _deferdemo package_, go-lab lessons. https://pkg.go.dev/github.com/ocrosby/go-lab/lessons/25-defer-and-cleanup
16. GeeksforGeeks — _Defer Keyword in Golang_. https://www.geeksforgeeks.org/defer-keyword-in-golang
