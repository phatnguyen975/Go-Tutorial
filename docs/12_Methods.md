<div align="center">
  <h1>Methods</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 17, 2026</sub>
</div>

## Table of Contents

1. [What Is a Method?](#1-what-is-a-method)
2. [Methods vs. Functions](#2-methods-vs-functions)
3. [Declaring a Method: Receiver Syntax](#3-declaring-a-method-receiver-syntax)
4. [Value Receivers vs. Pointer Receivers](#4-value-receivers-vs-pointer-receivers)
5. [How Go Decides Which Receiver Form to Use at the Call Site](#5-how-go-decides-which-receiver-form-to-use-at-the-call-site)
6. [Method Sets](#6-method-sets)
7. [Choosing Between Value and Pointer Receivers](#7-choosing-between-value-and-pointer-receivers)
8. [Methods on Non-Struct Named Types](#8-methods-on-non-struct-named-types)
9. [You Cannot Add Methods to Types You Don't Own](#9-you-cannot-add-methods-to-types-you-dont-own)
10. [Receiver Naming Conventions](#10-receiver-naming-conventions)
11. [Method Naming Conventions](#11-method-naming-conventions)
12. [Methods and Embedding: Promotion](#12-methods-and-embedding-promotion)
13. [Method Values](#13-method-values)
14. [Method Expressions](#14-method-expressions)
15. [Methods Satisfying Interfaces](#15-methods-satisfying-interfaces)
16. [Methods with Nil Receivers](#16-methods-with-nil-receivers)
17. [Getters and Setters (or the Lack Thereof)](#17-getters-and-setters-or-the-lack-thereof)
18. [Common Mistakes and Pitfalls](#18-common-mistakes-and-pitfalls)
19. [Best Practices Summary](#19-best-practices-summary)
20. [Use Case Summary Table](#20-use-case-summary-table)
21. [References](#21-references)

## 1. What Is a Method?

A method in Go is a function with a special extra parameter, called the **receiver**, that appears before the method name rather than in the normal parameter list. The receiver binds the function to a specific named type, so it can be called using dot-selector syntax on values of that type.

```go
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

circle := Circle{Radius: 5}
fmt.Println(circle.Area()) // called as a method, via dot syntax
```

Methods are how Go attaches behavior to types without needing classes — a struct (or any named type) plus a set of methods defined on it is Go's equivalent of an object with both data and behavior.

## 2. Methods vs. Functions

A method is really just a function with one difference: it declares a receiver. Everywhere else — parameters, return values, body — methods and functions are identical.

```go
// A plain function
func AreaOfCircle(c Circle) float64 {
    return math.Pi * c.Radius * c.Radius
}

// The equivalent method
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

AreaOfCircle(circle) // function call syntax
circle.Area()          // method call syntax
```

Both accomplish the same computation. The practical differences are:

- **Call syntax:** a method is invoked with `value.Method(...)`; a function is invoked with `Function(value, ...)`.
- **Interface satisfaction:** only methods (not free functions) count toward whether a type satisfies an interface — a type can only be substituted wherever an interface is expected if the required behaviors are declared as methods on it.
- **Discoverability and API design:** grouping related behavior as methods on a type keeps an editor's autocomplete (`value.` followed by a list of methods) useful, and signals "this behavior belongs to this type" more clearly than a same-named free function floating in the package.

**Guidance:** use a method when the operation is conceptually "an action a type can do" or "a property a type has," especially if it needs to participate in an interface. Use a plain function for operations that don't naturally belong to one specific type, or that operate across multiple unrelated types.

## 3. Declaring a Method: Receiver Syntax

The general form of a method declaration is:

```go
func (receiverName ReceiverType) MethodName(parameters) (returnTypes) {
    // method body
}
```

```go
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}
```

Inside the method body, the receiver (`r` here) behaves exactly like any other parameter — you can read its fields, and (depending on whether it's a value or pointer receiver) potentially modify them.

## 4. Value Receivers vs. Pointer Receivers

A method's receiver can be declared as either a plain value (`T`) or a pointer (`*T`), and the choice has real consequences.

### 4.1 Value Receiver

```go
func (c Circle) Scale(factor float64) {
    c.Radius = c.Radius * factor // modifies only the local copy
}
```

A value receiver receives a **copy** of the value the method was called on. Any modification inside the method affects only that copy — the original value the caller holds is completely untouched, exactly like passing a struct to an ordinary function by value.

### 4.2 Pointer Receiver

```go
func (c *Circle) Scale(factor float64) {
    c.Radius = c.Radius * factor // modifies the ORIGINAL Circle through the pointer
}
```

A pointer receiver receives the **address** of the value the method was called on. Modifications inside the method go through that pointer and are visible to the caller once the method returns.

```go
circle := Circle{Radius: 5}
circle.Scale(2) // with a pointer receiver: mutates circle in place
fmt.Println(circle.Radius) // 10
```

## 5. How Go Decides Which Receiver Form to Use at the Call Site

Go provides convenient, automatic conversions so that calling a method generally "just works," regardless of whether you have a value or a pointer, and regardless of whether the method has a value or pointer receiver — with one important exception involving addressability.

```go
c := Circle{Radius: 5}
c.Scale(2)   // c is a value; Scale has a pointer receiver
             // Go automatically rewrites this as (&c).Scale(2)

p := &Circle{Radius: 5}
p.Area()     // p is a pointer; Area has a value receiver
             // Go automatically rewrites this as (*p).Area()
```

- Calling a **pointer-receiver** method on an **addressable value** automatically takes that value's address for you.
- Calling a **value-receiver** method on a **pointer** automatically dereferences the pointer for you.

**The addressability catch:** the automatic `&c` conversion only works when `c` is _addressable_ — a local variable, a struct field of an addressable struct, an array element of an addressable array, and so on. It does **not** work for a value that isn't addressable, such as a struct literal used directly (`Circle{Radius: 5}.Scale(2)` is a compile error if `Scale` has a pointer receiver, since there's no variable to take the address of), or a value obtained from a map index expression (map values aren't addressable in Go).

```go
// Circle{Radius: 5}.Scale(2) // COMPILE ERROR if Scale has a pointer receiver: not addressable
```

## 6. Method Sets

The **method set** of a type is the complete list of methods that can be called on a value of that type. This concept determines both which methods you can call directly, and whether a type satisfies a given interface.

| Receiver type in declaration         | Included in method set of `T` | Included in method set of `*T` |
| ------------------------------------ | ----------------------------- | ------------------------------ |
| `func (t T) M()` (value receiver)    | Yes                           | Yes                            |
| `func (t *T) M()` (pointer receiver) | No                            | Yes                            |

In words: a pointer `*T`'s method set includes **both** value-receiver and pointer-receiver methods, while a plain value `T`'s method set includes **only** value-receiver methods. This asymmetry exists because a pointer-receiver method might need to modify the value or rely on it having a stable address, and the compiler cannot always safely produce an address for an arbitrary value (e.g., one stored in a map, or a non-addressable literal) — so the language simply doesn't extend `T`'s method set to include those methods; only `*T` gets them.

```go
type Greeter struct{ name string }

func (g Greeter) SayHello() string   { return "Hello, " + g.name }
func (g *Greeter) SetName(n string)  { g.name = n }

var g Greeter
// g.SetName("Alice") // fine at the call site (auto address-of), because g is addressable

type Namer interface {
    SetName(string)
}

var n Namer = g   // COMPILE ERROR: Greeter does not implement Namer (SetName has pointer receiver)
var n2 Namer = &g // OK: *Greeter's method set includes SetName
```

This is precisely why, when a type needs to satisfy an interface that includes a pointer-receiver method, you must use a `*T` value (not a plain `T`) wherever that interface is expected — this method-set rule is one of the most common sources of "does not implement interface" compiler errors for newcomers.

## 7. Choosing Between Value and Pointer Receivers

This is one of the most frequently asked design questions in Go. General guidance:

**Use a pointer receiver when:**

- The method needs to **modify** the receiver.
- The receiver is a **large struct**, and copying it on every method call would be wasteful.
- The type contains fields that should not be copied (e.g., a `sync.Mutex` — copying a struct containing a mutex that's already in use is a correctness bug, not just a performance concern).
- Other methods on the type already use pointer receivers (for consistency — see below).

**Use a value receiver when:**

- The method never needs to modify the receiver.
- The receiver is small (a handful of basic-type fields) and copying is cheap.
- The type represents an immutable value conceptually (e.g., a `Point`, a `Money` amount) where treating every operation as producing a new independent value makes sense.

**Consistency rule:** if _any_ method on a type genuinely needs a pointer receiver, it's conventional Go style to give **all** of that type's methods pointer receivers, even ones that don't strictly need to mutate anything. This avoids a confusing situation where a type's method set differs depending on whether some methods used `T` and others used `*T`, and ensures the type satisfies interfaces consistently regardless of which methods a given interface happens to require.

```go
// Consistent: every method on Counter uses a pointer receiver
type Counter struct{ count int }

func (c *Counter) Increment()   { c.count++ }
func (c *Counter) Value() int   { return c.count } // even though Value() doesn't mutate anything
```

## 8. Methods on Non-Struct Named Types

Methods aren't limited to structs — any **named type** defined in the current package can have methods, including named types based on basic types, slices, maps, and functions:

```go
type Celsius float64

func (c Celsius) ToFahrenheit() float64 {
    return float64(c)*9/5 + 32
}

type IntList []int

func (l IntList) Sum() int {
    total := 0
    for _, n := range l {
        total += n
    }
    return total
}

temp := Celsius(100)
fmt.Println(temp.ToFahrenheit()) // 212

nums := IntList{1, 2, 3}
fmt.Println(nums.Sum()) // 6
```

This is a common and powerful idiom: wrapping a basic type or a slice in a named type to attach meaningful, type-safe behavior directly to it, rather than writing free functions that take that type as a plain parameter.

## 9. You Cannot Add Methods to Types You Don't Own

Go enforces that a method's receiver type must be defined in the **same package** as the method itself. You cannot add a method directly to a type from another package (including built-in types like `int` or `string`, or types from the standard library or third-party packages):

```go
// func (t time.Time) MyExtraMethod() {} // COMPILE ERROR: cannot define new methods on non-local type time.Time
```

The idiomatic workaround is to define your own named type based on the foreign type, and attach the method to your new type instead:

```go
type MyTime time.Time

func (t MyTime) IsWeekend() bool {
    wd := time.Time(t).Weekday()
    return wd == time.Saturday || wd == time.Sunday
}
```

This restriction exists to prevent "method injection" conflicts — if any package could attach methods to any type, including types from unrelated packages, two different packages could define conflicting methods with the same name on the same foreign type, creating ambiguity and fragility across the whole ecosystem.

## 10. Receiver Naming Conventions

Idiomatic Go style names receivers with a **short, often one- or two-letter abbreviation** derived from the type name, consistently across every method on that type — not `this` or `self`, which are conventional in many other object-oriented languages but considered unidiomatic in Go.

```go
type Customer struct{ /* ... */ }

func (c Customer) FullName() string    { /* ... */ return "" } // consistently "c"
func (c *Customer) UpdateEmail(e string) { /* ... */ }          // consistently "c"
```

The reasoning behind this convention: the receiver's type is already stated explicitly in every method signature, so a long, descriptive receiver name would be redundant — a short, consistent abbreviation keeps signatures compact and lets the parts of the method body that actually _act_ on the receiver stand out more clearly. Consistency across a type's methods also matters: using different receiver names for different methods on the same type (`c` in one method, `cust` in another) is considered poor style and can be confusing to readers.

## 11. Method Naming Conventions

- Use **MixedCaps** (also called PascalCase for exported names, camelCase for unexported) rather than underscores — `CalculateTotal`, not `calculate_total`.
- Capitalize the first letter to **export** a method (make it visible outside the package); lowercase to keep it package-private.
- Honor Go's well-known method names and their expected signatures/meanings — `Read`, `Write`, `Close`, `Flush`, `String`, and similar names have canonical meanings across the ecosystem (largely from standard-library interfaces). Don't reuse one of these names for something unrelated, and conversely, if your type genuinely implements the same concept as one of these well-known methods, use the standard name and signature rather than inventing your own (e.g., name a string-conversion method `String() string` to satisfy `fmt.Stringer`, not `ToString()`).
- Name methods after the **action or behavior** they perform, not by restating the receiver's type name redundantly (`order.Cancel()`, not `order.CancelOrder()` — the receiver already tells you what's being cancelled).

## 12. Methods and Embedding: Promotion

When a struct embeds another type, the embedded type's methods are **promoted** to the outer struct, meaning they can be called directly on the outer struct as if they belonged to it:

```go
type Base struct{ id int }

func (b Base) Describe() string {
    return fmt.Sprintf("ID: %d", b.id)
}

type Widget struct {
    Base
    Name string
}

w := Widget{Base: Base{id: 1}, Name: "Gadget"}
fmt.Println(w.Describe()) // "ID: 1" — promoted from Base
```

**Crucially, this is not virtual dispatch.** When the promoted `Describe()` method runs, its receiver is bound to the embedded `Base` value, not to the outer `Widget`. If `Base.Describe()` internally called another method that `Widget` happens to also define with the same name, `Base`'s own version would still be the one invoked — there is no mechanism by which the embedded type's methods become aware of, or get overridden by, the outer type's methods when called through the embedded type itself. Go's embedding promotion is purely a convenience for reaching methods without spelling out the embedded field's name at each call site — it deliberately does not implement classical inheritance's polymorphic method resolution.

If the outer struct declares its **own** method with the same name, that new method simply shadows the promoted one for any call made directly on the outer type (`w.Describe()` would now call `Widget`'s own method); the embedded type's original method remains reachable by explicitly qualifying through the embedded field (`w.Base.Describe()`).

## 13. Method Values

A **method value** is created by selecting a method from a specific receiver value without calling it (no parentheses) — the result is an ordinary function value that has already "remembered" (bound) that particular receiver, so it can be called later without needing the receiver again:

```go
type Greeter struct{ prefix string }

func (g Greeter) Say(name string) string {
    return g.prefix + " " + name
}

g := Greeter{prefix: "Hello,"}
say := g.Say // method value: g is captured/bound right now

fmt.Println(say("Alice")) // "Hello, Alice" — say has type func(string) string
```

The receiver expression (`g` here) is evaluated and saved at the moment the method value is created — later calls to `say` always use that saved copy, even if `g` itself changes afterward. This makes method values useful for passing a bound behavior around as a plain function — for example, passing `logger.Log` as a callback, or `wg.Done` (from `sync.WaitGroup`) directly into `defer`.

## 14. Method Expressions

A **method expression** is written as `Type.Method` (or `(*Type).Method` for pointer-receiver methods) rather than `value.Method`. Unlike a method value, no specific receiver is bound yet — the resulting function value takes the receiver explicitly as its **first parameter**:

```go
type Greeter struct{ prefix string }

func (g Greeter) Say(name string) string {
    return g.prefix + " " + name
}

sayFn := Greeter.Say // method expression: type func(Greeter, string) string

g := Greeter{prefix: "Hi,"}
fmt.Println(sayFn(g, "Bob")) // "Hi, Bob" — receiver passed explicitly as the first argument
```

Method expressions are considerably rarer in everyday Go code than method values — they're mainly useful when you need a function that can be applied to a _variable_, not-yet-chosen receiver, such as building a table of operations that should each be applied to whatever receiver is supplied at call time.

## 15. Methods Satisfying Interfaces

A type's method set (see [Section 6](#6-method-sets)) determines exactly which interfaces it satisfies — an interface is satisfied when a type's method set includes every method the interface requires, matching name, parameters, and return types exactly.

```go
type Stringer interface {
    String() string
}

type Point struct{ X, Y int }

func (p Point) String() string {
    return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

var s Stringer = Point{1, 2} // Point satisfies Stringer because it has a String() string method
fmt.Println(s) // (1, 2)
```

Because interface satisfaction is checked against the method set, remember that if `String()` had been declared with a pointer receiver (`func (p *Point) String() string`), only `*Point` — not plain `Point` — would satisfy `Stringer`, which is a frequent source of "type does not implement interface" compile errors when a type's methods mix value and pointer receivers inconsistently.

## 16. Methods with Nil Receivers

Unlike calling a method on a `nil` pointer of an _interface_ type (which panics immediately, since there's no concrete type to dispatch to), calling a **pointer-receiver method** on a `nil` value of a **concrete pointer type** is perfectly legal in Go, as long as the method body doesn't actually dereference the nil pointer:

```go
type List struct {
    Value int
    Next  *List
}

func (l *List) Len() int {
    if l == nil {
        return 0 // safely handle the nil case before touching any field
    }
    return 1 + l.Next.Len()
}

var l *List // nil
fmt.Println(l.Len()) // 0 — works, because Len() checks for nil before dereferencing
```

This pattern is used deliberately in some standard-library and idiomatic Go code (recursive data structures like linked lists and trees are a classic example) to let a `nil` pointer represent "an empty instance" gracefully, rather than requiring every caller to check for `nil` before calling a method. It only works because the method itself checks `l == nil` before trying to access `l`'s fields — if it dereferenced `l` first, it would panic exactly like any other nil pointer dereference.

## 17. Getters and Setters (or the Lack Thereof)

Go culture generally discourages writing trivial getter/setter methods (`GetName()`, `SetName()`) purely as a matter of habit carried over from other languages. Because Go has no `private`/`public` keyword distinction beyond capitalization, an exported field can simply be accessed and set directly:

```go
// Not idiomatic Go, if Name has no special invariants to protect:
func (c *Customer) GetName() string      { return c.name }
func (c *Customer) SetName(n string)     { c.name = n }

// More idiomatic, if there's nothing to protect:
type Customer struct {
    Name string // exported field, accessed directly: customer.Name
}
```

**When a method genuinely earns its place:** if setting or getting a field needs to enforce an invariant (validation, derived/computed state, side effects, or the field must stay unexported for encapsulation reasons), a method is the right tool — just don't reach for `GetX`/`SetX` boilerplate reflexively when a plain exported field would do. When a getter is genuinely needed, Go convention also drops the `Get` prefix entirely — a method that returns a customer's balance is simply named `Balance()`, not `GetBalance()`.

## 18. Common Mistakes and Pitfalls

### 18.1 Expecting a Value Receiver to Mutate the Original

```go
type Counter struct{ n int }

func (c Counter) Increment() { c.n++ } // BUG: value receiver, mutates only the copy

c := Counter{}
c.Increment()
fmt.Println(c.n) // 0, not 1 — silent bug, no compile error
```

This compiles fine and produces no error or warning — it simply doesn't do what the caller likely intended. Fix: use a pointer receiver (`func (c *Counter) Increment()`) whenever the method needs to mutate state.

### 18.2 Mixing Value and Pointer Receivers Inconsistently

```go
type Account struct{ balance float64 }

func (a Account) Balance() float64     { return a.balance }
func (a *Account) Deposit(amt float64) { a.balance += amt }
```

While technically legal, this inconsistency can create confusing interface-satisfaction situations (only `*Account`, not `Account`, satisfies any interface requiring `Deposit`) and makes it less obvious at a glance whether a given method call mutates the receiver. Prefer uniformity across a type's methods once any one of them needs a pointer receiver.

### 18.3 Calling a Pointer-Receiver Method on a Non-Addressable Value

```go
type Circle struct{ Radius float64 }
func (c *Circle) Scale(f float64) { c.Radius *= f }

// Circle{Radius: 5}.Scale(2) // COMPILE ERROR: cannot take address of Circle{Radius: 5} literal
```

Fix: assign the literal to a variable first (`c := Circle{Radius: 5}; c.Scale(2)`), or construct it as a pointer from the start (`c := &Circle{Radius: 5}; c.Scale(2)`).

### 18.4 Assuming Embedded Methods Get Overridden Polymorphically

```go
type Base struct{}
func (b Base) Greet() string { return "Hi from " + b.who() }
func (b Base) who() string   { return "Base" }

type Derived struct{ Base }
func (d Derived) who() string { return "Derived" }

d := Derived{}
fmt.Println(d.Greet()) // "Hi from Base" — NOT "Hi from Derived"
```

`Base.Greet()` calls `Base.who()`, not `Derived.who()`, because embedding provides no virtual dispatch — `Base`'s methods only ever see `Base`'s own method set, regardless of what outer type embeds it. This is a common source of confusion for developers expecting classical-OOP-style method overriding.

### 18.5 Forgetting You Can't Add Methods to Foreign Types

```go
// func (s string) Shout() string { return strings.ToUpper(s) + "!" }
// COMPILE ERROR: cannot define new methods on non-local type string
```

Fix: define a named type based on the foreign type, and attach the method to that instead (`type Shout string; func (s Shout) Loud() string { ... }`).

### 18.6 Copying a Struct That Contains a Mutex (or Similar Type)

```go
import "sync"

type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c SafeCounter) Increment() { // BUG: value receiver copies the mutex too!
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

Copying a `sync.Mutex` (which happens automatically with a value receiver) produces an independent, useless copy of the lock — different callers end up locking different copies of the mutex, providing no actual mutual exclusion, and `go vet` will typically flag this. Any type containing a `sync.Mutex` (or similarly non-copyable synchronization primitives) should always use pointer receivers for every method, and should generally be passed around by pointer everywhere, not just within methods.

## 19. Best Practices Summary

1. **Use a pointer receiver whenever the method must mutate the receiver**, or when the receiver is large enough that copying it would be wasteful.
2. **Be consistent**: once any method on a type needs a pointer receiver, use pointer receivers for all of that type's methods.
3. **Name receivers with a short, consistent abbreviation** of the type name across every method — avoid `this`/`self`, and avoid varying the name method to method.
4. **Name methods after the action/behavior**, not by redundantly restating the receiver's type.
5. **Honor well-known method names and signatures** (`String`, `Read`, `Write`, `Close`, etc.) so your type integrates smoothly with standard-library expectations.
6. **Remember embedding promotion is not inheritance** — an embedded type's methods never see the outer type's overriding methods.
7. **Don't write reflexive `GetX`/`SetX` boilerplate** for fields with no invariants to protect — export the field directly instead, and drop the `Get` prefix on any getter you do write.
8. **Never give a value receiver to a type containing a `sync.Mutex`** (or any other type that must not be copied) — always use pointer receivers and pass such types by pointer.
9. **Define a named wrapper type when you need methods on a type you don't own** (a foreign package's type, or a built-in type) — Go does not allow attaching methods directly to non-local types.
10. **Use method values for bound callbacks** (passing `x.Method` as a function argument) and reserve method expressions for the rarer case of needing the receiver as an explicit, later-supplied argument.

## 20. Use Case Summary Table

| Technique                               | When to Use                                                                | Example Scenario                                                 |
| --------------------------------------- | -------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| Value receiver                          | Method never mutates the receiver; type is small                           | `func (p Point) Distance(other Point) float64`                   |
| Pointer receiver                        | Method mutates the receiver, or the receiver is large / must not be copied | `func (c *Counter) Increment()`, any type embedding `sync.Mutex` |
| Method on a non-struct named type       | Attach type-safe behavior to a basic type or slice/map alias               | `type Celsius float64` with a `ToFahrenheit()` method            |
| Named wrapper type + method             | Need a method on a type you don't own (foreign package, built-in type)     | `type MyTime time.Time` with an `IsWeekend()` method             |
| Method value (`x.Method`)               | Pass a bound behavior around as a plain callback                           | `defer wg.Done`, passing `logger.Log` into another function      |
| Method expression (`Type.Method`)       | Need the receiver supplied explicitly/later, e.g. in a table of operations | Building a dispatch table keyed by operation name                |
| Nil-safe pointer-receiver method        | A recursive/tree-like type where `nil` naturally represents "empty"        | `func (l *List) Len() int` returning 0 for a nil `*List`         |
| Embedding for method promotion          | Reuse a stable, cohesive type's behavior without repeating boilerplate     | Embedding a shared logging/base type across several structs      |
| Exported field instead of getter/setter | No invariant to enforce; direct access is simpler                          | `type Customer struct { Name string }`                           |

## 21. References

1. Go Team — _A Tour of Go: Methods_. https://go.dev/tour/methods/1
2. Go Team — _Effective Go_ (Methods, Embedding sections). https://go.dev/doc/effective_go
3. Go Language Specification — _Method declarations, Method sets_. https://go.dev/ref/spec#Method_declarations
4. codedamn — _Methods in Go_. https://codedamn.com/news/go/methods-in-go
5. go101.org — _Methods in Go_. https://go101.org/article/method.html
6. Basant C. — _Mastering Types, Methods, and Receivers in Go: The Ultimate Idiomatic Guide_, Medium. https://medium.com/@caring_smitten_gerbil_914/mastering-types-methods-and-receivers-in-go-the-ultimate-idiomatic-guide-753ee786c818
7. golang-nuts mailing list — _Naming convention for method receiver?_. https://groups.google.com/g/golang-nuts/c/73npFHhuL9k
8. Gabriel Anhaia — _Method Values vs Method Expressions in Go: A Distinction Worth Knowing_, DEV Community. https://dev.to/gabrielanhaia/method-values-vs-method-expressions-in-go-a-distinction-worth-knowing-4ke9
9. Boldly Go — _Method declarations_. https://boldlygo.tech/archive/2023-08-24-method-declarations/
10. Boldly Go — _Method values_. https://boldlygo.tech/archive/2023-11-02-method-values/
11. Boldly Go — _Method expressions, conclusion_. https://boldlygo.tech/archive/2023-11-01-method-expressions-conclusion/
12. Go standard library documentation — package `fmt` (`Stringer` interface and canonical method names). https://pkg.go.dev/fmt
13. Go standard library documentation — package `sync` (`Mutex` non-copyability). https://pkg.go.dev/sync
14. `go vet` documentation — copylocks check. https://pkg.go.dev/cmd/vet
