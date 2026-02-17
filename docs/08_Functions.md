<div align="center">
  <h1>Functions</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 17, 2026</sub>
</div>

## Declaring and Calling Functions

In Go, we declare a function using the `func` keyword. A function has a name, a list of comma-separated input parameters along with their types, the result type(s), and a body.

Following is an example of a simple function called `avg` that takes two input parameters of type `float64` and returns the average of the inputs. The result is also of type `float64`.

```go
func avg(x float64, y float64) float64 {
	return (x + y) / 2
}
```

Now, calling a function is very simple. You just need to pass the required number of parameters to the function like this:

```go
avg(6.56, 13.44)
```

Here is a complete example:

```go
package main

import "fmt"

func avg(x float64, y float64) float64 {
	return (x + y) / 2
}

func main() {
	x := 5.75
	y := 6.25

	result := avg(x, y)
	fmt.Printf("Average of %.2f and %.2f = %.2f", x, y, result)
}
```

**Function parameters and return type(s) are optional**

The input parameters and return type(s) are optional for a function. A function can be declared without any input and output.

The `main()` function is an example of such a function. Here is another example.

```go
func sayHello() {
	fmt.Println("Hello, World")
}
```

**Note:** In Go, a function must explicitly declare its return type if it returns any value. If the function does not return anything, no return type is specified. Go does not use the `void` keyword.

**You need to specify the type only once for multiple consecutive parameters of the same type**

If a function has two or more consecutive parameters of the same type, then it suffices to specify the type only once for the last parameter of that type.

For example, we can declare the `avg` function that we saw in the previous section like this as well:

```go
func avg(x, y float64) float64 { }
// Same as - func avg(x float64, y float64) float64 { }

func printPersonDetails(firstName, lastName string, age int) { }
// Same as - func printPersonDetails(firstName string, lastName string, age int) { }
```

## Functions with Multiple Return Values

Go functions are capable of returning multiple values. That's right! This is something that most programming languages don't support natively. But Go is different.

Let's say that you want to create a function that takes the **previous price** and the **current price** of a stock, and returns the **amount** by which the price has changed and the **percentage** of change.

Here is how you can implement such a function in Go.

```go
func getStockPriceChange(prevPrice, currentPrice float64) (float64, float64) {
	change := currentPrice - prevPrice
	percentChange := (change / prevPrice) * 100
	return change, percentChange
}
```

Simple! isn't it? You just need to specify the return types separated by comma inside parentheses, and then return multiple comma-separated values from the function.

Let's see a complete example with the `main()` function.

```go
package main

import (
	"fmt"
	"math"
)

func getStockPriceChange(prevPrice, currentPrice float64) (float64, float64) {
	change := currentPrice - prevPrice
	percentChange := (change / prevPrice) * 100
	return change, percentChange
}

func main() {
	prevStockPrice := 75000.0
	currentStockPrice := 100000.0

	change, percentChange := getStockPriceChange(prevStockPrice, currentStockPrice)

	if change < 0 {
		fmt.Printf("The Stock Price decreased by $%.2f which is %.2f%% of the prev price", math.Abs(change), math.Abs(percentChange))
	} else {
		fmt.Printf("The Stock Price increased by $%.2f which is %.2f%% of the prev price", change, percentChange)
	}
}
```

## Return an Error Value from a Function

Multiple return values are often used in Go to return an error from the function along with the result.

Let's see an example - The `getStockPriceChange` function that we saw in the previous section will return `±Inf` (Infinity) if the `prevPrice` is `0`. If you want to return an error instead, you can do so by adding another return value of type `error` and return the error value like so:

```go
func getStockPriceChangeWithError(prevPrice, currentPrice float64) (float64, float64, error) {
	if prevPrice == 0 {
		err := errors.New("Previous price cannot be zero")
		return 0, 0, err
	}
	change := currentPrice - prevPrice
	percentChange := (change / prevPrice) * 100
	return change, percentChange, nil
}
```

The `error` type is a built-in type in Go. Go programs use error values to indicate an abnormal situation. Don’t worry if you don't understand about `errors` for now. You'll learn more about error handling in a future chapter.

Following is a complete example demonstrating the above concept with a `main()` function.

```go
package main

import (
	"errors"
	"fmt"
	"math"
)

func getStockPriceChangeWithError(prevPrice, currentPrice float64) (float64, float64, error) {
	if prevPrice == 0 {
		err := errors.New("Previous price cannot be zero")
		return 0, 0, err
	}
	change := currentPrice - prevPrice
	percentChange := (change / prevPrice) * 100
	return change, percentChange, nil
}

func main() {
	prevStockPrice := 0.0
	currentStockPrice := 100000.0

	change, percentChange, err := getStockPriceChangeWithError(prevStockPrice, currentStockPrice)

	if err != nil {
		fmt.Println("Sorry! There was an error: ", err)
	} else {
		if change < 0 {
			fmt.Printf("The Stock Price decreased by $%.2f which is %.2f%% of the prev price", math.Abs(change), math.Abs(percentChange))
		} else {
			fmt.Printf("The Stock Price increased by $%.2f which is %.2f%% of the prev price", change, percentChange)
		}
	}
}
```

## Functions with Named Return Values

The return values of a function in Go may be named. Named return values behave as if you defined them at the top of the function.

Let's rewrite the `getStockPriceChange` function that we saw in the previous section with named return values.

```go
func getNamedStockPriceChange(prevPrice, currentPrice float64) (change, percentChange float64) {
	change = currentPrice - prevPrice
	percentChange = (change / prevPrice) * 100
	return change, percentChange
}
```

