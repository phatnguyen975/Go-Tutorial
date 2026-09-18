<div align="center">
  <h1>Goroutines</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 17, 2026</sub>
</div>

## Table of Contents

1. [What Is a Goroutine?](#1-what-is-a-goroutine)
2. [Goroutines vs. OS Threads](#2-goroutines-vs-os-threads)
3. [Starting a Goroutine](#3-starting-a-goroutine)
4. [The Go Scheduler: the GMP Model](#4-the-go-scheduler-the-gmp-model)
5. [GOMAXPROCS](#5-gomaxprocs)
6. [Concurrency vs. Parallelism](#6-concurrency-vs-parallelism)
7. [Synchronizing Goroutines with `sync.WaitGroup`](#7-synchronizing-goroutines-with-syncwaitgroup)
8. [Sharing Data Safely: `sync.Mutex` and `sync.RWMutex`](#8-sharing-data-safely-syncmutex-and-syncrwmutex)
9. [Channels: Communicating Between Goroutines](#9-channels-communicating-between-goroutines)
10. [The `select` Statement](#10-the-select-statement)
11. [Other Useful `sync` Primitives](#11-other-useful-sync-primitives)
12. [Cancellation and Timeouts with `context`](#12-cancellation-and-timeouts-with-context)
13. [Common Concurrency Patterns](#13-common-concurrency-patterns)
14. [Goroutine Leaks](#14-goroutine-leaks)
15. [Race Conditions and the Race Detector](#15-race-conditions-and-the-race-detector)
16. [Error Handling in Concurrent Code (`errgroup`)](#16-error-handling-in-concurrent-code-errgroup)
17. [Debugging and Observability](#17-debugging-and-observability)
18. [Best Practices Summary](#18-best-practices-summary)
19. [Common Mistakes](#19-common-mistakes)
20. [Use Case Summary Table](#20-use-case-summary-table)
21. [References](#21-references)

## 1. What Is a Goroutine?

A goroutine is a function or method that runs **concurrently** with the rest of your program, managed by the Go runtime rather than the operating system. The name comes from "Go" + "subroutine." Every Go program already has at least one goroutine — the one running `main()` — and you can launch more with the `go` keyword.

Goroutines are described as **lightweight threads of execution**: they are far cheaper to create than OS threads, which lets a single Go program comfortably run tens of thousands, or even millions, of goroutines at once.

**Key characteristics:**

- Created with the `go` keyword in front of a function call.
- Start with a very small stack (a few KB) that grows and shrinks dynamically as needed, unlike an OS thread's typically fixed, much larger stack.
- Scheduled cooperatively/preemptively by the **Go runtime scheduler**, not directly by the OS kernel.
- Communicate safely with each other primarily through **channels**, following Go's philosophy of _"Don't communicate by sharing memory; share memory by communicating."_

## 2. Goroutines vs. OS Threads

| Aspect                    | Goroutine                                          | OS Thread                         |
| ------------------------- | -------------------------------------------------- | --------------------------------- |
| Managed by                | Go runtime (user-space scheduler)                  | Operating system kernel           |
| Initial stack size        | ~2 KB, grows/shrinks dynamically                   | Typically 1–8 MB, often fixed     |
| Creation cost             | Very cheap (µs range)                              | Relatively expensive              |
| Context switch cost       | Cheap (no kernel trap needed)                      | Expensive (kernel-level switch)   |
| Typical count per process | Thousands to millions                              | Usually limited to a few thousand |
| Scheduling                | M:N — many goroutines mapped onto fewer OS threads | 1:1 with the kernel scheduler     |

Because goroutines are multiplexed onto a much smaller number of OS threads by the Go runtime, you get concurrency at a much lower memory and CPU cost than spawning an OS thread per task — this is exactly the kind of scaling problem goroutines were designed to solve for network services handling many simultaneous connections.

## 3. Starting a Goroutine

Launching a goroutine is done by prefixing a function call with the `go` keyword:

```go
package main

import (
    "fmt"
    "time"
)

func sayHello() {
    fmt.Println("Hello from a goroutine!")
}

func main() {
    go sayHello() // starts sayHello() as a new goroutine

    // main() itself runs in a goroutine and does not wait for others.
    time.Sleep(100 * time.Millisecond) // crude way to let sayHello() finish; avoid in real code
    fmt.Println("Hello from main!")
}
```

**Important:** `main()` does **not** automatically wait for goroutines it launches. If `main()` returns, the whole program exits immediately, even if other goroutines haven't finished — the `time.Sleep` above is a hack for illustration only; real code should use `sync.WaitGroup` or channels (see [Section 7](#7-synchronizing-goroutines-with-syncwaitgroup)).

You can also launch an anonymous function as a goroutine, which is very common for short-lived concurrent work:

```go
go func(msg string) {
    fmt.Println(msg)
}("running concurrently")
```

**Closure pitfall:** be careful when a goroutine's closure captures a loop variable — this is one of the most common beginner mistakes and is covered in [Section 19](#19-common-mistakes).

## 4. The Go Scheduler: the GMP Model

Go doesn't map goroutines directly onto OS threads one-to-one. Instead, the runtime uses **M:N scheduling**, meaning **M** goroutines are multiplexed across **N** OS threads. This is implemented through what is commonly called the **GMP model**:

| Symbol | Name      | Role                                                                                                                                                        |
| ------ | --------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **G**  | Goroutine | A lightweight, user-space unit of work — essentially your `go` statement's function plus its stack and scheduling metadata.                                 |
| **M**  | Machine   | An actual OS thread that executes goroutine code.                                                                                                           |
| **P**  | Processor | A logical scheduling context (not a physical CPU core) that holds a local run queue of runnable goroutines and pairs 1:1 with an M while executing Go code. |

**How it fits together:**

- The number of **P**s is controlled by `GOMAXPROCS` (defaults to the number of logical CPU cores).
- Each **P** owns a local run queue of goroutines ready to execute; there is also a smaller global run queue used as an overflow.
- An **M** (OS thread) must acquire a **P** before it can execute Go code; it then pulls goroutines from that P's local queue.
- When a P's local queue runs dry, it will **steal** goroutines from another P's queue — a technique called **work stealing** — which keeps all CPU cores busy without a centralized lock.
- If a goroutine performs a blocking system call, its M detaches from its P so that P can be handed to another M and keep executing other goroutines; the Go runtime's non-cooperative preemption also allows it to interrupt long-running goroutines that don't naturally yield, preventing one goroutine from starving the others.

**Why this matters practically:** you rarely need to think about M and P directly, but understanding that goroutines are cheap, cooperatively multiplexed units — not OS threads — explains why you can safely spawn thousands of them, and why blocking one goroutine (e.g., on I/O) doesn't necessarily block others.

## 5. GOMAXPROCS

`GOMAXPROCS` sets the maximum number of OS threads that can execute Go code simultaneously — effectively, the number of **P**s. Since Go 1.5, the default equals the number of logical CPU cores available to the process.

```go
import "runtime"

func main() {
    fmt.Println("Available CPUs:", runtime.NumCPU())
    fmt.Println("Current GOMAXPROCS:", runtime.GOMAXPROCS(0)) // 0 = query, don't change

    runtime.GOMAXPROCS(4) // explicitly set to 4
}
```

It can also be set via the `GOMAXPROCS` environment variable before starting the program.

**Practical guidance:**

- In most applications the default (number of CPU cores) is the right choice and should not be changed.
- Setting `GOMAXPROCS` higher than the number of available CPU cores generally does not help, and can hurt performance by forcing the OS to do more expensive context switching between threads than necessary.
- In containerized environments (Docker/Kubernetes) with CPU limits set below the host's core count, it is worth checking that Go correctly detects the container's CPU quota (modern Go versions and libraries like `go.uber.org/automaxprocs` address this).

## 6. Concurrency vs. Parallelism

These two terms are often conflated but mean different things:

- **Concurrency** is about _structure_: dealing with multiple tasks that are in progress at overlapping times, without necessarily running at the exact same instant. It's about managing many things at once.
- **Parallelism** is about _execution_: actually running multiple tasks at the exact same instant, which requires multiple CPU cores.

A Go program can be concurrent (structured as many goroutines) while running on a single CPU core — the scheduler interleaves them. It becomes parallel when `GOMAXPROCS` > 1 and there is real work to distribute across multiple cores simultaneously. Rob Pike's well-known formulation is that concurrency is about _dealing with_ lots of things at once, while parallelism is about _doing_ lots of things at once — concurrency is a way of structuring a program that may (but does not have to) enable parallel execution.

## 7. Synchronizing Goroutines with `sync.WaitGroup`

Since `main()` (or any calling goroutine) doesn't automatically wait for the goroutines it launches, you need explicit synchronization. `sync.WaitGroup` is the standard tool for "wait until N goroutines have finished."

```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // decrement the counter when this goroutine finishes
    fmt.Printf("Worker %d starting\n", id)
    // ... do work ...
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1) // increment the counter before launching
        go worker(i, &wg)
    }

    wg.Wait() // blocks until the counter returns to zero
    fmt.Println("All workers completed")
}
```

**Rules of thumb:**

- Call `Add()` **before** starting the goroutine, not inside it — otherwise `Wait()` might return before all goroutines have even registered.
- Always call `Done()` via `defer` at the very top of the goroutine's function, so it still runs even if the function panics or returns early.
- A `WaitGroup` must not be copied after first use — always pass it by pointer.

**Use case:** Fire off a fixed, known number of independent tasks (e.g., fetching several URLs concurrently) and wait for all of them to finish before continuing.

## 8. Sharing Data Safely: `sync.Mutex` and `sync.RWMutex`

Go's preferred concurrency style is "share memory by communicating" via channels, but sometimes protecting a simple shared variable (a counter, a map, a cache) with a lock is simpler and more idiomatic than building a channel-based solution. That's what `sync.Mutex` is for.

```go
package main

import (
    "fmt"
    "sync"
)

type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

func main() {
    counter := SafeCounter{}
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }

    wg.Wait()
    fmt.Println("Final count:", counter.Value()) // always 1000
}
```

Here, the `sync.Mutex` ensures only one goroutine can modify `count` at a time, preventing a **data race**. Without it, concurrent increments could overwrite each other and produce a wrong, non-deterministic final value.

**`sync.RWMutex`** is a variant that allows any number of concurrent readers, or a single writer, but not both at once — useful when reads vastly outnumber writes (e.g., an in-memory cache read very frequently but updated rarely):

```go
var mu sync.RWMutex

func read() {
    mu.RLock()
    defer mu.RUnlock()
    // ... read shared data ...
}

func write() {
    mu.Lock()
    defer mu.Unlock()
    // ... modify shared data ...
}
```

**Use case:** Protecting any piece of shared mutable state (counters, in-memory caches, connection pools, maps) that multiple goroutines read and write concurrently.

## 9. Channels: Communicating Between Goroutines

A **channel** is a typed conduit through which goroutines send and receive values, built directly into the language. Channels are the primary mechanism for goroutines to communicate and synchronize without directly sharing memory — encapsulated in Go's philosophy of preferring communication over shared state.

### 9.1 Creating and Using Channels

```go
ch := make(chan int)   // unbuffered channel of int
ch <- 42                // send a value into the channel (blocks until received)
value := <-ch           // receive a value from the channel (blocks until sent)
```

### 9.2 Unbuffered vs. Buffered Channels

- **Unbuffered channel** (`make(chan T)`): a send blocks until another goroutine is ready to receive, and vice versa. This provides strong synchronization — the two goroutines effectively "rendezvous" at that point.
- **Buffered channel** (`make(chan T, n)`): a send only blocks once the buffer is full; a receive only blocks once the buffer is empty. This decouples sender and receiver timing somewhat.

```go
buffered := make(chan string, 3)
buffered <- "a"
buffered <- "b"
buffered <- "c"
// A 4th send here would block until something is received.
```

### 9.3 Closing Channels and Ranging Over Them

A sender can `close()` a channel to signal that no more values will be sent. Receivers can detect this:

```go
func produce(ch chan<- int) {
    defer close(ch)
    for i := 0; i < 5; i++ {
        ch <- i
    }
}

func main() {
    ch := make(chan int)
    go produce(ch)

    for v := range ch { // automatically exits when ch is closed
        fmt.Println(v)
    }
}
```

**Rules:**

- Only the **sender** should close a channel, never the receiver.
- Sending on a closed channel causes a panic.
- Receiving from a closed channel returns the zero value immediately (you can check `v, ok := <-ch`; `ok` is `false` once the channel is closed and drained).
- Closing is optional — you only need to close a channel when receivers need to know "there's nothing more coming" (e.g., to end a `range` loop).

### 9.4 Directional Channels

Function signatures can restrict a channel to send-only (`chan<-`) or receive-only (`<-chan`), which the compiler enforces — a useful way to make intent explicit and prevent misuse:

```go
func produce(out chan<- int) { /* can only send to out */ }
func consume(in <-chan int)  { /* can only receive from in */ }
```

**Use case:** Channels shine for producer/consumer pipelines, fan-out/fan-in work distribution, signaling task completion, and implementing done/cancellation signals.

## 10. The `select` Statement

`select` lets a goroutine wait on multiple channel operations at once, proceeding with whichever is ready first — Go's equivalent of a `switch` for channels.

```go
select {
case msg1 := <-ch1:
    fmt.Println("received from ch1:", msg1)
case msg2 := <-ch2:
    fmt.Println("received from ch2:", msg2)
case <-time.After(2 * time.Second):
    fmt.Println("timeout waiting for a message")
default:
    fmt.Println("no channel ready right now") // makes select non-blocking
}
```

**Common uses of `select`:**

- **Timeouts:** race a channel receive against `time.After(...)` so a goroutine never blocks forever.
- **Cancellation:** race a channel receive against `ctx.Done()` (see [Section 12](#12-cancellation-and-timeouts-with-context)).
- **Multiplexing:** service whichever of several channels has data ready, without giving any one channel priority in code order (Go picks pseudo-randomly among ready cases).
- **Non-blocking operations:** adding a `default` case makes a channel send/receive attempt instantly instead of blocking.

## 11. Other Useful `sync` Primitives

| Primitive     | Purpose                                                                                                                                                                                  |
| ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `sync.Once`   | Ensures a function runs exactly once, no matter how many goroutines call it — ideal for lazy, thread-safe initialization (e.g., a singleton config or connection pool).                  |
| `sync.Cond`   | A condition variable that lets goroutines wait for, and be notified of, a particular condition becoming true; lower-level and less commonly needed than channels or `WaitGroup`.         |
| `sync.Pool`   | A pool of temporary, reusable objects to reduce garbage-collector pressure for frequently allocated/discarded objects (e.g., byte buffers in a hot path).                                |
| `sync/atomic` | Low-level atomic operations (`Add`, `Load`, `Store`, `CompareAndSwap`) on integers and pointers — faster than a mutex for simple counters, but easy to misuse for anything more complex. |

```go
var once sync.Once

func setup() {
    once.Do(func() {
        fmt.Println("This prints exactly once, even under concurrent calls")
    })
}
```

## 12. Cancellation and Timeouts with `context`

The `context` package is the standard way to carry **deadlines, cancellation signals, and request-scoped values** across API boundaries and between goroutines. It is critical in any production application for preventing goroutine leaks and ensuring clean shutdown.

### 12.1 Basic Cancellation

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("worker %d stopping: %v\n", id, ctx.Err())
            return
        default:
            // do a small unit of work
            time.Sleep(200 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    go worker(ctx, 1)

    time.Sleep(1 * time.Second)
    cancel() // signals ctx.Done() to close, worker exits
    time.Sleep(500 * time.Millisecond)
}
```

### 12.2 Timeouts and Deadlines

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel() // always call cancel to release resources, even if the timeout fires first

select {
case <-doWork(ctx):
    fmt.Println("work finished")
case <-ctx.Done():
    fmt.Println("timed out:", ctx.Err())
}
```

`context.WithDeadline` works similarly but takes an absolute `time.Time` instead of a duration.

### 12.3 Key Conventions

- Pass `context.Context` as the **first parameter** of any function that might block, do I/O, or spawn goroutines — by strong convention, named `ctx`.
- Always call the returned `cancel()` function (typically via `defer`), even if the context's deadline already expired — failing to do so leaks resources tied to the context.
- Never store a `context.Context` inside a struct; pass it explicitly through function calls.
- `context.Background()` is the root context, typically used in `main()`, tests, and top-level request handlers; `context.TODO()` signals "I haven't decided which context to use yet."

**Use case:** Any goroutine that could otherwise run indefinitely — a background worker, an HTTP handler's downstream calls, a long database query — should watch `ctx.Done()` so it can be told to stop cleanly when the caller no longer needs the result (client disconnected, request timed out, application shutting down).

## 13. Common Concurrency Patterns

### 13.1 Worker Pool

A fixed number of goroutines ("workers") pull tasks from a shared input channel and send results to an output channel — this bounds concurrency instead of spawning one goroutine per task, which matters when tasks are numerous or resource-intensive (e.g., DB connections, external API calls).

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        results <- j * 2 // pretend this is expensive work
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    for w := 1; w <= 3; w++ { // 3 workers
        go worker(w, jobs, results)
    }

    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= 9; a++ {
        fmt.Println(<-results)
    }
}
```

### 13.2 Fan-Out / Fan-In

**Fan-out** distributes work across multiple goroutines reading from the same channel; **fan-in** merges multiple result channels back into a single channel.

```go
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(channels))

    for _, c := range channels {
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(c)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### 13.3 Pipeline

Chain multiple stages together, each stage a goroutine reading from an input channel and writing to an output channel, forming a processing pipeline:

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

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// usage: for v := range square(generate(1, 2, 3, 4)) { fmt.Println(v) }
```

**Use case:** Stream processing where data flows through discrete transformation stages (e.g., read → parse → validate → write), each running concurrently.

## 14. Goroutine Leaks

A **goroutine leak** happens when a goroutine is started but never terminates — it stays blocked indefinitely (on a channel operation, a lock, or a blocking I/O call) and is never garbage collected, because the Go runtime cannot know a blocked goroutine will never be needed again. Unlike a typical memory leak, a leaked goroutine also holds onto whatever resources it was using — file descriptors, network connections, mutex locks — for as long as the process runs.

### 14.1 Common Causes

- **A goroutine sends on a channel with no receiver** (and vice versa), blocking forever with nothing to unblock it.
- **Forgetting to listen for a cancellation signal**, so a long-running or infinite loop never gets a chance to exit even when its work is no longer needed.
- **A `range` over a channel that is never closed**, leaving the consuming goroutine parked indefinitely.
- **Deadlocks between goroutines**, each waiting on a resource the other holds.
- **Unbounded goroutine creation** during request handling, without any cap, which can also exhaust memory even without a true "leak."

### 14.2 Example of a Leak

```go
func leaky() {
    ch := make(chan int) // unbuffered, no receiver
    go func() {
        ch <- 1 // blocks forever — nobody ever reads this
    }()
    // function returns; the goroutine above is now stuck permanently
}
```

### 14.3 Fixing It with `context`

```go
func fixed(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case ch <- 1:
        case <-ctx.Done(): // goroutine can exit even if nobody ever reads from ch
        }
    }()
}
```

### 14.4 Prevention Checklist

- Give every goroutine a guaranteed exit path — typically `ctx.Done()` inside a `select`.
- Close channels when the sending side is done, via `defer close(ch)`, so `range`-based readers terminate naturally.
- Use `select` with `time.After(...)` or a context timeout to avoid blocking on a channel operation forever.
- Bound concurrency explicitly (worker pools, semaphores) instead of spawning an unbounded number of goroutines per request.
- Monitor `runtime.NumGoroutine()` over time in long-running services, and use `pprof`'s goroutine profile to inspect what's actually blocked in production.
- In tests, use Uber's `go.uber.org/goleak` package (`defer goleak.VerifyNone(t)`) to catch leaked goroutines automatically.

## 15. Race Conditions and the Race Detector

A **data race** occurs when two or more goroutines access the same memory location concurrently, at least one of them is a write, and there is no synchronization between the accesses. The result is undefined behavior — output can vary from run to run, and can even corrupt memory in subtle ways.

```go
// UNSAFE: a classic data race
var counter int

func increment() {
    counter++ // read-modify-write, not atomic
}

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            increment() // race: many goroutines read/write `counter` unsynchronized
        }()
    }
    wg.Wait()
    fmt.Println(counter) // unpredictable — often less than 1000
}
```

Fixes: protect `counter` with a `sync.Mutex` (Section 8), use `sync/atomic`, or restructure so only one goroutine ever owns `counter` and others communicate with it via a channel.

### 15.1 The Built-in Race Detector

Go ships with a built-in race detector that instruments memory accesses and reports data races at runtime:

```bash
go run -race main.go
go test -race ./...
```

It's strongly recommended to run tests and, where feasible, staging/CI builds with `-race` enabled whenever concurrent code is written or modified — it will not catch every possible race (only ones actually exercised at runtime), but it catches a very large share of real bugs in practice.

## 16. Error Handling in Concurrent Code (`errgroup`)

Coordinating errors from multiple goroutines with a plain `WaitGroup` is awkward — there's no built-in way to know _which_ goroutine failed or to cancel the others. `golang.org/x/sync/errgroup` extends `WaitGroup` with error propagation and optional context cancellation:

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) error {
    g, ctx := errgroup.WithContext(ctx)

    for _, url := range urls {
        url := url // capture loop variable (see Section 19.1)
        g.Go(func() error {
            return fetch(ctx, url)
        })
    }

    return g.Wait() // returns the first non-nil error, if any
}
```

If any goroutine's function returns a non-nil error, `errgroup.WithContext`'s derived `ctx` is canceled, signaling the other goroutines to stop early — useful for "first error wins, cancel the rest" workflows like fetching multiple URLs where one failure should abort the whole batch.

**Use case:** Any batch of concurrent operations where you need to know if _any_ of them failed, and ideally want the remaining ones to stop early once a failure occurs.

## 17. Debugging and Observability

| Tool                             | What it does                                                                                                                                           |
| -------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `runtime.NumGoroutine()`         | Returns the current number of live goroutines — cheap to sample periodically as a leak indicator.                                                      |
| `net/http/pprof`                 | Exposes live profiling endpoints (`/debug/pprof/goroutine`, `/debug/pprof/block`, etc.) for inspecting what goroutines are doing in a running process. |
| `go tool trace`                  | Visualizes scheduler activity, goroutine execution, and blocking events over time — useful for diagnosing scheduling and latency issues.               |
| `go run -race` / `go test -race` | Built-in data race detector (Section 15).                                                                                                              |
| `go.uber.org/goleak`             | Third-party library for asserting "no goroutines leaked" at the end of a test.                                                                         |

## 18. Best Practices Summary

1. **Always have a plan for how a goroutine will stop.** Every goroutine should have a clear, reachable exit condition — don't fire off a goroutine "and forget it."
2. **Pass `context.Context` as the first argument** to any function that does I/O, blocks, or spawns further goroutines, and respect `ctx.Done()`.
3. **Use `sync.WaitGroup` (or `errgroup`) to wait for goroutines you launch** rather than `time.Sleep` hacks.
4. **Protect shared mutable state** with a `sync.Mutex`/`sync.RWMutex`, or avoid sharing it altogether by communicating over channels instead.
5. **Prefer bounded concurrency** (worker pools, semaphores via buffered channels) over spawning an unbounded goroutine per unit of work, especially per incoming HTTP request.
6. **Close channels from the sending side only**, and only when consumers actually need an end-of-stream signal.
7. **Run tests and CI with `-race` enabled** whenever concurrent code changes.
8. **Don't change `GOMAXPROCS` without a measured reason** — the default is correct for the overwhelming majority of programs.
9. **Keep goroutines' responsibilities small and well-defined** — it's easier to reason about correctness and shutdown when each goroutine does one clear thing.
10. **Monitor goroutine counts and use `pprof` in production** for long-running services, since leaks are invisible until they cause real damage (memory growth, exhausted connections, latency spikes).

## 19. Common Mistakes

### 19.1 Capturing the Loop Variable Incorrectly

```go
// BEFORE Go 1.22: classic bug
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i) // may print 5, 5, 5, 5, 5 — not 0,1,2,3,4
    }()
}
```

Prior to Go 1.22, `for` loop variables were reused across iterations, so a goroutine's closure could see a variable's _final_ value by the time it actually ran, instead of the value at the time `go` was called. The historical fix was to capture it explicitly:

```go
for i := 0; i < 5; i++ {
    i := i // shadow: create a new variable per iteration
    go func() {
        fmt.Println(i)
    }()
}
// or, equivalently, pass i as a function parameter:
for i := 0; i < 5; i++ {
    go func(i int) {
        fmt.Println(i)
    }(i)
}
```

**Note:** Go 1.22 changed loop semantics so that each iteration of a `for` loop gets its own copy of the loop variable, which eliminates this specific bug for code compiled with Go 1.22+ language version. Still, being explicit about what a goroutine's closure captures remains good practice, especially in codebases that must support older Go versions.

### 19.2 Deadlocks

A deadlock occurs when goroutines are mutually waiting on each other and none can proceed — for example, a channel send with no goroutine ever positioned to receive it.

```go
func main() {
    ch := make(chan int)
    ch <- 1 // deadlock: no other goroutine to receive; main blocks forever
}
```

Go's runtime can detect certain "all goroutines are asleep" deadlocks and will crash the program with `fatal error: all goroutines are asleep - deadlock!`, which is helpful during development.

### 19.3 Forgetting `wg.Add()` Placement

```go
// WRONG: Add() called inside the goroutine — race between Add and Wait
go func() {
    wg.Add(1)
    defer wg.Done()
    doWork()
}()
wg.Wait() // might return before Add(1) ever executes
```

Always call `Add()` in the launching goroutine, before starting the new one.

### 19.4 Treating Goroutines as Free

Because goroutines are cheap, it's tempting to spawn one per item in an unbounded loop (e.g., one per incoming request, one per row in a huge dataset) without any cap. Under load, this can exhaust memory or downstream resources (DB connections, file descriptors) even though no individual goroutine is "leaked" in the strict sense — bound concurrency explicitly instead.

### 19.5 Using `sync.Mutex` and Channels Interchangeably Without Thought

Both can solve synchronization problems, but they fit different situations: a mutex is usually simpler for protecting a piece of shared state accessed from many places, while channels are usually clearer for passing ownership of data between goroutines or signaling events. Mixing both without a clear reason can make code harder to reason about.

## 20. Use Case Summary Table

| Technique                           | When to Use                                                                   | Example Scenario                                                    |
| ----------------------------------- | ----------------------------------------------------------------------------- | ------------------------------------------------------------------- |
| Plain `go func()`                   | Fire a short, independent task that doesn't need to report back               | Logging, sending a metric, best-effort notification                 |
| `sync.WaitGroup`                    | Wait for a known, fixed number of goroutines to finish                        | Fetching N URLs concurrently, then aggregating results              |
| `sync.Mutex` / `RWMutex`            | Protect shared mutable state accessed by multiple goroutines                  | In-memory cache, shared counter, connection pool state              |
| Unbuffered channel                  | Strict rendezvous / hand-off between exactly one sender and one receiver      | Signaling "task done," passing ownership of a single value          |
| Buffered channel                    | Decouple producer/consumer timing, or implement a semaphore                   | Worker pool job queue, limiting concurrent DB connections           |
| `select` + `time.After`             | Prevent a goroutine from blocking forever                                     | Network call with a timeout                                         |
| `context.WithCancel/Timeout`        | Propagate cancellation/deadlines across goroutines and API boundaries         | HTTP request handling, background workers, graceful shutdown        |
| Worker pool pattern                 | Bound concurrency for a large or unbounded stream of tasks                    | Processing a queue of jobs with limited DB/API connections          |
| Fan-out/fan-in                      | Parallelize independent work, then merge results                              | Concurrently querying multiple services and combining responses     |
| `errgroup`                          | Run several goroutines, fail fast on the first error, get all errors' context | Concurrent multi-source data fetch that should abort on any failure |
| `go test -race` / `-race` builds    | Catch unsynchronized shared-memory access                                     | Any test suite covering concurrent code                             |
| `goleak`, `pprof`, `NumGoroutine()` | Detect and diagnose goroutines that never terminate                           | Long-running services suspected of leaking goroutines               |

## 21. References

1. Go Team — _Effective Go: Concurrency_. https://go.dev/doc/effective_go
2. Go Team — _Go Documentation_. https://go.dev/doc/
3. Go Team — _Go Wiki: LearnConcurrency_. https://go.dev/wiki/LearnConcurrency
4. Rob Pike — _Share Memory By Communicating_ (Go Blog / codewalk), referenced via Go Wiki LearnConcurrency. https://go.dev/wiki/LearnConcurrency
5. GetStream Engineering Blog — _Goroutines in Go: A Practical Guide to Concurrency_. https://getstream.io/blog/goroutines-go-concurrency-guide/
6. Brandon Wofford — _The Comprehensive Guide to Concurrency in Golang_, Medium. https://bwoff.medium.com/the-comprehensive-guide-to-concurrency-in-golang-aaa99f8bccf6
7. Sarvesh Sharma — _Understanding the Go Scheduler: The GMP Model Explained_, Medium. https://medium.com/@sharmasarvesh826/understanding-the-go-scheduler-the-gmp-model-explained-dee532c15c5f
8. Ram Prawesh Kumar — _Go Scheduler Deep Dive: Understanding the GMP Model and Tracing_, Medium. https://medium.com/@rpglearnai/go-scheduler-deep-dive-understanding-the-gmp-model-and-tracing-823a1be20133
9. Priyanka Guha — _The Go Scheduler Explained: Deep Dive into the G-M-P Concurrency Model_, Medium. https://medium.com/@priyankaguha.2012/the-go-scheduler-explained-deep-dive-into-the-g-m-p-concurrency-model-dfaf40a7508c
10. DEV Community (debianbaker) — _Inside the Go Scheduler: How GMP Model Powers Millions of Goroutines_. https://dev.to/debianbaker/inside-the-go-scheduler-how-gmp-model-powers-millions-of-goroutines-940
11. Leapcell — _Unveiling Go's Scheduler Secrets: The G-M-P Model in Action_. https://leapcell.io/blog/unveiling-go-s-scheduler-secrets-the-g-m-p-model-in-action
12. Educative — _The Go Scheduler_. https://www.educative.io/courses/advanced-techniques-in-go-programming/np/the-go-scheduler
13. OneUptime Engineering Blog — _How to Avoid Common Goroutine Leaks in Go_. https://oneuptime.com/blog/post/2026-01-07-go-goroutine-leaks/view
14. Serif Colakel — _Go Concurrency Mastery: Preventing Goroutine Leaks with Context, Timeout & Cancellation Best Practices_, Medium. https://medium.com/@serifcolakel/go-concurrency-mastery-preventing-goroutine-leaks-with-context-timeout-cancellation-best-ade0573d7532
15. Or Ben Shmueli — _Handling Goroutine Leaks — When Concurrency Goes Wrong in Go_, Medium. https://medium.com/@orbens/handling-goroutine-leaks-when-concurrency-goes-wrong-in-go-4c5d8ff36b10
16. Sogol Hedayatmanesh — _Goroutine Leaks in Go: Root Causes, Real-World Examples, and Ironclad Detection Strategies_, Medium. https://medium.com/@sogol.hedayatmanesh/goroutine-leaks-in-go-root-causes-real-world-examples-and-ironclad-detection-strategies-435c938d66ed
17. Gopher (author handle) — _Preventing Goroutine Leaks in Go: Diagnosis and Defensive Coding Techniques_, Medium. https://medium.com/@gane18/preventing-goroutine-leaks-in-go-diagnosis-and-defensive-coding-techniques-2ef1144cc99d
18. Serif Colakel — _Detecting and Preventing Goroutine Leaks in Production_, Medium. https://medium.com/@serifcolakel/detecting-and-preventing-goroutine-leaks-in-production-35f9cbc61f16
19. knowledgelib.io — _How to Detect and Fix Goroutine Leaks in Go_. https://knowledgelib.io/software/debugging/go-goroutine-leak/2026
20. The Golang Blueprint — _Context Leaks: The Hidden Technical Debt in Go Systems_. https://thegolangblueprint.substack.com/p/context-leaks-the-hidden-technical
21. Go standard library documentation — `sync` package. https://pkg.go.dev/sync
22. Go standard library documentation — `context` package. https://pkg.go.dev/context
23. `golang.org/x/sync/errgroup` documentation. https://pkg.go.dev/golang.org/x/sync/errgroup
24. `go.uber.org/goleak` documentation. https://pkg.go.dev/go.uber.org/goleak
25. Go Team — _Go 1.22 Release Notes_ (for-loop variable scoping change). https://go.dev/doc/go1.22
