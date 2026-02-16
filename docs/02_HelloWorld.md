<div align="center">
  <h1>Hello World</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 16, 2026</sub>
</div>

## Hello World Program

When we start learning a new programming language, we typically start by writing the classic `"Hello, World"` program.

So let's write the `"Hello, World"` program in Go and understand how it works. Open your favorite text editor, create a new file named `hello.go`, and type in the following code

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World")
}
```

You can run the above program using `go run` command like so:

```bash
$ go run hello.go
Hello, World
```

The `go run` command does two things - It first compiles the program to machine code and then runs the compiled code.

If however, you want to produce a standalone binary executable from your Go source that can be run without using the Go tool, then use the `go build` command:

```bash
$ go build hello.go
$ ls
hello		hello.go
```

You may now run the built binary file like this:

```bash
$ ./hello
Hello, World
```

## How It Works?

When you write a program in Go you will have a `main` package defined with a `main` func inside it. Packages are ways of grouping up related Go code together.

The `func` keyword is how you define a function with a name and a body.

With `import "fmt"` we are importing a package which contains the `Println` function that we use to print.

## How To Test?

How do you test this? It is good to separate your "domain" code from the outside world (side-effects). The `fmt.Println` is a side effect (printing to stdout) and the string we send in is our domain.

So let's separate these concerns so it's easier to test.

```go
package main

import "fmt"

func Hello() string {
    return "Hello, World"
}

func main() {
    fmt.Println(Hello())
}
```

We have created a new function again with `func` but this time we've added another keyword `string` in the definition. This means this function returns a `string`.

Now create a new file called `hello_test.go` where we are going to write a test for our `Hello` function.

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello()
    want := "Hello, World"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

Before explaining, let's just run the code. Run `go test` in your terminal. It should've passed! Just to check, try deliberately breaking the test by changing the `want` string.

Notice how you have not had to pick between multiple testing frameworks and then figure out how to install. Everything you need is built in to the language and the syntax is the same as the rest of the code you will write.

### Writing Tests

Writing a test is just like writing a function, with a few rules:

- It needs to be in a file with a name like `xxx_test.go`.
- The test function must start with the word `Test`.
- The test function takes one argument only `t *testing.T`.

For now it's enough to know that your `t` of type `*testing.T` is your "hook" into the testing framework so you can do things like `t.Fail()` when you want to fail.

### `t.Errorf`

We are calling the `Errorf` method on our `t` which will print out a message and fail the test. The `f` stands for format which allows us to build a string with values inserted into the placeholder values `%q`. When you made the test fail it should be clear how it works.

You can read more about the placeholder strings in the [fmt go doc](https://golang.org/pkg/fmt/#hdr-Printing). For tests `%q` is very useful as it wraps your values in double quotes.

### Go doc

Another quality of life feature of Go is the documentation. You can launch the docs locally by running `godoc -http :8000`. If you go to [localhost:8000/pkg](localhost:8000/pkg) you will see all the packages installed on your system.

The vast majority of the standard library has excellent documentation with examples. Navigating to [http://localhost:8000/pkg/testing](http://localhost:8000/pkg/testing) would be worthwhile to see what's available to you.

If you don't have `godoc` command, then maybe you are using the newer version of Go (1.14 or later) which is no longer including `godoc`. You can manually install it with `go get golang.org/x/tools/cmd/godoc`.

### Hello, YOU

Now that we have a test we can iterate on our software safely.

In the last example we wrote the test after the code had been written just so you could get an example of how to write a test and declare a function. From this point on we will be writing tests first.

Our next requirement is to let us specify the recipient of the greeting.

Let's start by capturing these requirements in a test. This is basic test driven development and allows us to make sure our test is actually testing what we want. When you retrospectively write tests there is the risk that your test may continue to pass even if the code doesn't work as intended.

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello("Chris")
    want := "Hello, Chris"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

Now run `go test`, you should have a compilation error.

```go
./hello_test.go:6:18: too many arguments in call to Hello
    have (string)
    want ()
