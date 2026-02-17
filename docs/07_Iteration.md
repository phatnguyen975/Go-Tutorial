<div align="center">
  <h1>Iteration</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 17, 2026</sub>
</div>

## For Loop

To do stuff repeatedly in Go, you'll need `for`. In Go there are no `while`, `do`, `until` keywords, you can only use `for`.

```go
for initialization; condition; increment {
	// Body
}
```

The `initialization` statement is executed exactly once before the first iteration of the loop. In each iteration, the `condition` is checked. If the `condition` evaluates to `true` then the body of the loop is executed, otherwise, the loop terminates. The `increment` statement is executed at the end of every iteration.

```go
package main

import "fmt"

func main() {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", i)
	}
}
```

Unlike other languages like `C`, `C++`, and `Java`, Go's `for` loop doesn't contain parentheses, and the curly braces are mandatory.

**Note:** Both `initialization` and `increment` statements in the `for` loop are optional and can be omitted.

- Omitting the `initialization` statement:

```go
package main

import "fmt"

func main() {
    i := 2
    for ;i <= 10; i += 2 {
        fmt.Printf("%d ", i)
    }
}
```

- Omitting the `increment` statement:

```go
package main

import "fmt"

func main() {
	i := 2
	for ;i <= 20; {
		fmt.Printf("%d ", i)
		i *= 2
	}
}
```

**Note:** You can also omit the semicolons from the `for` loop in the above example and write it like this:

```go
package main

import "fmt"

func main() {
	i := 2
	for i <= 20 {
		fmt.Printf("%d ", i)
		i *= 2
	}
}
```

The above `for` loop is similar to a `while` loop in other languages. Go doesn't have a `while` loop because we can easily represent a `while` loop using `for`.

- You can also omit the `condition` from the `for` loop in Go. This will give you an **infinite** loop.

```go
package main

import "fmt"

func main() {
	for {
        fmt.Println("loop")
        break
	}
}
```

- Another way of accomplishing the basic **"do this N times"** iteration is `range` over an integer.

```go
package main

import "fmt"

func main() {
    for i := range 3 {
        fmt.Println("range", i)
    }
}
```

## `break` Statement

You can use `break` statement to break out of a loop before its normal termination.

```go
package main

import "fmt"

func main() {
	for num := 1; num <= 100; num++ {
		if num % 3 == 0 && num % 5 == 0 {
			fmt.Printf("First positive number divisible by both 3 and 5 is %d", num)
			break
		}
	}
}
```

## `continue` Statement

The `continue` statement is used to stop running the loop body midway and continue to the next iteration of the loop.

```go
package main

import "fmt"

func main() {
	for num := 1; num <= 10; num++ {
		if num % 2 == 0 {
			continue;
		}
		fmt.Printf("%d ", num)
	}
}
```
