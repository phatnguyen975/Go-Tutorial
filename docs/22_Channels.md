<div align="center">
  <h1>Channels</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 17, 2026</sub>
</div>

## Table of Contents

1. [What Is a Channel?](#1-what-is-a-channel)
2. [Declaring and Creating Channels](#2-declaring-and-creating-channels)
3. [Sending and Receiving](#3-sending-and-receiving)
4. [Unbuffered vs. Buffered Channels](#4-unbuffered-vs-buffered-channels)
5. [The Four Channel Axioms](#5-the-four-channel-axioms)
6. [Closing Channels](#6-closing-channels)
7. [Who Should Close a Channel?](#7-who-should-close-a-channel)
8. [Ranging Over a Channel](#8-ranging-over-a-channel)
9. [Directional (Typed-Direction) Channels](#9-directional-typed-direction-channels)
10. [The `select` Statement](#10-the-select-statement)
11. [Nil Channels: A Surprisingly Useful Tool](#11-nil-channels-a-surprisingly-useful-tool)
12. [Channels as Synchronization Primitives](#12-channels-as-synchronization-primitives)
13. [Common Channel Patterns](#13-common-channel-patterns)
14. [Channels vs. Mutexes: When to Use Which](#14-channels-vs-mutexes-when-to-use-which)
15. [Common Mistakes and Pitfalls](#15-common-mistakes-and-pitfalls)
16. [Best Practices Summary](#16-best-practices-summary)
17. [Use Case Summary Table](#17-use-case-summary-table)
18. [References](#18-references)

## 1. What Is a Channel?

A **channel** is a built-in Go type that provides a typed conduit through which goroutines send and receive values, and — critically — synchronize their execution around those sends and receives. Channels are the concrete embodiment of Go's concurrency philosophy, rooted in Tony Hoare's Communicating Sequential Processes (CSP) model: instead of protecting shared memory with locks, goroutines pass ownership of data to one another through a channel, encapsulated in the phrase _"Don't communicate by sharing memory; share memory by communicating."_

A channel is declared with a specific element type, and only values of that type can flow through it:

```go
var ch chan int       // a channel that carries int values
var strCh chan string // a channel that carries string values
```

## 2. Declaring and Creating Channels

A `var` declaration alone produces a **nil channel** (see [Section 11](#11-nil-channels-a-surprisingly-useful-tool)) — to get a usable channel you must initialize it with `make`:

```go
ch := make(chan int)        	 // unbuffered channel of int
buffered := make(chan int, 10) // buffered channel with capacity 10
```

The optional second argument to `make` sets the buffer capacity; omitting it (or passing `0`) creates an unbuffered channel.

## 3. Sending and Receiving

The `<-` operator is used for both sending and receiving, with the channel on one side and the arrow pointing in the direction the data flows:

```go
ch <- 42        // send: put 42 into ch
value := <-ch   // receive: take a value out of ch and assign it

<-ch            // receive and discard the value (still synchronizes/blocks)
```

Both operations can also appear in an assignment with a second, boolean result that reports whether the channel is still open (covered in [Section 6](#6-closing-channels)):

```go
value, ok := <-ch
```

## 4. Unbuffered vs. Buffered Channels

### 4.1 Unbuffered Channels

An **unbuffered channel** (`make(chan T)`) has zero capacity. A send on it blocks until another goroutine is ready to receive at that exact moment, and vice versa — the two goroutines effectively "rendezvous." This gives unbuffered channels strong synchronization guarantees: successfully sending a value on an unbuffered channel guarantees the receiving goroutine has already started receiving it.

```go
ch := make(chan string)

go func() {
    ch <- "done" // blocks until main() receives
}()

msg := <-ch // blocks until the goroutine sends
fmt.Println(msg)
```

### 4.2 Buffered Channels

A **buffered channel** (`make(chan T, n)`) has a fixed-size internal queue. A send only blocks once the buffer is full; a receive only blocks once the buffer is empty. This decouples the sender's and receiver's timing to some degree.

```go
ch := make(chan int, 2)
ch <- 1 // does not block — buffer has room
ch <- 2 // does not block — buffer now full
// ch <- 3 // would block here until something is received
fmt.Println(len(ch), cap(ch)) // 2 2 — current length and capacity
```

`len(ch)` returns how many elements are currently queued; `cap(ch)` returns the buffer's total capacity.

**When to use which:**

- **Unbuffered** — when you need a strict guarantee that the receiver has taken the value before the sender continues (a true hand-off / synchronization point).
- **Buffered** — when you want to smooth out small timing mismatches between producers and consumers, or implement a bounded queue/semaphore (see [Section 13](#13-common-channel-patterns)). A buffered channel is **not** a substitute for proper synchronization — it only shifts the point at which blocking happens; it doesn't eliminate the need to think about backpressure or completion.

## 5. The Four Channel Axioms

These four rules, popularized by Dave Cheney's "Channel Axioms," summarize almost everything you need to remember about channel blocking/panic behavior:

| #      | Axiom                                                                                                                          |
| ------ | ------------------------------------------------------------------------------------------------------------------------------ |
| **A1** | A send to a **nil** channel blocks forever.                                                                                    |
| **A2** | A receive from a **nil** channel blocks forever.                                                                               |
| **A3** | A send to a **closed** channel panics.                                                                                         |
| **A4** | A receive from a **closed** channel returns the channel's zero value immediately (once any buffered values have been drained). |

```go
// A1 & A2 — nil channel
var c chan int
c <- 1   // blocks forever (deadlock if nothing else is running)
<-c      // also blocks forever

// A3 — send to closed channel
ch := make(chan int)
close(ch)
ch <- 1 // panic: send on closed channel

// A4 — receive from closed channel
ch2 := make(chan int, 2)
ch2 <- 10
ch2 <- 20
close(ch2)
fmt.Println(<-ch2) // 10
fmt.Println(<-ch2) // 20
fmt.Println(<-ch2) // 0 (zero value, returned immediately, no blocking)
```

Two more rules worth memorizing alongside these: **closing a nil channel panics**, and **closing an already-closed channel panics**.

Knowing these axioms by heart makes a large share of channel-related bugs (deadlocks, unexpected panics) immediately recognizable.

## 6. Closing Channels

`close(ch)` marks a channel so that no further values can be sent on it. It is a **signal**, not a requirement — most channels in a program never need to be closed explicitly, especially ones that live for the whole program's lifetime.

You need to close a channel specifically when the **receiver must be told there is nothing more coming** — most commonly so a `for range` loop over the channel can terminate (see [Section 8](#8-ranging-over-a-channel)).

```go
func produce(ch chan<- int) {
    defer close(ch) // signal "no more values" once this function returns
    for i := 0; i < 5; i++ {
        ch <- i
    }
}
```

**Panics to avoid:**

- Sending to a closed channel → panic (Axiom A3).
- Closing a channel twice → panic.
- Closing a nil channel → panic.

## 7. Who Should Close a Channel?

This is one of the most common points of confusion for newcomers. The rule of thumb is simple:

> **Only the sender should close a channel. The receiver should never close it.**

The reasoning: the sender is the only party that knows for certain when it has finished sending. If a receiver closes a channel while a sender might still try to send on it, the sender's next send will panic (Axiom A3) — a race condition between closing and sending.

### 7.1 Single Sender (Simple Case)

```go
func producer(ch chan<- int) {
    defer close(ch)
    for i := 0; i < 3; i++ {
        ch <- i
    }
}
```

This is the most straightforward and most common case: the goroutine that owns and writes to the channel is also the one that closes it.

### 7.2 Multiple Senders

When several goroutines send on the same channel, none of them individually knows when _all_ of them are done, so none should close the channel directly — doing so causes a race between "the first goroutine to finish closes the channel" and "the other goroutines are still trying to send," which panics. Instead, use a coordinator: typically a `sync.WaitGroup` tracking all senders, with a separate goroutine that waits for the group to finish and only then closes the channel.

```go
func fanIn(sources ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(sources))

    for _, src := range sources {
        go func(src <-chan int) {
            defer wg.Done()
            for v := range src {
                out <- v
            }
        }(src)
    }

    go func() {
        wg.Wait()  // wait until every sender goroutine has finished
        close(out) // only now is it safe to close
    }()

    return out
}
```

### 7.3 When Not to Close at All

If a channel is used for the entire lifetime of the program (e.g., a control/command channel that's simply garbage collected when the program exits), it's perfectly acceptable to never close it. Channels are not like file handles — there's no resource leak from an unclosed channel by itself, as long as nothing is left permanently blocked waiting on it (which _would_ be a goroutine leak: a goroutine stuck forever on a channel operation with no way to unblock, which keeps its stack and any resources it holds alive for as long as the program runs).

## 8. Ranging Over a Channel

`for range` over a channel receives values until the channel is both closed and drained, then exits the loop automatically — it is the idiomatic way to consume a stream of values from a channel:

```go
func producer(ch chan<- int) {
    defer close(ch)
    for i := 0; i < 5; i++ {
        ch <- i
    }
}

func main() {
    ch := make(chan int)
    go producer(ch)

    for v := range ch {
        fmt.Println(v)
    }
    fmt.Println("channel closed and drained")
}
```

This is exactly equivalent, under the hood, to a manual loop using the comma-ok form:

```go
for {
    v, ok := <-ch
    if !ok {
        break // channel closed and drained
    }
    // use v
}
```

**Important:** if the channel is never closed and no more values ever arrive, a `for range` loop over it will block forever — a classic goroutine leak (see [Section 15](#15-common-mistakes-and-pitfalls)).

## 9. Directional (Typed-Direction) Channels

A channel type can be restricted to send-only (`chan<- T`) or receive-only (`<-chan T`) in a function signature. This is enforced by the compiler and is a valuable way to make a function's intent explicit and prevent accidental misuse (e.g., a consumer function accidentally closing or sending on a channel it should only read from):

```go
func send(out chan<- int, v int) {
    out <- v
    // <-out // compile error: cannot receive from a send-only channel
}

func receive(in <-chan int) int {
    return <-in
    // in <- 5 // compile error: cannot send to a receive-only channel
}
```

A regular bidirectional channel (`chan T`) is automatically convertible to a directional one when passed as an argument, but not the other way around — this asymmetry is what lets you safely narrow a channel's capabilities as it's passed deeper into a call chain.

## 10. The `select` Statement

`select` lets a goroutine wait on multiple channel operations simultaneously, proceeding with whichever one becomes ready first. If multiple cases are ready at once, Go picks one at random (uniformly) among them — this is a deliberate design choice to avoid starving any particular case.

```go
select {
case v := <-ch1:
    fmt.Println("received from ch1:", v)
case ch2 <- 99:
    fmt.Println("sent 99 to ch2")
case <-time.After(2 * time.Second):
    fmt.Println("timed out waiting")
default:
    fmt.Println("nothing ready right now")
}
```

**Key behaviors:**

- A `select` with **no `default`** blocks until at least one case can proceed. `select` never picks a case that would block — a case is only eligible once its channel operation can actually complete immediately.
- A `select` **with a `default`** case never blocks: if no other case is ready right away, `default` runs immediately, making the operation non-blocking.
- An empty `select {}` with no cases at all blocks forever — occasionally used deliberately to keep `main()` alive while other goroutines run, though a proper synchronization mechanism (`WaitGroup`, signal channel) is usually clearer.

**Typical uses:**

1. **Timeouts** — race a receive against `time.After(...)`.
2. **Cancellation** — race a receive against `ctx.Done()`.
3. **Multiplexing** — service whichever of several input channels has data.
4. **Non-blocking probes** — check if a channel has a value ready without committing to wait for it, via `default`.

## 11. Nil Channels: A Surprisingly Useful Tool

A **nil channel** is a channel variable that was declared but never initialized with `make` (or explicitly set to `nil`). Per the channel axioms, both sending to and receiving from a nil channel block forever — which sounds useless, but is actually a valuable technique **inside a `select` statement**: a nil channel's case in a `select` is never chosen (since it would block forever), which lets you dynamically "disable" a case.

```go
func merge(a, b <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for a != nil || b != nil {
            select {
            case v, ok := <-a:
                if !ok {
                    a = nil // disable this case once `a` is drained
                    continue
                }
                out <- v
            case v, ok := <-b:
                if !ok {
                    b = nil // disable this case once `b` is drained
                    continue
                }
                out <- v
            }
        }
    }()
    return out
}
```

Here, setting `a` or `b` to `nil` after it closes effectively removes that `select` case from consideration, without needing extra boolean flags or restructuring the loop. This is a common, idiomatic pattern for merging multiple channels that finish at different times (see also [Section 13.2](#132-fan-in)).

## 12. Channels as Synchronization Primitives

Beyond carrying data, an empty-struct channel (`chan struct{}`) is a very common idiom purely for **signaling** — since `struct{}` occupies zero bytes, it communicates "an event happened" with no payload:

```go
done := make(chan struct{})

go func() {
    // ... do work ...
    close(done) // broadcast "I'm finished" to anyone listening
}()

<-done // block until the goroutine signals completion
```

Closing a channel (rather than sending on it) is a particularly useful signaling idiom because **closing broadcasts to all current and future receivers simultaneously** — every goroutine blocked on `<-done`, no matter how many, unblocks at once, and any goroutine that receives from `done` _afterward_ gets the zero value immediately (Axiom A4). This is exactly how `context.Context`'s `Done()` channel works internally.

### 12.1 The Semaphore Pattern

A buffered channel's fixed capacity makes it a natural fit for implementing a **counting semaphore** — limiting how many goroutines can do something concurrently:

```go
var sem = make(chan struct{}, 3) // at most 3 concurrent operations

func limitedWork() {
    sem <- struct{}{}        // acquire a slot (blocks if 3 are already in use)
    defer func() { <-sem }() // release the slot

    // ... do the actual work ...
}
```

**Use case:** capping concurrent outbound HTTP requests, database connections, or file handles so a burst of goroutines doesn't overwhelm a downstream resource.

## 13. Common Channel Patterns

### 13.1 Worker Pool

A fixed number of goroutines pull tasks from a shared input channel and push results to an output channel, bounding concurrency instead of spawning one goroutine per task:

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        results <- j * j
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs) // tells workers there are no more jobs; their range loops exit

    for a := 0; a < 9; a++ {
        fmt.Println(<-results)
    }
}
```

### 13.2 Fan-In

Merges several input channels into a single output channel, typically using a `WaitGroup` to know when it is safe to close the merged output (see [Section 7.2](#72-multiple-senders)).

### 13.3 Fan-Out

The inverse: a single source channel is read by several worker goroutines concurrently, distributing the work among them (this is effectively what the worker pool above does with `jobs`).

### 13.4 Pipeline

Channels chain multiple processing stages together, each stage a goroutine that reads from an input channel and writes to an output channel:

```go
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func double(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * 2
        }
    }()
    return out
}

// usage: for v := range double(generate(1, 2, 3)) { fmt.Println(v) } // 2 4 6
```

**Use case:** streaming data transformations (read → parse → filter → write), each stage running as its own goroutine so the whole pipeline overlaps work instead of processing everything sequentially.

### 13.5 Timeout Pattern

```go
select {
case res := <-resultCh:
    fmt.Println("got result:", res)
case <-time.After(3 * time.Second):
    fmt.Println("operation timed out")
}
```

### 13.6 Done/Cancellation Channel (Pre-`context` Idiom)

Before or alongside `context.Context`, a simple `chan struct{}` that gets closed is a common way to broadcast "stop" to many goroutines at once:

```go
func worker(done <-chan struct{}, work <-chan int) {
    for {
        select {
        case <-done:
            return
        case w, ok := <-work:
            if !ok {
                return
            }
            process(w)
        }
    }
}
```

In modern Go code, `context.Context`'s `Done()` channel is generally preferred over a hand-rolled `done` channel because it also carries deadlines, cancellation reasons (`ctx.Err()`), and request-scoped values in a standardized way.

## 14. Channels vs. Mutexes: When to Use Which

Both channels and `sync.Mutex` can coordinate access to shared state, but they suit different situations, and Go's own FAQ notes that structuring a program so a single goroutine owns a piece of data at a time is generally good practice regardless of which mechanism enforces it.

| Use a **channel** when...                                               | Use a **mutex** when...                                                          |
| ----------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| You're passing _ownership_ of a value from one goroutine to another     | Multiple goroutines need to read/write the _same_ piece of shared state in place |
| You're coordinating a pipeline or workflow of independent stages        | The critical section is small, simple, and local (e.g., a counter, a map)        |
| You need to broadcast an event to many goroutines at once (via `close`) | The data doesn't naturally "flow" anywhere — it just needs protecting            |
| You want built-in support for timeouts/cancellation via `select`        | You want lower overhead for a very hot, simple, short critical section           |

Neither is strictly "more idiomatic" than the other — Go provides both because different problems fit each shape more naturally. A common guideline is to reach for channels when data is moving between goroutines, and for a mutex when the goroutines are converging on shared state that doesn't logically move anywhere.

## 15. Common Mistakes and Pitfalls

### 15.1 Deadlock: Sending With No Receiver

```go
func main() {
    ch := make(chan int)
    ch <- 1 // deadlock: nothing will ever receive this
}
```

Go's runtime can detect this specific "all goroutines are asleep" situation and crashes with `fatal error: all goroutines are asleep - deadlock!` — helpful for catching the bug immediately during development, though it won't catch every deadlock scenario (e.g., ones involving blocked I/O).

### 15.2 Closing a Channel From the Receiver

```go
// WRONG: receiver closes the channel it's reading from
func consumer(ch chan int) {
    for v := range ch {
        fmt.Println(v)
    }
    close(ch) // if the sender tries to send afterward, it panics
}
```

Only the sender (or a coordinator that knows all senders are done) should close a channel.

### 15.3 Closing a Channel Twice, or From Multiple Goroutines

```go
close(ch)
close(ch) // panic: close of closed channel
```

When multiple goroutines might all reach "I'm the last one, I'll close it," you need explicit coordination (`sync.Once`, a `WaitGroup`-based coordinator, or an atomic flag) — never let more than one goroutine call `close` on the same channel unconditionally.

### 15.4 Forgetting a Channel Can Leak a Goroutine

```go
func leaky() {
    ch := make(chan int) // unbuffered, nobody will ever receive
    go func() {
        ch <- 1 // blocks forever
    }()
}
```

If nothing ever receives from `ch`, the sending goroutine is stuck permanently — a goroutine leak. Always make sure a channel operation has a way to unblock, whether through a proper receiver, a `select` with a cancellation case, or a buffer sized to avoid blocking in the first place.

### 15.5 Using an Unbuffered Channel Where a Buffered One Was Needed (or vice versa)

Reasoning incorrectly about buffer size is a frequent source of unexpected blocking. For example, sending N values on an unbuffered channel from the same goroutine that will later try to receive them — without a separate goroutine to receive concurrently — deadlocks, because the first send can't complete until something else receives:

```go
ch := make(chan int)
ch <- 1 // deadlock: same goroutine can't also be there to receive
v := <-ch
```

### 15.6 Assuming `select` Case Order Matters

`select` does **not** try cases in the order they're written when more than one is ready — it picks uniformly at random among the ready cases. Relying on source-code order for priority is a common but incorrect assumption; if you need priority between cases, you must implement it explicitly (e.g., a nested `select` that checks the higher-priority channel first with a `default`).

### 15.7 Not Draining a Channel Before Assuming a Goroutine Exited

When a producer goroutine writes to a buffered channel and then exits, values may still be sitting in the buffer waiting to be read; forgetting this can lead to code that assumes "no more values" prematurely, or to a goroutine leak if the buffer is never drained and something is still blocked trying to write to a full buffer.

## 16. Best Practices Summary

1. **Decide channel ownership up front.** Know which goroutine (or coordinator) is responsible for closing a channel before you write the code, not after a panic in production.
2. **Only the sender closes; the receiver never does.** With multiple senders, use a `WaitGroup`-based coordinator to close safely.
3. **Close a channel only when receivers need an end-of-stream signal**, typically to let a `for range` loop terminate — don't close channels reflexively.
4. **Always give a channel operation an escape hatch.** Wrap blocking sends/receives in a `select` alongside `ctx.Done()` or `time.After(...)` when there's any chance the other side might never show up.
5. **Prefer directional channel types (`chan<-`, `<-chan`) in function signatures** to make intent clear and let the compiler catch misuse.
6. **Use `chan struct{}` for pure signaling** — it makes the intent (an event, not data) explicit and costs nothing extra.
7. **Remember the four channel axioms** — they explain almost every channel-related panic or deadlock you'll encounter.
8. **Size buffered channels deliberately**, not arbitrarily — a buffer is a design decision (how much slack between producer and consumer is acceptable), not just a way to silence blocking.
9. **Use channels for moving/ownership-transferring data between goroutines; use a mutex for protecting shared state that doesn't move.**
10. **Test concurrent channel code with `go test -race`** to catch synchronization bugs the type system can't catch for you.

## 17. Use Case Summary Table

| Technique                     | When to Use                                                                   | Example Scenario                                                |
| ----------------------------- | ----------------------------------------------------------------------------- | --------------------------------------------------------------- |
| Unbuffered channel            | Need a strict hand-off/rendezvous between exactly one sender and one receiver | Signaling "goroutine finished," passing a single computed value |
| Buffered channel              | Decouple producer/consumer timing, implement a queue                          | Worker pool job queue, batching results                         |
| `close(ch)` + `range`         | Tell a consumer "no more values are coming"                                   | Producer streaming a finite sequence of items                   |
| Directional channel params    | Restrict a function to only send or only receive                              | Public API functions that should not misuse a channel           |
| `select` with `time.After`    | Prevent a goroutine from blocking indefinitely                                | Network call with a timeout                                     |
| `select` with `ctx.Done()`    | Respect cancellation/deadlines                                                | Any long-running worker or downstream call                      |
| `select` with `default`       | Non-blocking probe of a channel                                               | Polling for available work without waiting                      |
| Nil channel in `select`       | Dynamically disable a `select` case                                           | Merging multiple channels that finish at different times        |
| `chan struct{}` + `close`     | Broadcast an event to many goroutines at once                                 | Shutdown signal, "ready" signal                                 |
| Buffered channel as semaphore | Bound concurrent access to a limited resource                                 | Cap concurrent outbound HTTP requests or DB connections         |
| Fan-in / fan-out              | Distribute or merge concurrent work across goroutines                         | Parallel processing of an input stream                          |
| Pipeline of channels          | Chain independent processing stages                                           | Multi-step streaming data transformation                        |

## 18. References

1. Go Team — _Effective Go: Concurrency (channels)_. https://go.dev/doc/effective_go
2. Go Team — _A Tour of Go: Concurrency_. https://go.dev/tour/concurrency
3. Go Team — _Go Wiki: LearnConcurrency_. https://go.dev/wiki/LearnConcurrency
4. Dave Cheney — _Channel Axioms_. https://dave.cheney.net/2014/03/19/channel-axioms
5. Usman Mahmood — _Go Channel Axioms_. https://www.usman.me.uk/2015/12/channel-axioms/
6. udhos — _Golang Concurrency Tricks_. https://udhos.github.io/golang-concurrency-tricks/
7. Suryansh Shrivastava — _Understanding Channel Closure in Go: Who Should Close It?_, Medium. https://medium.com/@suryanshshrivastava_75738/understanding-channel-closure-in-go-who-should-close-it-299fb022e583
8. ByteSizeGo — _Go channels: send, receive, close, and range without surprises_. https://www.bytesizego.com/blog/go-channels-send-receive-close-and-range-without-surprises
9. Leo Lara — _Closing a Go channel written by several goroutines_, DEV Community. https://dev.to/leolara/closing-a-go-channel-written-by-several-goroutines-52j2
10. go101.org — _Channel Closing Principles_, referenced via Leo Lara's article above. https://go101.org/article/channel-closing.html
11. Educative — _Closing Go Channels_. https://www.educative.io/courses/advanced-techniques-in-go-programming/closing-go-channels
12. Educative — _Synchronization of goroutines_. https://www.educative.io/courses/the-way-to-go/synchronization-of-goroutines
13. Educative — _Semaphore Pattern_. https://www.educative.io/module/page/wjB3xQCAR51LoORMY/10370001/5765062970572800/6154041951780864
14. Pratik Pandey — _Golang Concepts: Nil Channels_. https://pratikpandey.substack.com/p/golang-concepts-nil-channels
15. Nicolas Lepage — _Go channels in JS (3/5): Closing_, DEV Community. https://dev.to/zenika/go-channels-in-js-3-5-closing-4l6l
16. OpenZiti Engineering Blog — _Golang Aha! Moments: Channels_. https://blog.openziti.io/golang-aha-moments-channels
17. GetStream Engineering Blog — _Goroutines in Go: A Practical Guide to Concurrency_. https://getstream.io/blog/goroutines-go-concurrency-guide/
18. Go standard library documentation — _Package `context`_. https://pkg.go.dev/context
19. Go FAQ — _What operations are atomic? What about mutexes?_. https://go.dev/doc/faq#What_operations_are_atomic_What_about_mutexes
20. Go Language Specification — _Channel types, Send statements, Receive operator, Select statements_. https://go.dev/ref/spec
