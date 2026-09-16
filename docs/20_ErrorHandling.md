<div align="center">
  <h1>Error Handling</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 16, 2026</sub>
</div>

## Table of Contents

1. [Philosophy: Errors Are Values](#1-philosophy-errors-are-values)
2. [The `error` Interface](#2-the-error-interface)
3. [Basic Error Handling Pattern](#3-basic-error-handling-pattern)
4. [Creating Errors](#4-creating-errors)
5. [Sentinel Errors](#5-sentinel-errors)
6. [Custom Error Types](#6-custom-error-types)
7. [Error Wrapping (Go 1.13+)](#7-error-wrapping-go-113)
8. [Inspecting Errors: `errors.Is` and `errors.As`](#8-inspecting-errors-errorsis-and-errorsas)
9. [Multiple Errors: `errors.Join` (Go 1.20+)](#9-multiple-errors-errorsjoin-go-120)
10. [Panic and Recover vs. Errors](#10-panic-and-recover-vs-errors)
11. [Error Handling with `defer`](#11-error-handling-with-defer)
12. [Best Practices](#12-best-practices)
13. [Common Mistakes](#13-common-mistakes)
14. [Third-Party Error Libraries (Historical Context)](#14-third-party-error-libraries-historical-context)
15. [Use Case Summary Table](#15-use-case-summary-table)
16. [References](#16-references)

## 1. Philosophy: Errors Are Values

Unlike Java, C++, JavaScript, or Python, Go does **not** use exceptions (`try`/`catch`/`throw`) as the primary mechanism for error handling. Instead, Go treats errors as **ordinary values** that are returned alongside a function's normal result. The caller is responsible for explicitly checking and handling that value.

This design choice is intentional: Go's conventions push developers to check for errors explicitly at the point they occur, rather than throwing an exception that might be caught (or missed) somewhere far up the call stack.

The practical consequence is a lot of code shaped like:

```go
result, err := doSomething()
if err != nil {
    // handle the error
    return err
}
// use result
```

This is more verbose than a `try/catch` block, but it has a key benefit: **the flow of control is always explicit and visible in the code**. You can look at any function and see exactly where it might fail and what happens when it does — nothing is "thrown" silently up an invisible call stack.

**Key takeaway:** Go deliberately trades brevity for explicitness and traceability.

## 2. The `error` Interface

At the core of Go's error handling is a single, tiny built-in interface:

```go
type error interface {
    Error() string
}
```

Any type that implements a method `Error() string` satisfies the `error` interface. This is why errors are so flexible in Go — they can be a simple string-based value, or a fully custom struct carrying structured data (error codes, HTTP status, a wrapped cause, etc.), as long as it has an `Error()` method.

The zero value of an `error` is `nil`, which by convention means "no error occurred."

```go
var err error // err == nil
```

**Important idiom:** Never use the _presence_ of a non-nil error struct pointer as a substitute for checking `err != nil` in unexpected ways — always compare the `error` interface value to `nil` directly (more on the "typed nil" pitfall in [Common Mistakes](#13-common-mistakes)).

## 3. Basic Error Handling Pattern

Because Go supports multiple return values, the idiomatic pattern is to return the result **and** an `error` as the last return value:

```go
package main

import (
    "errors"
    "fmt"
)

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}
```

Conventions to follow:

- The error is always the **last** return value.
- A `nil` error means success; any non-`nil` error means the caller must decide how to react.
- Check errors **immediately** after the call that could produce them — don't defer the check to "later."

**Use case:** This pattern is used everywhere in Go — file I/O, network calls, database queries, JSON parsing, etc. Any operation that can fail returns `(value, error)`.

## 4. Creating Errors

There are two standard-library ways to create a simple error:

### 4.1 `errors.New`

Creates a new error from a static string:

```go
err := errors.New("something went wrong")
```

### 4.2 `fmt.Errorf`

Creates a formatted error message — useful when you want to embed dynamic values:

```go
name := "config.yaml"
err := fmt.Errorf("failed to open file %q", name)
```

`fmt.Errorf` becomes especially powerful when combined with the `%w` verb for **error wrapping** (covered in [Section 7](#7-error-wrapping-go-113)).

**Use case:** Use `errors.New`/`fmt.Errorf` for one-off, local errors that don't need to be identified programmatically elsewhere in your codebase.

## 5. Sentinel Errors

A **sentinel error** is a predefined, exported error _value_ that callers can compare against using `==` (or, more robustly, `errors.Is`). The convention is to name them `ErrXxx` and declare them as package-level variables created with `errors.New`.

```go
package store

import "errors"

var ErrNotFound = errors.New("item not found")

func GetItem(id string) (string, error) {
    // ... lookup logic
    if notFound {
        return "", ErrNotFound
    }
    return item, nil
}
```

Caller side:

```go
item, err := store.GetItem("123")
if errors.Is(err, store.ErrNotFound) {
    // handle "not found" case specifically
}
```

A well-known real-world example from the standard library is `sql.ErrNoRows` in `database/sql`, and `io.EOF` in the `io` package.

**Use case:** Sentinel errors are ideal when the caller needs to make a decision based on **error identity** — "did this specific, known failure happen?" — such as distinguishing "record not found" from "database connection failed."

**Limitation:** Sentinel errors can't carry additional structured data (e.g., which ID was not found). For that, you need a **custom error type**.

## 6. Custom Error Types

When an error needs to carry extra context (fields), define a struct that implements the `error` interface:

```go
package main

import "fmt"

type ValidationError struct {
    Field string
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on field %q: %s", e.Field, e.Msg)
}

func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{Field: "age", Msg: "must not be negative"}
    }
    return nil
}
```

Callers can extract the concrete type using a type assertion or, preferably, `errors.As` (Go 1.13+):

```go
err := validateAge(-5)

var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Println("Field with problem:", ve.Field)
}
```

This pattern — checking a concrete error type to make decisions — has existed since early Go. As the Go blog notes about a similar built-in example, `*QueryError`:

Code can inspect a concrete error type like `*QueryError` to make decisions based on what actually went wrong — this pattern of looking inside a container error is generally called "unwrapping." The standard library's `os.PathError` is another well-known example of an error type that contains another error.

**Use case:** Use custom error types when consumers of your package need structured information about the failure — for example, an HTTP client library that needs to expose the failed status code, or a parser that needs to report a line/column number.

**Naming convention:** Name custom error types `XxxError` (e.g., `ValidationError`, `QueryError`), mirroring the `ErrXxx` convention for sentinel values.

## 7. Error Wrapping (Go 1.13+)

### 7.1 The Problem

Before Go 1.13, adding context to an error while bubbling it up the call stack (e.g., "failed while opening config: <original error>") usually meant using `fmt.Errorf("...: %v", err)`. This **loses the original error's identity** — the caller can no longer detect the underlying cause programmatically, only read the message as a string.

### 7.2 The Solution: `%w`

Go 1.13 introduced the `%w` verb for `fmt.Errorf`, which wraps an error while preserving a link to it:

When `%w` is used, the error returned by `fmt.Errorf` gains an `Unwrap` method that returns the wrapped argument (which must itself be an `error`). In every other respect, `%w` behaves exactly like `%v`.

```go
if err != nil {
    return fmt.Errorf("decompress %v: %w", name, err)
}
```

### 7.3 The `Unwrap` Convention

An error that contains another error can implement an `Unwrap` method that returns that underlying error. If calling `e1.Unwrap()` returns `e2`, we say `e1` wraps `e2`, and that `e1` can be unwrapped to reach `e2`.

```go
type wrappedError struct {
    msg string
    err error
}

func (w *wrappedError) Error() string { return w.msg }
func (w *wrappedError) Unwrap() error { return w.err }
```

You rarely need to write this yourself — `fmt.Errorf("...%w", err)` does it for you — but understanding the mechanism helps when building custom wrapping types.

**Use case:** Wrap an error every time you "bubble it up" through a function boundary and want to add context about _what operation_ was being attempted, while still letting callers detect the _original_ cause with `errors.Is`/`errors.As`.

```go
func loadConfig(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("loadConfig: %w", err)
    }
    // ...
    return nil
}
```

## 8. Inspecting Errors: `errors.Is` and `errors.As`

Go 1.13 also added two functions to the `errors` package specifically to **traverse the wrap chain**:

### 8.1 `errors.Is` — Identity Check

`errors.Is(err, target)` reports whether `err` — or any error it wraps, recursively — matches `target` (via `==` or a custom `Is` method). Use it to check for a **sentinel error**.

```go
var ErrPermission = errors.New("permission denied")

func doSomething() error {
    // ...
    return fmt.Errorf("%w", ErrPermission)
}

if err := doSomething(); errors.Is(err, ErrPermission) {
    // handle the permission error, regardless of how deep it was wrapped
}
```

The Go blog explicitly recommends wrapping sentinel errors so `errors.Is` always works even after future refactors:

The recommended approach is to always return an error that wraps the sentinel, so that callers reliably detect it through `errors.Is`, regardless of how the internal implementation changes over time.

### 8.2 `errors.As` — Type Check

`errors.As(err, &target)` reports whether `err` — or any error it wraps — matches the type of `target`, and if so, assigns the matched error to `target`.

```go
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Println(ve.Field) // access struct fields directly
}
```

### 8.3 `errors.Unwrap`

Returns the single wrapped error one level down (or `nil` if there is none). Both `errors.Is` and `errors.As` use `Unwrap` internally to walk the chain:

`errors.Is` checks whether a specific error is present anywhere in the chain of wrapped errors, unwrapping one level at a time until it finds a match or runs out of errors to check.

### 8.4 Comparison Table

| Function                  | Purpose                                  | Returns                        | Matches by                                    |
| ------------------------- | ---------------------------------------- | ------------------------------ | --------------------------------------------- |
| `errors.Is(err, target)`  | "Is this (or a cause) exactly `target`?" | `bool`                         | Value equality (`==`) or custom `Is()` method |
| `errors.As(err, &target)` | "Is this (or a cause) of type `T`?"      | `bool`, and populates `target` | Concrete type match                           |
| `errors.Unwrap(err)`      | Get the next error down the chain        | `error` (or `nil`)             | N/A                                           |

**Rule of thumb:** Use `errors.Is` for sentinel values; use `errors.As` for custom error types whose fields you need to inspect.

## 9. Multiple Errors: `errors.Join` (Go 1.20+)

### 9.1 Motivation

Sometimes an operation can fail for more than one reason at once — e.g., validating a form with three invalid fields, or running several independent cleanup steps where more than one might fail. Prior to Go 1.20 there was no standard way to combine several errors into one `error` value while still preserving each one's identity for `errors.Is`/`errors.As`. Developers often reached for third-party packages such as `hashicorp/go-multierror` or `uber-go/multierr`.

### 9.2 `errors.Join`

Go 1.20 added `errors.Join`, plus support for **multiple `%w` verbs** in `fmt.Errorf`, and updated `errors.Is`/`errors.As` to traverse a **tree** of wrapped errors (not just a linear chain):

Go 1.20 expanded error wrapping so that a single error can wrap multiple other errors at once. An error achieves this by implementing an `Unwrap` method that returns a `[]error` instead of a single `error`. Both `errors.Is` and `errors.As` were updated to traverse this multi-error structure, `fmt.Errorf` now accepts multiple `%w` verbs in one call, and the new `errors.Join` function returns an error that wraps an entire list of errors.

Basic usage:

```go
err1 := errors.New("err1")
err2 := errors.New("err2")
err := errors.Join(err1, err2)

fmt.Println(err) // prints "err1\nerr2"

errors.Is(err, err1) // true
errors.Is(err, err2) // true
```

### 9.3 Real Use Case: Batch Processing

```go
func processAll(items []Item) error {
    var err error
    for _, item := range items {
        if e := process(item); e != nil {
            err = errors.Join(err, e)
        }
    }
    return err
}
```

`errors.Join` gracefully ignores `nil` arguments, so this accumulator pattern (`err = errors.Join(err, e)`) works correctly even on the first iteration when `err` is still `nil`.

### 9.4 Caveat: `errors.Unwrap` Returns `nil` for Joined Errors

Because `errors.Join` implements `Unwrap() []error` (plural) instead of the classic `Unwrap() error` (singular), calling the older, single-error `errors.Unwrap()` function on a joined error returns `nil` — this is intentional, to preserve backward compatibility:

For error types implementing the plural `Unwrap() []error` form, the older single-error `errors.Unwrap()` function always returns `nil` — a backward-compatibility trade-off worth keeping in mind when introducing multi-error wrapping into existing code.

In practice this rarely matters because idiomatic code checks errors with `errors.Is`/`errors.As`, both of which were updated to correctly traverse this tree structure.

**Use case:** Use `errors.Join` when you must report _all_ failures from an operation at once (form validation, running multiple independent cleanup/close steps, batch jobs) rather than stopping at the first error.

## 10. Panic and Recover vs. Errors

Go has a separate mechanism, `panic`/`recover`, that superficially resembles exceptions — but it is **not** meant to be used for ordinary error handling.

- **`error`** — for expected, recoverable failure conditions that are part of normal program flow (file not found, invalid input, network timeout).
- **`panic`** — for _unrecoverable_ programmer errors or conditions that should stop execution immediately (e.g., an out-of-bounds slice access, a nil pointer dereference, or a deliberate `panic()` call for "this should never happen" invariant violations).
- **`recover`** — called inside a `deferred` function to regain control after a panic, typically used at a boundary (like an HTTP handler or goroutine entry point) to prevent a single failure from crashing the entire process.

```go
func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
        }
    }()
    result = a / b // panics on division by zero
    return
}
```

**Use case for `recover`:** Guarding the top level of a goroutine or an HTTP server middleware so that an unexpected panic in one request doesn't take down the whole server, converting it into a logged 500 error instead. It should **not** be used as a substitute for `if err != nil` checks in everyday logic.

## 11. Error Handling with `defer`

`defer` is commonly combined with error handling to guarantee cleanup (closing files, releasing locks, closing DB connections) regardless of how a function exits.

```go
func readConfig(path string) ([]byte, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("open config: %w", err)
    }
    defer f.Close()

    data, err := io.ReadAll(f)
    if err != nil {
        return nil, fmt.Errorf("read config: %w", err)
    }
    return data, nil
}
```

A more advanced pattern captures the error from a deferred cleanup call (e.g., `f.Close()` can itself fail) using a **named return value**:

```go
func writeAll(path string, data []byte) (err error) {
    f, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("create file: %w", err)
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = fmt.Errorf("close file: %w", cerr)
        }
    }()

    _, err = f.Write(data)
    return err
}
```

**Use case:** Any function that acquires a resource (file handle, DB connection, mutex, HTTP response body) should `defer` its release right after a successful acquisition, and — when the release itself can fail — surface that failure without silently swallowing it.

## 12. Best Practices

1. **Check errors immediately.** Don't let a function continue with stale/zero data after an unhandled error.
2. **Add context when propagating (bubbling up), but avoid double-handling.** A common rule of thumb is to wrap an error with at least the current function's name (or the operation being attempted) every time it is passed back up the call chain, so the final message reads like a trail of breadcrumbs rather than a bare, context-free message.
3. **Handle an error exactly once.** Either you `log` it and stop propagating it, or you `return` it (possibly wrapped) — never both. Logging _and_ returning the same error typically causes duplicate log entries further up the stack.
4. **Prefer `errors.Is`/`errors.As` over direct type assertions or `==`** when the error may have been wrapped anywhere along the call chain, since these functions traverse the whole chain.
5. **Use sentinel errors for identity, custom types for data.** If the caller only needs a yes/no ("was it a 'not found' error?"), a sentinel is enough. If the caller needs structured details (field name, status code, retry-after duration), use a custom type.
6. **Be careful about wrapping in public/library APIs.** Wrapping always preserves the original message, which can leak internal details. Because wrapping keeps the original error message intact, doing so indiscriminately can expose internal implementation details that raise security, privacy, or UX concerns. In such cases it can be better to handle the error internally and return a fresh, sanitized error rather than wrapping the original — this is especially relevant for REST APIs and open-source libraries that shouldn't leak internal details to third-party callers.
7. **Never expose raw internal errors to end users** in security-sensitive contexts (e.g., don't return raw SQL error text to an HTTP client) — log the detailed error internally and return a generic, safe message externally.
8. **Keep error messages lowercase and without trailing punctuation**, since they're often wrapped into larger sentences (Go convention: `"failed to open file"`, not `"Failed to open file."`).
9. **Name sentinel errors `ErrXxx` and custom error types `XxxError`** for consistency with the standard library.
10. **Use `errors.Join` for aggregating independent failures**, not for adding sequential context — that's what `%w` is for.

## 13. Common Mistakes

### 13.1 Comparing Wrapped Errors with `==`

```go
// WRONG — breaks once the error is wrapped anywhere in the chain
if err == ErrNotFound { ... }

// RIGHT — traverses the wrap chain
if errors.Is(err, ErrNotFound) { ... }
```

### 13.2 The "Typed Nil" Trap

A classic Go gotcha: returning a `nil` pointer of a concrete error type through an `error` interface produces a **non-nil interface value**.

```go
type MyError struct{}
func (e *MyError) Error() string { return "my error" }

func doWork() error {
    var e *MyError = nil
    return e // returns a non-nil `error` interface, even though the underlying pointer is nil!
}

err := doWork()
fmt.Println(err == nil) // false! Surprising.
```

**Fix:** Only return a literal `nil`, not a typed nil variable, when there is no error:

```go
func doWork() error {
    var e *MyError
    if somethingFailed {
        e = &MyError{}
    }
    if e != nil {
        return e
    }
    return nil
}
```

### 13.3 Ignoring Errors with `_`

```go
data, _ := os.ReadFile("config.yaml") // silently discards failure information
```

Only discard an error deliberately and when you can justify it (e.g., a `Close()` call where failure is genuinely irrelevant) — and prefer to still log it in that case.

### 13.4 Wrapping the Same Error Multiple Times Unnecessarily

Excessive `fmt.Errorf("...: %w", err)` at every single call frame can make error messages extremely long and redundant. Add context only where it's genuinely useful (i.e., where you know something the caller of _your_ function doesn't).

### 13.5 Using `panic` for Ordinary/Expected Failures

Reserve `panic` for truly exceptional, programmer-error conditions. Using it for expected failures (like "user not found") breaks Go idioms and forces every caller to use `recover`, defeating the purpose of the `error` interface.

## 14. Third-Party Error Libraries (Historical Context)

Before Go 1.13 standardized wrapping, the community widely used **[`github.com/pkg/errors`](https://github.com/pkg/errors)** by Dave Cheney, which provided `errors.Wrap`, `errors.Cause`, and automatic stack-trace capture. Many existing/legacy Go codebases still use it.

Before Go 1.13 shipped in 2019, the standard library offered very little dedicated tooling for errors — essentially just `errors.New` and `fmt.Errorf` — which is why many older codebases still rely on third-party libraries such as `pkg/errors`.

For **aggregating multiple errors**, before `errors.Join` existed (Go < 1.20), common choices were `hashicorp/go-multierror` and `uber-go/multierr`. As of Go 1.20, the standard library's `errors.Join` offers similar functionality for combining multiple errors, and the `go-multierror` maintainers themselves now recommend it for new projects — reserving `go-multierror` for cases needing its extra features, like custom formatting, the concurrent-safe `Group` pattern, or utilities such as `Append`, `Flatten`, and `Prefix`.

**Recommendation for new projects (Go 1.20+):** Prefer the standard `errors` package (`errors.New`, `fmt.Errorf` with `%w`, `errors.Is`, `errors.As`, `errors.Join`) unless you specifically need a feature the standard library lacks, such as automatic stack traces (`pkg/errors`) or the concurrent error-collection `Group` pattern (`go-multierror`).

## 15. Use Case Summary Table

| Technique                             | When to Use                                                                         | Example Scenario                                                                   |
| ------------------------------------- | ----------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `errors.New` / `fmt.Errorf` (no `%w`) | Local, one-off error; caller doesn't need to inspect it programmatically            | `errors.New("invalid input")`                                                      |
| Sentinel error (`ErrXxx`)             | Caller needs to check _identity_ of a known failure case                            | `sql.ErrNoRows`, `io.EOF`, `ErrNotFound`                                           |
| Custom error type (`XxxError`)        | Caller needs _structured data_ about the failure                                    | HTTP client exposing status code; validation error with field name                 |
| `fmt.Errorf("...: %w", err)`          | Add context while bubbling an error up the call stack, preserving inspectability    | `"loadConfig: %w"` wrapping an `os.Open` failure                                   |
| `errors.Is`                           | Check whether a specific sentinel error occurred anywhere in the chain              | Retry logic on a network timeout sentinel                                          |
| `errors.As`                           | Extract a specific error type's fields anywhere in the chain                        | Reading `ValidationError.Field` after several wraps                                |
| `errors.Join`                         | Combine multiple independent failures into one value                                | Form validation with several invalid fields; batch job errors                      |
| `panic` / `recover`                   | Truly unrecoverable programmer errors; top-level safety net (e.g., HTTP middleware) | Recovering from a panic in a goroutine so one bad request doesn't crash the server |
| `defer` + error capture               | Guarantee cleanup and optionally surface cleanup failures                           | Closing a file/DB connection, capturing `Close()` errors via named return          |

## 16. References

1. Go Team — _Error handling and Go_, The Go Blog. https://go.dev/blog/error-handling-and-go
2. Go Team — _Working with Errors in Go 1.13_, The Go Blog. https://go.dev/blog/go1.13-errors
3. Go Team — _syntactic support for error handling_, The Go Blog. https://go.dev/blog/error-syntax
4. Go 1.20 Release Notes — _Wrapping multiple errors_. https://go.dev/doc/go1.20 (summarized via secondary sources below)
5. Endorama — _Golang 1.20 - multiple errors_. https://endorama.dev/2023/golang-1.20-multiple-errors/
6. Lukáš Zapletal — _New in Go 1.20: wrapping multiple errors_. https://lukas.zapletalovi.com/posts/2022/wrapping-multiple-errors/
7. Earthly Blog — _Effective Error Handling in Golang_. https://earthly.dev/blog/golang-errors/
8. JetBrains GoLand Blog — _How to Handle Errors in Go_. https://blog.jetbrains.com/go/2026/09/02/how-to-handle-errors-in-go/
9. JetBrains Blog — _Best Practices for Secure Error Handling in Go_. https://blog.jetbrains.com/go/2026/03/02/secure-go-error-handling-best-practices/
10. Peter Bourgon — _Programming with errors_. https://peter.bourgon.org/blog/2019/09/11/programming-with-errors.html
11. Ganesh — _Untangling Go Errors: Wrap, Is, As, and When to Use Them_, Medium. https://medium.com/@ganeshpant/untangling-go-errors-wrap-is-as-and-when-to-use-them-e7674362e139
12. HashiCorp — _go-multierror_ (README, migration notes to `errors.Join`). https://github.com/hashicorp/go-multierror
13. Datadog Engineering Blog — _A practical guide to error handling in Go_. https://www.datadoghq.com/blog/go-error-handling/
14. `pkg/errors` by Dave Cheney. https://github.com/pkg/errors
15. Go standard library documentation — `errors` package. https://pkg.go.dev/errors
16. Go standard library documentation — `fmt` package (`%w` verb). https://pkg.go.dev/fmt
