<div align="center">
  <h1>Conditional Statements</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 18, 2026</sub>
</div>

## Table of Contents

1. [Overview: Go's Branching Constructs](#1-overview-gos-branching-constructs)
2. [The Basic `if` Statement](#2-the-basic-if-statement)
3. [No Parentheses, Mandatory Braces](#3-no-parentheses-mandatory-braces)
4. [`else` and `else if`](#4-else-and-else-if)
5. [The `if` Statement's Init Clause](#5-the-if-statements-init-clause)
6. [The Common Idiom: `if err != nil`](#6-the-common-idiom-if-err--nil)
7. [Early Returns to Reduce Nesting](#7-early-returns-to-reduce-nesting)
8. [No Ternary Operator](#8-no-ternary-operator)
9. [The `switch` Statement: Overview](#9-the-switch-statement-overview)
10. [Expression Switch: Basic Form](#10-expression-switch-basic-form)
11. [Implicit Break: No Fallthrough by Default](#11-implicit-break-no-fallthrough-by-default)
12. [Multiple Values in a Single `case`](#12-multiple-values-in-a-single-case)
13. [The `switch` Statement's Init Clause](#13-the-switch-statements-init-clause)
14. [Expressionless Switch (`switch true`)](#14-expressionless-switch-switch-true)
15. [The `fallthrough` Keyword](#15-the-fallthrough-keyword)
16. [`break` Inside a `switch`](#16-break-inside-a-switch)
17. [`default` Can Appear Anywhere](#17-default-can-appear-anywhere)
18. [Case Expressions Don't Need to Be Constants](#18-case-expressions-dont-need-to-be-constants)
19. [Type Switch](#19-type-switch)
20. [`switch` vs. `if`/`else`: When to Use Which](#20-switch-vs-ifelse-when-to-use-which)
21. [Common Mistakes and Pitfalls](#21-common-mistakes-and-pitfalls)
22. [Best Practices Summary](#22-best-practices-summary)
23. [Use Case Summary Table](#23-use-case-summary-table)
24. [References](#24-references)

## 1. Overview: Go's Branching Constructs

Go has two primary constructs for conditional branching: **`if`/`else`** and **`switch`**. (A third construct, `select`, handles branching specifically over channel operations and is a distinct topic from general-purpose conditionals.) Both `if` and `switch` share a few distinctive Go conventions worth internalizing up front: no parentheses around the condition, mandatory braces around every branch's body, and an optional short initialization statement that scopes a variable to just that conditional.

## 2. The Basic `if` Statement

```go
age := 20

if age >= 18 {
    fmt.Println("You are an adult")
}
```

The condition must evaluate to a `bool` — unlike C, Python (pre-3), or JavaScript, Go does **not** allow arbitrary non-boolean values (like an integer or a pointer) to be implicitly treated as "truthy" or "falsy" in a condition. A condition must be an actual boolean expression:

```go
// if age { ... } // COMPILE ERROR: age (type int) is not a boolean value
if age != 0 { ... } // must be explicit
```

## 3. No Parentheses, Mandatory Braces

Two syntax rules distinguish Go's `if` from many C-family languages:

- **No parentheses** are required (or idiomatically used) around the condition.
- **Braces `{ }` are always mandatory**, even for a single-statement body — unlike C or Java, where braces can be omitted for a one-line body.

```go
if age >= 18 {       // no parentheses around the condition
    fmt.Println("adult")
}

// if age >= 18 fmt.Println("adult") // COMPILE ERROR: missing braces
```

The opening brace must also appear on the **same line** as the `if` (or `else`) keyword — Go's automatic semicolon insertion makes a brace on its own line a syntax error in this position, so `gofmt` always formats it this way.

## 4. `else` and `else if`

```go
score := 75

if score >= 90 {
    fmt.Println("A")
} else if score >= 80 {
    fmt.Println("B")
} else if score >= 70 {
    fmt.Println("C")
} else {
    fmt.Println("F")
}
```

Each condition is checked in order, top to bottom, and the first one that evaluates to `true` has its block executed — the rest are skipped entirely. If none of the conditions match and an `else` is present, its block runs as the fallback.

The `else` keyword (and any `else if`) must appear on the **same line** as the preceding block's closing brace, not on its own line — this, too, is a consequence of Go's semicolon-insertion rules and is what `gofmt` always enforces:

```go
if x > 0 {
    fmt.Println("positive")
} else { // must be on the same line as the closing brace above
    fmt.Println("not positive")
}
```

## 5. The `if` Statement's Init Clause

An `if` statement can include a short **init statement** before the condition, separated by a semicolon — any variable declared there is scoped **only** to the `if`/`else if`/`else` chain it belongs to, not to the surrounding code:

```go
if err := doSomething(); err != nil {
    fmt.Println("error:", err)
}
// err is not accessible here — its scope ended with the if statement
```

This is an extremely common and idiomatic Go pattern, especially for functions returning `(value, error)` or `(value, ok)` pairs, since it lets the "check the result" logic live in the same line as the call that produced it, without polluting the surrounding scope with a variable that's only relevant to this one check:

```go
if value, ok := myMap["key"]; ok {
    fmt.Println("found:", value)
} else {
    fmt.Println("not found")
}
```

Note that a variable declared in the init clause **is** visible across the entire `if`/`else if`/`else` chain that follows it (including any `else` branches), even though it's invisible outside that chain entirely.

## 6. The Common Idiom: `if err != nil`

Because Go functions conventionally return an `error` as their final return value, and Go has no exceptions, the pattern of checking that error immediately after a call — very often combined with the init-clause form from the previous section — appears constantly throughout idiomatic Go code:

```go
data, err := os.ReadFile("config.json")
if err != nil {
    log.Fatal(err)
}
// use data here, knowing err was nil
```

This repetitive-looking pattern is a deliberate, core part of Go's design philosophy: making the possibility of failure, and the point where it's handled, fully explicit and visible in the code, rather than implicit or hidden behind an exception mechanism that could be silently skipped or caught far away from where the failure actually occurred.

## 7. Early Returns to Reduce Nesting

A widely followed Go style convention is to handle error/failure conditions with an early `return` (or `continue`/`break`, in a loop) at the top of a function, keeping the function's "normal," successful path unindented and easy to follow, rather than nesting the success path inside an `else` block:

```go
// Less idiomatic: success path nested inside an else
func process(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    } else {
        // ... a lot of logic, indented one level deeper than necessary ...
        return nil
    }
}

// More idiomatic: early return keeps the success path unindented
func process(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }
    // ... the same logic, now at the function's base indentation level ...
    return nil
}
```

This pattern — sometimes called a "guard clause" style — scales especially well when a function has several sequential steps that can each fail, since each check stays flat rather than progressively nesting deeper with every additional validation.

## 8. No Ternary Operator

Go deliberately does **not** provide a ternary conditional operator (`condition ? a : b`), unlike C, Java, JavaScript, and many other languages. This was an intentional design choice — the Go team's stated position is that the ternary operator is frequently abused to cram complex, hard-to-read logic onto a single line, and that a full `if`/`else` is clearer in the overwhelming majority of cases, even for simple choices.

```go
// No such syntax in Go:
// max := a > b ? a : b

// Idiomatic Go equivalent:
max := a
if b > a {
    max = b
}
```

For very simple cases used often enough to justify a helper, a small function can substitute:

```go
func maxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}

max := maxInt(a, b)
```

(Since Go 1.21, the standard library's built-in `max`/`min` functions cover this exact numeric case directly, without needing a hand-written helper at all — but the general pattern above applies to any type or condition a ternary might otherwise have expressed.)

## 9. The `switch` Statement: Overview

Go's `switch` provides multi-way branching, comparing an expression (or a type, in a type switch) against a list of `case` clauses. It comes in two main forms:

1. **Expression switch** — compares a value against a list of case values.
2. **Type switch** — compares the dynamic type stored in an interface value against a list of case types (covered in [Section 19](#19-type-switch)).

Go's `switch` differs from `switch` in C, Java, or PHP in one especially important way covered in [Section 11](#11-implicit-break-no-fallthrough-by-default): cases do **not** fall through to the next case by default.

## 10. Expression Switch: Basic Form

```go
day := 3

switch day {
case 1:
    fmt.Println("Monday")
case 2:
    fmt.Println("Tuesday")
case 3:
    fmt.Println("Wednesday")
default:
    fmt.Println("Some other day")
}
// Output: Wednesday
```

The switch expression (`day`) is evaluated once, then compared against each `case` value in order, top to bottom; the first matching case's block runs, and the rest are skipped. If nothing matches and a `default` clause is present, it runs as the fallback (a `default` is optional — if there is no match and no `default`, execution simply falls through the entire switch without running anything).

## 11. Implicit Break: No Fallthrough by Default

This is the single most important behavioral difference from C-family `switch` statements: **each Go `case` automatically stops (breaks) after its block runs** — execution never falls through into the next case unless `fallthrough` is used explicitly.

```go
switch c {
case '&':
    esc = "&amp;"
case '\'':
    esc = "&#39;"
case '<':
    esc = "&lt;"
default:
    panic("unrecognized escape character")
}
```

No explicit `break` statements are needed at the end of each case (unlike C/Java, where forgetting one is a classic bug) — Go's `switch` was deliberately designed with the far more commonly desired behavior (stop after the matching case) as the default, precisely to avoid that entire class of accidental-fallthrough bugs.

## 12. Multiple Values in a Single `case`

A single `case` can list several comma-separated values — the case matches if the switch expression equals **any one** of them:

```go
func isVowel(c rune) bool {
    switch c {
    case 'a', 'e', 'i', 'o', 'u':
        return true
    }
    return false
}
```

This is a common, concise way to group several distinct values that should all lead to the same branch, without needing several separate `case` clauses (which, without `fallthrough`, wouldn't share a body anyway) or an equivalent chain of `||` conditions in an `if`.

## 13. The `switch` Statement's Init Clause

Like `if`, `switch` supports an optional init statement before the (optional) switch expression, separated by a semicolon, scoping any declared variable to the entire `switch` block:

```go
switch hour := time.Now().Hour(); {
case hour < 12:
    fmt.Println("Good morning!")
case hour < 17:
    fmt.Println("Good afternoon!")
default:
    fmt.Println("Good evening!")
}
```

Note the syntax here: `hour := time.Now().Hour();` is the init statement, and the switch expression itself is left empty (covered next) — this combination is common for the "cleaner if/else chain" style of `switch`.

## 14. Expressionless Switch (`switch true`)

When a `switch` has **no expression** at all, it's treated as equivalent to `switch true` — each `case` is then a full boolean expression, and the first one that evaluates to `true` matches. This gives a `switch` the ability to replace a chain of `if`/`else if`/`else` conditions with what's often cleaner, more readable syntax:

```go
func unhex(c byte) byte {
    switch {
    case '0' <= c && c <= '9':
        return c - '0'
    case 'a' <= c && c <= 'f':
        return c - 'a' + 10
    case 'A' <= c && c <= 'F':
        return c - 'A' + 10
    }
    return 0
}
```

This is functionally identical to writing the same conditions as an `if`/`else if` chain, but many Go developers find it reads more clearly once there are more than two or three conditions to check, since each case's condition lines up visually rather than nesting.

## 15. The `fallthrough` Keyword

To deliberately force execution to continue into the **next** case's block — the opposite of Go's default behavior — use the `fallthrough` keyword as the last statement in a case:

```go
switch num := 2; num {
case 1:
    fmt.Println("One")
    fallthrough
case 2:
    fmt.Println("Two")
    fallthrough
case 3:
    fmt.Println("Three")
default:
    fmt.Println("Other")
}
// Output:
// Two
// Three
// Other
```

Even though `num` is `2`, execution starts at `case 2` (the actual match), prints "Two", then `fallthrough` unconditionally jumps into `case 3`'s body next — **without** re-evaluating `case 3`'s own condition/value at all. `fallthrough` is an unconditional jump to the very next case's body, not a "keep checking further cases" mechanism.

**Rules governing `fallthrough`:**

- It must be the **last statement** in its case's block — nothing can follow it within that case.
- It **cannot be used in the final case** of a switch (there being no subsequent case to fall through into) — this is a compile-time error.
- It works the same way in both expression switches and (with some restrictions) type switches, though its use in type switches is considerably rarer.

`fallthrough` should be used sparingly — because it deliberately breaks from Go's normal (and generally clearer) "stop after the matching case" behavior, overusing it can make a switch statement's control flow noticeably harder to trace at a glance.

## 16. `break` Inside a `switch`

Since each case already stops automatically, an explicit `break` inside a `switch` case might look redundant — but it remains useful for exiting a case **early**, before reaching the end of that case's block:

```go
switch command {
case "echo":
    fmt.Println(args)
case "cat":
    if len(args) == 0 {
        fmt.Println("Usage: cat <filename>")
        break // exit this case early, skipping the rest of its logic
    }
    printFile(args[0])
default:
    fmt.Println("Unknown command")
}
```

**Important distinction inside a loop:** a `break` inside a `switch` that is itself nested inside a `for` loop only exits the `switch`, **not** the enclosing loop — if you need to break out of the outer loop from within a nested switch case, you need a labeled `break` referencing the loop's label (covered under loops), since an unlabeled `break` always targets only its innermost enclosing `switch`/`select`/loop.

## 17. `default` Can Appear Anywhere

Unlike some languages that require `default` to be the last clause, Go allows `default` to appear **anywhere** within a `switch`'s case list — it still only runs when no other case matches, regardless of its position in the source:

```go
switch x {
default:
    fmt.Println("no match")
case 1:
    fmt.Println("one")
case 2:
    fmt.Println("two")
}
```

By strong convention (and for readability), `default` is almost always written **last**, even though the language permits placing it elsewhere — deviating from this convention without a specific reason can make a switch statement harder for readers to scan.

## 18. Case Expressions Don't Need to Be Constants

Unlike C or Java, where `switch` case labels must be compile-time constants, Go's `case` expressions can be **arbitrary expressions**, evaluated at runtime, including function calls:

```go
func Foo(n int) int {
    fmt.Println(n)
    return n
}

switch Foo(2) {
case Foo(1), Foo(2), Foo(3):
    fmt.Println("matched")
}
```

Case expressions in a single `case` clause with multiple comma-separated values are evaluated left to right, and evaluation stops as soon as a match is found (short-circuiting any remaining, unevaluated expressions in that same case) — but be aware that any expression with side effects (like the `Foo` calls above, which print as a side effect) will actually execute during the switch's evaluation, which is worth keeping in mind if case expressions aren't simple, pure values.

## 19. Type Switch

A **type switch** is a special form of `switch` that compares the **dynamic type** stored in an interface value against a list of case types, rather than comparing a value:

```go
func describe(i any) {
    switch v := i.(type) {
    case int:
        fmt.Printf("int: %d\n", v)
    case string:
        fmt.Printf("string: %q\n", v)
    case bool:
        fmt.Printf("bool: %t\n", v)
    case nil:
        fmt.Println("nil value")
    default:
        fmt.Printf("unknown type: %T\n", v)
    }
}
```

The syntax `switch v := i.(type)` is specific to type switches — `i.(type)` is only valid inside this exact construct, nowhere else in the language. Inside each case, the variable `v` automatically takes on the concrete type of that case (an `int` inside `case int:`, a `string` inside `case string:`, and so on), which the compiler enforces, so no manual type assertion is needed within the branch.

### 19.1 Multiple Types in One Case

A type switch case can also list several types together, separated by commas — in that situation, however, the case variable's type inside the block remains the original interface type (`any`, or whatever the switch expression's static type was), since the compiler can't narrow it down to one specific type when several are listed together:

```go
switch v := i.(type) {
case int, int64, float64:
    fmt.Printf("some kind of number: %v\n", v) // v's static type here is still `any`
default:
    fmt.Printf("something else: %v\n", v)
}
```

### 19.2 The `case nil` Branch

As shown in the first example, `case nil` specifically matches when the interface value itself is `nil` (no concrete type at all) — this is distinct from matching a value whose concrete type happens to be some nil-capable type like a nil pointer, which the discussion of interfaces elsewhere covers as the "typed nil" subtlety.

**Use case:** handling several possible concrete types flowing through a generic (`any`-typed) pipeline cleanly — for example, dispatching different behavior for different AST node types in a parser, encoding arbitrary values into a custom format, or handling several specific, expected concrete error types when `errors.As` isn't a natural fit for multi-branch dispatch.

## 20. `switch` vs. `if`/`else`: When to Use Which

Both constructs can often express the same logic, and the choice is largely a matter of readability:

| Prefer `if`/`else` when...                                                                              | Prefer `switch` when...                                                                    |
| ------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| There are only one or two conditions to check                                                           | There are three or more mutually exclusive conditions/values to compare                    |
| The conditions are complex, compound boolean expressions that don't naturally line up as parallel cases | The conditions naturally compare one single value/expression against several possibilities |
| You need to branch on a genuinely unrelated set of conditions, not a single dimension                   | You're branching on a type (type switch is the _only_ way to do this cleanly)              |

```go
// Reasonable as if/else — few, non-parallel conditions
if user == nil {
    return errors.New("no user")
} else if !user.Active {
    return errors.New("user inactive")
}

// Cleaner as a switch — many parallel comparisons against one value
switch status {
case "pending":
    // ...
case "active":
    // ...
case "suspended":
    // ...
case "closed":
    // ...
default:
    // ...
}
```

## 21. Common Mistakes and Pitfalls

### 21.1 Expecting C-Style Fallthrough by Default

```go
switch n {
case 1:
    fmt.Println("one")
case 2: // WRONG ASSUMPTION: expecting this to also run if n == 1, like in C
    fmt.Println("two")
}
```

Go's `switch` stops after the matching case by default — use `fallthrough` explicitly if you genuinely want the next case's block to run too, per [Section 15](#15-the-fallthrough-keyword).

### 21.2 Putting a Statement After `fallthrough`

```go
switch n {
case 1:
    fallthrough
    fmt.Println("unreachable") // COMPILE ERROR: fallthrough must be the last statement in the case
}
```

`fallthrough` must be the final statement in its case block — nothing else can follow it.

### 21.3 Using `fallthrough` in the Last Case

```go
switch n {
case 1:
    // ...
case 2:
    fallthrough // COMPILE ERROR: cannot fallthrough final case in switch
}
```

There's no subsequent case to fall through to from the last one — this is a compile-time error.

### 21.4 Assuming a Multi-Type `case` Narrows the Variable's Type

```go
switch v := i.(type) {
case int, string:
    fmt.Println(v + 1) // COMPILE ERROR: v's type here is still `any`, not int or string specifically
}
```

When a type-switch case lists multiple types, the case variable keeps the switch expression's original (usually interface) type inside that block — see [Section 19.1](#191-multiple-types-in-one-case).

### 21.5 Forgetting `else` Must Be on the Same Line as the Preceding `}`

```go
if x > 0 {
    fmt.Println("positive")
}
else { // COMPILE ERROR: syntax error — else must be on the same line as the closing brace
    fmt.Println("not positive")
}
```

Go's semicolon-insertion rules require `else`/`else if` to immediately follow the prior block's closing brace on the same line — `gofmt` always formats it this way, so this is rarely hand-written incorrectly when using standard tooling, but worth knowing if you encounter it.

### 21.6 Treating a Non-Boolean Value as a Condition

```go
var count int = 5
// if count { ... } // COMPILE ERROR: count (type int) is not a boolean
```

Unlike some languages, Go requires an actual `bool` expression in an `if`/`switch`-without-expression condition — there's no implicit truthiness for numbers, strings, pointers, or slices.

### 21.7 Relying on an Init-Clause Variable Outside Its Scope

```go
if err := doSomething(); err != nil {
    // handle err
}
fmt.Println(err) // COMPILE ERROR: err is not defined here — its scope ended with the if statement
```

A variable declared in an `if`/`switch` init clause is scoped strictly to that conditional statement (including any `else`/`case` branches) — it does not leak into the surrounding code.

## 22. Best Practices Summary

1. **Prefer early returns (guard clauses) over nested `else` blocks** for handling error/failure conditions, keeping a function's main logic at a consistent, shallow indentation level.
2. **Use the `if`/`switch` init-clause form** for variables that are only relevant to that one conditional check, to keep their scope as narrow as possible.
3. **Reach for `switch` once you have three or more mutually exclusive conditions/values to compare**, especially when they naturally compare one single expression against several possibilities.
4. **Don't add unnecessary `break` statements at the end of `switch` cases** — Go's implicit break already covers that; reserve explicit `break` for exiting a case early, before its natural end.
5. **Use `fallthrough` sparingly and deliberately**, since it deviates from Go's normal, generally clearer default behavior — and remember it unconditionally executes the next case's body without re-checking that case's own condition.
6. **Use a type switch, not a chain of type assertions, when branching on more than one or two possible concrete types.**
7. **Keep `default` as the last clause**, even though Go technically allows it anywhere, for readability and convention.
8. **Don't try to force a ternary-style one-liner** — write out a full `if`/`else`, or a small named helper function for a frequently repeated simple choice.
9. **Be aware that case expressions with side effects will actually execute** during a switch's evaluation, in order, stopping at the first match.
10. **Remember an unlabeled `break` inside a `switch` nested in a loop only exits the `switch`**, not the enclosing loop — use a labeled `break` on the loop itself if that's what's needed.

## 23. Use Case Summary Table

| Technique                               | When to Use                                                      | Example Scenario                                                               |
| --------------------------------------- | ---------------------------------------------------------------- | ------------------------------------------------------------------------------ |
| Basic `if`/`else`                       | A single condition, or a small number of non-parallel conditions | Validating one precondition before proceeding                                  |
| `if` with init clause                   | A check that's only relevant to that one conditional             | `if err := doX(); err != nil { ... }`                                          |
| Early return / guard clause             | Handling failure conditions without nesting the success path     | Validating several sequential preconditions in a function                      |
| Expression `switch`                     | Comparing one value against several discrete possibilities       | Branching on a status string, a day-of-week constant                           |
| Expressionless `switch` (`switch true`) | Many parallel boolean conditions, as a cleaner if/else chain     | Categorizing a score into a letter grade                                       |
| Multiple values in one `case`           | Several distinct values that should share one branch             | Checking if a character is a vowel                                             |
| `fallthrough`                           | Deliberately cascading into the next case's logic                | Sequentially building up behavior across adjacent cases (used sparingly)       |
| Type switch                             | Branching on the concrete type stored in an interface value      | Handling different AST node types, dispatching on several known concrete types |
| `switch` over `if`/`else` chain         | Three or more mutually exclusive comparisons against one value   | Command dispatch, state-machine-style branching                                |

## 24. References

1. Go Team — _A Tour of Go: If_, and the following pages on `if`/`switch`. https://go.dev/tour/flowcontrol/5
2. Go Team — _Effective Go_ (If, Switch, Type switch sections). https://go.dev/doc/effective_go
3. Go Language Specification — _If statements, Switch statements, Expression switches, Type switches_. https://go.dev/ref/spec#If_statements
4. Go Wiki — _Switch_. https://go.dev/wiki/Switch
5. YourBasic Go — _5 switch statement patterns_ and _4 basic if-else statement patterns_. https://yourbasic.org/golang/switch-statement/
6. ZetCode — _Go switch statement_. https://zetcode.com/golang/switch/
7. GeeksforGeeks — _Switch Statement in Go_. https://www.geeksforgeeks.org/go-language/switch-statement-in-go/
8. Leapcell — _Understanding `fallthrough` in Go's `switch` Statements_. https://leapcell.io/blog/understanding-fallthrough-in-go-switch-statements
9. Go Team — _Go 1.21 Release Notes_ (built-in `min`/`max` functions). https://go.dev/doc/go1.21
