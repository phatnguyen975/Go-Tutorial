<div align="center">
  <h1>Packages</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 18, 2026</sub>
</div>

## Table of Contents

1. [Packages vs. Modules: The Two Levels of Organization](#1-packages-vs-modules-the-two-levels-of-organization)
2. [What Is a Package?](#2-what-is-a-package)
3. [The `package` Declaration](#3-the-package-declaration)
4. [Package Naming Conventions](#4-package-naming-conventions)
5. [The Special `main` Package](#5-the-special-main-package)
6. [Exported vs. Unexported Identifiers: Scope Between Packages](#6-exported-vs-unexported-identifiers-scope-between-packages)
7. [Importing Packages](#7-importing-packages)
8. [Import Aliasing](#8-import-aliasing)
9. [Blank Imports](#9-blank-imports)
10. [Grouped and Sorted Imports](#10-grouped-and-sorted-imports)
11. [What Is a Module?](#11-what-is-a-module)
12. [Creating a Module: `go mod init`](#12-creating-a-module-go-mod-init)
13. [Anatomy of `go.mod`](#13-anatomy-of-gomod)
14. [The Role of `go.sum`](#14-the-role-of-gosum)
15. [Adding, Updating, and Removing Dependencies](#15-adding-updating-and-removing-dependencies)
16. [Idiomatic Project Layout](#16-idiomatic-project-layout)
17. [The `internal` Directory: Enforced Package Privacy](#17-the-internal-directory-enforced-package-privacy)
18. [Multiple Packages and Multiple Commands in One Module](#18-multiple-packages-and-multiple-commands-in-one-module)
19. [Package Initialization: the `init` Function](#19-package-initialization-the-init-function)
20. [Semantic Versioning and Major Version Suffixes](#20-semantic-versioning-and-major-version-suffixes)
21. [Go Workspaces (`go.work`)](#21-go-workspaces-gowork)
22. [Vendoring Dependencies](#22-vendoring-dependencies)
23. [Circular Imports Are Not Allowed](#23-circular-imports-are-not-allowed)
24. [Common Package-Management Commands](#24-common-package-management-commands)
25. [Common Mistakes and Pitfalls](#25-common-mistakes-and-pitfalls)
26. [Best Practices Summary](#26-best-practices-summary)
27. [Use Case Summary Table](#27-use-case-summary-table)
28. [References](#28-references)

## 1. Packages vs. Modules: The Two Levels of Organization

Go organizes code at two distinct levels, and it's important to keep them separate conceptually:

- A **package** is the basic unit of code organization within Go — a directory of `.go` files that share a namespace, all declaring the same `package` name at the top of each file.
- A **module** is a versioned collection of one or more packages that are released, distributed, and depended-upon together, defined by a `go.mod` file at the root of the module's directory tree.

A single module very commonly contains many packages (one per subdirectory), and a project you're working on is normally exactly one module, unless you're deliberately using a multi-module workspace (see [Section 21](#21-go-workspaces-gowork)).

## 2. What Is a Package?

Every `.go` file belongs to exactly one package, declared at the top of the file. All files in the same directory must declare the **same** package name — a directory cannot mix files from two different packages (with one narrow exception: `_test` package suffixes for external test packages).

```go
// greet.go
package greet

func Hello(name string) string {
    return "Hello, " + name
}
```

```go
// greet_helpers.go — same directory, same package
package greet

func shout(s string) string {
    return strings.ToUpper(s) + "!"
}
```

A package's name is tied to the **directory** it lives in, not to any particular file name — this is why Go's convention is "organize code by directory," not by individual file naming schemes.

## 3. The `package` Declaration

The `package` clause must be the first non-comment line of every Go source file:

```go
package mypackage
```

Every file in the directory repeats this exact same declaration. When another piece of code imports this package, it refers to it by its **import path** (derived from the module path plus the directory's location within the module — see [Section 7](#7-importing-packages)), but within the code itself, the package's declared short name (`mypackage`) is what's used to qualify its exported identifiers (`mypackage.SomeFunction()`).

## 4. Package Naming Conventions

Go has strong, widely followed conventions for package names, distinct from many other languages:

- **Short, lowercase, single-word names** — `http`, `json`, `time`, `sort`. Avoid `mixedCaps` or `under_scores` in package names.
- **No plurals** — `strings` (a standard-library exception, largely historical) aside, prefer `buf` over `buffers`, `user` over `users`, unless a plural genuinely reads better and is well established.
- **Name the package for what it _provides_, not what it _contains_**. A package of utility functions for strings is named `strings`, not `stringutils` or `util` — this makes call sites read naturally (`strings.Contains(...)`), since the package name becomes a prefix at every call site.
- **Avoid generic, catch-all names** like `common`, `util`, `helpers`, or `misc`. These tend to accumulate unrelated code over time and provide callers no signal about what's actually inside — a more specific name (`validate`, `retry`, `ratelimit`) is almost always better and keeps related functionality cohesive.
- **The package name is not repeated in its own exported identifiers.** Since callers already write `pkgname.Thing`, a type or function inside package `user` should be named `User`, `New`, `Validate` — not `UserUser`, `NewUser` when already inside `package user`, or `ValidateUser` — the package name already supplies that context.

```go
// Inside package user:
type User struct{ /* ... */ }       // good: called as user.User
func New(name string) *User { ... } // good: called as user.New(...)

// Avoid:
type UserAccount struct{ /* ... */ } // redundant: called as user.UserAccount
```

## 5. The Special `main` Package

A package named `main` is treated specially by the Go toolchain: it's the entry point of an executable program, rather than an importable library. A `main` package must also declare a `func main()` with no parameters and no return values — this is where program execution begins.

```go
// cmd/myapp/main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello from the entry point")
}
```

Every standalone Go program (as opposed to a library meant to be imported by other code) needs exactly one `main` package with exactly one `main()` function in its build. A module can define multiple different `main` packages, in different directories, to produce multiple separate executables from one codebase (see [Section 18](#18-multiple-packages-and-multiple-commands-in-one-module)).

## 6. Exported vs. Unexported Identifiers: Scope Between Packages

Go determines visibility purely through **capitalization**, with no separate `public`/`private` keywords:

- An identifier (function, type, variable, constant, struct field, interface method) whose name starts with an **uppercase** letter is **exported** — visible and usable from any package that imports the package it's declared in.
- An identifier whose name starts with a **lowercase** letter is **unexported** — visible only within the package it's declared in, including across every file in that same package directory, but invisible to any importing package.

```go
package account

type Account struct {
    ID      string // exported — accessible as account.Account{}.ID from other packages
    balance float64 // unexported — only visible inside package account itself
}

func New(id string) *Account { // exported: account.New(...)
    return &Account{ID: id}
}

func validateID(id string) bool { // unexported: only callable from within package account
    return len(id) > 0
}
```

```go
package main

import "myapp/account"

func main() {
    a := account.New("acc-1")   // OK: New is exported
    fmt.Println(a.ID)             // OK: ID is exported
    // fmt.Println(a.balance)    // COMPILE ERROR: balance is unexported, invisible outside package account
    // account.validateID("x")  // COMPILE ERROR: validateID is unexported
}
```

**Scope within a package:** every identifier — exported or not — is visible across **all files in the same package directory**, without needing any import at all between them, since they all share one namespace. Unexported identifiers are only hidden from _other_ packages, not from sibling files within the same package.

**This is Go's primary encapsulation mechanism.** With no `private`/`protected`/`public` keywords, capitalization plus the package boundary is the entirety of Go's visibility model — a deliberate simplicity that keeps the rule easy to remember (and easy to see at a glance in any piece of code, without needing to check a separate declaration).

## 7. Importing Packages

An `import` statement brings another package's exported identifiers into scope, qualified by that package's name:

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    s := strings.ToUpper("hello")
    fmt.Println(s) // HELLO
}
```

### 7.1 Import Paths

The string inside an `import` statement is an **import path** — for standard-library packages, it's simply the package's location relative to the standard library root (`"fmt"`, `"net/http"`, `"encoding/json"`). For packages within your own module or a third-party module, the import path is the **module path** (declared in that module's `go.mod`) plus the package's directory path relative to the module root:

```go
import (
    "myapp/internal/auth"           // a package in the current module
    "github.com/gorilla/mux"        // a third-party package from a different module
)
```

**Important distinction:** the string used in the `import` statement (the import path) can differ from the short package **name** actually declared inside those files with `package ...` — by convention they usually match the last path segment, but this isn't a hard requirement, and code always refers to the package by its declared `package` name, not by the last segment of the import path, when the two happen to differ.

### 7.2 Unused Imports Are a Compile Error

Go treats an imported-but-never-used package as a compile-time error, not a warning — this is a deliberate language design choice to keep import lists accurate and prevent unused-dependency drift:

```go
import "fmt" // COMPILE ERROR if fmt is never actually referenced anywhere in the file
```

## 8. Import Aliasing

An import can be given a different local name than its package's own declared name, using an alias directly before the import path:

```go
import (
    "fmt"
    myjson "encoding/json" // aliased: referred to as myjson.Marshal(...), not json.Marshal(...)
)
```

**Common reasons to alias an import:**

- **Resolving a name collision** when two different packages you need to import happen to share the same default package name.
- **Improving clarity** when a package's default name is ambiguous or unclear in context.
- **Following a project-wide convention** for a specific commonly-used package (less common, but seen in some codebases).

```go
import (
    v1 "myapp/api/v1"
    v2 "myapp/api/v2"
)

// now referred to as v1.Something and v2.Something, resolving what would otherwise be
// two conflicting imports both trying to use the same default package name
```

## 9. Blank Imports

Prefixing an import with an underscore (`_`) imports a package purely for its **side effects** — running its `init()` function(s) (see [Section 19](#19-package-initialization-the-init-function)) — without making any of its exported identifiers directly accessible, and without triggering the "unused import" compile error that a normal, unused import would cause:

```go
import (
    "database/sql"
    _ "github.com/lib/pq" // registers the "postgres" driver with database/sql via its init() function, but is never referenced by name
)

func main() {
    db, err := sql.Open("postgres", connectionString) // works because the pq driver registered itself
}
```

This pattern is common for **database drivers** (which register themselves with `database/sql` via `init()`), and for packages whose only purpose at a given import site is to trigger some registration or setup logic, not to expose functions/types you'll call directly.

## 10. Grouped and Sorted Imports

Idiomatic Go groups imports into blocks — conventionally, standard-library imports first, then a blank line, then third-party and local module imports — and within each group, imports are sorted alphabetically:

```go
import (
    "fmt"
    "net/http"
    "strings"

    "github.com/gorilla/mux"

    "myapp/internal/auth"
    "myapp/internal/config"
)
```

The `gofmt`/`goimports` tools automatically sort and group imports according to this convention (and `goimports` additionally adds/removes imports automatically based on what the file actually uses), so in practice most Go developers rely on their editor or a pre-commit hook to enforce this formatting rather than doing it by hand.

## 11. What Is a Module?

A **module** is a collection of one or more related Go packages that are versioned and released together as a single unit — it's the granularity at which Go's dependency management (via the `go` command and `go.mod`) actually operates. A module is defined by having a `go.mod` file at the root of its directory tree; every package inside that tree, in every subdirectory, belongs to that one module (unless a subdirectory itself contains another `go.mod`, making it a separate, nested module — a pattern that's possible but relatively uncommon outside of monorepos with deliberately independent sub-projects).

Modules were introduced in **Go 1.11** as the successor to the earlier `GOPATH`-based dependency workflow, specifically to give Go reliable, reproducible dependency version management — before modules, all Go source across every project on a machine typically lived under one shared `$GOPATH/src` directory with no per-project versioning of dependencies at all.

## 12. Creating a Module: `go mod init`

A new module is initialized with:

```bash
go mod init github.com/someuser/myproject
```

This creates a `go.mod` file in the current directory, declaring the module's **path** — conventionally the repository location where the code will be hosted (a GitHub path, for instance), which also becomes the import-path prefix for every package inside this module.

```
myproject/
  go.mod         # module github.com/someuser/myproject
  main.go        # package main
  greet/
    greet.go     # package greet, importable as "github.com/someuser/myproject/greet"
```

**The module path does not have to correspond to a real, publicly hosted repository** for purely local or private development — but if the module will ever be published for others to `go get`, using its actual intended repository location as the module path is required, since that's exactly how the `go` command locates and downloads it.

## 13. Anatomy of `go.mod`

A typical `go.mod` file looks like this:

```
module github.com/someuser/myproject

go 1.22

require (
    github.com/gorilla/mux v1.8.1
    github.com/stretchr/testify v1.9.0
)

require (
    github.com/davecgh/go-spew v1.1.1 // indirect
    github.com/pmezard/go-difflib v1.0.0 // indirect
)
```

| Directive     | Meaning                                                                                                                                                                                                              |
| ------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `module`      | Declares the module's path — both its identity and the import-path prefix for all packages inside it.                                                                                                                |
| `go`          | Declares the minimum Go language version the module requires/targets, affecting which language features and standard-library behaviors are available.                                                                |
| `require`     | Lists direct and indirect (transitive) dependencies, each pinned to a specific version.                                                                                                                              |
| `// indirect` | A comment (not a directive) automatically added by tooling to mark a dependency that isn't imported directly by this module's own code, but is needed because one of the module's direct dependencies depends on it. |
| `replace`     | (Not shown above) Substitutes a dependency's source — often pointing to a local filesystem path or a fork — useful during development or when patching a dependency.                                                 |
| `exclude`     | (Not shown above) Explicitly excludes a specific version of a dependency from being selected, even if some other dependency would otherwise pull it in.                                                              |

`go.mod` is normally maintained automatically by the `go` command (`go get`, `go mod tidy`, and so on) rather than being hand-edited line by line, though manual edits (especially adding a `replace` directive) are common and fully supported.

## 14. The Role of `go.sum`

`go.sum` is an auto-generated companion file that records the **cryptographic checksums** of every dependency's module content (and, historically, its `go.mod` file), across every version that has ever been resolved for the build at some point. Its purpose is **integrity verification**: when the `go` command downloads a dependency, it checks the downloaded content's checksum against what's recorded in `go.sum` — if they don't match (indicating the dependency's published content was altered, corrupted, or tampered with since it was first recorded), the build fails loudly rather than silently using untrusted code.

```
github.com/gorilla/mux v1.8.1 h1:TuBL49tXwgrFYWhqrNgrUNEY92u81SPhu7sTdzQEiWY=
github.com/gorilla/mux v1.8.1/go.mod h1:DVbg23sWSpFRCP0SfiEN6jmj59UnW/n46BH5rLB71So=
```

`go.sum` should be **committed to version control alongside `go.mod`** — it's precisely what gives Go's dependency resolution its reproducibility guarantee: anyone building the project later gets the exact same, verified dependency content, not whatever happens to be available at build time.

## 15. Adding, Updating, and Removing Dependencies

### 15.1 The Implicit Way: Just Import and Run `go mod tidy`

The simplest way to add a dependency is often to simply add the `import` statement in your code and then run:

```bash
go mod tidy
```

`go mod tidy` scans all the imports actually used across the module's source code, then adds any missing entries to `go.mod`/`go.sum` and removes any entries that are no longer actually needed by any file — it's the standard, idiomatic way to keep these two files accurate and minimal, and is commonly run after any change to imports.

### 15.2 The Explicit Way: `go get`

```bash
go get github.com/gorilla/mux          # add/update to the latest version
go get github.com/gorilla/mux@v1.8.1    # add/update to a specific version
go get github.com/gorilla/mux@latest    # explicitly request the latest version
go get github.com/gorilla/mux@none      # remove the dependency entirely
```

`go get` explicitly modifies `go.mod` and downloads the requested version (updating `go.sum` accordingly) in one step, without needing to have already written the corresponding `import` statement first.

### 15.3 `go install` for Executables

Since Go 1.17, `go install` (rather than `go get`) is the dedicated command for building and installing an executable binary from a module, without affecting the current directory's own `go.mod`:

```bash
go install github.com/some/tool@latest
```

This places the compiled binary in `$GOBIN` (or `$GOPATH/bin` by default), making it directly runnable from the command line, separate from adding something as a source-level dependency of your own module.

## 16. Idiomatic Project Layout

While Go does not mandate any specific directory structure beyond "one package per directory," a widely converged-upon layout has emerged across the community for typical application projects:

```
myapp/
├── go.mod
├── go.sum
├── cmd/                      # entry points — one subdirectory per executable
│   ├── server/
│   │   └── main.go           # package main — the HTTP server binary
│   └── worker/
│       └── main.go           # package main — a background worker binary
├── internal/                 # private application code, not importable outside this module
│   ├── handler/
│   │   ├── handler.go
│   │   └── handler_test.go
│   ├── store/
│   │   ├── postgres.go
│   │   └── store.go
│   └── config/
│       └── config.go
└── pkg/                       # code intended to be importable by OTHER modules/projects
    └── validate/
        ├── validate.go
        └── validate_test.go
```

| Directory   | Purpose                                                                                                                                                                                                                                                                                                                                                                      |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cmd/`      | One subdirectory per runnable binary the module produces, each containing its own `package main` and `main()`. Keeps multiple entry points cleanly separated.                                                                                                                                                                                                                |
| `internal/` | Application-specific code that should never be imported by anything outside this module — enforced by the compiler, not just convention (see [Section 17](#17-the-internal-directory-enforced-package-privacy)). The overwhelming majority of a typical application's own code lives here.                                                                                   |
| `pkg/`      | Code that's genuinely meant to be reused and imported by other, separate projects/modules. This directory is a convention, not a compiler-enforced boundary — unlike `internal`, anything here is importable by anyone. Many Go projects, especially applications rather than libraries, don't need a `pkg/` directory at all if nothing is meant to be externally reusable. |

**This layout is a convention, not a requirement enforced by the Go toolchain** (aside from `internal`'s special compiler-enforced behavior) — small projects, single-package libraries, and simple CLIs are perfectly idiomatic without adopting the full `cmd`/`internal`/`pkg` structure. Adopt it when a project's size and shape actually benefits from the separation it provides, not reflexively for every project regardless of scale.

## 17. The `internal` Directory: Enforced Package Privacy

Unlike `pkg/` (a pure convention), a directory literally named `internal` is given **special, compiler-enforced meaning** by the Go toolchain: any package under a path containing an `internal` directory segment can only be imported by code that lives inside the directory tree **rooted at the parent of that `internal` directory** — code outside that tree gets a compile error if it tries to import it.

```
project-root/
  internal/
    auth/
      auth.go        # package auth, importable ONLY by code under project-root/
  main.go             # can import "myproject/internal/auth" — inside the allowed tree
```

A module elsewhere on the filesystem, or even a different part of a larger monorepo outside `project-root`, attempting `import "myproject/internal/auth"` gets a build error — the Go toolchain itself enforces this boundary, not merely a documentation convention.

```
project-root-directory/
  go.mod                    # module github.com/someuser/modname
  modname.go                 # package modname — the root package
  auth/
    auth.go                  # package auth — publicly importable
    token/
      token.go                # package token — publicly importable
  hash/
    hash.go                   # package hash — publicly importable
  internal/
    trace/
      trace.go                 # package trace — importable ONLY within this module's tree
```

With this layout, external users of the `modname` module could import `github.com/someuser/modname/auth` or `github.com/someuser/modname/auth/token`, but never `github.com/someuser/modname/internal/trace`.

**Why this matters for library authors:** it's the standard way to expose a clean, deliberately curated public API (the packages outside `internal/`) while keeping true implementation details free to be refactored, renamed, or restructured without ever being considered a breaking change for external consumers, since nothing outside the module could have depended on them in the first place. For this reason, it's recommended to keep as much of a module's code under `internal/` as reasonably possible, exposing only what genuinely needs to be a public, stable API surface.

## 18. Multiple Packages and Multiple Commands in One Module

A single module can contain many importable packages, each its own subdirectory, and can also produce multiple separate executables:

```
myapp/
  go.mod                    # module github.com/someuser/myapp
  internal/
    ...                      # shared internal packages used by both commands below
  server/
    main.go                  # package main — produces one binary: "server"
  worker/
    main.go                  # package main — produces a second, separate binary: "worker"
```

Each `main.go` here is its own independent `package main`, and running `go build ./...` (or building each directory individually) produces two separate executables, both able to import and share the same `internal/` packages between them.

## 19. Package Initialization: the `init` Function

Any package (including `main`) may declare one or more functions named `init`, each taking no arguments and returning nothing:

```go
package config

var Settings map[string]string

func init() {
    Settings = loadDefaultSettings() // runs automatically before main() or any other code in this package is used
}
```

**Rules governing `init`:**

- A package can have **multiple** `init` functions, even within a single file, and they run in the order they appear in the source (across multiple files in a package, the order follows the order the Go compiler presents the files, which is typically alphabetical by filename, though relying on cross-file ordering is generally discouraged).
- `init` functions run **automatically** — they are never called explicitly by name, and in fact cannot be, since `init` is not an ordinary callable identifier from outside the package.
- All of a package's dependencies are fully initialized (their own `init` functions have already run, and their package-level variables are already set) **before** that package's own `init` functions run — Go resolves this ordering automatically based on the import graph.
- For a `main` package, all `init` functions across every imported package (transitively) run before `func main()` itself begins executing.

**Common, legitimate uses of `init`:** registering a database driver or codec with a shared registry (as with blank imports in [Section 9](#9-blank-imports)), validating that some required package-level configuration is well-formed at startup, or performing genuinely one-time setup that must happen before any of the package's exported functionality is used.

**A caution:** overusing `init` for complex, order-dependent, or side-effect-heavy logic makes a program's startup behavior harder to trace and test (since `init` runs implicitly, with no way to pass in different behavior per test or per call site) — many experienced Go developers prefer explicit initialization functions that calling code invokes deliberately, reserving `init` for the narrower cases above.

## 20. Semantic Versioning and Major Version Suffixes

Go Modules follow **semantic versioning** (`vMAJOR.MINOR.PATCH`, e.g., `v1.4.2`) for released module versions, and the Go tooling has a specific, distinctive rule for major version changes: starting at **major version 2 and above**, the module's path itself must include the major version as a suffix, both in `go.mod`'s `module` line and in every importer's import path:

```
module github.com/someuser/mymodule/v2
```

```go
import "github.com/someuser/mymodule/v2"
```

This convention — called **semantic import versioning** — exists so that two incompatible major versions of the same module (say, `v1` and `v2`) can be imported **side by side** within the same build if needed (for instance, during a gradual migration), since their import paths are, from Go's perspective, genuinely different paths pointing at different code. Versions `v0` and `v1` do not require this suffix; it only becomes necessary starting from `v2`.

## 21. Go Workspaces (`go.work`)

Introduced in **Go 1.18**, workspaces let you work across **multiple modules simultaneously** — for example, developing a library and an application that depends on it, in the same local checkout, with changes to the library immediately visible to the application without needing to publish a new version or add manual `replace` directives to `go.mod`.

```bash
go work init ./mymodule ./myapp
```

This creates a `go.work` file:

```
go 1.22

use (
    ./mymodule
    ./myapp
)
```

While a `go.work` file is present, the `go` command resolves imports across all the listed modules together, letting local, in-progress changes in one module be picked up immediately by another module in the same workspace during local development — `go.work` is typically **not committed to version control**, since it reflects a developer's local, personal workspace setup, not the module's own published dependency graph.

## 22. Vendoring Dependencies

`go mod vendor` copies the full source of every dependency into a `vendor/` directory inside the module, producing a self-contained snapshot that doesn't require network access (or a module proxy) to build:

```bash
go mod vendor
```

```
myapp/
  go.mod
  go.sum
  vendor/
    github.com/gorilla/mux/
      ...
    modules.txt
  main.go
```

Once a `vendor/` directory is present, builds can use `-mod=vendor` (or Go will use it automatically in many configurations when the directory exists and is consistent with `go.mod`) to build entirely from the vendored copies rather than fetching from the network. Vendoring is less commonly needed since the introduction of Go's module proxy and checksum database (which already provide reliable, cached, verified dependency resolution for most setups), but it remains useful for fully air-gapped build environments, strict reproducibility/auditability requirements, or organizations with policies against builds fetching code over the network at build time.

## 23. Circular Imports Are Not Allowed

Go's compiler forbids **import cycles** — if package `A` imports package `B`, package `B` cannot (directly or transitively, through any chain of imports) import package `A` back:

```go
// package a imports package b
// package b imports package a
// → COMPILE ERROR: import cycle not allowed
```

This is a deliberate, enforced constraint, not just a style guideline — Go's build model requires a package's dependencies to form a strict, acyclic graph, partly because it simplifies compilation (each package can be compiled once its dependencies are already compiled) and partly because it pushes toward cleaner, more layered designs.

**When a cycle seems necessary**, it's almost always a sign that some piece of shared functionality should be extracted into a **third package** that both original packages can depend on, rather than depending on each other directly:

```
Before (cyclic, not allowed):
  package order  →  imports  →  package customer
  package customer  →  imports  →  package order

After (acyclic, valid):
  package order     →  imports  →  package shared
  package customer  →  imports  →  package shared
```

## 24. Common Package-Management Commands

| Command                     | Purpose                                                                                             |
| --------------------------- | --------------------------------------------------------------------------------------------------- |
| `go mod init <module-path>` | Create a new module, generating `go.mod`.                                                           |
| `go mod tidy`               | Add missing and remove unused dependencies in `go.mod`/`go.sum`, based on what's actually imported. |
| `go get <path>[@version]`   | Add or update a dependency to a specific version (or `@latest`, or `@none` to remove it).           |
| `go install <path>@version` | Build and install an executable binary, without modifying the current module's `go.mod`.            |
| `go build ./...`            | Build all packages in the module (the `./...` pattern matches every package recursively).           |
| `go list -m all`            | List the current module and all of its resolved dependencies.                                       |
| `go mod why <path>`         | Explain why a particular dependency is needed (which import chain pulls it in).                     |
| `go mod vendor`             | Copy all dependencies' source into a local `vendor/` directory.                                     |
| `go mod verify`             | Verify that locally cached dependency content matches the checksums recorded in `go.sum`.           |
| `go work init <dirs...>`    | Create a `go.work` file for a multi-module local workspace.                                         |

## 25. Common Mistakes and Pitfalls

### 25.1 Editing `go.sum` by Hand

`go.sum` is machine-generated and should never be hand-edited — always let `go mod tidy`, `go get`, or `go mod verify` manage its contents; manual edits risk breaking the exact checksum-matching mechanism the file exists to provide.

### 25.2 Forgetting to Run `go mod tidy` After Changing Imports

Adding a new `import` without running `go mod tidy` (or an equivalent `go get`) can leave `go.mod` missing a required entry, which surfaces as a build failure for anyone else pulling the latest code (or in CI) until `go mod tidy` is run and the resulting changes committed.

### 25.3 Using a Generic, Catch-All Package Name

```go
package utils // vague — accumulates unrelated code over time, provides no signal to readers
```

Prefer a name that describes what the package actually provides, as discussed in [Section 4](#4-package-naming-conventions).

### 25.4 Expecting `internal` to Work Without the Exact Directory Name

The special import-restriction behavior is triggered specifically by a directory literally named `internal` — a directory named `private`, `impl`, or anything else does not get this compiler-enforced protection, even if it's intended to serve the same purpose.

### 25.5 Forgetting the Major Version Suffix for v2+ Modules

```
module github.com/someuser/mymodule // missing /v2 suffix, but the module is tagged v2.0.0
```

Once a module releases a `v2` (or higher) tag, its `go.mod` module path and every importer's import path must include the matching `/v2` (or `/v3`, etc.) suffix, per [Section 20](#20-semantic-versioning-and-major-version-suffixes) — omitting this leads to confusing version-resolution errors for consumers.

### 25.6 Committing a `go.work` File Meant Only for Local Development

A `go.work` file that points at local filesystem paths (e.g., `./mymodule`) is specific to one developer's local checkout layout — committing it can break builds for anyone else (or CI) whose local directory structure doesn't match, since `go.work` isn't typically part of a module's own published configuration.

### 25.7 Creating an Import Cycle by Accident

A design where package `A`'s types need to reference package `B`'s types, and vice versa, commonly leads to an accidental cyclic-import attempt as the code grows — resolve it by extracting the shared pieces both packages need into a separate, lower-level package neither depends back on, per [Section 23](#23-circular-imports-are-not-allowed).

## 26. Best Practices Summary

1. **Name packages short, lowercase, and for what they provide**, avoiding generic catch-all names like `utils` or `common`.
2. **Default to unexported (lowercase) identifiers**, and only capitalize (export) what genuinely needs to be part of the package's public API.
3. **Run `go mod tidy` after any change to imports**, and commit both `go.mod` and `go.sum` to version control.
4. **Never hand-edit `go.sum`** — let the tooling manage it.
5. **Put the overwhelming majority of an application's own code under `internal/`**, exposing only a deliberately curated, minimal public API in any packages outside it.
6. **Use `cmd/` for entry points** when a module produces more than one executable, keeping each `main` package cleanly separated.
7. **Reserve `pkg/` for code genuinely meant to be imported by other, separate projects** — many applications don't need it at all.
8. **Keep `init` functions narrow and side-effect-focused** (driver registration, startup validation) rather than complex, order-dependent application logic.
9. **Follow semantic import versioning** — add the `/v2`, `/v3`, etc. suffix to both the module path and every import path once a module's major version reaches 2 or higher.
10. **Resolve an apparent need for circular imports by extracting a shared third package**, rather than trying to work around the compiler's cycle restriction.

## 27. Use Case Summary Table

| Technique                                 | When to Use                                                                           | Example Scenario                                                           |
| ----------------------------------------- | ------------------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `go mod init <path>`                      | Starting any new Go project meant to manage its own dependencies                      | Beginning a new application or library                                     |
| `go mod tidy`                             | After adding/removing/changing any import in the codebase                             | Routine maintenance, before every commit that touches imports              |
| `go get <path>@<version>`                 | Adding a dependency, or pinning/upgrading to a specific version                       | Adding a new third-party library, upgrading past a known bug fix           |
| `internal/` directory                     | Code that must never be importable outside this module                                | The bulk of an application's own business logic and implementation details |
| `cmd/` directory                          | A module that builds more than one executable                                         | A server binary and a separate CLI/worker binary sharing internal packages |
| `pkg/` directory                          | Code specifically meant for reuse by other, separate modules                          | A genuinely standalone, reusable validation or parsing library             |
| Import aliasing                           | Two imports would otherwise share the same default package name                       | Importing two different versions/variants of a similarly-named API         |
| Blank import (`_ "pkg"`)                  | Triggering a package's `init()` side effects without using its identifiers directly   | Registering a database driver with `database/sql`                          |
| `replace` directive in `go.mod`           | Developing against a local fork/patch of a dependency, or resolving a temporary issue | Testing an unreleased fix to a dependency before it's published            |
| `go.work` workspace                       | Developing across multiple local modules simultaneously                               | Working on a library and the application consuming it, in the same session |
| `go mod vendor`                           | Fully offline, air-gapped, or strictly audited build environments                     | CI/CD in a network-restricted environment                                  |
| Semantic import versioning (`/v2`, `/v3`) | Releasing a breaking-change major version of a public module                          | Publishing a v2 of a library with an incompatible API                      |

## 28. References

1. Go Team — _Organizing a Go module_. https://go.dev/doc/modules/layout
2. Go Team — _Go Modules Reference_. https://go.dev/ref/mod
3. Go Team — _Tutorial: Create a Go module_. https://go.dev/doc/tutorial/create-module
4. Go Team — _Using Go Modules_, The Go Blog. https://go.dev/blog/using-go-modules
5. Go Team — _Go Modules: v2 and Beyond_, The Go Blog (semantic import versioning). https://go.dev/blog/v2-go-modules
6. Go Team — _Go Workspaces (go.work)_. https://go.dev/ref/mod#workspaces
7. Go Team — _Tutorial: Getting started with multi-module workspaces_. https://go.dev/doc/tutorial/workspaces
8. Go Language Specification — _Packages, Import declarations, Package initialization_. https://go.dev/ref/spec#Packages
9. Go Team — _Effective Go_ (Package names, Exported names). https://go.dev/doc/effective_go
10. Go standard library documentation — _cmd/go_ command reference (`go get`, `go install`, `go mod`, etc.). https://pkg.go.dev/cmd/go
11. golinuxcloud — _Golang packaging, package structure, naming conventions_. https://www.golinuxcloud.com/golang-packaging/
12. Benita Emudianughe — _Go Beyond Basics: A Complete Guide to Go Packages, Modules, and Project Setup_, Medium. https://medium.com/@emusbeny/go-beyond-basics-a-complete-guide-to-go-packages-modules-and-project-setup-9ae082fbf3cd
13. ferztyle — _Go Packages and Modules explained_, DEV Community. https://dev.to/ferztyle/go-packages-and-modules-explained-4j4c
14. Earthly Blog — _Understanding Go Package Management and Modules_. https://earthly.dev/blog/go-modules/
15. ZetCode — _Working with Modules in Go_. https://zetcode.com/golang/module/
