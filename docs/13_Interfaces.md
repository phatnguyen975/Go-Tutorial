<div align="center">
  <h1>Interfaces</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 17, 2026</sub>
</div>

## Table of Contents

1. [What Is an Interface?](#1-what-is-an-interface)
2. [Implicit Satisfaction: Go's Structural Typing](#2-implicit-satisfaction-gos-structural-typing)
3. [Declaring and Using Interfaces](#3-declaring-and-using-interfaces)
4. [Interface Internals: Type + Value](#4-interface-internals-type--value)
5. [Method Sets: Value vs. Pointer Receivers](#5-method-sets-value-vs-pointer-receivers)
6. [The Empty Interface (`interface{}` / `any`)](#6-the-empty-interface-interface--any)
7. [Type Assertions](#7-type-assertions)
8. [Type Switches](#8-type-switches)
9. [Interface Embedding and Composition](#9-interface-embedding-and-composition)
10. [The Nil Interface Pitfall (Typed Nil)](#10-the-nil-interface-pitfall-typed-nil)
11. [Comparing Interface Values](#11-comparing-interface-values)
12. [Compile-Time Interface Satisfaction Checks](#12-compile-time-interface-satisfaction-checks)
13. [Well-Known Standard Library Interfaces](#13-well-known-standard-library-interfaces)
14. [Design Philosophy: Accept Interfaces, Return Structs](#14-design-philosophy-accept-interfaces-return-structs)
15. [Interfaces and Testing (Mocking)](#15-interfaces-and-testing-mocking)
16. [Interfaces vs. Generics](#16-interfaces-vs-generics)
17. [Performance Considerations](#17-performance-considerations)
18. [Common Mistakes and Pitfalls](#18-common-mistakes-and-pitfalls)
19. [Best Practices Summary](#19-best-practices-summary)
20. [Use Case Summary Table](#20-use-case-summary-table)
21. [References](#21-references)

## 1. What Is an Interface?

An interface in Go is a type that defines a **set of method signatures** — a contract describing behavior, without specifying how that behavior is implemented or what concrete data backs it. An interface itself carries no fields and no method bodies; it purely says "anything that can do these things satisfies me."

```go
type Namer interface {
    Name() string
}
```

Any concrete type — a struct, a named primitive type, a map, a function type, even another interface — that implements all the methods declared by an interface **automatically** satisfies it. Go is not a classic class-based, inheritance-driven object-oriented language; interfaces are how it achieves polymorphism and abstraction instead.

## 2. Implicit Satisfaction: Go's Structural Typing

The single most distinctive feature of Go's interfaces, compared to languages like Java or C#, is that **satisfaction is implicit**. A type never declares "I implement interface X" — the compiler simply checks whether the type's method set includes every method the interface requires, and if so, treats it as satisfying that interface.

```go
type Greeter interface {
    Greet() string
}

type EnglishGreeter struct{}
func (e EnglishGreeter) Greet() string { return "Hello!" }

type SpanishGreeter struct{}
func (s SpanishGreeter) Greet() string { return "¡Hola!" }

func main() {
    var g Greeter
    g = EnglishGreeter{}
    fmt.Println(g.Greet()) // Hello!
    g = SpanishGreeter{}
    fmt.Println(g.Greet()) // ¡Hola!
}
```

Notice that neither `EnglishGreeter` nor `SpanishGreeter` mentions `Greeter` anywhere in its own definition. This is sometimes loosely compared to "duck typing" ("if it walks like a duck and quacks like a duck…"), though it's more precisely called **structural typing**, because it is fully checked at compile time — unlike true duck typing in dynamically typed languages, an unsatisfied interface is a compile error, not a runtime surprise.

### 2.1 Why This Matters

- **Decoupling:** the package defining an interface and the package providing an implementation need not know about each other at all. You can define your own interface in your own package, and any third-party type — even one from a library whose author never heard of your interface — satisfies it automatically as long as its methods line up.
- **Retroactive interfaces:** you can introduce a new interface after the fact to describe behavior that existing types already happen to have, without touching those types.
- **Small, focused interfaces are natural:** because there's no cost to "declaring" satisfaction, Go culture strongly favors small interfaces (often just one or two methods) that describe a single capability, rather than large, all-encompassing ones.
- **Encourages "accept interfaces, return structs"** — a design idiom covered in [Section 14](#14-design-philosophy-accept-interfaces-return-structs).

### 2.2 Trade-offs

Implicit satisfaction isn't free of downsides:

- **No explicit statement of intent.** The compiler enforces that the methods exist, but not that a type was _designed_ for a particular interface — a type could accidentally satisfy an interface with methods that don't behave as that interface's contract expects. In practice this is rare with well-named, small interfaces.
- **Discoverability.** For a large interface, it isn't always obvious at a glance which concrete types satisfy it — you may need an IDE or `grep` to find implementations, since there's no explicit "implements" list anywhere.

## 3. Declaring and Using Interfaces

An interface's zero value is `nil`. A `nil` interface has no methods that can be safely called (calling one panics) until it's assigned a concrete value:

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

func describe(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    shapes := []Shape{
        Rectangle{Width: 3, Height: 4},
        Circle{Radius: 5},
    }
    for _, s := range shapes {
        describe(s) // works uniformly for any Shape, regardless of concrete type
    }
}
```

This is Go's version of polymorphism: `describe` doesn't know or care whether it received a `Rectangle` or a `Circle` — it only relies on the `Shape` contract.

## 4. Interface Internals: Type + Value

Understanding how an interface value is represented internally clarifies many of Go's more surprising behaviors (especially around `nil`). Conceptually, a non-empty interface value is a pair:

```
(type descriptor, value)
```

- The **type word** records the concrete type currently stored in the interface (e.g., `*MyError`, `Rectangle`).
- The **value word** holds (or points to) the actual data.

An interface value is `nil` **only when both parts are unset** — i.e., no concrete type has ever been assigned. This two-part representation is exactly why a "typed nil" (a nil pointer of some concrete type, stored in an interface) is _not_ the same as a nil interface — see [Section 10](#10-the-nil-interface-pitfall-typed-nil) for the full explanation and its practical consequences.

## 5. Method Sets: Value vs. Pointer Receivers

Whether a type satisfies an interface depends on its **method set**, which in turn depends on whether its methods use value receivers or pointer receivers.

| Receiver type                        | Method set of `T`        | Method set of `*T` |
| ------------------------------------ | ------------------------ | ------------------ |
| `func (t T) M()` (value receiver)    | Includes `M`             | Includes `M`       |
| `func (t *T) M()` (pointer receiver) | Does **not** include `M` | Includes `M`       |

In other words:

- A value of type `T` only has access to methods declared with a **value** receiver.
- A pointer `*T` has access to **both** value-receiver and pointer-receiver methods.

```go
type Animal interface {
    Speak() string
}

type Dog struct{}
func (d *Dog) Speak() string { return "Woof!" } // pointer receiver

func main() {
    var d Dog
    // var a Animal = d  // COMPILE ERROR: Dog does not implement Animal (Speak has pointer receiver)
    var a Animal = &d    // OK: *Dog implements Animal
    fmt.Println(a.Speak())
}
```

**Why this rule exists:** if `M` has a pointer receiver, it may mutate the value it's called on, or may rely on the value having a stable address; the compiler cannot always safely take the address of an arbitrary value (e.g., a value stored inside a map, or an unaddressable literal), so Go does not automatically promote a value's method set to include pointer-receiver methods.

**Practical guidance:**

- If **any** method of a type needs a pointer receiver (to mutate the receiver, or to avoid copying a large struct), it's generally best to make **all** of that type's methods use pointer receivers, for consistency — this also avoids subtly satisfying (or failing to satisfy) an interface with only some of its methods "promoted."
- When passing something into a function expecting an interface, remember to pass a `&T` if `T`'s relevant methods use pointer receivers.

## 6. The Empty Interface (`interface{}` / `any`)

An interface with **zero methods**, `interface{}`, is satisfied by every type, since there are no requirements to meet at all. Since Go 1.18, the standard library defines `any` as an alias for `interface{}`, and it is now the idiomatic spelling:

```go
func describe(i any) {
    fmt.Printf("value: %v, type: %T\n", i, i)
}

describe(42)
describe("hello")
describe(true)
describe(Rectangle{Width: 2, Height: 3})
```

**Historical context:** before Go 1.18 introduced generics, `interface{}` was the closest thing to "generic" code — containers like a homemade stack or a JSON decoder used `interface{}` to accept or return arbitrary types, at the cost of losing compile-time type safety (you had to type-assert values back out) and some runtime performance overhead (boxing values into an interface, and needing a type assertion or reflection to use them concretely).

**Modern guidance:** for genuinely generic algorithms and data structures (a type-safe stack, a `Map`/`Filter`/`Reduce` helper, etc.), prefer Go's built-in **generics** (type parameters, Go 1.18+) over `any`, since generics give you compile-time type checking and avoid the runtime cost of boxing/unboxing. Reserve `any`/`interface{}` for genuinely dynamic cases: heterogeneous collections, JSON/config decoding into unknown shapes, `fmt`-style formatting functions, and similar situations where the type truly isn't known ahead of time. See [Section 16](#16-interfaces-vs-generics) for a fuller comparison.

## 7. Type Assertions

A **type assertion** extracts the concrete value stored inside an interface, when you already have reason to believe (or want to check) what concrete type it holds:

```go
var i any = "hello"

s := i.(string)         // "unsafe" form: panics if i does not hold a string
fmt.Println(s)

s, ok := i.(string)     // "comma-ok" form: ok is false instead of panicking
if ok {
    fmt.Println("it's a string:", s)
}

n, ok := i.(int) // ok == false, n == 0 (zero value), no panic
```

**Rule of thumb:** always use the two-result "comma-ok" form unless you are certain (and willing to crash if wrong) that the interface holds the asserted type — for example, immediately after a code path that guarantees it. Using the single-result form in general-purpose code is a common source of unexpected panics in production.

You can also assert to another **interface** type, not just a concrete type, which is useful for checking optional capabilities:

```go
type Closer interface {
    Close() error
}

func maybeClose(w io.Writer) {
    if c, ok := w.(Closer); ok {
        c.Close()
    }
}
```

This pattern — checking whether a value _additionally_ satisfies some other, more specific interface — is common in the standard library (for example, `http.ResponseWriter` implementations are sometimes checked for an optional `http.Flusher` or `http.Hijacker` interface).

## 8. Type Switches

A **type switch** lets you branch on the concrete type stored in an interface value across several cases, more cleanly than a chain of type assertions:

```go
func describe(i any) {
    switch v := i.(type) {
    case int:
        fmt.Printf("int: %d\n", v)
    case string:
        fmt.Printf("string: %q\n", v)
    case bool:
        fmt.Printf("bool: %t\n", v)
    case []int:
        fmt.Printf("slice of int, length %d\n", len(v))
    case nil:
        fmt.Println("nil value")
    default:
        fmt.Printf("unknown type: %T\n", v)
    }
}
```

Inside each `case`, the variable `v` automatically takes on the type of that case (e.g., `v` is an `int` inside `case int:`), which the compiler enforces — no manual assertion needed within the branch.

**Use case:** handling several possible concrete types that flow through a generic pipeline — for example, encoding arbitrary values to JSON, dispatching different AST node types in a parser, or routing errors by their concrete type when `errors.As` isn't a fit for multi-branch logic.

## 9. Interface Embedding and Composition

Interfaces can be composed from other interfaces by embedding them, which is how the standard library builds larger contracts out of small, focused ones:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
    Reader // embeds Reader's method set
    Writer // embeds Writer's method set
}
```

A type satisfies `ReadWriter` only if it implements **both** `Read` and `Write`. This composition pattern is central to Go's standard library design (`io.Reader`, `io.Writer`, `io.Closer`, and their combinations `io.ReadWriter`, `io.ReadCloser`, `io.ReadWriteCloser`, etc.), and it reflects the broader Go philosophy of building larger behavior out of small, single-purpose pieces rather than designing large monolithic interfaces up front.

## 10. The Nil Interface Pitfall (Typed Nil)

This is one of the most well-known "gotchas" in Go, and understanding it thoroughly is essential.

### 10.1 The Problem

```go
type MyError struct{}
func (e *MyError) Error() string { return "my error" }

func doWork() error {
    var e *MyError // e is nil, but of type *MyError
    if false {
        e = &MyError{}
    }
    return e // returns a NON-nil error interface!
}

func main() {
    err := doWork()
    fmt.Println(err == nil) // false — surprising!
}
```

### 10.2 Why This Happens

Recall from [Section 4](#4-interface-internals-type--value) that a non-empty interface value is internally a `(type, value)` pair, and it only equals `nil` when **both** the type and value are unset. When `doWork` returns the nil `*MyError` pointer `e` as an `error`, Go must convert the concrete `*MyError` into the `error` interface. That conversion sets the type word to `*MyError` and the value word to `nil` — the resulting interface value is `(type=*MyError, value=nil)`, which is **not** the zero pair, and therefore `err == nil` is `false`, even though the underlying pointer genuinely is `nil`.

Confusingly, `fmt.Println(err)` in this scenario often prints `<nil>`, because it calls `Error()` on the underlying pointer and the method happens to handle a nil receiver gracefully (or `fmt` catches the panic) — but that's the _pointer's_ nil-ness being displayed, not a statement about the interface itself.

### 10.3 The Fix

Never return a typed nil variable through an interface-typed return value when you mean "no error" / "nothing here." Return the literal `nil` explicitly:

```go
func doWork() error {
    var e *MyError
    if somethingFailed {
        e = &MyError{}
        return e // fine: e is genuinely non-nil here
    }
    return nil // return literal nil, not the typed variable
}
```

### 10.4 Where This Bites in Practice

This bug most often shows up in:

- Functions that declare a typed error/interface variable early, conditionally populate it, and return it directly at the end instead of branching to return a literal `nil`.
- Repository/service layers that return a domain interface, where a "happy path" constructs a typed pointer that ends up nil.
- Assigning any nil-capable concrete value (a nil pointer, nil slice, nil map, nil channel, nil function) into an `any`/`interface{}` variable — the interface is non-nil once a concrete type is attached, even though the underlying value looks empty:

```go
var b *int = nil
var a any = b
fmt.Println(b == nil) // true — b itself is a plain nil pointer
fmt.Println(a == nil) // false — a is a non-nil interface wrapping a nil *int
```

### 10.5 General Defensive Guidance

- Be explicit: return literal `nil`, not a nil-valued typed variable, whenever a function's contract is "return nil to mean no error / nothing here."
- When accepting an interface parameter and needing to check "is there really nothing here," be aware that a plain `== nil` check may not be enough if the caller could have handed you a typed nil; if this matters, you may need reflection (`reflect.ValueOf(x).IsNil()`) — though needing this is usually itself a sign the API should be redesigned to avoid the ambiguity in the first place.
- Prefer initializing variables to clear zero values or via constructors, and avoid mixing "nil pointer" and "nil interface" as if they always meant the same thing.

## 11. Comparing Interface Values

Two interface values are equal (`==`) if they hold identical dynamic types **and** their dynamic values are equal under `==`. If the dynamic type is not comparable (e.g., a slice, map, or function stored in an interface), comparing two interface values that hold that type **panics at runtime**, even though the comparison compiles (since the static type is just `any`/some interface):

```go
var a any = []int{1, 2, 3}
var b any = []int{1, 2, 3}
fmt.Println(a == b) // panic: comparing uncomparable type []int
```

This is a subtle risk whenever interface values (especially `any`) are used as map keys or compared directly — make sure the concrete types that will flow through are actually comparable, or guard the comparison (e.g., with `reflect.DeepEqual` for structural comparison, understanding that it has different semantics and performance characteristics than `==`).

## 12. Compile-Time Interface Satisfaction Checks

Because satisfaction is implicit, it's easy to accidentally break it — for example, by renaming a method or changing its signature — without any error until something tries to use the type as that interface, potentially far away in the codebase. A common idiom forces the compiler to check satisfaction immediately, at the point the type is defined:

```go
type FileWriter struct{ /* ... */ }
func (f *FileWriter) Write(p []byte) (int, error) { /* ... */ return 0, nil }

var _ io.Writer = (*FileWriter)(nil) // compile-time check: does *FileWriter implement io.Writer?
```

`var _ io.Writer = (*FileWriter)(nil)` creates a nil `*FileWriter` pointer and attempts to assign it to a variable of the blank identifier `_` typed as `io.Writer`. The assignment itself is discarded at runtime, but the **compiler** must verify that `*FileWriter` satisfies `io.Writer` for the assignment to type-check at all — if it doesn't, you get an immediate, clear compile error right next to the type definition, instead of a confusing error somewhere else later.

**Use case:** place this line right after defining any type that's meant to implement a specific interface, especially for types implementing interfaces from external packages (like `io.Writer`, `sort.Interface`, or your own domain interfaces), so refactors that accidentally break the contract are caught immediately.

## 13. Well-Known Standard Library Interfaces

Several small, widely used interfaces form the backbone of idiomatic Go APIs:

| Interface         | Method(s)                                            | Purpose                                                                                                  |
| ----------------- | ---------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| `error`           | `Error() string`                                     | The standard error contract used throughout Go for reporting failures.                                   |
| `fmt.Stringer`    | `String() string`                                    | Lets `fmt` (and anything using `%v`/`%s`) print a custom, human-readable representation of a type.       |
| `io.Reader`       | `Read(p []byte) (n int, err error)`                  | Anything that can be read from as a stream of bytes.                                                     |
| `io.Writer`       | `Write(p []byte) (n int, err error)`                 | Anything that can be written to as a stream of bytes.                                                    |
| `io.Closer`       | `Close() error`                                      | Anything that holds a resource that must be released.                                                    |
| `sort.Interface`  | `Len() int`, `Less(i, j int) bool`, `Swap(i, j int)` | Lets `sort.Sort` order any custom collection.                                                            |
| `context.Context` | `Done()`, `Err()`, `Deadline()`, `Value()`           | Carries deadlines, cancellation signals, and request-scoped values across API boundaries and goroutines. |

Example: implementing `fmt.Stringer` for custom formatting.

```go
type Point struct{ X, Y int }

func (p Point) String() string {
    return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

fmt.Println(Point{3, 4}) // prints "(3, 4)" instead of "{3 4}"
```

Studying these interfaces is a good way to internalize idiomatic Go interface design: each one is small (often a single method), focused on one capability, and composable with others.

## 14. Design Philosophy: Accept Interfaces, Return Structs

A widely repeated Go idiom, often summarized as **"accept interfaces, return concrete types"** (sometimes phrased "accept interfaces, return structs"), guides how to design function and constructor signatures:

- **Function parameters** should be interfaces when possible — this keeps the function flexible about what it can work with, and decoupled from any specific implementation. A function that takes an `io.Reader` works with a file, a network connection, an in-memory buffer, or a test fixture, without modification.
- **Return values** should generally be concrete types (structs, or pointers to structs) — this gives callers full access to everything the concrete type offers, doesn't hide capabilities behind an artificially narrow interface, and avoids forcing every caller to go through the interface's necessarily generic contract.

```go
// Good: accepts the general io.Writer interface
func WriteHeader(w io.Writer, name string) error {
    _, err := fmt.Fprintf(w, "Name: %s\n", name)
    return err
}

// Good: returns a concrete type, giving the caller full access to *Client's capabilities
func NewClient(baseURL string) *Client {
    return &Client{baseURL: baseURL}
}
```

**Why this matters for third-party code specifically:** because Go's interfaces are satisfied implicitly, you can define your own narrow interface in your own package, and any external type — even from a library whose author never anticipated your interface — will satisfy it automatically as long as its method signatures line up. This lets you decouple your code from a specific dependency without wrapping it in adapter boilerplate, which is a much heavier requirement in languages with explicit `implements` declarations.

**A related warning: avoid interface pollution.** Don't define an interface (especially in the producer/library package) just because "it might be useful later" or purely to enable mocking. Wait until there is a genuine need for abstraction (multiple real implementations, or a real requirement to swap implementations in tests) before introducing an interface — small interfaces are best created at the point of _use_ (the consumer), not preemptively by the producer.

## 15. Interfaces and Testing (Mocking)

Because Go interfaces are implicitly satisfied and typically small, they are a natural seam for substituting real dependencies with test doubles:

```go
type EmailSender interface {
    Send(to, subject, body string) error
}

// Production implementation
type SMTPSender struct{ /* ... */ }
func (s *SMTPSender) Send(to, subject, body string) error { /* real SMTP call */ return nil }

// Test double — no mocking library needed
type FakeSender struct {
    sent []string
}
func (f *FakeSender) Send(to, subject, body string) error {
    f.sent = append(f.sent, to)
    return nil
}

func TestNotifyUser(t *testing.T) {
    fake := &FakeSender{}
    NotifyUser(fake, "user@example.com") // pass the fake in place of a real SMTPSender
    if len(fake.sent) != 1 {
        t.Errorf("expected 1 email sent, got %d", len(fake.sent))
    }
}
```

Because satisfying an interface requires no special declaration, a hand-written "fake" struct like `FakeSender` above is often all you need — no mocking framework required, though popular tools like `gomock` and `testify/mock` exist for generating mocks automatically from larger interfaces, which can save time when an interface has many methods or you need fine-grained call/argument assertions.

**Use case:** define a small interface for any external dependency your code needs to call (a database, an email service, a payment gateway, the system clock) at the point where your code _consumes_ it, then substitute a lightweight fake or generated mock in tests instead of hitting the real dependency.

## 16. Interfaces vs. Generics

Since Go 1.18, **type parameters** (generics) provide an alternative way to write code that works across multiple types, and it's important to know which tool fits which problem:

|               | Interfaces                                                                                    | Generics                                                                                    |
| ------------- | --------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| Purpose       | Abstract over **behavior** (what a type can _do_)                                             | Abstract over **the type itself** (write one algorithm for many types)                      |
| Type checking | Dynamic dispatch at runtime; concrete type resolved via the interface's method set            | Resolved at compile time; the compiler generates/type-checks code per instantiation         |
| Typical use   | Swappable implementations, dependency injection, polymorphic collections of different types   | Type-safe containers and algorithms (a generic `Stack[T]`, `Map`/`Filter`/`Reduce` helpers) |
| Runtime cost  | Slight overhead from the interface's dynamic dispatch and possible heap allocation ("boxing") | No dynamic dispatch overhead; behaves like hand-written code for each type                  |
| Example       | `io.Writer`, `sort.Interface`, `error`                                                        | `func Max[T cmp.Ordered](a, b T) T`                                                         |

Note that generics and interfaces aren't mutually exclusive — Go's generics use interfaces as **type constraints** to describe what operations a type parameter must support:

```go
type Number interface {
    ~int | ~int64 | ~float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}
```

**Guidance:** reach for a plain interface when you need runtime polymorphism (different concrete types handled uniformly through one variable or slice, often with genuinely different behavior per type). Reach for generics when you're writing an algorithm or data structure whose _logic_ is identical for every type it supports, and you want compile-time type safety without runtime boxing.

## 17. Performance Considerations

Storing a value in an interface is not entirely free:

- **Boxing:** when a concrete value doesn't already live on the heap, storing it in an interface can force a heap allocation so the interface's value word can hold a pointer to it — this is one reason very hot code paths sometimes avoid unnecessary interface conversions.
- **Dynamic dispatch:** calling a method through an interface involves an indirect call (looking up the method via the interface's type descriptor) rather than a direct, staticly-resolved call — normally negligible, but measurable in extremely hot loops.
- **Comparability panics:** as noted in [Section 11](#11-comparing-interface-values), comparing interface values holding uncomparable concrete types panics at runtime — a correctness risk more than a raw performance one, but worth remembering when interfaces are used generically (e.g., as map keys).

In the overwhelming majority of application code, these costs are irrelevant compared to the clarity and flexibility interfaces provide — this is a consideration for narrow, profiled hot paths, not a reason to avoid interfaces generally.

## 18. Common Mistakes and Pitfalls

### 18.1 Returning a Typed Nil as an Interface

Covered in depth in [Section 10](#10-the-nil-interface-pitfall-typed-nil) — the single most common and confusing interface-related bug in Go.

### 18.2 Forgetting the Method-Set Rule for Pointer Receivers

```go
type Speaker interface{ Speak() string }
type Dog struct{}
func (d *Dog) Speak() string { return "Woof" }

var s Speaker = Dog{} // COMPILE ERROR: Dog does not implement Speaker (Speak has pointer receiver)
```

Fix: use `&Dog{}` (or store the value as a pointer from the start) whenever a type's relevant methods use pointer receivers.

### 18.3 Designing Interfaces Too Early or Too Large

Introducing an interface before there's a genuine need for abstraction (only one implementation will ever exist, or it exists purely "for testability" without a real fake ever being written) adds indirection without benefit. Likewise, cramming many unrelated methods into one interface makes it harder for types to satisfy, harder to mock meaningfully, and works against Go's preference for small, focused interfaces.

### 18.4 Ignoring the `ok` in Type Assertions

```go
v := i.(string) // panics if i does not hold a string
```

Prefer the comma-ok form (`v, ok := i.(string)`) unless a panic is genuinely the desired behavior for a type mismatch (e.g., you've already established the type elsewhere and a mismatch would indicate a real bug).

### 18.5 Comparing Interfaces Holding Uncomparable Types

As shown in [Section 11](#11-comparing-interface-values), comparing two `any` values that happen to hold slices, maps, or functions panics — this can surface unexpectedly deep inside generic code (e.g., using an `any` as a map key) long after the code was written.

### 18.6 Treating Implicit Satisfaction as License to Skip Documentation

Because there's no explicit "implements" declaration, it can be unclear from a type's own source file which interfaces it's meant to satisfy. Document the intended interfaces a type implements (in comments, and ideally with a compile-time assertion per [Section 12](#12-compile-time-interface-satisfaction-checks)) so future maintainers don't have to reverse-engineer the relationship.

## 19. Best Practices Summary

1. **Keep interfaces small.** Favor one or two methods that describe a single capability, following the standard library's own style (`io.Reader`, `io.Writer`, `fmt.Stringer`).
2. **Define interfaces at the point of consumption**, not preemptively in the producing package — let real usage patterns drive what abstraction is actually needed.
3. **Accept interfaces as parameters; return concrete types.** This keeps functions flexible for callers while giving callers full access to whatever a constructor actually builds.
4. **Be consistent with receiver types** — if any method on a type needs a pointer receiver, use pointer receivers for all of that type's methods, to avoid subtle, partial interface satisfaction.
5. **Never return a typed-nil variable as an interface value when you mean "no value."** Return the literal `nil` explicitly.
6. **Use the comma-ok form of type assertions** in general-purpose code; reserve the single-result form for cases where a mismatch would indicate a genuine bug.
7. **Add a compile-time satisfaction check** (`var _ SomeInterface = (*SomeType)(nil)`) right after defining a type meant to implement a specific interface.
8. **Prefer generics over `any`** for algorithms/containers whose logic is identical across types and where compile-time type safety matters; reserve `any` for genuinely dynamic data.
9. **Watch for uncomparable concrete types** flowing through `any`/interface values that might be compared or used as map keys.
10. **Document which interfaces a type is meant to satisfy**, since implicit satisfaction provides no such record automatically.

## 20. Use Case Summary Table

| Technique                                      | When to Use                                                                                               | Example Scenario                                                                    |
| ---------------------------------------------- | --------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------- |
| Small, single-method interface                 | Abstract over one capability so multiple implementations can be swapped in                                | `io.Writer` for logging to a file, buffer, or network socket interchangeably        |
| Empty interface (`any`)                        | Genuinely dynamic/unknown data                                                                            | Decoding arbitrary JSON, heterogeneous collections, `fmt`-style formatting          |
| Type assertion (comma-ok)                      | Check whether a value additionally satisfies a more specific interface, or recover a concrete type safely | Checking if an `io.Writer` also implements `io.Closer`                              |
| Type switch                                    | Branch behavior across several known possible concrete types                                              | Encoding logic for different AST node types                                         |
| Interface embedding                            | Compose smaller behaviors into a larger contract                                                          | Building `io.ReadWriteCloser` from `Reader` + `Writer` + `Closer`                   |
| Compile-time assertion (`var _ I = (*T)(nil)`) | Guarantee a type satisfies an interface, catch regressions immediately                                    | Any type meant to implement a public interface (your own or the standard library's) |
| Interface-based test doubles                   | Substitute real dependencies with fakes/mocks in tests                                                    | Faking an email/payment/database dependency in unit tests                           |
| Generics with interface constraints            | Type-safe algorithms/containers shared across many types                                                  | A generic `Stack[T]`, a numeric `Sum[T Number]` function                            |
| `fmt.Stringer`                                 | Customize how a type prints with `%v`/`%s`                                                                | Pretty-printing a domain type like `Point`, `Money`, or `UserID`                    |
| `sort.Interface`                               | Sort any custom collection with `sort.Sort`                                                               | Sorting a slice of custom structs by an arbitrary field                             |

## 21. References

1. Go Team — _A Tour of Go: Methods and Interfaces_. https://go.dev/tour/methods/10
2. Go Team — _Effective Go_. https://go.dev/doc/effective_go
3. Go Language Specification — _Interface types_. https://go.dev/ref/spec#Interface_types
4. Leapcell — _Elegant Interface Implementation in Go: The Beauty of Implicit Contracts_. https://leapcell.io/blog/elegant-interface-implementation-in-go-the-beauty-of-implicit-contracts
5. Sling Academy — _How Interfaces Work in Go: Implicit Implementation Explained_. https://www.slingacademy.com/article/how-interfaces-work-in-go-implicit-implementation-explained/
6. Educative — _Satisfying Interfaces_ and _What is an Interface?_, from _The Way To Go_. https://www.educative.io/courses/the-way-to-go/what-is-an-interface
7. Calvin McLean — _The Magic of Interfaces in Go_, DEV Community. https://dev.to/calvinmclean/the-magic-of-interfaces-in-go-5gga
8. Java Code Geeks — _Go's Interface Satisfaction: Why Explicit Implementation Declarations Are Considered Harmful_. https://www.javacodegeeks.com/2026/02/gos-interface-satisfaction-why-explicit-implementation-declarations-are-considered-harmful.html
9. ByteSizeGo — _Interfaces in Go_. https://www.bytesizego.com/blog/golang-interfaces
10. Gabor Koos — _Go Interfaces - Beyond the Basics_. https://blog.gaborkoos.com/posts/2025-08-13-Go-Interfaces-Beyond-the-Basics/ (also mirrored at https://dev.to/gkoos/go-interfaces-beyond-the-basics-305j)
11. Gabriel Anhaia — _The Strange Case of Go's nil Interface Comparison_, DEV Community. https://dev.to/gabrielanhaia/the-strange-case-of-gos-nil-interface-comparison-2deh
12. Ben Meehan — _Understanding Go's Typed Nil_, Medium. https://medium.com/@ben.meehan_27368/understanding-nil-in-go-interfaces-typed-nil-and-common-pitfalls-6b1154718e00
13. golang/go Issue #25496 — _Method called with NIL pointer receiver when using interfaces_. https://github.com/golang/go/issues/25496
14. Leapcell — _Go Interface Pitfalls: When nil != nil_, DEV Community. https://dev.to/leapcell/go-interface-pitfalls-when-nil-nil-4ai1
15. Alexander Obregon — _Interface Nil Checks in Go_. https://alexanderobregon.substack.com/p/interface-nil-checks-in-go
16. DoltHub Blog — _Much Ado About Nil Things: More Go Pitfalls_. https://www.dolthub.com/blog/2023-09-08-much-ado-about-nil-things/
17. Moksh S — _Interfaces in Go: It Looks Like Duck Typing, But It's Not (Exactly)_, Medium. https://medium.com/@moksh.9/interfaces-in-go-it-looks-like-duck-typing-but-its-not-exactly-9cfbcb09ad06
18. Go standard library documentation — packages `io`, `sort`, `fmt`, `context`. https://pkg.go.dev/io, https://pkg.go.dev/sort, https://pkg.go.dev/fmt, https://pkg.go.dev/context
19. Go Team — _Type Parameters Proposal / Generics documentation_. https://go.dev/doc/tutorial/generics
20. Go Blog — _Why Generics?_. https://go.dev/blog/why-generics
