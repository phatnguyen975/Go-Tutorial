<div align="center">
  <h1>Packages</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 17, 2026</sub>
</div>

## Go Package

In the most basic terms, A package is nothing but a directory inside your Go workspace containing one or more Go source files, or other Go packages.

Every Go source file belongs to a package. To declare a source file to be part of a package, we use the following syntax:

```go
package <packagename>
```

The above package declaration must be the first line of code in your Go source file. All the functions, types, and variables defined in your Go source file become part of the declared package.

You can choose to export a member defined in your package to outside packages, or keep them private to the same package. Other packages can import and reuse the functions or types that are exported from your package.

Almost all the code that we have seen so far in this tutorial series include the following line:

```go
import "fmt"
```

`fmt` is a core library package that contains functionalities related to formatting and printing output or reading input from various I/O sources. It exports functions like `Println()`, `Printf()`, `Scan()`, etc, for other packages to reuse.

**Packaging functionalities in this way has the following benefits**

- It reduces naming conflicts. You can have the same function names in different packages. This keeps our function names short and concise.
- It organizes related code together so that it is easier to find the code you want to reuse.
- It speeds up the compilation process by only requiring recompilation of smaller parts of the program that has actually changed. Although we use the `fmt` package, we don't need to recompile it every time we change our program.

### The `main` Package and `main()` Function

Go programs start running in the `main` package. It is a special package that is used with programs that are meant to be executable.

By convention, executable programs (the ones with the `main` package) are called **Commands**. Others are called simply **Packages**.

The `main()` function is a special function that is the entry point of an executable program. Let's see an example of an executable program in Go.

```go
// Package declaration
package main

// Importing packages
import (
    "fmt"
    "time"
    "math"
    "math/rand"
)

func main() {
    // Finding the Max of two numbers
    fmt.Println(math.Max(73.15, 92.46))

    // Calculate the square root of a number
    fmt.Println(math.Sqrt(225))

    // Printing the value of '𝜋'
    fmt.Println(math.Pi)

    // Epoch time in milliseconds
    epoch := time.Now().Unix()
    fmt.Println(epoch)

    // Generating a random integer between 0 to 100
    rand.Seed(epoch)
    fmt.Println(rand.Intn(100))
}
```

### Importing Packages

There are two ways to import packages in Go:

```go
// Multiple import statements
import "fmt"
import "time"
import "math"
import "math/rand"

// Factored import statements
import (
    "fmt"
    "time"
    "math"
    "math/rand"
)
```

Go's convention is that - the package name is the same as the last element of the import path. For example, the name of the package imported as `math/rand` is `rand`. It is imported with path `math/rand` because It is nested inside the `math` package as a subdirectory.

### Exported vs Unexported Names

> Anything (variable, type, or function) that starts with a capital letter is exported, and visible outside the package.

> Anything that does not start with a capital letter is not exported, and is visible only inside the same package.

When you import a package, you can only access its exported names.

```go
package main

import (
    "fmt"
    "math"
)

func main() {
    // MaxInt64 is an exported name
    fmt.Println("Max value of int64: ", int64(math.MaxInt64))

    // Phi is an exported name
    fmt.Println("Value of Phi (ϕ): ", math.Phi)

    // pi starts with a small letter, so it is not exported
    fmt.Println("Value of Pi (𝜋): ", math.pi)
}
```

```text
# Output
./exported_names.go:16:38: cannot refer to unexported name math.pi
./exported_names.go:16:38: undefined: math.pi
```

To fix the above error, you need to change `math.pi` to `math.Pi`.

## Creating and Managing Custom Packages

Until now, We have only written code in the `main` package and used functionalities imported from Go's core library packages.

Let's create a sample Go project that has multiple custom packages with a bunch of source code files and see how the same concept of package declaration, imports, and exports apply to custom packages as well.

Fire up your terminal and create a directory for our Go project:

```bash
$ mkdir packer
```

Next, we'll create a [Go module](https://blog.golang.org/using-go-modules) and make the project directory the root of the module.

