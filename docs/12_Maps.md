<div align="center">
  <h1>Maps</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 18, 2026</sub>
</div>

A map is an unordered collection of key-value pairs. It maps keys to values. The keys are **unique** within a map while the values may not be.

The map data structure is used for fast lookups, retrieval, and deletion of data based on keys. It is one of the most used data structures in computer science.

## Declaring a Map

A map is declared using the following syntax:

```go
var m map[KeyType]ValueType
```

For example, Here is how you can declare a map of `string` keys to `int` values:

```go
var m map[string]int
```

The zero value of a map is `nil`. A `nil` map has no keys. Moreover, any attempt to add keys to a `nil` map will result in a runtime error.

Let's see an example:

```go
package main

import "fmt"

func main() {
    var m map[string]int
    fmt.Println(m)
    if m == nil {
        fmt.Println("m is nil")
    }

    // Attempting to add keys to a nil map will result in a runtime error
    // m["one hundred"] = 100
}
```

If you uncomment the statement `m["one hundred"] = 100`, the program will generate the following error:

```bash
panic: assignment to entry in nil map
```

## Initializing a Map

### Initializing a Map Using The Built-in `make()` Function

You can initialize a map using the built-in `make()` function. You just need to pass the type of the map to the `make()` function as in the example below. The function will return an initialized and ready to use map.

```go
var m = make(map[string]int)
```

Let's see a complete example:

```go
package main

import "fmt"

func main() {
    var m = make(map[string]int)

    fmt.Println(m)

    if m == nil {
        fmt.Println("m is nil")
    } else {
        fmt.Println("m is not nil")
    }

    // make() function returns an initialized and ready to use map
    // Since it is initialized, you can add new keys to it
    m["one hundred"] = 100
    fmt.Println(m)
}
```

### Initializing a Map Using Literal

A map literal is a very convenient way to initialize a map with some data. You just need to pass the key-value pairs separated by colon inside curly braces like this:

```go
var m = map[string]int{
    "one": 1,
    "two": 2,
    "three": 3,
}
```

**Note:** The last trailing comma is necessary, otherwise, you’ll get a compiler error.

Let's check out a complete example:

```go
package main

import "fmt"

func main() {
    var m = map[string]int{
        "one":   1,
        "two":   2,
        "three": 3,
        "four":  4,
        "five":  5,
    }

  	fmt.Println(m)
}
```

You can also create an empty map using a map literal by leaving the curly braces empty.

```go
var m = map[string]int{}
```

The above statement is functionally identical to using the `make()` function.

## Adding Items (key-value pairs) to a Map

You can add new items to an initialized map using the following syntax:

```go
m[key] = value
```

The following example initializes a map using the `make()` function and adds some new items to it:

```go
package main

import "fmt"

func main() {
    // Initializing a map
    var tinderMatch = make(map[string]string)

    // Adding keys to a map
    tinderMatch["Rajeev"] = "Angelina"
    tinderMatch["James"] = "Sophia"
    tinderMatch["David"] = "Emma"

    fmt.Println(tinderMatch)

    /*
      Adding a key that already exists will simply override
      the existing key with the new value
    */
    tinderMatch["Rajeev"] = "Jennifer"
    fmt.Println(tinderMatch)
}
```

**Note:** If you try to add a key that already exists in the map, then it will simply be overridden by the new value.

## Retrieving the Value Associated with a Given Key in a Map

You can retrieve the value assigned to a key in a map using the syntax `m[key]`. If the key exists in the map, you'll get the assigned value. Otherwise, you'll get the zero value of the map's value type.

Let's check out an example to understand this:

```go
package main

import "fmt"

func main() {
    var personMobileNo = map[string]string{
        "John":  "+33-8273658526",
        "Steve": "+1-8579822345",
        "David": "+44-9462834443",
    }

    var mobileNo = personMobileNo["Steve"]
    fmt.Println("Steve's Mobile No: ", mobileNo)

    // If a key doesn't exist in the map, we get the zero value of the value type
    mobileNo = personMobileNo["Jack"]
    fmt.Println("Jack's Mobile No: ", mobileNo)
}
```

In the above example, since the key `"Jack"` doesn't exist in the map, we get the zero value of the map's value type. Since the map's value type is `string`, we get `""`.

Unlike other languages, we do not get a runtime error in Go if the key doesn't exist in the map.

