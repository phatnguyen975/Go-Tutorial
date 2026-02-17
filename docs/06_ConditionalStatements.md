<div align="center">
  <h1>Conditional Statements</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 17, 2026</sub>
</div>

## If Statement

The `if` statements are used to specify whether a block of code should be executed or not depending on a given condition.

Following is the syntax of `if` statements in Go:

```go
if (condition) {
    // Code to be executed if the condition is true
}
```

**Note:** You can omit the parentheses `()` from an `if` statement in Go, but the curly braces `{}` are mandatory.

```go
var y = -1

if y < 0 {
    fmt.Printf("%d is negative", y)
}
```

You can combine multiple conditions using short circuit operators `&&` and `||` like so:

```go
var age = 21

if age >= 17 && age <= 30 {
    fmt.Println("My age is between 17 and 30")
}
```

**Note:** There is no **ternary if** in Go, so you'll need to use a full `if` statement even for basic conditions.

## If-Else Statement

An `if` statement can be combined with an `else` block. The `else` block is executed if the condition specified in the `if` statement is false.

```go
if condition {
    // Code to be executed if the condition is true
} else {
    // Code to be executed if the condition is false
}
```

Here is a simple example:

```go
package main

import "fmt"

func main() {
    var age = 18

    if age >= 18 {
        fmt.Println("You're eligible to vote!")
    } else {
        fmt.Println("You're not eligible to vote!")
    }
}
```

## If-Else-If Statement

The `if` statements can also have multiple `else if` parts making a chain of conditions like this:

```go
package main

import "fmt"

func main() {
    var BMI = 21.0

    if BMI < 18.5 {
        fmt.Println("You are underweight");
    } else if BMI >= 18.5 && BMI < 25.0 {
        fmt.Println("Your weight is normal");
    } else if BMI >= 25.0 && BMI < 30.0 {
        fmt.Println("You're overweight")
    } else {
        fmt.Println("You're obese")
    }
}
```

## If with a Short Statement

An `if` statement in Go can also contain a [short declaration statement](https://www.callicoder.com/golang-variables-zero-values-type-inference/#short-declaration) preceding the conditional expression:

```go
if n := 10; n % 2 == 0 {
    fmt.Printf("%d is even", n)
}
```

The variable declared in the short statement is only available inside the `if` block and it's `else` or `else-if` branches.

```go
if n := 15; n % 2 == 0 {
    fmt.Printf("%d is even", n)
} else {
    fmt.Printf("%d is odd", n)
}
```

**Note:** If you're using a short statement, then you can't use parentheses. So the following code will generate a syntax :

```go
if (n := 15; n % 2 == 0) { // Syntax error

}
```

## Switch Statement

A `switch` statement takes an expression and matches it against a list of possible cases. Once a match is found, it executes the block of code specified in the matched case.

Here is a simple example of `switch` statement:

```go
package main

import "fmt"

func main() {
    var dayOfWeek = 6

    switch dayOfWeek {
        case 1: fmt.Println("Monday")
        case 2: fmt.Println("Tuesday")
        case 3: fmt.Println("Wednesday")
        case 4: fmt.Println("Thursday")
        case 5: fmt.Println("Friday")
        case 6: {
            fmt.Println("Saturday")
            fmt.Println("Weekend. Yaay!")
        }
        case 7: {
            fmt.Println("Sunday")
            fmt.Println("Weekend. Yaay!")
        }
        default: fmt.Println("Invalid day")
    }
}
```

Go evaluates all the switch cases one by one from top to bottom until a case succeeds. Once a case succeeds, it runs the block of code specified in that case and then stops (it doesn’t evaluate any further cases).

This is contrary to other languages like `C`, `C++`, and `Java`, where you explicitly need to insert a `break` statement after the body of every case to stop the evaluation of cases that follow.

If none of the cases succeed, then the `default` case is executed.

### Switch with a Short Statement

Just like `if`, `switch` can also contain a short declaration statement preceding the conditional expression. So you could also write the previous `switch` example like this:

```go
switch dayOfWeek := 6; dayOfWeek {
    case 1: fmt.Println("Monday")
    case 2: fmt.Println("Tuesday")
    case 3: fmt.Println("Wednesday")
    case 4: fmt.Println("Thursday")
    case 5: fmt.Println("Friday")
    case 6: {
        fmt.Println("Saturday")
        fmt.Println("Weekend. Yaay!")
    }
    case 7: {
        fmt.Println("Sunday")
        fmt.Println("Weekend. Yaay!")
    }
    default: fmt.Println("Invalid day")
}
```

The only difference is that the variable declared by the short statement (`dayOfWeek`) is only available inside the `switch` block.

### Combining Multiple Switch Cases

You can combine multiple `switch` cases into one like so:

```go
package main

import "fmt"

func main() {
    switch dayOfWeek := 5; dayOfWeek {
        case 1, 2, 3, 4, 5:
            fmt.Println("Weekday")
        case 6, 7:
            fmt.Println("Weekend")
        default:
            fmt.Println("Invalid Day")
    }
}
```

This comes handy when you need to run a common logic for multiple cases.

### Switch with No Expression

In Go, the expression that we specify in the `switch` statement is optional. A `switch` statement without an expression is same as `switch true`. It evaluates all the cases one by one, and runs the first case that evaluates to `true`.

```go
package main

import "fmt"

func main() {
	var BMI = 21.0

	switch {
		case BMI < 18.5:
			fmt.Println("You're underweight")
		case BMI >= 18.5 && BMI < 25.0:
			fmt.Println("Your weight is normal")
		case BMI >= 25.0 && BMI < 30.0:
			fmt.Println("You're overweight")
		default:
			fmt.Println("You're obese")
	}
}
```

Switch without an expression is simply a concise way of writing `if-else-if` chains.
