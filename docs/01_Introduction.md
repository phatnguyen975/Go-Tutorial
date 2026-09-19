<div align="center">
  <h1>Introduction to Go</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 19, 2026</sub>
</div>

## Table of Contents

1. [What Is Go?](#1-what-is-go)
2. [A Brief History of Go](#2-a-brief-history-of-go)
3. [Why Was Go Created?](#3-why-was-go-created)
4. [Go's Core Design Philosophy](#4-gos-core-design-philosophy)
5. [Key Features of Go](#5-key-features-of-go)
6. [What Is Go Commonly Used For?](#6-what-is-go-commonly-used-for)
7. [Who Uses Go?](#7-who-uses-go)
8. [Installing Go](#8-installing-go)
9. [Verifying the Installation](#9-verifying-the-installation)
10. [Setting Up Your First Project](#10-setting-up-your-first-project)
11. [Writing "Hello, World!"](#11-writing-hello-world)
12. [Understanding the Program, Line by Line](#12-understanding-the-program-line-by-line)
13. [Running the Program](#13-running-the-program)
14. [Building an Executable Binary](#14-building-an-executable-binary)
15. [The Go Toolchain: Essential Commands](#15-the-go-toolchain-essential-commands)
16. [Code Formatting with `gofmt`](#16-code-formatting-with-gofmt)
17. [A Slightly Expanded Example](#17-a-slightly-expanded-example)
18. [Common Beginner Mistakes](#18-common-beginner-mistakes)
19. [Where to Go from Here](#19-where-to-go-from-here)
20. [References](#20-references)

## 1. What Is Go?

**Go** (often called **Golang**, mainly because of its `golang.org`/`go.dev` domain and to make it easier to search for online) is an open-source, general-purpose programming language. It is:

- **Statically typed** — every variable's type is known and checked at compile time.
- **Compiled** — Go source code is compiled directly into a native machine-code binary, not interpreted or run inside a virtual machine.
- **Garbage-collected** — memory management is handled automatically by the runtime; there's no manual `malloc`/`free`.
- **Concurrent by design** — it has built-in language-level support for concurrent programming through goroutines and channels.
- **Simple and minimal** — the language deliberately has a small set of keywords and features, favoring clarity and consistency over flexibility for its own sake.

Its official website describes Go as an open-source programming language that makes it simple to build reliable and efficient software.

## 2. A Brief History of Go

Go's design began in **September 2007** at Google, created by three engineers: **Robert Griesemer**, **Rob Pike**, and **Ken Thompson**. Each brought decades of relevant experience — Ken Thompson co-created Unix and the B programming language (a direct ancestor of C); Rob Pike had worked on the Plan 9 operating system and earlier concurrent languages like Newsqueak and Limbo, which directly influenced Go's goroutines and channels; Robert Griesemer had worked on Google's V8 JavaScript engine and the Java HotSpot virtual machine, bringing deep compiler and runtime expertise.

Key milestones:

- **September 2007** — Design work begins at Google.
- **November 10, 2009** — Go is publicly announced and released as an open-source project.
- **March 2012** — **Go 1.0** is released, establishing the language's first stable, backward-compatibility-guaranteed specification.
- **2015 and beyond** — Go steadily gains adoption, particularly in cloud infrastructure, networking tools, and backend services; landmark projects like Docker and Kubernetes are written in Go, which significantly boosted its visibility.
- **March 2022 (Go 1.18)** — Generics (type parameters) are added, one of the most significant language additions since Go 1.0.
- **Ongoing** — Go continues to be actively developed by Google and a large community of open-source contributors, with a new minor release roughly every six months.

## 3. Why Was Go Created?

By the mid-2000s, Google's engineering teams were working with enormous codebases (tens of millions of lines of code, maintained by thousands of engineers) primarily in C++ and Java. Several recurring pain points motivated Go's creation:

- **Slow compilation.** Large C++ builds could take many minutes — sometimes described as stretching to an hour or more for a major internal binary — significantly slowing down everyday development.
- **Language complexity.** C++ had accumulated decades of features, making it powerful but difficult to use consistently and safely across a huge team of engineers with varying experience levels.
- **Poor built-in support for concurrency and networked, multicore systems.** Google's infrastructure was increasingly built on machines with many CPU cores and large networked/distributed systems, and the mainstream languages of the time offered comparatively little direct language-level help for writing concurrent code correctly.
- **A gap between "fast to write" and "fast to run."** Dynamically typed, interpreted languages like Python were quick and pleasant to write in, but lacked the raw performance and compile-time safety of a statically typed, compiled language like C++.

Go's creators set out to combine the **ease and readability** of a dynamic, interpreted language with the **performance and safety** of a statically typed, compiled one — while adding first-class support for concurrency and keeping compilation fast, even for very large codebases.

## 4. Go's Core Design Philosophy

A few guiding principles run consistently through Go's design, and understanding them helps make sense of many of the language's specific choices:

- **Simplicity over cleverness.** Go deliberately omits many features found in other modern languages (classical inheritance, operator overloading, a ternary operator, exceptions) in favor of a small number of orthogonal, composable building blocks.
- **Readability matters.** Code is read far more often than it's written, and Go's syntax and conventions (like `gofmt`'s single, non-negotiable formatting style) are designed to make any Go codebase familiar-looking to any Go developer.
- **Explicit over implicit.** Error handling via explicit return values (rather than exceptions), and no implicit type conversions, are both examples of Go favoring code where the possibility of failure or a change in representation is always visible.
- **Fast compilation.** Go's dependency model and language design were specifically shaped to allow very fast builds, even for enormous codebases — a direct response to the slow C++ build times that frustrated its creators.
- **Practical, not academic.** Go was designed "in the service of software engineering," in Rob Pike's own words — solving real, everyday problems of building and maintaining large software systems, rather than exploring novel programming-language theory for its own sake.

## 5. Key Features of Go

| Feature                             | What It Means                                                                                                                                                                                   |
| ----------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Statically & strongly typed**     | Types are checked at compile time; no implicit conversions between types.                                                                                                                       |
| **Compiled to native machine code** | Produces a single, self-contained, fast-starting binary with no separate runtime/interpreter needed.                                                                                            |
| **Garbage collected**               | Automatic memory management — no manual allocation/deallocation.                                                                                                                                |
| **Built-in concurrency**            | Goroutines (lightweight, cheaply-created concurrent functions) and channels (typed conduits for communicating between them) are part of the language itself, not a library bolted on afterward. |
| **Fast compilation**                | Go's build times remain fast even for very large projects, by design.                                                                                                                           |
| **Simple, small language**          | A short, memorable list of keywords; deliberately no classical inheritance, no exceptions, no generics-heavy complexity (generics were added later, but kept intentionally restrained).         |
| **Rich standard library**           | Batteries-included support for HTTP servers/clients, JSON, cryptography, text processing, and much more, without needing third-party packages for many common tasks.                            |
| **Cross-compilation**               | A single Go installation can compile binaries for other operating systems and CPU architectures (e.g., building a Windows binary from a Mac) without extra toolchains.                          |
| **Integrated tooling**              | Formatting (`gofmt`), dependency management (`go mod`), testing (`go test`), and more are built directly into the `go` command, with no separate tools to install for these basics.             |
| **Static binaries**                 | A compiled Go program typically has no external runtime dependency to install on the target machine — you can usually just copy the binary and run it.                                          |

## 6. What Is Go Commonly Used For?

- **Backend web services and APIs** — Go's standard library has strong built-in HTTP support, and its performance and concurrency model suit request-heavy server workloads well.
- **Cloud infrastructure and DevOps tooling** — Docker, Kubernetes, Terraform, and many other widely used infrastructure tools are written in Go.
- **Command-line tools (CLIs)** — Go compiles to a single, dependency-free binary, which makes distributing a CLI tool to users straightforward.
- **Networking and distributed systems** — built-in concurrency primitives make Go a natural fit for services that handle many simultaneous connections.
- **Microservices** — fast startup time, low memory footprint, and simple deployment (a single binary) suit containerized microservice architectures well.

## 7. Who Uses Go?

Go is used by a wide range of organizations for backend and infrastructure workloads, including Google itself (its creator and primary steward), and it underpins major open-source projects like Docker and Kubernetes, both of which have significantly shaped how the broader software industry deploys and orchestrates applications. Many other companies across cloud computing, fintech, and internet services also use Go for backend systems, though the specific list of companies and their exact usage changes over time — for current, specific case studies, checking Go's own website or recent industry surveys is more reliable than relying on any fixed list.

## 8. Installing Go

Go can be installed on Windows, macOS, and Linux. The general steps:

1. Visit the official downloads page: **https://go.dev/dl/**
2. Download the installer appropriate for your operating system:
   - **Windows**: an `.msi` installer
   - **macOS**: a `.pkg` installer
   - **Linux**: a `.tar.gz` archive (extracted to `/usr/local/go`, with `/usr/local/go/bin` added to your `PATH`)
3. Run the installer (Windows/macOS) or follow the extraction and `PATH`-configuration steps (Linux) as described on the download page.

Package managers are also commonly used:

```bash
# macOS, via Homebrew
brew install go

# Ubuntu/Debian, via apt (may lag behind the latest release)
sudo apt install golang-go

# Windows, via winget
winget install GoLang.Go
```

**Tip:** the official downloads page always reflects the current, correct installation instructions and latest release for each platform — since Go releases a new version roughly every six months, checking there directly is more reliable than relying on a fixed version number written down elsewhere.

## 9. Verifying the Installation

After installing, open a terminal (Command Prompt, PowerShell, or a Unix shell) and run:

```bash
go version
```

This should print the installed Go version, something like:

```
go version go1.23.0 darwin/arm64
```

If instead you get a "command not found" (or similar) error, Go's `bin` directory likely isn't in your system's `PATH` environment variable — revisiting the installation instructions for your specific operating system usually resolves this.

You can also check where Go's environment is configured:

```bash
go env
```

This prints a list of Go-related environment variables (like `GOPATH`, `GOROOT`, and `GOOS`/`GOARCH`), which can be useful for troubleshooting installation issues.

## 10. Setting Up Your First Project

Modern Go organizes code into **modules**, tracked by a `go.mod` file. To start a new project:

```bash
mkdir hello-go
cd hello-go
go mod init hello-go
```

This creates a `go.mod` file with contents similar to:

```
module hello-go

go 1.23
```

You don't strictly need a module to run a single, quick throwaway file with `go run` in some setups, but creating a module is the standard, idiomatic way to start any real project, and is required as soon as you want to add any external dependencies.

## 11. Writing "Hello, World!"

Create a new file named `main.go` inside your project directory, with the following contents:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

That's the complete, runnable program — four lines of actual code (aside from the blank line), and every Go program you'll ever write shares this same basic shape at its entry point.

## 12. Understanding the Program, Line by Line

Here is the full program again, with each line numbered for reference:

```go
1  package main
2
3  import "fmt"
4
5  func main() {
6      fmt.Println("Hello, World!")
7  }
```

- **Line 1 — `package main`**: Every Go source file begins with a `package` declaration. `main` is a special package name: it tells the Go toolchain that this package produces a standalone, runnable **executable**, rather than a library meant to be imported by other code.
- **Line 3 — `import "fmt"`**: This imports the `fmt` package (short for "format") from Go's standard library, which provides functions for formatted input and output — printing to the console, reading input, formatting strings, and so on. You can only use identifiers from a package (like `fmt.Println`) after importing that package.
- **Line 5 — `func main() {`**: This declares the **entry point** of the program. Every executable Go program must have exactly one function named `main`, inside a package named `main` — this is where program execution begins. It takes no parameters and returns no values. The opening brace `{` starts the function's body.
- **Line 6 — `fmt.Println("Hello, World!")`**: This calls the `Println` function from the `fmt` package, passing it the string `"Hello, World!"`. `fmt.Println` prints its arguments to standard output (typically your terminal), followed by a newline. The dot syntax (`fmt.Println`) is how Go accesses an **exported** (capitalized) identifier from an imported package — `Println` starts with a capital letter, which is exactly what makes it visible and usable outside the `fmt` package itself.
- **Line 7 — `}`**: Closes the body of the `main` function, matching the opening brace on line 5. Once execution reaches this point, the program exits.

**A note on formatting:** Go requires the opening brace `{` of a function (and of `if`, `for`, and similar blocks) to be on the **same line** as the declaration — writing `func main()` and `{` on separate lines is a syntax error, due to how Go's compiler automatically inserts semicolons at the end of certain lines. This is one of several small syntax rules `gofmt` (covered in [Section 16](#16-code-formatting-with-gofmt)) takes care of automatically, so you rarely need to think about it once your editor is set up correctly.

## 13. Running the Program

With `main.go` saved, run it directly with:

```bash
go run main.go
```

Output:

```
Hello, World!
```

`go run` **compiles and immediately executes** the program in one step, without leaving a permanent binary file behind — it's the fastest way to try out a small program or script-like piece of code while developing, since you don't need a separate "build" step every time you want to test a change.

If your project has multiple files or you want to run everything in the current module's main package, you can also use:

```bash
go run .
```

## 14. Building an Executable Binary

To produce a standalone, distributable executable file (rather than running it directly), use `go build`:

```bash
go build
```

This compiles the program and produces a binary named after the module (on this example, `hello-go`, or `hello-go.exe` on Windows) in the current directory. You can then run that binary directly:

```bash
./hello-go        # macOS/Linux
hello-go.exe      # Windows
```

**This is one of Go's most appreciated practical features:** the resulting binary is a single, self-contained file with no separate runtime or interpreter needed on the machine it runs on — you can copy it to another machine (with a compatible OS/architecture) and run it immediately, with nothing else to install.

You can also specify the output binary's name explicitly:

```bash
go build -o myapp
```

## 15. The Go Toolchain: Essential Commands

The single `go` command is the entry point to nearly everything you'll do day-to-day:

| Command                    | Purpose                                                                                                                                                  |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `go run <file(s)/package>` | Compile and immediately run a program, without leaving a binary behind.                                                                                  |
| `go build`                 | Compile the current module's package(s) into an executable (or check that they compile, for a library).                                                  |
| `go fmt`                   | Automatically reformat source files to Go's standard style (a thin wrapper around `gofmt`).                                                              |
| `go vet`                   | Analyze code for common mistakes and suspicious constructs the compiler doesn't catch (e.g., a `Printf` format string that doesn't match its arguments). |
| `go test`                  | Run the tests in the current package (functions named `TestXxx` in `_test.go` files).                                                                    |
| `go mod init <name>`       | Initialize a new module, creating `go.mod`.                                                                                                              |
| `go mod tidy`              | Add missing and remove unused dependencies based on actual imports.                                                                                      |
| `go get <package>`         | Add or update a dependency.                                                                                                                              |
| `go install <package>`     | Build and install an executable to your Go binary directory.                                                                                             |
| `go doc <package/symbol>`  | Show documentation for a package, type, or function, right from the terminal.                                                                            |
| `go version`               | Print the installed Go version.                                                                                                                          |
| `go env`                   | Print Go's environment configuration variables.                                                                                                          |

## 16. Code Formatting with `gofmt`

Go ships with an opinionated, automatic code formatter, `gofmt` (usually invoked via `go fmt` for the current module), which enforces one single, canonical formatting style for indentation, spacing, brace placement, and more:

```bash
go fmt ./...
```

Because every Go developer's tooling applies the exact same formatting rules, Go code from wildly different authors and projects tends to look remarkably consistent — this eliminates entire categories of style debates (tabs vs. spaces, brace placement, and so on) that consume time in many other language communities. Most editors and IDEs with Go support run `gofmt` automatically every time you save a file, so in practice you rarely invoke it manually.

## 17. A Slightly Expanded Example

To show a few more fundamentals in context, here's a slightly larger "hello world"-style program that takes a name and greets it, introducing variables and a simple function:

```go
package main

import "fmt"

func greet(name string) string {
    return "Hello, " + name + "!"
}

func main() {
    name := "Gopher" // Go developers are informally called "Gophers"
    message := greet(name)
    fmt.Println(message)
}
```

Running this with `go run main.go` prints:

```
Hello, Gopher!
```

This small example already demonstrates several core Go ideas: a user-defined function (`greet`) with a typed parameter and return value, short variable declaration (`:=`), and calling a function to build a value used later — patterns that scale up directly into much larger, real-world Go programs.

## 18. Common Beginner Mistakes

### 18.1 Forgetting the `package main` Declaration

```go
// missing "package main" at the top
import "fmt"

func main() {
    fmt.Println("Hello")
}
```

Every Go file needs a `package` declaration as its first non-comment line — an executable program's entry-point file needs specifically `package main`.

### 18.2 Placing the Opening Brace on Its Own Line

```go
func main()
{  // COMPILE ERROR: syntax error, unexpected semicolon (from automatic insertion)
    fmt.Println("Hello")
}
```

The opening brace must be on the same line as the preceding declaration — this is enforced by Go's automatic semicolon insertion rules, and `gofmt` always places it correctly for you.

### 18.3 Importing a Package but Never Using It

```go
import (
    "fmt"
    "os" // COMPILE ERROR if os is never referenced anywhere in the file
)
```

An unused import is a compile-time error in Go, not just a warning — remove any import you're not actually using, or use `goimports` (an extended version of `gofmt`) to have this handled automatically.

### 18.4 Trying to Run a File That Isn't Part of `package main`

```bash
go run mylib.go
# error: go run: cannot run non-main package
```

`go run` (and producing a runnable binary with `go build`) only works for files belonging to `package main` with a `func main()` — a library package is meant to be imported, not run directly.

### 18.5 Expecting Semicolons to Be Required

```go
fmt.Println("Hello");  // valid but unnecessary — gofmt will remove this
```

Go's compiler automatically inserts semicolons at the end of lines according to specific rules, so they're essentially never written by hand in idiomatic Go code — `gofmt` strips out any that are typed explicitly.

## 19. Where to Go from Here

With Go installed and your first program running, natural next steps for continued learning include:

- Exploring Go's basic building blocks: **variables, constants, and data types**.
- Learning Go's **control flow** — conditional statements and loops.
- Understanding **functions**, and how Go's multiple return values and error-handling conventions work.
- Getting comfortable with **structs and methods**, Go's approach to organizing data and behavior without classical classes.
- Diving into Go's standout feature: **concurrency**, via goroutines and channels.
- Learning how Go organizes larger codebases through **packages and modules**.

The official Go documentation (**https://go.dev/doc/**) and the interactive **A Tour of Go** (**https://go.dev/tour/**) are both excellent, authoritative starting points for continuing to explore the language hands-on, directly in a browser, with no local installation required.

## 20. References

1. Go Team — _Go.dev homepage_. https://go.dev/
2. Go Team — _Download and install_. https://go.dev/doc/install
3. Go Team — _A Tour of Go_. https://go.dev/tour/
4. Go Team — _Effective Go_. https://go.dev/doc/effective_go
5. Go Team — _Getting Started_. https://go.dev/doc/tutorial/getting-started
6. Wikipedia — _Go (programming language)_. https://en.wikipedia.org/wiki/Go_(programming_language)
7. AlgoMaster — _History of Go_. https://algomaster.io/learn/go/history-of-go
8. CodiLime — _The Go programming language — everything you should know_. https://codilime.com/blog/what-is-go-language/
9. Ardan Labs — _16 Years of Go: A Programming Language Built to Last_. https://www.ardanlabs.com/news/2025/16-years-of-go-a-programming-language-built-to-last/
10. Albert Bausili Fernández — _The Birth of Go_. https://albertbf.com/articles/the-birth-of-go/
11. Go Team — _Go 1.18 Release Notes_ (generics). https://go.dev/doc/go1.18
12. Go Team — _gofmt documentation_. https://pkg.go.dev/cmd/gofmt