But what if you want to check for the existence of a key? In the above example, the map would return `""` even if the key `"Jack"` existed with the value `""`. So how do we distinguish between cases where a key exists with the value equal to the zero value of the value type, and the absence of a key?

### Checking If a Key Exists in a Map

When you retrieve the value assigned to a given key using the syntax `map[key]`, it returns an additional boolean value as well which is true if the key exists in the map, and false if it doesn't exist.

So you can check for the existence of a key in a map by using the following two-value assignment:

```go
value, ok := m[key]
```

The boolean variable ok will be `true` if the key exists, and `false` otherwise.

Consider the following map for example. It maps employeeIds to names:

```go
var employees = map[int]string{
    1001: "Rajeev",
    1002: "Sachin",
    1003: "James",
}
```

Accessing the key `1001` will return `"Rajeev"` and `true`, since the key `1001` exists in the map:

```go
name, ok := employees[1001] // "Rajeev", true
```

However, if you try to access a key that doesn't exist, then the map will return an empty string `""` (zero value of strings), and `false`.

```go
name, ok := employees[1010] // "", false
```

If you just want to check for the existence of a key without retrieving the value associated with that key, then you can use an `_` (underscore) in place of the first value.

```go
_, ok := employees[1005]
```

Now let's check out a complete example:

```go
package main

import "fmt"

func main() {
    var employees = map[int]string{
        1001: "John",
        1002: "Steve",
        1003: "Maria",
    }

    printEmployee(employees, 1001)
    printEmployee(employees, 1010)

    if isEmployeeExists(employees, 1002) {
        fmt.Println("EmployeeId 1002 found")
    }
}

func printEmployee(employees map[int]string, employeeId int) {
    if name, ok := employees[employeeId]; ok {
        fmt.Printf("name = %s, ok = %v\n", name, ok)
    } else {
        fmt.Printf("EmployeeId %d not found\n", employeeId)
    }
}

func isEmployeeExists(employees map[int]string, employeeId int) bool {
    _, ok := employees[employeeId]
    return ok
}
```

In the above example, I've used a short declaration in the `if` statement to initialize the `name` and `ok` values, and then test the boolean value `ok`. It makes the code more concise.

## Deleting a Key from a Map

You can delete a key from a map using the built-in `delete()` function. The syntax looks like this:

```go
delete(map, key)
```

The `delete()` function doesn't return any value. Also, it doesn't do anything if the key doesn't exist in the map.

Here is a complete example:

```go
package main

import "fmt"

func main() {
    var fileExtensions = map[string]string{
        "Python": ".py",
        "C++":    ".cpp",
        "Java":   ".java",
        "Golang": ".go",
        "Kotlin": ".kt",
    }

    fmt.Println(fileExtensions)

    delete(fileExtensions, "Kotlin")

    // delete() doesn't do anything if the key doesn't exist
    delete(fileExtensions, "Javascript")

    fmt.Println(fileExtensions)
}
```

To **remove all** key/value pairs from a map, use the `clear()` builtin.

```go
clear(m)
fmt.Println("map:", m)
```

## Maps are Reference Types

Maps are reference types. When you assign a map to a new variable, they both refer to the same underlying data structure. Therefore changes done by one variable will be visible to the other.

```go
package main

import "fmt"

func main() {
    var m1 = map[string]int{
        "one":   1,
        "two":   2,
        "three": 3,
        "four":  4,
        "five":  5,
    }

    var m2 = m1
    fmt.Println("m1 = ", m1)
    fmt.Println("m2 = ", m2)

    m2["ten"] = 10
    fmt.Println("\nm1 = ", m1)
    fmt.Println("m2 = ", m2)
}
```

**Note:** The same concept applies when you pass a map to a function. Any changes done to the map inside the function is also visible to the caller.

## Iterating over a Map

You can iterate over a map using `range` form of the `for` loop. It gives you the key, value pair in every iteration.

```go
package main

import "fmt"

func main() {
    var personAge = map[string]int{
        "Rajeev": 25,
        "James":  32,
        "Sarah":  29,
    }

    for name, age := range personAge {
        fmt.Println(name, age)
    }
}
```

**Note:** A map is an **unordered** collection and therefore the iteration order of a map is not guaranteed to be the same every time you iterate over it.