```

When using a statically typed language like Go it is important to listen to the compiler. The compiler understands how your code should snap together and work so you don't have to.

In this case the compiler is telling you what you need to do to continue. We have to change our function `Hello` to accept an argument.

Edit the `Hello` function to accept an argument of type `string`.

```go
func Hello(name string) string {
    return "Hello, World"
}
```

If you try and run your tests again your `hello_test.go` will fail to compile because you're not passing an argument. Send in `"World"` to make it pass.

```go
func main() {
    fmt.Println(Hello("World"))
}
```

Now when you run your tests you should see something like:

```bash
hello_test.go:10: got 'Hello, world' want 'Hello, Chris''
```

We finally have a compiling program but it is not meeting our requirements according to the test. Let's make the test pass by using the `name` argument and concatenate it with `Hello`.

```go
func Hello(name string) string {
    return "Hello, " + name
}
```

When you run the tests they should now pass. Normally as part of the TDD cycle we should now **refactor**. We can now refactor our code.

```go
const englishHelloPrefix = "Hello, "

func Hello(name string) string {
    return englishHelloPrefix + name
}
```

After refactoring, re-run your tests to make sure you haven't broken anything. Constants should improve performance of your application as it saves you creating the `"Hello, "` string instance every time Hello is called.

To be clear, the performance boost is incredibly negligible for this example! But it's worth thinking about creating constants to capture the meaning of values and sometimes to aid performance.

### Hello, World... Again

The next requirement is when our function is called with an empty string it defaults to printing `"Hello, World"`, rather than `"Hello, "`.

Start by writing a new failing test.

```go
func TestHello(t *testing.T) {

    t.Run("saying hello to people", func(t *testing.T) {
        got := Hello("Chris")
        want := "Hello, Chris"

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    })

    t.Run("say 'Hello, World' when an empty string is supplied", func(t *testing.T) {
        got := Hello("")
        want := "Hello, World"

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    })

}
```

Here we are introducing another tool in our testing arsenal, subtests. Sometimes it is useful to group tests around a "thing" and then have subtests describing different scenarios.

A benefit of this approach is you can set up shared code that can be used in the other tests.

There is repeated code when we check if the message is what we expect.

Refactoring is not just for the production code!

It is important that your tests are clear specifications of what the code needs to do.

We can and should refactor our tests.

```go
func TestHello(t *testing.T) {

    assertCorrectMessage := func(t *testing.T, got, want string) {
        t.Helper()
        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    }

    t.Run("saying hello to people", func(t *testing.T) {
        got := Hello("Chris")
        want := "Hello, Chris"
        assertCorrectMessage(t, got, want)
    })

    t.Run("empty string defaults to 'World'", func(t *testing.T) {
        got := Hello("")
        want := "Hello, World"
        assertCorrectMessage(t, got, want)
    })

}
```

What have we done here?

We've refactored our assertion into a function. This reduces duplication and improves readability of our tests. In Go you can declare functions inside other functions and assign them to variables. You can then call them, just like normal functions. We need to pass in `t *testing.T` so that we can tell the test code to fail when we need to.

`t.Helper()` is needed to tell the test suite that this method is a helper. By doing this when it fails the line number reported will be in our function call rather than inside our test helper. This will help other developers track down problems easier. If you still don't understand, comment it out, make a test fail and observe the test output.

Now that we have a well-written failing test, let's fix the code, using an `if`.

```go
const englishHelloPrefix = "Hello, "

