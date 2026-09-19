<div align="center">
  <h1>Maps</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 17, 2026</sub>
</div>

## Table of Contents

1. [What Is a Map?](#1-what-is-a-map)
2. [Declaring and Creating Maps](#2-declaring-and-creating-maps)
3. [The Nil Map Trap](#3-the-nil-map-trap)
4. [Map Literals](#4-map-literals)
5. [Adding and Updating Entries](#5-adding-and-updating-entries)
6. [Reading Entries and the Comma-Ok Idiom](#6-reading-entries-and-the-comma-ok-idiom)
7. [Deleting Entries](#7-deleting-entries)
8. [Checking the Number of Entries](#8-checking-the-number-of-entries)
9. [Iterating Over a Map](#9-iterating-over-a-map)
10. [Iteration Order Is Randomized](#10-iteration-order-is-randomized)
11. [Valid Key Types: Comparability](#11-valid-key-types-comparability)
12. [Struct Values in Maps Are Not Addressable](#12-struct-values-in-maps-are-not-addressable)
13. [Maps Behave Like Reference Types](#13-maps-behave-like-reference-types)
14. [Maps Are Not Comparable (Except to `nil`)](#14-maps-are-not-comparable-except-to-nil)
15. [Maps Are Not Safe for Concurrent Use](#15-maps-are-not-safe-for-concurrent-use)
16. [`sync.Map`: A Concurrency-Safe Alternative](#16-syncmap-a-concurrency-safe-alternative)
17. [Using a Map as a Set](#17-using-a-map-as-a-set)
18. [Nested Maps](#18-nested-maps)
19. [Sizing Maps with `make`'s Capacity Hint](#19-sizing-maps-with-makes-capacity-hint)
20. [Clearing a Map](#20-clearing-a-map)
21. [The Standard Library's `maps` Package](#21-the-standard-librarys-maps-package)
22. [Performance Characteristics](#22-performance-characteristics)
23. [Common Mistakes and Pitfalls](#23-common-mistakes-and-pitfalls)
24. [Best Practices Summary](#24-best-practices-summary)
25. [Use Case Summary Table](#25-use-case-summary-table)
26. [References](#26-references)

## 1. What Is a Map?

A **map** is Go's built-in associative data structure — an unordered collection of key-value pairs, where each key is unique and maps to exactly one value. Internally, a map is implemented as a **hash table**, which is why lookups, insertions, and deletions are all extremely fast on average (constant time, O(1)), regardless of how many entries the map holds.

```go
var ages map[string]int // a map from string keys to int values
```

The type of a map is written `map[KeyType]ValueType`. Maps are one of Go's three built-in "reference-like" types, alongside slices and channels — internally, a map variable holds a pointer to the actual hash table data structure, which has important consequences covered in [Section 13](#13-maps-behave-like-reference-types).

## 2. Declaring and Creating Maps

There are two primary ways to obtain a usable map:

### 2.1 `make`

```go
ages := make(map[string]int) // an empty, ready-to-use map
ages["Alice"] = 30
```

`make` is the standard, idiomatic way to create a map you intend to populate — it allocates and initializes the underlying hash table so the map is immediately safe to write to.

### 2.2 Map Literal

```go
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
}
```

Covered in more detail in [Section 4](#4-map-literals).

### 2.3 Zero Value: `nil`

```go
var ages map[string]int // ages is nil — declared but not initialized
```

A `nil` map behaves differently from an empty map created with `make` in one crucial way, covered next.

## 3. The Nil Map Trap

This is the single most important beginner pitfall with Go maps: **a `nil` map can be read from safely, but writing to it panics.**

```go
var m map[string]int // nil map

v := m["anything"] // fine — reading from a nil map returns the value type's zero value
fmt.Println(v)       // 0

m["key"] = 1 // PANIC: assignment to entry in nil map
```

This asymmetry — safe to read, unsafe to write — trips up many newcomers, since a `nil` map often "looks fine" until the exact moment code tries to insert into it. Reading is safe because a lookup on an empty/nonexistent hash table can simply report "not found" (the zero value), but a write needs an actual, initialized hash table data structure to insert into, which a `nil` map does not have.

**Fix:** always initialize a map with `make` (or a literal) before writing to it:

```go
var m map[string]int
if m == nil {
    m = make(map[string]int)
}
m["key"] = 1 // now safe
```

Or, more simply, just initialize it with `make` at declaration time whenever you intend to write to it:

```go
m := make(map[string]int)
m["key"] = 1
```

**A common real-world source of this bug:** a struct field of map type that's zero-valued (`nil`) because the struct was created with a plain literal or `new()` without explicitly initializing that field — the map field silently starts out `nil` until something writes to it, exactly like any other zero-valued map variable.

```go
type Config struct {
    Settings map[string]string // nil until explicitly initialized
}

c := Config{} // c.Settings is nil
// c.Settings["key"] = "value" // PANIC
```

## 4. Map Literals

A map literal lets you declare and populate a map in one expression:

```go
capitals := map[string]string{
    "France": "Paris",
    "Japan":  "Tokyo",
    "Egypt":  "Cairo",
}
```

An empty map literal (`map[string]int{}`) is a valid, non-nil, immediately writable map — this is a subtle but important distinction from `var m map[string]int` (which is `nil`):

```go
m1 := map[string]int{}       // non-nil, empty, safe to write to immediately
var m2 map[string]int         // nil — writing to it panics until initialized
```

Both `m1` and `m2` report `len(m) == 0` and print the same way, but only `m1` can be written to without first calling `make` — this is a common source of confusion, since the two look nearly identical until a write is attempted.

## 5. Adding and Updating Entries

Both inserting a brand-new key and updating an existing key's value use the exact same assignment syntax:

```go
m := make(map[string]int)

m["apples"] = 5   // insert: "apples" doesn't exist yet, so this creates it
m["apples"] = 10  // update: "apples" already exists, so this overwrites the value
```

There is no separate "insert" vs. "update" operation in Go's map syntax — assignment handles both uniformly, depending on whether the key was already present.

## 6. Reading Entries and the Comma-Ok Idiom

### 6.1 The Basic Form

```go
value := m["apples"]
fmt.Println(value) // 10, or the zero value if "apples" isn't in the map
```

**The core ambiguity:** a plain lookup like this returns the value type's zero value both when the key genuinely maps to that zero value, and when the key doesn't exist in the map at all — there is no way to distinguish "key exists and its value is 0" from "key does not exist" using this form alone.

### 6.2 The Comma-Ok Idiom

To resolve this ambiguity, a map lookup can return a second, boolean result indicating whether the key was actually present:

```go
value, ok := m["bananas"]
if !ok {
    fmt.Println("bananas is not in the map")
} else {
    fmt.Println("bananas:", value)
}
```

`ok` is `true` if the key exists (regardless of what its value is, even if that value happens to be the zero value), and `false` if the key is absent (in which case `value` is set to the value type's zero value).

```go
m := map[string]int{"zero-value-key": 0}

v1, ok1 := m["zero-value-key"] // v1 = 0, ok1 = true  — key exists, value happens to be 0
v2, ok2 := m["missing-key"]     // v2 = 0, ok2 = false — key doesn't exist at all
```

**Use case:** any time the zero value is a meaningful, legitimate value your map might actually store (0 for numbers, `""` for strings, `false` for booleans), the comma-ok form is the only reliable way to distinguish "present with that value" from "absent entirely."

## 7. Deleting Entries

The built-in `delete` function removes a key (and its associated value) from a map:

```go
m := map[string]int{"a": 1, "b": 2}
delete(m, "a")
fmt.Println(m) // map[b:2]
```

`delete` is always safe to call, even if the key doesn't exist in the map — it simply does nothing in that case, without panicking or returning an error:

```go
delete(m, "nonexistent-key") // no-op, perfectly safe
```

`delete` is also explicitly safe to call **during** a `range` loop over the same map — deleting the current entry, or an entry not yet visited, is well-defined behavior (see [Section 9.2](#92-deleting-during-iteration)).

## 8. Checking the Number of Entries

The built-in `len` function returns the number of key-value pairs currently in a map:

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
fmt.Println(len(m)) // 3
```

`len` on a `nil` map returns `0` (this is one of the safe, read-only operations a `nil` map supports, alongside plain lookups).

## 9. Iterating Over a Map

`for range` is the idiomatic way to visit every key-value pair in a map:

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}

for key, value := range m {
    fmt.Println(key, value)
}

for key := range m { // keys only
    fmt.Println(key)
}

for _, value := range m { // values only
    fmt.Println(value)
}
```

### 9.1 Modifying the Map's Values During Iteration

Updating the value of an **existing** key you're currently iterating over is well-defined and safe:

```go
for k := range m {
    m[k] = m[k] * 2 // safe: updating an existing key's value during iteration
}
```

However, whether a **newly inserted** key (one that didn't exist when the `range` started) gets visited during that same iteration is explicitly left **unspecified** by the language — it might be visited, or it might not, and this can even vary between runs. Avoid inserting new keys into a map while iterating over it if you need predictable behavior; collect the keys/values to insert separately, and apply them after the loop finishes.

### 9.2 Deleting During Iteration

Deleting the current or a not-yet-visited entry during a `range` loop over the same map is explicitly documented as safe and well-defined:

```go
for k, v := range m {
    if v < 0 {
        delete(m, k) // safe, even for the entry currently being visited
    }
}
```

## 10. Iteration Order Is Randomized

**Go deliberately does not guarantee any particular order** when iterating over a map with `range` — and in practice, the Go runtime actively randomizes the starting point of each iteration specifically to prevent code from accidentally coming to rely on a particular order that isn't actually guaranteed:

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k, v := range m {
    fmt.Println(k, v)
}
// Running this multiple times can (and often does) print the entries in a different order each time.
```

**If you need a predictable, deterministic order** (for example, printing a map's contents in a stable, readable way, or producing reproducible output for tests), extract the keys into a slice, sort that slice, and then iterate using the sorted key order:

```go
import "sort"

keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)

for _, k := range keys {
    fmt.Println(k, m[k]) // now printed in a stable, sorted order
}
```

(Go 1.21+ can simplify the key-collection step using the standard library's `maps.Keys` function, discussed in [Section 21](#21-the-standard-librarys-maps-package).)

## 11. Valid Key Types: Comparability

A map's key type must be **comparable** — that is, it must support the `==`/`!=` operators. This includes basic types (numbers, strings, booleans), pointers, channels, interfaces (as long as their dynamic value is also comparable at runtime), arrays of comparable element types, and structs made entirely of comparable fields.

**Not valid as map keys:** slices, maps, and functions — none of these support `==` in Go, so attempting to declare a map with one of them as the key type is a compile-time error:

```go
// var m map[[]int]string // COMPILE ERROR: invalid map key type []int
```

### 11.1 Structs as Map Keys

A struct with only comparable fields makes an excellent map key when you need to key by a combination of values:

```go
type Point struct {
    X, Y int
}

visited := make(map[Point]bool)
visited[Point{1, 2}] = true

if visited[Point{1, 2}] {
    fmt.Println("already visited (1, 2)")
}
```

This is a common and idiomatic pattern for problems like grid/graph traversal, where the "key" is naturally a coordinate pair or similar composite value, rather than a single primitive.

## 12. Struct Values in Maps Are Not Addressable

A struct value **stored directly as a map's value type** (not behind a pointer) cannot be mutated in place through the map index expression, because map values are not addressable in Go:

```go
type Point struct{ X, Y int }

m := map[string]Point{"origin": {0, 0}}

// m["origin"].X = 5 // COMPILE ERROR: cannot assign to struct field m["origin"].X in map
```

**Fix 1 — read, modify, write back:**

```go
p := m["origin"]
p.X = 5
m["origin"] = p
```

**Fix 2 — store pointers to structs instead:**

```go
m2 := map[string]*Point{"origin": {0, 0}}
m2["origin"].X = 5 // works: dereferencing a pointer's field is addressable
```

Storing pointers also avoids copying the full struct on every read, which matters if the struct is large — though it does mean multiple parts of your program can now observe (and potentially race on) the same shared struct through different references, so this trade-off should be made deliberately.

## 13. Maps Behave Like Reference Types

Although Go's "everything is passed by value" rule holds universally (see the discussion of pointers for the full explanation), a map variable internally holds a **pointer to its underlying hash table**. Copying a map variable, or passing it to a function, copies this small internal descriptor — but the copy still points to the **same** underlying data:

```go
func addEntry(m map[string]int) {
    m["new"] = 1 // mutates the SAME underlying hash table the caller has
}

original := map[string]int{"a": 1}
addEntry(original)
fmt.Println(original) // map[a:1 new:1] — the caller sees the change
```

This is why you almost never need an explicit `*map[K]V` just to let a function mutate a map's existing contents — the map variable's internal pointer already provides that. An explicit pointer to a map is only needed in the rare case where a function must **reassign** the caller's map variable itself (e.g., assigning it to an entirely new map value, or setting it to `nil`), not merely add/modify/remove entries within it.

## 14. Maps Are Not Comparable (Except to `nil`)

Unlike some other composite types, two maps **cannot** be compared with `==` — the only comparison Go allows for a map is checking whether it equals `nil`:

```go
m1 := map[string]int{"a": 1}
m2 := map[string]int{"a": 1}

// fmt.Println(m1 == m2) // COMPILE ERROR: map can only be compared to nil

fmt.Println(m1 == nil) // false — valid comparison
```

If you need to check whether two maps hold the same keys and values, use `reflect.DeepEqual` (for a general, recursive structural comparison) or the standard library's `maps.Equal` function (Go 1.21+, discussed in [Section 21](#21-the-standard-librarys-maps-package)), which is both clearer in intent and typically faster than `reflect.DeepEqual` for this specific comparison.

```go
import "maps"

fmt.Println(maps.Equal(m1, m2)) // true
```

## 15. Maps Are Not Safe for Concurrent Use

**This is one of the most important practical facts about Go maps to internalize.** A plain Go map has **no built-in synchronization** — if multiple goroutines access the same map concurrently, and at least one of them is writing, the behavior is officially undefined, and in practice the Go runtime actively detects this at run time and crashes the program deliberately, rather than allowing it to silently corrupt data:

```go
m := make(map[int]int)

go func() {
    for i := 0; i < 1000; i++ {
        m[i] = i // write
    }
}()

go func() {
    for i := 0; i < 1000; i++ {
        _ = m[i] // read
    }
}()

// fatal error: concurrent map read and map write
```

This crash is **intentional runtime behavior**, not a bug — the Go team has explained that concurrent, unsynchronized access can corrupt a map's internal hash-table structure in ways that could otherwise produce silent, hard-to-diagnose data corruption; the runtime chooses to fail loudly and immediately instead. Critically, this specific fatal error **cannot be caught with `recover()`** — it's a runtime-level fatal error, not a normal Go panic, so it always terminates the entire program.

### 15.1 Fixing It with a Mutex

The standard fix is to protect the map with a `sync.Mutex` (or `sync.RWMutex` if reads vastly outnumber writes):

```go
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]string
}

func (s *SafeMap) Get(key string) (string, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    v, ok := s.m[key]
    return v, ok
}

func (s *SafeMap) Set(key, value string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.m[key] = value
}
```

This is generally the preferred approach for most concurrent-map use cases in application code, since it's explicit, simple to reason about, and works well when the map is accessed as a whole unit with clear read/write boundaries.

## 16. `sync.Map`: A Concurrency-Safe Alternative

The standard library also provides `sync.Map`, a map implementation with built-in internal synchronization, safe for concurrent use without any external locking:

```go
var m sync.Map

m.Store("key", 42)

value, ok := m.Load("key")
if ok {
    fmt.Println(value) // 42
}

m.Delete("key")

m.Range(func(key, value any) bool {
    fmt.Println(key, value)
    return true // return false to stop iterating early
})
```

`sync.Map` uses `any` for its keys and values (it's not generic), which sacrifices some compile-time type safety compared to a regular typed map guarded by your own mutex.

**When `sync.Map` is actually a good fit:** the standard library documentation notes it's optimized for two specific access patterns — (1) when a given key is written once but read many times, with the set of keys growing mostly monotonically (a cache that only grows), or (2) when many goroutines each work with disjoint sets of keys, with little contention between them. For the more general case of a map with frequent, interleaved reads and writes across arbitrary keys from many goroutines, a plain map guarded by a `sync.Mutex`/`sync.RWMutex` is often simpler to reason about and can perform just as well or better — `sync.Map` is a specialized tool, not a universal drop-in replacement for "a map used by multiple goroutines."

## 17. Using a Map as a Set

Go has no built-in "set" type, but a `map[T]bool` or, more idiomatically, `map[T]struct{}` is the standard way to represent one — since a set only cares about key _presence_, not any associated value:

```go
seen := make(map[string]struct{})

seen["alice"] = struct{}{}
seen["bob"] = struct{}{}

if _, ok := seen["alice"]; ok {
    fmt.Println("alice has been seen")
}

delete(seen, "bob")
```

Using `struct{}` (which occupies zero bytes) instead of `bool` communicates intent more precisely — there's no meaningful "false" state to a value being present in a set, so a `bool` would misleadingly suggest there is one. `map[T]bool` is also seen in real code and works fine functionally, but `map[T]struct{}` is generally considered the more idiomatic and slightly more memory-efficient choice for genuine set semantics.

## 18. Nested Maps

A map's value type can itself be another map, allowing multi-level, nested lookups:

```go
inventory := make(map[string]map[string]int) // warehouse -> (item -> quantity)

inventory["warehouse-A"] = make(map[string]int) // inner maps must be initialized individually
inventory["warehouse-A"]["widgets"] = 100

fmt.Println(inventory["warehouse-A"]["widgets"]) // 100
```

**Important:** creating the outer map does **not** automatically initialize any inner maps — each inner map is its own independent `nil` map until explicitly created with `make` (or a literal), exactly following the nil-map rule from [Section 3](#3-the-nil-map-trap). Forgetting to initialize an inner map before writing to it is a common source of the "assignment to entry in nil map" panic in nested-map code:

```go
// inventory["warehouse-B"]["gadgets"] = 5
// PANIC: inventory["warehouse-B"] is nil (never initialized), so writing into it panics
```

A safe pattern for "get or create" on a nested map:

```go
func addItem(inventory map[string]map[string]int, warehouse, item string, qty int) {
    if inventory[warehouse] == nil {
        inventory[warehouse] = make(map[string]int)
    }
    inventory[warehouse][item] += qty
}
```

## 19. Sizing Maps with `make`'s Capacity Hint

`make(map[K]V, n)` accepts an optional second argument — a **hint** for how many entries the map is expected to eventually hold:

```go
m := make(map[string]int, 1000) // hints that ~1000 entries are expected
```

This is purely a performance optimization: it lets the runtime pre-allocate an appropriately sized internal hash table up front, reducing the number of incremental resize/rehash operations that would otherwise happen as the map grows organically one insertion at a time. It does **not** limit the map's actual size — you can insert more or fewer entries than the hint, and the map will grow (or simply have some unused extra capacity) as needed regardless. Providing an accurate hint is a worthwhile, low-effort optimization whenever you know roughly how many entries a map will end up holding in advance (for example, when building a map from a slice or query result of a known length).

## 20. Clearing a Map

Go 1.21 introduced the built-in `clear` function, which removes all entries from a map (or zeroes a slice) in one call, working for any map type:

```go
m := map[string]int{"a": 1, "b": 2}
clear(m)
fmt.Println(len(m)) // 0
```

`clear(m)` is equivalent to (but generally more efficient than) manually deleting every key in a loop, and — unlike reassigning `m = make(map[string]int)` — it reuses the map's existing underlying storage rather than allocating an entirely new hash table, which can matter for maps that are cleared and refilled repeatedly in a hot loop.

## 21. The Standard Library's `maps` Package

Go 1.21 also introduced the `maps` package, providing common generic operations over maps without needing a hand-written loop for each one:

```go
import "maps"

m := map[string]int{"a": 1, "b": 2}

m2 := maps.Clone(m)          // a shallow copy of the map
equal := maps.Equal(m, m2)    // structural equality check
maps.DeleteFunc(m, func(k string, v int) bool {
    return v < 2 // deletes all entries whose value is less than 2
})

for k := range maps.Keys(m) { // Keys returns an iterator (Go 1.23+ range-over-func); see note below
    fmt.Println(k)
}
```

**Note on `maps.Keys`/`maps.Values`:** in Go 1.21 and 1.22, these returned an `iter.Seq`-like sequence usable with the newer range-over-function support; from Go 1.23 onward, `for range` directly over these iterator-returning functions became fully supported language syntax. Check the documentation for the exact Go version you're targeting if you rely on these specific functions, since their exact signatures evolved slightly as Go's range-over-func iterator support matured.

**Practical takeaway:** before writing a manual loop to clone, compare, filter, or extract keys/values from a map, check whether the `maps` package (Go 1.21+) already provides it — these implementations are maintained as part of the standard library and are generally well-optimized.

## 22. Performance Characteristics

- **Average-case O(1)** for lookups, insertions, and deletions, thanks to the underlying hash table implementation — this holds regardless of how many entries the map contains, as long as the hash function distributes keys reasonably evenly.
- **Growth involves occasional rehashing.** As a map grows past its current internal capacity, Go's runtime allocates a larger backing structure and redistributes existing entries — an operation with higher cost than a typical insertion, but amortized across many insertions so the _average_ cost per insertion remains low.
- **Providing a capacity hint via `make(map[K]V, n)`** (see [Section 19](#19-sizing-maps-with-makes-capacity-hint)) avoids some of these incremental resizes when the eventual size is known in advance.
- **Struct keys have a cost proportional to their size**, since hashing and comparing a struct key involves every one of its fields — very large struct keys can meaningfully slow down map operations compared to a simple primitive key like a string or int.
- **Maps generally have more memory overhead per entry** than an equivalently sized slice or array, due to the hash table's internal bucket structure — for very performance- or memory-sensitive code with a small, fixed, known set of keys, a slice with linear search (or a specialized data structure) can sometimes outperform a map, though this is a micro-optimization relevant mainly to hot paths with small N.

## 23. Common Mistakes and Pitfalls

### 23.1 Writing to a Nil Map

```go
var m map[string]int
m["key"] = 1 // panic: assignment to entry in nil map
```

Always initialize with `make` (or a literal) before writing — see [Section 3](#3-the-nil-map-trap).

### 23.2 Confusing "Zero Value" with "Key Absent"

```go
count := m["missing-key"] // returns 0 whether "missing-key" exists with value 0, or doesn't exist at all
```

Use the comma-ok form (`value, ok := m[key]`) whenever this ambiguity matters — see [Section 6](#6-reading-entries-and-the-comma-ok-idiom).

### 23.3 Relying on Iteration Order

```go
for k, v := range m {
    // assuming this always visits keys in insertion order, or alphabetical order, etc. — WRONG
}
```

Map iteration order is intentionally randomized. Sort keys explicitly (see [Section 10](#10-iteration-order-is-randomized)) whenever deterministic output matters.

### 23.4 Concurrent Access Without Synchronization

```go
// Multiple goroutines reading/writing the same plain map with no lock
// → fatal error: concurrent map read and map write (crashes the whole program, not recoverable)
```

Protect shared maps with a `sync.Mutex`/`sync.RWMutex`, or use `sync.Map` for its specific supported access patterns — see [Sections 15](#15-maps-are-not-safe-for-concurrent-use) and [16](#16-syncmap-a-concurrency-safe-alternative).

### 23.5 Trying to Mutate a Struct Field Through a Map Index

```go
type Point struct{ X, Y int }
m := map[string]Point{"p": {1, 2}}
// m["p"].X = 5 // COMPILE ERROR
```

Read, modify, and write back the whole struct, or store pointers to structs instead — see [Section 12](#12-struct-values-in-maps-are-not-addressable).

### 23.6 Forgetting to Initialize Inner Maps in a Nested Map

```go
outer := make(map[string]map[string]int)
// outer["a"]["b"] = 1 // PANIC: outer["a"] is nil until explicitly initialized
```

Always check for and initialize the inner map before writing into it — see [Section 18](#18-nested-maps).

### 23.7 Comparing Maps with `==`

```go
// if m1 == m2 { ... } // COMPILE ERROR: maps can only be compared to nil
```

Use `reflect.DeepEqual` or `maps.Equal` (Go 1.21+) instead — see [Section 14](#14-maps-are-not-comparable-except-to-nil).

### 23.8 Assuming a Newly Inserted Key Will (or Won't) Be Visited Mid-Iteration

```go
for k := range m {
    if someCondition {
        m["new-key"] = 1 // whether this gets visited later in THIS SAME loop is unspecified
    }
}
```

Avoid inserting new keys during iteration if you need predictable results — collect changes separately and apply them after the loop, as noted in [Section 9.1](#91-modifying-the-maps-values-during-iteration).

## 24. Best Practices Summary

1. **Always initialize a map with `make` or a literal before writing to it** — a `nil` map panics on write, even though it can be read safely.
2. **Use the comma-ok idiom whenever the value type's zero value is a meaningful possible value** — it's the only reliable way to distinguish "present with zero value" from "absent."
3. **Never assume any particular iteration order** — sort keys explicitly whenever deterministic output is required.
4. **Protect maps shared across goroutines with a `sync.Mutex`/`sync.RWMutex`**, or use `sync.Map` specifically for its documented supported access patterns (mostly-write-once keys, or largely disjoint per-goroutine key sets).
5. **Use `map[T]struct{}` for set semantics**, communicating "presence only" more precisely than `map[T]bool`.
6. **Provide a capacity hint (`make(map[K]V, n)`) when the eventual size is roughly known**, to avoid unnecessary incremental resizing.
7. **Store pointers to structs (not structs directly) as map values** when you need to mutate individual struct fields in place through the map.
8. **Always explicitly initialize inner maps in a nested map structure** before writing to them — outer initialization does not cascade to inner maps.
9. **Use `reflect.DeepEqual` or the standard library's `maps.Equal`** to compare two maps' contents, since `==` is not defined for maps beyond comparing to `nil`.
10. **Check the standard library's `maps` package (Go 1.21+) and the built-in `clear` function** before hand-writing common map operations like cloning, filtering, or bulk-clearing.

## 25. Use Case Summary Table

| Technique                        | When to Use                                                           | Example Scenario                                                              |
| -------------------------------- | --------------------------------------------------------------------- | ----------------------------------------------------------------------------- |
| `make(map[K]V)`                  | Any map you intend to write to                                        | General-purpose key-value storage                                             |
| Map literal                      | A map whose initial contents are known upfront                        | A lookup table of constants (`map[string]string{"US": "United States", ...}`) |
| Comma-ok idiom (`v, ok := m[k]`) | Distinguishing "present with zero value" from "absent"                | Counting occurrences where 0 is a valid legitimate count                      |
| `delete(m, k)`                   | Removing an entry, including safely during iteration                  | Expiring a cache entry, removing a processed item from a work queue           |
| Sorted key iteration             | Deterministic, reproducible output                                    | Printing a config map's contents in a stable order for logs/tests             |
| Struct as map key                | Composite/multi-field lookup keys                                     | Grid coordinates, composite (host, port) keys                                 |
| Pointer values in a map          | Need to mutate struct fields in place through the map                 | A registry of mutable objects keyed by ID                                     |
| `sync.Mutex` + map               | General-purpose concurrent map access with interleaved reads/writes   | A shared, frequently updated cache accessed by many goroutines                |
| `sync.Map`                       | Write-once/read-many keys, or largely disjoint per-goroutine key sets | A cache of computed results that only grows, rarely overwritten               |
| `map[T]struct{}`                 | Set semantics — presence-only membership testing                      | Deduplicating a list, tracking "already visited" nodes                        |
| Nested maps                      | Multi-level lookups                                                   | Inventory by warehouse then item, permissions by role then resource           |
| `make(map[K]V, n)` capacity hint | The approximate final size is known in advance                        | Building a map from a slice or query result of known length                   |
| `clear(m)` (Go 1.21+)            | Reusing a map's storage while removing all entries                    | Resetting a per-iteration scratch map in a hot loop                           |
| `maps` package (Go 1.21+)        | Cloning, comparing, or filtering maps                                 | `maps.Clone`, `maps.Equal`, `maps.DeleteFunc` instead of hand-written loops   |

## 26. References

1. Go Team — _A Tour of Go: Maps_. https://go.dev/tour/moretypes/19
2. Go Team — _Go maps in action_, The Go Blog. https://go.dev/blog/maps
3. Go Language Specification — _Map types_. https://go.dev/ref/spec#Map_types
4. Go standard library documentation — package `sync` (`Map`, `Mutex`, `RWMutex`). https://pkg.go.dev/sync
5. Go standard library documentation — package `maps`. https://pkg.go.dev/maps
6. Go standard library documentation — package `reflect` (`DeepEqual`). https://pkg.go.dev/reflect
7. Go Team — _Go 1.21 Release Notes_ (`maps` package, built-in `clear`). https://go.dev/doc/go1.21
8. frihk_ian — _Maps in Golang_, DEV Community. https://dev.to/frihk_ian/maps-in-golang-4jif
9. IdeaToCode — _Fatal Error: Concurrent Map Read and Write in Go: Understanding and Solving the Panic_, Medium. https://medium.com/@ideatocode.tech/fatal-error-concurrent-map-read-and-write-in-go-understanding-and-solving-the-panic-95e8abe88a26
10. Luan Figueredo — _Concurrent map access in Go_, Medium. https://medium.com/@luanrubensf/concurrent-map-access-in-go-a6a733c5ffd1
11. Divesh Kumar Chordia — _Go Maps Aren't Concurrency-Safe — Here's What That Really Means_, Medium. https://medium.com/@diveshkumarchordia/go-maps-and-concurrency-why-you-need-synchronization-and-how-to-get-it-right-6cf5cae2f26b
12. Vincent — _Go: Concurrency Access with Maps — Part III_, Medium (_A Journey With Go_). https://medium.com/a-journey-with-go/go-concurrency-access-with-maps-part-iii-8c0a0e4eb27e
13. golang/go Issue #41674 — discussion on the intentional design of the concurrent map access fatal error. https://github.com/golang/go/issues/41674
14. golang.design — _Go: Under the Hood, Chapter 11.7: Concurrency-Safe Hash Table (sync.Map)_. https://golang.design/under-the-hood/en/part3concurrency/ch11sync/map/