**Note:** Go module is Go's new dependency management system. A module is a collection of Go packages stored in a directory with a `go.mod` file at its root. The `go.mod` file defines the module's path, which is also the import path used while importing packages that are part of this module.

Let's initialize a Go module by typing the following commands:

```bash
$ cd packer
$ go mod init github.com/username/packer
```

Let's now create some source files and place them in different packages inside our project. The following image displays all the packages and the source files:

<p align="center">
  <img src="https://www.callicoder.com/static/1fff108fd759e383c7ed482fdc33f730/b06ae/go-custom-package-organization.png" style="width:80%;" alt="Packer">
</p>

**numbers/prime.go**

```go
package numbers

import "math"

// Checks if a number is prime or not
func IsPrime(num int) bool {
    for i := 2; i <= int(math.Floor(math.Sqrt(float64(num)))); i++ {
        if num % i == 0 {
            return false
        }
    }
    return num > 1
}
```

**strings/reverse.go**

```go
package strings

// Reverses a string
func Reverse(s string) string {
    runes := []rune(s)
    reversedRunes := reverseRunes(runes)
    return string(reversedRunes)
}
```

**strings/reverse_runes.go**

```go
package strings

// Reverses an array of runes
// This function is not exported (It is only visible inside the 'strings' package)
func reverseRunes(r []rune) []rune {
    for i, j := 0, len(r) - 1; i < j; i, j = i + 1, j - 1 {
        r[i], r[j] = r[j], r[i]
    }
    return r
}
```

**strings/greeting/texts.go (nested package)**

```go
// Nested Package
package greeting

// Exported
const  (
    WelcomeText = "Hello, World to Golang"
    MorningText = "Good Morning"
    EveningText = "Good Evening"
)

// Not exported (only visible inside the 'greeting' package)
var loremIpsumText = `Lorem ipsum dolor sit amet, consectetur adipiscing elit,
                      sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
                      Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris
                      nisi ut aliquip ex ea commodo consequat.`
```

**main.go (entry point of our program)**

```go
package main

import (
    "fmt"
    str "strings" // Package Alias
    "github.com/username/packer/numbers"
    "github.com/username/packer/strings"
    "github.com/username/packer/strings/greeting" // Importing a nested package
)

func main() {
    fmt.Println(numbers.IsPrime(19))

    fmt.Println(greeting.WelcomeText)

    fmt.Println(strings.Reverse("username"))

    fmt.Println(str.Count("Go is Awesome. I love Go", "Go"))
}
```

### Things to Note

- **Import Paths:** All import paths are relative to the module's path `github.com/username/packer`.

```go
import (
    "github.com/username/packer/numbers"
    "github.com/username/packer/strings"
    "github.com/username/packer/strings/greeting"
)
```

- **Package Alias:** You can use package alias to resolve conflicts between different packages of the same name, or just to give a short name to the imported package.

```go
import (
    str "strings"	// Package Alias
)
```

- **Nested Package:** You can nest a package inside another. It's as simple as creating a subdirectory.

```text
packer
    strings          # Package
        greeting     # Nested Package
            texts.go
```

A nested package can be imported similar to a root package. Just provide its path relative to the module's path `github.com/username/packer`.

```go
import (
    "github.com/callicoder/packer/strings/greeting"
)
```

## Adding 3rd Party Packages

Adding 3rd party packages to your project is very easy with Go modules. You can just import the package to any of the source files in your project, and the next time you build/run the project, Go automatically downloads it for you.

```go
package main

import (
    "fmt"
    "rsc.io/quote"
)

func main() {
    fmt.Println(quote.Go())
}
```

Go will also add this new dependency to the `go.mod` file.

### Manually Installing Packages

You can use `go get` command to download 3rd party packages from remote repositories.

```bash
go get -u github.com/jinzhu/gorm
```

The above command fetches the `gorm` package from Github and adds it as a dependency to your `go.mod` file.

That's it. You can now import and use the above package in your program like this:

```go
import "github.com/jinzhu/gorm"
```