Notice how we changed `:=` (short declarations) with `=` (assignments) in the function body. This is because Go itself defines all the named return values and makes them available for use in the function. Since they are already defined, you can't define them again using short declarations.

Named return values allow you to use the so-called **Naked return** (a return statement without any argument). When you specify a return statement without any argument, it returns the named return values by default. So you can write the above function like this as well:

```go
func getNamedStockPriceChange(prevPrice, currentPrice float64) (change, percentChange float64) {
	change = currentPrice - prevPrice
	percentChange = (change / prevPrice) * 100
	return
}
```

Let's use the above function in a complete example with the `main()` function and verify the output:

```go
package main

import (
	"fmt"
	"math"
)

func getNamedStockPriceChange(prevPrice, currentPrice float64) (change, percentChange float64) {
	change = currentPrice - prevPrice
	percentChange = (change / prevPrice) * 100
	return
}

func main() {
	prevStockPrice := 100000.0
	currentStockPrice := 90000.0

	change, percentChange := getNamedStockPriceChange(prevStockPrice, currentStockPrice)

	if change < 0 {
		fmt.Printf("The Stock Price decreased by $%.2f which is %.2f%% of the prev price", math.Abs(change), math.Abs(percentChange))
	} else {
		fmt.Printf("The Stock Price increased by $%.2f which is %.2f%% of the prev price", change, percentChange)
	}
}
```

Named return values improve the readability of your functions. Using meaningful names would let the consumers of your function know what the function will return just by looking at its signature.

The naked return statements are good for short functions. But don't use them if your functions are long. They can harm the readability. You should explicitly specify the return values in longer functions.

## Blank Identifier

Sometimes you may want to ignore some of the results from a function that returns multiple values.

For example, let's say that you're using the `getStockPriceChange` function that we defined in the previous section, but you're only interested in the amount of change, not the percentage change.

Now, you might just declare local variables and store all the values returned from the function like this:

```go
change, percentChange := getStockPriceChange(prevStockPrice, currentStockPrice)
```

But in that case, you'll be forced to use the `percentChange` variable because Go doesn't allow creating variables that you never use.

So what's the solution? Well, you can use a blank identifier instead:

```go
change, _ := getStockPriceChange(prevStockPrice, currentStockPrice)
```

The blank identifier is used to tell Go that you don't need this value. The following example demonstrates this concept:

```go
package main

import (
	"fmt"
	"math"
)

func getStockPriceChange(prevPrice, currentPrice float64) (float64, float64) {
	change := currentPrice - prevPrice
	percentChange := (change / prevPrice) * 100
	return change, percentChange
}

func main() {
	prevStockPrice := 80000.0
	currentStockPrice := 120000.0

	change, _ := getStockPriceChange(prevStockPrice, currentStockPrice)

	if change < 0 {
		fmt.Printf("The Stock Price decreased by $%.2f", math.Abs(change))
	} else {
		fmt.Printf("The Stock Price increased by $%.2f", change)
	}
}
```

## Variadic Functions

[Variadic functions](https://en.wikipedia.org/wiki/Variadic_function) can be called with any number of trailing arguments. For example, `fmt.Println` is a common variadic function.

Here's a function that will take an arbitrary number of ints as arguments.

```go
package main

import "fmt"

func sum(nums ...int) {
    fmt.Print(nums, " ")
    total := 0

    for _, num := range nums {
        total += num
    }
    fmt.Println(total)
}

func main() {
    // Variadic functions can be called in the usual way with individual arguments
    sum(1, 2)
    sum(1, 2, 3)

    // If you already have multiple args in a slice, apply them to a variadic function using 'func(slice...)' like this
    nums := []int{1, 2, 3, 4}
    sum(nums...)
}
```

Within the function, the type of `nums` is equivalent to `[]int`. We can call `len(nums)`, iterate over it with `range`, etc.

## Closures

Go supports [anonymous functions](https://en.wikipedia.org/wiki/Anonymous_function), which can form [closures](<https://en.wikipedia.org/wiki/Closure_(computer_science)>). Anonymous functions are useful when you want to define a function inline without having to name it.

```go
package main

import "fmt"

func intSeq() func() int {
    i := 0
    return func() int {
        i++
        return i
    }
}

func main() {
    nextInt := intSeq()

    // See the effect of the closure by calling 'nextInt' a few times
    fmt.Println(nextInt())
    fmt.Println(nextInt())
    fmt.Println(nextInt())

    // To confirm that the state is unique to that particular function, create and test a new one
    newInts := intSeq()
    fmt.Println(newInts())
}
```

This function `intSeq` returns another function, which we define anonymously in the body of `intSeq`. The returned function closes over the variable `i` to form a closure.

We call `intSeq`, assigning the result (a function) to `nextInt`. This function value captures its own `i` value, which will be updated each time we call `nextInt`.

## Recursion

Go supports recursive functions. Here's a classic example.

```go
package main

import "fmt"

// This fact function calls itself until it reaches the base case of fact(0)
func fact(n int) int {
    if n == 0 {
        return 1
    }
    return n * fact(n - 1)
}

func main() {
    fmt.Println(fact(7))

    // Anonymous functions can also be recursive, but this requires explicitly
    // declaring a variable with 'var' to store the function before it's defined
    var fib func(n int) int

    fib = func(n int) int {
        if n < 2 {
            return n
        }

        // Since 'fib' was previously declared, Go knows which function to call with 'fib' here
        return fib(n - 1) + fib(n - 2)
    }

    fmt.Println(fib(7))
}
```
