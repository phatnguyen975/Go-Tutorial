<div align="center">
  <h1>Defer</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>June 28, 2026</sub>
</div>

## What is Defer?

Defer is used to ensure that a function call is performed later in a program’s execution, usually for purposes of cleanup. defer is often used where e.g. `ensure` and `finally` would be used in other languages.

The definition might seem complex but it's pretty simple to understand by means of an example.

```go
package main

import (
    "fmt"
    "time"
)

func totalTime(start time.Time) {
    fmt.Printf("Total time taken %f seconds", time.Since(start).Seconds())
}

func test() {
    start := time.Now()
    defer totalTime(start)
    time.Sleep(2 * time.Second)
    fmt.Println("Sleep complete")
}

func main() {
    test()
}
```

The above is a simple program which illustrates the use of `defer`. In the above program, `defer` is used to find out the total time taken for the execution of the `test()` function. The start time of the `test()` function execution is passed as argument to `defer totalTime(start)`. This defer call is executed just before `test()` returns. `totalTime` prints the difference between `start` and the current time using `time.Since`. To simulate some computation happening in `test()`, a 2 second `sleep` is added.

```bash
Sleep complete
Total time taken 2.000000 seconds
```

The output correlates to the 2 second sleep added. Before the `test()` function returns, `totalTime` is called and it prints the total time taken for `test()` to execute.

## Arguments Evaluation

The arguments of a deferred function are evaluated when the `defer` statement is executed and not when the actual function call is done.

Let's understand this by means of an example.

```go
package main

import "fmt"

func displayValue(a int) {
	fmt.Println("value of a in deferred function", a)
}

func main() {
	a := 5
	defer displayValue(a)
	a = 10
	fmt.Println("value of a before deferred function call", a)
}
```

In the program above `a` initially has a value of `5`. When the `defer` statement is executed, the value of `a` is `5` and hence this will be the argument to the `displayValue` function which is deferred. We change the value of `a` to `10`. The next line prints the value of `a`. This program outputs,

```bash
value of a before deferred function call 10
value of a in deferred function 5
```

From the above output it can be understood that although the value of `a` changes to `10` after the defer statement is executed, the actual deferred function call `displayValue(a)` still prints `5`.

## Deferred Methods

Defer is not restricted only to functions. It is perfectly legal to defer a method call too. Let’s write a small program to test this.

```go
package main

import "fmt"

type person struct {
	firstName string
	lastName string
}

func (p person) fullName() {
	fmt.Printf("%s %s",p.firstName,p.lastName)
}

func main() {
	p := person {
		firstName: "John",
		lastName: "Smith",
	}
	defer p.fullName()
	fmt.Printf("Welcome ")
}
```

In the above program we have deferred a method call. The rest of the program is self explanatory. This program outputs,

```bash
Welcome John Smith
```

## Multiple defer calls are placed in stack

When a function has multiple defer calls, they are pushed to a stack and executed in Last In First Out (LIFO) order.

We will write a small program which prints a string in reverse using a stack of defers.

```go
package main

import "fmt"

func main() {
	str := "Gopher"
	fmt.Printf("Original String: %s\n", string(str))
	fmt.Printf("Reversed String: ")
	for _, v := range str {
		defer fmt.Printf("%c", v)
	}
}
```

In the program above, the `for range` loop, iterates the string and calls `defer fmt.Printf("%c", v)`. These deferred calls will be added to a stack.

```bash
Original String: Gopher
Reversed String: rehpoG
```

## Practical Use of `defer`

Defer is used in places where a function call should be executed irrespective of the code flow. Let’s understand this with the example of a program which makes use of `WaitGroup`. We will first write the program without using `defer` and then we will modify it to use `defer` and understand how useful defer is.

```go
package main

import (
	"fmt"
	"sync"
)

type rect struct {
	length int
	width  int
}

func (r rect) area(wg *sync.WaitGroup) {
	if r.length < 0 {
		fmt.Printf("rect %v's length should be greater than zero\n", r)
		wg.Done()
		return
	}
	if r.width < 0 {
		fmt.Printf("rect %v's width should be greater than zero\n", r)
		wg.Done()
		return
	}
	area := r.length * r.width
	fmt.Printf("rect %v's area %d\n", r, area)
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	r1 := rect{-67, 89}
	r2 := rect{5, -67}
	r3 := rect{8, 9}
	rects := []rect{r1, r2, r3}
	for _, v := range rects {
		wg.Add(1)
		go v.area(&wg)
	}
	wg.Wait()
	fmt.Println("All go routines finished executing")
}
```

In the program above, we have created a `rect` struct and a method `area` on `rect` which calculates the area of the rectangle. This method checks whether the length and width of the rectangle is less than zero. If it is so, it prints a corresponding error message else it prints the area of the rectangle.

The `main` function creates 3 variables `r1`, `r2` and `r3` of type `rect`. They are then added to the `rects` slice. This slice is then iterated using a `for range` loop and the `area` method is called as a concurrent `Goroutine`. The `WaitGroup wg` is used to ensure that the main function is waiting until all Goroutines finish executing. This `WaitGroup` is passed to the `area` method as an argument and the `area` method calls `wg.Done()` to notify the main function that the `Goroutine` is done with its job. If you notice closely, you can see that these calls happen just before the `area` method returns. `wg.Done()` should be called before the method returns irrespective of the path the code flow takes and hence these calls can be effectively replaced by a single `defer` call.

Let’s rewrite the above program using defer.

In the program below, we have removed the 3 `wg.Done()` calls in the above program and replaced it with a single defer `wg.Done()` call. This makes the code more simple and readable.

```go
package main

import (
	"fmt"
	"sync"
)

type rect struct {
	length int
	width  int
}

func (r rect) area(wg *sync.WaitGroup) {
	defer wg.Done()
	if r.length < 0 {
		fmt.Printf("rect %v's length should be greater than zero\n", r)
		return
	}
	if r.width < 0 {
		fmt.Printf("rect %v's width should be greater than zero\n", r)
		return
	}
	area := r.length * r.width
	fmt.Printf("rect %v's area %d\n", r, area)
}

func main() {
	var wg sync.WaitGroup
	r1 := rect{-67, 89}
	r2 := rect{5, -67}
	r3 := rect{8, 9}
	rects := []rect{r1, r2, r3}
	for _, v := range rects {
		wg.Add(1)
		go v.area(&wg)
	}
	wg.Wait()
	fmt.Println("All go routines finished executing")
}
```

This program outputs,

```bash
rect {8 9}'s area 72
rect {-67 89}'s length should be greater than zero
rect {5 -67}'s width should be greater than zero
All go routines finished executing
```

There is one more advantage of using defer in the above program. Let’s say we add another return path to the `area` method using a new `if` condition. If the call to `wg.Done()` was not deferred, we have to be careful and ensure that we call `wg.Done()` in this new return path. But since the call to `wg.Done()` is defered, we need not worry about adding new return paths to this method.