func Hello(name string) string {
    if name == "" {
        name = "World"
    }
    return englishHelloPrefix + name
}
```

If we run our tests we should see it satisfies the new requirement and we haven't accidentally broken the other functionality.

### Discipline

Let's go over the cycle again

- Write a test
- Make the compiler pass
- Run the test, see that it fails and check the error message is meaningful
- Write enough code to make the test pass
- Refactor

On the face of it this may seem tedious but sticking to the feedback loop is important.

Not only does it ensure that you have relevant tests, it helps ensure you design good software by refactoring with the safety of tests.

Seeing the test fail is an important check because it also lets you see what the error message looks like. As a developer it can be very hard to work with a codebase when failing tests do not give a clear idea as to what the problem is.

By ensuring your tests are fast and setting up your tools so that running tests is simple you can get in to a state of flow when writing your code.

By not writing tests you are committing to manually checking your code by running your software which breaks your state of flow and you won't be saving yourself any time, especially in the long run.

## Keep Going! More Requirements

Goodness me, we have more requirements. We now need to support a second parameter, specifying the language of the greeting. If a language is passed in that we do not recognise, just default to English.

We should be confident that we can use TDD to flesh out this functionality easily!

Write a test for a user passing in Spanish. Add it to the existing suite.

```go
t.Run("in Spanish", func(t *testing.T) {
    got := Hello("Elodie", "Spanish")
    want := "Hola, Elodie"
    assertCorrectMessage(t, got, want)
})
```

Remember not to cheat! Test first. When you try and run the test, the compiler should complain because you are calling `Hello` with two arguments rather than one.

```bash
./hello_test.go:27:19: too many arguments in call to Hello
    have (string, string)
    want (string)
```

Fix the compilation problems by adding another string argument to `Hello`.

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }
    return englishHelloPrefix + name
}
```

When you try and run the test again it will complain about not passing through enough arguments to `Hello` in your other tests and in `hello.go`.

```bash
./hello.go:15:19: not enough arguments in call to Hello
    have (string)
    want (string, string)
```

Fix them by passing through empty strings. Now all your tests should compile and pass, apart from our new scenario.

```bash
hello_test.go:29: got 'Hello, Elodie' want 'Hola, Elodie'
```

We can use `if` here to check the language is equal to `"Spanish"` and if so change the message.

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == "Spanish" {
        return "Hola, " + name
    }

    return englishHelloPrefix + name
}
```

The tests should now pass.

Now it is time to refactor. You should see some problems in the code, "magic" strings, some of which are repeated. Try and refactor it yourself, with every change make sure you re-run the tests to make sure your refactoring isn't breaking anything.

```go
const spanish = "Spanish"
const englishHelloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "

func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == spanish {
        return spanishHelloPrefix + name
    }

    return englishHelloPrefix + name
}
```

### French

- Write a test asserting that if you pass in `"French"` you get `"Bonjour, "`.
- See it fail, check the error message is easy to read.
- Do the smallest reasonable change in the code.

You may have written something that looks roughly like this:

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == spanish {
        return spanishHelloPrefix + name
    }

    if language == french {
        return frenchHelloPrefix + name
    }

    return englishHelloPrefix + name
}
```

When you have lots of `if` statements checking a particular value it is common to use a `switch` statement instead. We can use `switch` to refactor the code to make it easier to read and more extensible if we wish to add more language support later.

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    prefix := englishHelloPrefix

    switch language {
    case french:
        prefix = frenchHelloPrefix
    case spanish:
        prefix = spanishHelloPrefix
    }

    return prefix + name
}
```

Write a test to now include a greeting in the language of your choice and you should see how simple it is to extend our amazing function.

You could argue that maybe our function is getting a little big. The simplest refactor for this would be to extract out some functionality into another function.

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
    switch language {
    case french:
        prefix = frenchHelloPrefix
    case spanish:
        prefix = spanishHelloPrefix
    default:
        prefix = englishHelloPrefix
    }
    return
}
```

A few new concepts:

- In our function signature we have made a named return value `(prefix string)`.
- This will create a variable called `prefix` in your function.
  - It will be assigned the "zero" value. This depends on the type, for example ints are `0` and for strings it is `""`.
    - You can return whatever it's set to by just calling `return` rather than `return prefix`.
  - This will display in the Go Doc for your function so it can make the intent of your code clearer.
- `default` in the switch case will be branched to if none of the other `case` statements match.
- The function name starts with a lowercase letter. In Go public functions start with a capital letter and private ones start with a lowercase. We don't want the internals of our algorithm to be exposed to the world, so we made this function private.
