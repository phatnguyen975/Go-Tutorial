<div align="center">
  <h1>Structs</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 17, 2026</sub>
</div>

## Table of Contents

1. [What Is a Struct?](#1-what-is-a-struct)
2. [Declaring a Struct Type](#2-declaring-a-struct-type)
3. [Creating and Initializing Struct Values](#3-creating-and-initializing-struct-values)
4. [Accessing and Modifying Fields](#4-accessing-and-modifying-fields)
5. [Structs Are Value Types](#5-structs-are-value-types)
6. [Pointers to Structs](#6-pointers-to-structs)
7. [Anonymous Structs](#7-anonymous-structs)
8. [Struct Comparison and Comparability](#8-struct-comparison-and-comparability)
9. [Struct Tags](#9-struct-tags)
10. [Struct Embedding (Composition)](#10-struct-embedding-composition)
11. [Field and Method Promotion](#11-field-and-method-promotion)
12. [Name Conflicts and Shadowing in Embedding](#12-name-conflicts-and-shadowing-in-embedding)
13. [Embedding Interfaces in Structs](#13-embedding-interfaces-in-structs)
14. [Exported vs. Unexported Fields](#14-exported-vs-unexported-fields)
15. [Empty Structs (`struct{}`)](#15-empty-structs-struct)
16. [Struct Alignment and Memory Layout](#16-struct-alignment-and-memory-layout)
17. [Constructors: The `NewXxx` Convention](#17-constructors-the-newxxx-convention)
18. [Common Mistakes and Pitfalls](#18-common-mistakes-and-pitfalls)
19. [Best Practices Summary](#19-best-practices-summary)
20. [Use Case Summary Table](#20-use-case-summary-table)
21. [References](#21-references)

## 1. What Is a Struct?

A **struct** is a composite data type that groups together zero or more named fields (also called members), each with its own type, under a single new type name. Structs are Go's primary tool for modeling structured data — a record, an entity, a configuration bundle, anything with multiple related pieces of information that belong together.

```go
type Person struct {
    Name string
    Age  int
}
```

Go does not have classes in the traditional object-oriented sense. Structs, combined with methods (functions bound to a type — a topic large enough to deserve its own dedicated treatment) and interfaces, are how Go achieves the equivalent of encapsulating data and behavior, favoring **composition over inheritance** as its core organizing principle.

## 2. Declaring a Struct Type

A struct type is declared with the `type` keyword, a name, and the `struct` keyword followed by a field list in braces:

```go
type Rectangle struct {
    Width  float64
    Height float64
}

type Employee struct {
    ID       int
    Name     string
    Salary   float64
    IsActive bool
}
```

Fields of the same type can be grouped on one line:

```go
type Point struct {
    X, Y int
}
```

A struct can contain fields of any type — basic types, other structs, slices, maps, pointers, interfaces, functions, or even channels:

```go
type Order struct {
    ID       string
    Items    []string
    Metadata map[string]string
    Customer *Person // pointer to another struct
}
```

## 3. Creating and Initializing Struct Values

There are several ways to create a struct value:

### 3.1 Zero Value

Declaring a struct variable without initialization gives every field its own zero value (`0` for numbers, `""` for strings, `false` for booleans, `nil` for pointers/slices/maps/interfaces/channels/functions):

```go
var p Person
fmt.Println(p) // {  0}  — Name is "", Age is 0
```

### 3.2 Struct Literal (Positional)

```go
p := Person{"Alice", 30} // fields set in declaration order
```

This form is fragile — it breaks if the struct's field order or count ever changes — so it's generally discouraged outside of very small, stable, well-known structs.

### 3.3 Struct Literal (Keyed) — Preferred

```go
p := Person{
    Name: "Alice",
    Age:  30,
}
```

Keyed literals are self-documenting, allow fields to be listed in any order, and allow omitting fields you don't need to set (they get their zero value automatically):

```go
p := Person{Name: "Bob"} // Age defaults to 0
```

### 3.4 Using `new()`

```go
p := new(Person) // p is *Person, all fields zero-valued
p.Name = "Carol"
```

### 3.5 Address of a Literal

```go
p := &Person{Name: "Dave", Age: 25} // p is *Person directly
```

This is the most common pattern for constructor-style functions that need to return a pointer.

## 4. Accessing and Modifying Fields

Fields are accessed with the dot (`.`) operator:

```go
p := Person{Name: "Alice", Age: 30}
fmt.Println(p.Name) // Alice
p.Age = 31          // modify a field
fmt.Println(p.Age)  // 31
```

If you have a pointer to a struct, the dot operator works identically — Go automatically dereferences the pointer for you:

```go
ptr := &p
fmt.Println(ptr.Name) // Alice — no need to write (*ptr).Name
ptr.Age = 32          // mutates the original struct through the pointer
```

## 5. Structs Are Value Types

Assigning a struct, or passing it to a function, copies **every field**. This is a fundamental and consistent rule in Go — unlike slices or maps, a struct does not implicitly share underlying storage.

```go
type Point struct{ X, Y int }

func modify(p Point) {
    p.X = 999 // modifies only the local copy
}

func main() {
    a := Point{X: 1, Y: 2}
    b := a       // b is a completely independent COPY of a
    b.X = 100

    fmt.Println(a.X) // 1 — unaffected by changes to b
    fmt.Println(b.X) // 100

    modify(a)
    fmt.Println(a.X) // still 1 — modify() only changed its own copy
}
```

For small structs, this copying is cheap and often preferable for its simplicity and safety (no risk of unintended shared mutation). For large structs, or when a function genuinely needs to mutate the caller's original struct, pass a pointer instead (see [Section 6](#6-pointers-to-structs)).

## 6. Pointers to Structs

A pointer to a struct (`*Person`) lets a function operate on the original struct rather than a copy, and avoids copying every field when the struct is large:

```go
func birthday(p *Person) {
    p.Age++ // mutates the caller's actual struct
}

func main() {
    alice := Person{Name: "Alice", Age: 30}
    birthday(&alice)
    fmt.Println(alice.Age) // 31
}
```

This pattern — pass a pointer when you need to mutate, pass a value when you don't — is one of the most common design decisions in everyday Go code, and it directly parallels the choice between value and pointer receivers on methods.

## 7. Anonymous Structs

A struct type can be declared **inline**, without a separate `type` declaration, when a one-off grouping of fields is needed in a single place:

```go
person := struct {
    Name string
    Age  int
}{
    Name: "Eve",
    Age:  28,
}

fmt.Println(person.Name) // Eve
```

Anonymous structs are useful in a few recurring situations:

- **Table-driven tests**, where each test case is a small, local bundle of inputs and expected outputs that only matters within that one test function.
- **Decoding a subset of a JSON payload**, when you only need a few fields and don't want to declare (and maintain) a full named type just for that.
- **A composite map key**, when you need to key a map on a combination of values and don't want a separate named type just for the key.

```go
tests := []struct {
    input    int
    expected int
}{
    {input: 2, expected: 4},
    {input: 3, expected: 9},
}

for _, tc := range tests {
    if got := square(tc.input); got != tc.expected {
        t.Errorf("square(%d) = %d, want %d", tc.input, got, tc.expected)
    }
}
```

**Limitation:** because an anonymous struct literal has no type name, you cannot declare methods on it — a method's receiver must always refer to a named type. If a shape of data needs behavior (methods) or is reused across multiple places in a codebase, give it a proper named type instead.

## 8. Struct Comparison and Comparability

Two struct values can be compared with `==`/`!=` **if and only if every field's type is itself comparable**. Comparable field types include numbers, strings, booleans, pointers, channels, interfaces, arrays of comparable element types, and other structs made entirely of comparable fields. Fields of type slice, map, or function make the whole struct **not** comparable — attempting `==` on such a struct is a compile-time error, not a runtime one.

```go
type Point struct{ X, Y int }

p1 := Point{1, 2}
p2 := Point{1, 2}
fmt.Println(p1 == p2) // true — same field values

type Container struct {
    Items []int // slice — not comparable
}

// c1 := Container{[]int{1}}
// c2 := Container{[]int{1}}
// c1 == c2 // COMPILE ERROR: struct containing []int cannot be compared
```

When a struct isn't comparable with `==` because of a slice/map/function field, and you need to check "are these logically equal," use `reflect.DeepEqual` (or write your own field-by-field comparison), understanding that `reflect.DeepEqual` has different semantics (structural, recursive comparison) and a real performance cost compared to `==`.

## 9. Struct Tags

A **struct tag** is an optional, raw string literal attached after a field's type, conventionally holding `key:"value"` pairs that libraries (most commonly encoding packages like `encoding/json`, `encoding/xml`, and ORMs) read via reflection to customize how a field is treated:

```go
type User struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    Email     string `json:"email,omitempty"`
    Password  string `json:"-"` // "-" tells encoding/json to always skip this field
}
```

Common tag conventions:

| Tag                     | Purpose                                                                                    | Example                     |
| ----------------------- | ------------------------------------------------------------------------------------------ | --------------------------- |
| `json:"name"`           | Rename the field for JSON encoding/decoding                                                | `json:"first_name"`         |
| `json:"name,omitempty"` | Omit the field from JSON output when it's the zero value                                   | `json:"email,omitempty"`    |
| `json:"-"`              | Always exclude the field from JSON                                                         | `json:"-"`                  |
| `xml:"name"`            | Equivalent renaming for `encoding/xml`                                                     | `xml:"Name"`                |
| `db:"column_name"`      | Map a field to a database column (used by various SQL libraries)                           | `db:"user_id"`              |
| `validate:"required"`   | Declarative validation rules (used by validation libraries like `go-playground/validator`) | `validate:"required,email"` |

Tags are plain strings as far as the Go compiler is concerned — the compiler does not interpret or enforce them at all. Their meaning is entirely defined by whichever library reads them via reflection (`reflect.StructTag`), so an unsupported or misspelled tag simply does nothing rather than causing a compile error, which is worth double-checking carefully since typos in tags fail silently.

**Important caveat with embedding:** tags do not automatically merge or compose across an embedded struct's fields — how an embedded struct's tagged fields are treated (flattened into the parent's JSON object, or nested under their own key) depends on the specific library's own rules; `encoding/json`, for example, flattens anonymous (embedded) struct fields into the parent object by default.

## 10. Struct Embedding (Composition)

Go does not have class inheritance. Instead, it provides **embedding**: including one struct type inside another **without giving it an explicit field name** (an "anonymous field"). This is Go's primary mechanism for composition and code reuse.

```go
type Animal struct {
    Name string
}

func (a Animal) Describe() string {
    return "Animal: " + a.Name
}

type Dog struct {
    Animal // embedded field — no field name, just the type
    Breed  string
}

func main() {
    d := Dog{
        Animal: Animal{Name: "Rex"},
        Breed:  "Labrador",
    }

    fmt.Println(d.Name)       // "Rex" — promoted field access
    fmt.Println(d.Describe()) // "Animal: Rex" — promoted method
    fmt.Println(d.Animal.Name) // "Rex" — can also access explicitly via the embedded type's name
}
```

The embedded field's implicit name is the type name itself (`Animal`), which is why `d.Animal` works to access the embedded value directly, even though it wasn't given an explicit field name.

**This is not inheritance.** `Dog` does not "become" an `Animal` in the way a subclass becomes its superclass — there is no polymorphic substitutability implied (`Dog` does not automatically satisfy an interface requiring an `Animal` parameter, for example). It is purely a **composition** mechanism: `Dog` _has_ an `Animal` inside it, and Go's compiler generates convenient shorthand ("promotion," see [Section 11](#11-field-and-method-promotion)) for reaching the embedded value's fields and methods without needing to write `d.Animal.Name` every time.

**Good fits for embedding:** shared behavior across otherwise-unrelated types (a common logging helper, a mutex for synchronization), wrapper types that extend a base implementation, and composing two cohesive values where promotion genuinely makes call sites clearer. Embedding too many types, or types with a very wide exported surface, can obscure where fields and methods actually come from, and can widen a struct's effective exported API by accident — when the relationship should be explicit and visible to readers, use a regular **named field** instead of embedding:

```go
type Dog struct {
    Pet Animal // named field: explicit, no promotion — d.Pet.Name required
    Breed string
}
```

## 11. Field and Method Promotion

When a struct embeds another type, that type's exported fields and methods are **promoted** to the outer struct — meaning they become accessible through the outer struct directly, as if they were declared on it, as long as there's no ambiguity (see [Section 12](#12-name-conflicts-and-shadowing-in-embedding)).

```go
type Base struct {
    ID int
}

func (b Base) Describe() string {
    return fmt.Sprintf("ID: %d", b.ID)
}

type Widget struct {
    Base // embedded
    Name string
}

w := Widget{Base: Base{ID: 42}, Name: "Gadget"}
fmt.Println(w.ID)         // 42 — promoted field
fmt.Println(w.Describe()) // "ID: 42" — promoted method
```

**A key subtlety:** when a promoted method executes, its receiver is still the _embedded_ type, not the outer type — calling `w.Describe()` runs `Base.Describe()` with `b` bound to `w.Base`, not to `w` itself. This is precisely why embedding is not the same as classical inheritance with virtual dispatch: the embedded type's methods have no awareness of, and cannot be "overridden" from the perspective of, the outer type calling into them polymorphically. If the outer type defines its own method with the same name, that new method simply shadows the promoted one for calls made directly on the outer type — but the embedded type's own method still only ever sees the embedded value when called through promotion, or through its own type's method set.

Promotion also applies when a struct embeds a **pointer** to another struct (`*Base` instead of `Base`), which additionally allows the embedded value to be `nil` and shared across multiple owners, at the cost of needing to ensure it's initialized before use.

## 12. Name Conflicts and Shadowing in Embedding

If the outer struct declares a field or method with the **same name** as one coming from an embedded type, the outer struct's own field/method takes precedence — the embedded type's version is "shadowed," but still reachable by qualifying it explicitly through the embedded field's name:

```go
type Base struct {
    tag string
}

func (b Base) DescribeTag() string {
    return "Base tag is " + b.tag
}

type Container struct {
    Base
    tag string // shadows Base.tag
}

func (co Container) DescribeTag() string {
    return "Container tag is " + co.tag // shadows Base.DescribeTag
}

b := Base{tag: "b's tag"}
co := Container{Base: b, tag: "co's tag"}

fmt.Println(b.DescribeTag())        // "Base tag is b's tag"
fmt.Println(co.DescribeTag())       // "Container tag is co's tag"
fmt.Println(co.Base.DescribeTag())  // "Base tag is b's tag" — explicit access to the shadowed method
```

If **two or more** embedded types at the same nesting depth both have a field or method of the same name, and the outer struct doesn't resolve the conflict with its own field/method of that name, accessing that name through the outer struct directly is a compile-time ambiguity error — you must qualify which embedded type's version you mean.

## 13. Embedding Interfaces in Structs

A struct can also embed an **interface** type (not just another struct). This is a common pattern for building flexible wrapper types, adapters, and partial implementations:

```go
type Logger interface {
    Log(msg string)
}

type Service struct {
    Logger // embedded interface field
    Name   string
}

func main() {
    svc := Service{
        Logger: someLoggerImplementation{},
        Name:   "OrderService",
    }
    svc.Log("service started") // promoted from the embedded Logger
}
```

This pattern lets `Service` satisfy any interface requiring a `Log` method (via promotion) without `Service` implementing `Log` itself — the call is simply forwarded to whatever concrete `Logger` was assigned. It's frequently used to build test doubles or decorators that only need to override a subset of an interface's methods, while delegating everything else to an embedded real implementation.

## 14. Exported vs. Unexported Fields

Go's visibility rules apply to struct fields exactly as they do to any other identifier: a field name starting with an **uppercase** letter is exported (accessible from other packages); a field starting with a **lowercase** letter is unexported (package-private).

```go
type Account struct {
    ID      string // exported — visible outside the package
    balance float64 // unexported — only accessible within this package
}
```

Unexported fields are a common way to enforce invariants: external code cannot set `balance` directly and must go through exported methods (if any are provided) that can validate or control how it changes. This is one of the main tools Go offers for encapsulation, since the language has no `private`/`public` keywords — visibility is entirely determined by capitalization and the package boundary.

**Caveat:** libraries that use reflection (like `encoding/json`) generally cannot see or set unexported fields at all — only exported fields participate in JSON marshaling/unmarshaling, struct tag processing, and similar reflection-based tooling by default.

## 15. Empty Structs (`struct{}`)

`struct{}` is a struct type with **zero fields**, and it is guaranteed to occupy zero bytes of memory. It's a common idiom for values that carry no data at all — only their _presence_ matters:

```go
seen := make(map[string]struct{})
seen["alice"] = struct{}{}

if _, ok := seen["alice"]; ok {
    fmt.Println("already seen")
}
```

Using `map[string]struct{}` as a "set" (rather than `map[string]bool`) communicates intent more precisely — the value genuinely carries no information, whereas a `bool` value implies there might be a meaningful "false" state, when in practice a set only ever cares about key presence. This same idiom (`chan struct{}`) is also common for pure signaling channels, where only the fact that a value was sent (or the channel closed) matters, not any particular payload.

## 16. Struct Alignment and Memory Layout

The order in which fields are declared can affect a struct's total memory size, due to **alignment padding** the compiler inserts so each field starts at an address suitable for its type (e.g., an 8-byte `int64` typically needs to start at an address that's a multiple of 8).

```go
type Inefficient struct {
    A bool    // 1 byte
    B int64   // 8 bytes — needs 7 bytes of padding before it for alignment
    C bool    // 1 byte
    // total: significantly more than 1+8+1 = 10 bytes, due to padding
}

type Efficient struct {
    B int64  // 8 bytes
    A bool   // 1 byte
    C bool   // 1 byte
    // total: much closer to the sum of field sizes, less padding needed
}
```

Ordering fields from largest to smallest alignment requirement generally minimizes total padding. This is a micro-optimization relevant mainly to structs allocated in very large quantities (e.g., millions of elements in a slice) — for most application-level structs, the difference is negligible and readability/logical grouping of fields should take priority. You can inspect a struct's actual size with `unsafe.Sizeof` if you need to verify this in a specific hot path.

## 17. Constructors: The `NewXxx` Convention

Go has no built-in constructor syntax (no special method that always runs on creation). Instead, the idiomatic pattern is a plain function, conventionally named `NewXxx`, that builds and returns a properly initialized value (often a pointer):

```go
type Client struct {
    baseURL string
    timeout time.Duration
}

func NewClient(baseURL string) *Client {
    return &Client{
        baseURL: baseURL,
        timeout: 30 * time.Second, // sensible default
    }
}

c := NewClient("https://api.example.com")
```

This pattern is especially valuable when a struct has unexported fields that need internal setup or validation, or when you want to guarantee certain invariants (a non-empty default, a required dependency) are always satisfied before the value is used — something a bare struct literal, built directly by any caller, cannot enforce.

## 18. Common Mistakes and Pitfalls

### 18.1 Forgetting Structs Copy on Assignment/Pass

```go
type Config struct{ Debug bool }

func enableDebug(c Config) {
    c.Debug = true // only changes the local copy
}

cfg := Config{}
enableDebug(cfg)
fmt.Println(cfg.Debug) // false — unaffected
```

Fix: pass a pointer (`*Config`) if the function needs to mutate the caller's struct.

### 18.2 Comparing Structs That Aren't Comparable

```go
type Data struct {
    Tags []string // slice — not comparable
}

// d1 == d2 // COMPILE ERROR
```

Fix: use `reflect.DeepEqual`, or write a manual field-by-field equality method if `==` semantics genuinely aren't available.

### 18.3 Relying on Positional Struct Literals

```go
p := Person{"Alice", 30} // breaks silently if fields are reordered later
```

Fix: use keyed literals (`Person{Name: "Alice", Age: 30}`), which remain correct even if the struct's field order changes.

### 18.4 Assuming Embedding Provides Polymorphism

```go
type Animal struct{ Name string }
type Dog struct{ Animal }

func feed(a Animal) { /* ... */ }

d := Dog{Animal: Animal{Name: "Rex"}}
// feed(d) // COMPILE ERROR: Dog is not an Animal, even though it embeds one
feed(d.Animal) // must explicitly pass the embedded value
```

Embedding gives convenient field/method promotion, not an "is-a" substitutability relationship — a `Dog` is never automatically usable wherever an `Animal` is expected.

### 18.5 Mutating a Struct Field Obtained From a Map

```go
type Point struct{ X, Y int }
m := map[string]Point{"origin": {0, 0}}

// m["origin"].X = 5 // COMPILE ERROR: cannot assign to struct field in map value (not addressable)
```

Map values are not addressable in Go, so you can't mutate a struct field in place through a map index expression. Fix: either store pointers (`map[string]*Point`), or read the struct out, modify the copy, and write it back:

```go
p := m["origin"]
p.X = 5
m["origin"] = p
```

### 18.6 Unintentionally Widening a Public API via Embedding

Embedding a type with a large exported method/field surface inside a public struct silently makes all of that surface part of your own type's public API, which can be surprising to consumers and constrains future refactoring (removing or changing the embedded type becomes a breaking change for anyone relying on the promoted members). Be deliberate about what you embed in exported types.

## 19. Best Practices Summary

1. **Prefer keyed struct literals** (`Person{Name: "x", Age: 1}`) over positional ones for anything beyond the smallest, most stable structs.
2. **Pass small structs by value, large or mutable-in-place structs by pointer** — mirror this decision consistently with how methods on the same type choose their receivers.
3. **Use unexported fields plus exported constructor/accessor functions** when you need to enforce invariants that a bare struct literal wouldn't guarantee.
4. **Reach for embedding to reuse behavior across cohesive, stable types**, not as a substitute for interfaces or as an attempt to simulate class inheritance.
5. **Use a named field instead of embedding** whenever the relationship should be visible and explicit to readers, or when embedding would unintentionally widen your type's public surface.
6. **Use anonymous structs** for genuinely one-off, local groupings (table-driven test cases, partial JSON decoding, composite map keys) — give data a named type as soon as it's reused or needs methods.
7. **Double-check struct tags carefully** — a typo silently does nothing rather than causing a compile error, since tags are just strings interpreted by reflection.
8. **Use `struct{}`** for set-like maps and pure signaling channels, to communicate "presence only, no payload" more precisely than `bool`.
9. **Don't over-optimize field ordering for alignment** unless profiling shows it matters for a specific, large-scale hot path — readability and logical grouping usually matter more.
10. **Remember map values aren't addressable** — read, modify, and write back a struct copy, or store pointers in the map, when you need to mutate struct fields held in a map.

## 20. Use Case Summary Table

| Technique                     | When to Use                                                                           | Example Scenario                                         |
| ----------------------------- | ------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| Named struct + keyed literal  | Any reusable, multi-field data model                                                  | `type User struct{...}`, `User{Name: "x", Email: "y"}`   |
| Struct passed by value        | Small struct, no mutation needed, simplicity/safety preferred                         | A `Point{X, Y}` passed into a math function              |
| Struct passed by pointer      | Large struct, or the function must mutate the caller's data                           | Filling in default fields of a large `Config`            |
| Anonymous struct              | One-off grouping used in exactly one place                                            | Table-driven test cases, a composite map key             |
| Struct tags                   | Mapping a Go type to an external format/schema                                        | `json:"..."`, `db:"..."`, `validate:"..."`               |
| Struct embedding              | Reusing behavior/data from a stable, cohesive type                                    | Embedding a `sync.Mutex`, wrapping a base implementation |
| Named field (no embedding)    | The relationship should be explicit, or embedding would widen the public API too much | `type Dog struct { Pet Animal }`                         |
| Embedded interface            | Partial implementation / delegation of a contract                                     | A decorator that overrides only one method of a `Logger` |
| `struct{}`                    | Presence-only data — sets, pure signaling                                             | `map[string]struct{}`, `chan struct{}`                   |
| `NewXxx` constructor function | Enforcing required initialization/invariants                                          | `func NewClient(baseURL string) *Client`                 |

## 21. References

1. Go Team — _A Tour of Go: Structs_. https://go.dev/tour/moretypes/2
2. Go Team — _Effective Go_. https://go.dev/doc/effective_go
3. Go Language Specification — _Struct types_. https://go.dev/ref/spec#Struct_types
4. Kiran Adhikari — _Golang Interview: Struct_, Medium. https://medium.com/@kiruu1238/golang-interview-struct-758db20ef4c9
5. GoLinuxCloud — _Golang struct embedding — composition, promoted fields, methods, interfaces_. https://www.golinuxcloud.com/golang-embedding/
6. GoLinuxCloud — _Go anonymous struct — inline syntax, JSON, receivers, embedding_. https://www.golinuxcloud.com/golang-anonymous-structs/
7. Go FAQ (gofaq.org) — _Anonymous structs_. https://www.gofaq.org/en/anonymous-structs/
8. Go FAQ (gofaq.org) — _How to Use Struct Tags in Go (json, db, yaml, validate)_. https://www.gofaq.org/en/how-to-use-struct-tags-in-go-json-db-yaml-validate/
9. Eli Bendersky — _Embedding in Go, Part 1: structs in structs_. https://eli.thegreenplace.net/2020/embedding-in-go-part-1-structs-in-structs
10. Leapcell — _Understanding Struct Embedding in Go_. https://leapcell.io/blog/understanding-struct-embedding-in-go
11. LabEx — _How to access embedded struct fields_. https://labex.io/tutorials/go-how-to-access-embedded-struct-fields-437892
12. go101.org — _Type Embedding_. https://go101.org/article/type-embedding.html
13. Go standard library documentation — package `encoding/json` (struct tag behavior). https://pkg.go.dev/encoding/json
14. Go standard library documentation — package `reflect` (`StructTag`, `DeepEqual`). https://pkg.go.dev/reflect
15. Go standard library documentation — package `unsafe` (`Sizeof`, alignment). https://pkg.go.dev/unsafe
