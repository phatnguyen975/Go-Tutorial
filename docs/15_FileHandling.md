<div align="center">
  <h1>File Handling</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>September 18, 2026</sub>
</div>

## Table of Contents

1. [Overview: Which Packages Are Involved](#1-overview-which-packages-are-involved)
2. [Opening and Closing Files](#2-opening-and-closing-files)
3. [Reading an Entire File at Once](#3-reading-an-entire-file-at-once)
4. [Writing an Entire File at Once](#4-writing-an-entire-file-at-once)
5. [File Flags and Permissions with `os.OpenFile`](#5-file-flags-and-permissions-with-osopenfile)
6. [Reading Files Incrementally (Streaming)](#6-reading-files-incrementally-streaming)
7. [Reading Line by Line with `bufio.Scanner`](#7-reading-line-by-line-with-bufioscanner)
8. [Writing Incrementally and Buffered Writes](#8-writing-incrementally-and-buffered-writes)
9. [Copying Data Between Files with `io.Copy`](#9-copying-data-between-files-with-iocopy)
10. [Appending to a File](#10-appending-to-a-file)
11. [File Metadata: `os.Stat` and `os.FileInfo`](#11-file-metadata-osstat-and-osfileinfo)
12. [Checking Whether a File Exists](#12-checking-whether-a-file-exists)
13. [Working with Directories](#13-working-with-directories)
14. [Walking a Directory Tree](#14-walking-a-directory-tree)
15. [Path Manipulation with `path/filepath`](#15-path-manipulation-with-pathfilepath)
16. [Renaming, Removing, and Copying Files](#16-renaming-removing-and-copying-files)
17. [Temporary Files and Directories](#17-temporary-files-and-directories)
18. [Seeking Within a File](#18-seeking-within-a-file)
19. [Embedding Files at Compile Time (`embed`)](#19-embedding-files-at-compile-time-embed)
20. [Error Handling Conventions for File Operations](#20-error-handling-conventions-for-file-operations)
21. [Common Mistakes and Pitfalls](#21-common-mistakes-and-pitfalls)
22. [Best Practices Summary](#22-best-practices-summary)
23. [Use Case Summary Table](#23-use-case-summary-table)
24. [References](#24-references)

## 1. Overview: Which Packages Are Involved

Go's standard library spreads file-handling functionality across a small number of focused packages, rather than one monolithic "file I/O" package:

| Package         | Responsibility                                                                                                                                                                    |
| --------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `os`            | The primary entry point: opening, creating, reading, writing, and removing files; directory operations; file metadata; environment/process-level file access.                     |
| `io`            | Defines the core streaming interfaces (`io.Reader`, `io.Writer`, `io.Closer`, etc.) and generic helpers like `io.Copy` that work across any type satisfying them, not just files. |
| `bufio`         | Buffered wrappers around readers/writers (`bufio.Reader`, `bufio.Writer`, `bufio.Scanner`) for more efficient, line-oriented, or chunked I/O.                                     |
| `path/filepath` | Platform-independent path manipulation and directory-tree traversal.                                                                                                              |
| `io/fs`         | Defines a read-only filesystem abstraction (`fs.FS`) and shared filesystem-related types/errors (like `fs.ErrNotExist`) used across `os`, `embed`, and other packages.            |
| `embed`         | Compile-time embedding of files directly into the compiled binary.                                                                                                                |

**Historical note:** before Go 1.16, many convenience functions for whole-file reads/writes lived in a separate `io/ioutil` package (`ioutil.ReadFile`, `ioutil.WriteFile`, `ioutil.ReadDir`, etc.). Go 1.16 moved this functionality into `os` and `io` directly (`os.ReadFile`, `os.WriteFile`, `os.ReadDir`), and `io/ioutil` is now considered legacy — new code should use the `os`/`io` equivalents.

The type at the center of everything is **`*os.File`**, the concrete type returned by `os.Open`, `os.Create`, and `os.OpenFile`. It implements `io.Reader`, `io.Writer`, `io.Closer`, `io.Seeker`, and related interfaces, which is precisely why file handles can be passed directly into any function that accepts one of those general-purpose interfaces (like `bufio.NewScanner`, `io.Copy`, or `json.NewDecoder`).

## 2. Opening and Closing Files

### 2.1 `os.Open` — Read-Only

```go
package main

import (
    "log"
    "os"
)

func main() {
    file, err := os.Open("data.txt") // opens read-only
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // use file...
}
```

### 2.2 `os.Create` — Write, Truncate or Create

```go
file, err := os.Create("output.txt") // creates if missing, truncates if it exists, opens read-write
if err != nil {
    log.Fatal(err)
}
defer file.Close()
```

`os.Create` is shorthand for `os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)` — it always truncates an existing file to zero length, so use it only when you genuinely want to start the file fresh.

### 2.3 The Safe Close Pattern

**Always check the error before deferring `Close()`.** Closing a `nil` `*os.File` panics, so the defer must come only after confirming the open succeeded:

```go
file, err := os.Open("data.txt")
if err != nil {
    log.Fatal(err) // return/handle before touching file at all
}
defer file.Close() // only reached if err == nil
```

`Close()` itself can fail (e.g., a final flush to disk fails, or a network filesystem drops the connection) and returns an `error`. For most read-only operations, a `defer file.Close()` without checking that error's return value is common and acceptable — but for writes where data integrity really matters, check the error explicitly instead of relying purely on `defer`:

```go
func writeImportantData(path string, data []byte) (err error) {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = cerr // surface a close failure if nothing else already failed
        }
    }()

    _, err = f.Write(data)
    return err
}
```

## 3. Reading an Entire File at Once

For files you know are reasonably sized (config files, small data files), `os.ReadFile` reads the whole content into memory in one call, handling opening and closing internally:

```go
data, err := os.ReadFile("config.json")
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(data))
```

`os.ReadFile` returns a `[]byte` containing the entire file's contents. This is the simplest option, but it is **not appropriate for very large files**, since the whole file must fit in memory at once — for large or unbounded-size files, prefer the streaming approaches in [Section 6](#6-reading-files-incrementally-streaming) onward.

## 4. Writing an Entire File at Once

Symmetrically, `os.WriteFile` creates (or truncates) a file and writes a full byte slice to it in one call:

```go
data := []byte("Hello, Gophers!\n")
err := os.WriteFile("output.txt", data, 0644)
if err != nil {
    log.Fatal(err)
}
```

The third argument (`0644`) is the file's permission bits, used only when the file doesn't already exist (see [Section 5](#5-file-flags-and-permissions-with-osopenfile)). Like `os.ReadFile`, this is convenient for small, complete writes, but streaming is preferable for large amounts of data.

## 5. File Flags and Permissions with `os.OpenFile`

`os.OpenFile` is the fully general form underlying both `os.Open` and `os.Create`, giving explicit control over how the file is opened and what permissions a newly created file should get:

```go
file, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
if err != nil {
    log.Fatal(err)
}
defer file.Close()
```

### 5.1 Common Flags (combined with bitwise OR `|`)

| Flag          | Meaning                                                                             |
| ------------- | ----------------------------------------------------------------------------------- |
| `os.O_RDONLY` | Open read-only                                                                      |
| `os.O_WRONLY` | Open write-only                                                                     |
| `os.O_RDWR`   | Open for both reading and writing                                                   |
| `os.O_APPEND` | Append writes to the end of the file, rather than the current offset                |
| `os.O_CREATE` | Create the file if it doesn't already exist                                         |
| `os.O_EXCL`   | Used with `O_CREATE`: fail if the file already exists (atomic "create only if new") |
| `os.O_TRUNC`  | Truncate the file to zero length when opened                                        |
| `os.O_SYNC`   | Make writes synchronous (bypass OS write buffering — rarely needed)                 |

### 5.2 Permission Bits

The final argument to `os.OpenFile` (and `os.WriteFile`, `os.MkdirAll`, etc.) sets Unix-style **file permission bits**, expressed as an octal literal — three digits representing the owner, group, and others' permissions, each a sum of read (4), write (2), and execute (1):

| Value  | Meaning                                                                                            |
| ------ | -------------------------------------------------------------------------------------------------- |
| `0644` | Owner: read/write; Group and Others: read-only — a very common default for regular files           |
| `0600` | Owner: read/write; Group and Others: no access — appropriate for files containing secrets          |
| `0755` | Owner: read/write/execute; Group and Others: read/execute — common for directories and executables |
| `0777` | Everyone: full read/write/execute — rarely appropriate; avoid unless genuinely required            |

This value is only applied when the operation actually creates a new file — opening an existing file with `os.OpenFile` never changes its existing permissions, and the actual resulting permissions are also filtered through the operating system's `umask` on Unix-like systems, which can further restrict the effective bits. On Windows, most of these Unix-style permission bits have limited or no effect, since Windows uses a different, ACL-based permission model.

## 6. Reading Files Incrementally (Streaming)

For large files, read data in fixed-size chunks instead of loading everything into memory at once, using the raw `Read` method (satisfying `io.Reader`) in a loop:

```go
package main

import (
    "fmt"
    "io"
    "log"
    "os"
)

func main() {
    file, err := os.Open("largefile.dat")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    buffer := make([]byte, 4096) // 4KB chunks
    totalBytes := 0

    for {
        n, err := file.Read(buffer)
        if n > 0 {
            totalBytes += n
            // process buffer[:n] here
        }
        if err == io.EOF {
            break // normal, expected end of file
        }
        if err != nil {
            log.Fatal(err)
        }
    }

    fmt.Println("Total bytes read:", totalBytes)
}
```

**Key detail:** `Read` can return a non-zero `n` **and** a non-nil `err` in the same call (for example, returning the last few bytes together with `io.EOF`) — always process `buffer[:n]` before checking the error, and treat `io.EOF` as the normal, expected signal that there's nothing left to read, not as a failure to report.

For finer control over buffering size and behavior, `bufio.NewReaderSize(file, size)` wraps a file in a buffered reader that reduces the number of underlying system calls:

```go
reader := bufio.NewReaderSize(file, 32*1024) // 32KB internal buffer
```

## 7. Reading Line by Line with `bufio.Scanner`

For text files where you want to process content line by line (logs, CSVs, config files), `bufio.Scanner` is the idiomatic tool:

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
)

func main() {
    file, err := os.Open("notes.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text() // current line, without the trailing newline
        fmt.Println(line)
    }

    if err := scanner.Err(); err != nil { // check for a scan error AFTER the loop
        log.Fatal(err)
    }
}
```

**Important:** `scanner.Scan()` returns `false` both when the input is exhausted normally _and_ when a genuine error occurred — always call `scanner.Err()` after the loop finishes to distinguish "reached the end successfully" (returns `nil`) from "something went wrong while scanning" (returns the actual error).

`bufio.Scanner` defaults to splitting input into lines (`bufio.ScanLines`), but can be configured with `scanner.Split(...)` to instead split by words (`bufio.ScanWords`), runes (`bufio.ScanRunes`), or a custom `bufio.SplitFunc` for specialized formats. Its default internal buffer has a maximum token size (64KB by default) — for files with unusually long individual lines, increase it with `scanner.Buffer(buf, maxSize)` before scanning, or switch to `bufio.Reader.ReadString('\n')` / `ReadBytes` for unbounded line lengths.

## 8. Writing Incrementally and Buffered Writes

### 8.1 Direct Writes

```go
file, err := os.OpenFile("example.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
if err != nil {
    log.Fatal(err)
}
defer file.Close()

_, err = file.WriteString("Hello, Gophers!\n")
if err != nil {
    log.Fatal(err)
}
```

Every call to `file.Write`/`file.WriteString` typically triggers a separate underlying system call, which can be inefficient for many small writes.

### 8.2 Buffered Writes with `bufio.Writer`

Wrapping a file in `bufio.Writer` accumulates writes in an in-memory buffer, flushing to the underlying file in larger batches, which significantly reduces the number of system calls for write-heavy code:

```go
file, err := os.Create("large_output.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

writer := bufio.NewWriter(file)
defer writer.Flush() // CRITICAL: without this, buffered data may never reach the file

for i := 0; i < 1000; i++ {
    fmt.Fprintf(writer, "Line %d\n", i)
}

if err := writer.Flush(); err != nil { // flush explicitly if you need to check for errors before the function ends
    log.Fatal(err)
}
```

**Critical rule:** data written through a `bufio.Writer` sits in memory until the buffer fills up or `Flush()` is called explicitly. Forgetting to flush before the program exits (or before the underlying file is closed) is a very common bug that silently drops buffered-but-unflushed data. `defer writer.Flush()` is a common safety net, but for cases where you need to actually detect a flush failure, call `Flush()` explicitly and check its returned error too, since a deferred call's error return is easy to ignore by accident.

## 9. Copying Data Between Files with `io.Copy`

`io.Copy(dst, src)` efficiently streams all data from any `io.Reader` to any `io.Writer`, using an internally managed buffer, without needing to load the whole source into memory:

```go
package main

import (
    "io"
    "log"
    "os"
)

func main() {
    src, err := os.Open("source.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer src.Close()

    dst, err := os.Create("destination.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer dst.Close()

    bytesWritten, err := io.Copy(dst, src)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Copied %d bytes\n", bytesWritten)
}
```

Because `io.Copy` operates purely in terms of the `io.Reader`/`io.Writer` interfaces, it works identically whether the source/destination are files, network connections, in-memory buffers, or any other type satisfying those interfaces — this is one of the clearest illustrations of why Go's small, composable interfaces are valuable in practice.

## 10. Appending to a File

To add content to the end of an existing file without erasing what's already there, open it with the `os.O_APPEND` flag (combined with `os.O_CREATE` if the file might not exist yet):

```go
file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
    log.Fatal(err)
}
defer file.Close()

if _, err := file.WriteString("new log entry\n"); err != nil {
    log.Fatal(err)
}
```

`O_APPEND` ensures every write goes to the current end of the file, atomically with respect to other appenders on most operating systems — this differs from opening the file normally and manually seeking to the end, which has a small race window between seeking and writing if multiple processes write concurrently.

## 11. File Metadata: `os.Stat` and `os.FileInfo`

`os.Stat(path)` retrieves metadata about a file or directory without opening it, returning an `fs.FileInfo` (aliased historically as `os.FileInfo`):

```go
info, err := os.Stat("data.txt")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Name:", info.Name())
fmt.Println("Size (bytes):", info.Size())
fmt.Println("Permissions:", info.Mode())
fmt.Println("Last modified:", info.ModTime())
fmt.Println("Is directory:", info.IsDir())
```

If you already have an open `*os.File`, you can call `.Stat()` directly on it instead of re-resolving the path:

```go
info, err := file.Stat()
```

`os.Lstat` is a variant that, on systems supporting symbolic links, returns information about the **link itself** rather than following it to the linked target — relevant when you specifically need to detect or inspect symlinks rather than transparently following them.

## 12. Checking Whether a File Exists

The idiomatic way to check for a file's existence relies on interpreting the error from `os.Stat`, rather than any dedicated "exists" function (Go's standard library deliberately doesn't provide one, since a check followed by a separate operation always has an inherent race condition against other processes anyway):

```go
func fileExists(path string) bool {
    _, err := os.Stat(path)
    if err == nil {
        return true
    }
    if errors.Is(err, os.ErrNotExist) {
        return false
    }
    // some other error (permissions, I/O failure, etc.) — treat cautiously
    return false
}
```

**Modern idiom:** use `errors.Is(err, os.ErrNotExist)` (or the equivalent `errors.Is(err, fs.ErrNotExist)`) rather than the older `os.IsNotExist(err)` function directly, since `errors.Is` correctly traverses any error-wrapping that might have occurred, whereas `os.IsNotExist` performs a more limited, non-wrapping-aware check. Both approaches recognize the same underlying condition, but `errors.Is` composes more reliably with wrapped errors from other layers of your application.

**Important caveat:** checking existence and then acting on it (e.g., "if it doesn't exist, create it") is inherently racy in a concurrent or multi-process environment — another process could create or remove the file between your check and your subsequent action. When atomicity actually matters, prefer flags like `os.O_CREATE|os.O_EXCL` (which atomically fails if the file already exists) over a separate existence check followed by a separate create.

## 13. Working with Directories

### 13.1 Creating Directories

```go
err := os.Mkdir("newdir", 0755)          // creates one directory; fails if parent doesn't exist
err := os.MkdirAll("path/to/newdir", 0755) // creates all necessary parent directories too, like `mkdir -p`
```

`os.MkdirAll` does **not** return an error if the directory already exists — it's safe to call even when you're not sure the path exists yet.

### 13.2 Listing Directory Contents

```go
entries, err := os.ReadDir(".")
if err != nil {
    log.Fatal(err)
}

for _, entry := range entries {
    fmt.Println(entry.Name(), entry.IsDir())
}
```

`os.ReadDir` (Go 1.16+) returns a slice of `fs.DirEntry`, which is more efficient than the older `ioutil.ReadDir` because it doesn't need to call `Stat` on every single entry unless you actually ask for full `FileInfo` via `entry.Info()`.

### 13.3 Removing Directories

```go
err := os.Remove("emptydir")     // removes a single empty directory (or a single file)
err := os.RemoveAll("some/path") // recursively removes a directory and everything inside it
```

**Warning:** `os.RemoveAll` is powerful and unforgiving — it recursively deletes everything under the given path with no confirmation and no recycle bin. Double-check any path passed to it, especially if built from user input or configuration.

## 14. Walking a Directory Tree

`filepath.WalkDir` (Go 1.16+) recursively visits every file and directory under a root path, calling a callback for each:

```go
package main

import (
    "fmt"
    "io/fs"
    "log"
    "path/filepath"
)

func main() {
    err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err // propagate errors encountered while walking (e.g., permission denied)
        }
        if d.IsDir() {
            fmt.Println("DIR: ", path)
        } else {
            fmt.Println("FILE:", path)
        }
        return nil // return a non-nil error to stop the walk early, or filepath.SkipDir to skip a directory
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

`filepath.WalkDir` is preferred over the older `filepath.Walk` because it uses `fs.DirEntry` (which avoids an extra `Stat` syscall per entry in many cases) rather than the heavier `os.FileInfo` that `filepath.Walk` requires for every entry, making it noticeably faster on large directory trees.

## 15. Path Manipulation with `path/filepath`

Building file paths by concatenating strings with `/` is not portable (Windows uses `\`) and is error-prone (double slashes, trailing slashes, `..` segments). The `path/filepath` package provides portable, correct path handling:

```go
import "path/filepath"

full := filepath.Join("dir", "subdir", "file.txt") // "dir/subdir/file.txt" (or "dir\subdir\file.txt" on Windows)

dir := filepath.Dir("/home/user/file.txt")   // "/home/user"
base := filepath.Base("/home/user/file.txt") // "file.txt"
ext := filepath.Ext("archive.tar.gz")         // ".gz"

abs, err := filepath.Abs("relative/path.txt") // resolves to an absolute path based on the current working directory
```

**Guidance:** always use `filepath.Join` instead of manual string concatenation (`+` or `fmt.Sprintf` with `/`) to build paths that need to work correctly across operating systems.

## 16. Renaming, Removing, and Copying Files

```go
err := os.Rename("old.txt", "new.txt") // renames (or moves, on the same filesystem) a file
err := os.Remove("unwanted.txt")        // deletes a single file
```

**Go has no dedicated "copy a file" function in the standard library.** The idiomatic approach is to open the source, create the destination, and use `io.Copy` (as shown in [Section 9](#9-copying-data-between-files-with-iocopy)) — this streaming approach also naturally works for files too large to fit comfortably in memory. Remember to also copy over any needed metadata explicitly (like the original file's permission bits, via `os.Chmod` on the destination) if that matters for your use case, since `io.Copy` only copies the byte content, not the source file's metadata.

`os.Rename` behaves atomically on most operating systems when the source and destination are on the same filesystem/volume, but may fail (or silently fall back to a non-atomic copy-then-delete on some platforms/tools, though not in the standard library's own implementation) when moving across different filesystems/volumes — check the returned error rather than assuming success across arbitrary paths.

## 17. Temporary Files and Directories

For scratch files/directories that should not collide with other concurrent processes and are meant to be short-lived, use the dedicated temp-file helpers rather than hand-rolling unique names:

```go
tmpFile, err := os.CreateTemp("", "myapp-*.tmp") // "" uses the OS default temp directory
if err != nil {
    log.Fatal(err)
}
defer os.Remove(tmpFile.Name()) // clean up when done
defer tmpFile.Close()

tmpDir, err := os.MkdirTemp("", "myapp-*")
if err != nil {
    log.Fatal(err)
}
defer os.RemoveAll(tmpDir)
```

The `*` in the pattern (`"myapp-*.tmp"`) is replaced with a random string to guarantee a unique name, avoiding collisions between concurrent runs of the same program or between multiple goroutines. These are the modern replacements for the deprecated `ioutil.TempFile`/`ioutil.TempDir`.

## 18. Seeking Within a File

`*os.File` implements `io.Seeker`, allowing the read/write position ("offset") within an open file to be moved explicitly, rather than always reading/writing sequentially from wherever the last operation left off:

```go
file, err := os.Open("data.bin")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

// Move to byte offset 100 from the start of the file
_, err = file.Seek(100, io.SeekStart)

// Move 10 bytes forward from the current position
_, err = file.Seek(10, io.SeekCurrent)

// Move to 50 bytes before the end of the file
_, err = file.Seek(-50, io.SeekEnd)
```

**Use case:** random-access binary file formats (databases, custom binary containers, certain media formats) where you need to jump directly to a known offset rather than reading everything sequentially up to that point.

## 19. Embedding Files at Compile Time (`embed`)

The `embed` package (Go 1.16+) lets you bundle static files directly into your compiled binary at build time, so they don't need to exist as separate files at runtime — useful for templates, static web assets, default configuration, and similar bundled resources.

```go
package main

import (
    "embed"
    "fmt"
)

//go:embed templates/*.html
var templatesFS embed.FS

func main() {
    data, err := templatesFS.ReadFile("templates/index.html")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))
}
```

`embed.FS` implements the read-only `fs.FS` interface, so embedded content can be used anywhere an `fs.FS` is accepted (for example, `html/template.ParseFS`, or serving static assets via `http.FileServerFS`), without needing the files to be physically present on disk alongside the deployed binary — this is particularly valuable for producing a single, self-contained deployable binary.

## 20. Error Handling Conventions for File Operations

Most `os` package functions that fail return a `*os.PathError` (or similar), which wraps the underlying operating-system error together with the operation name and the path involved — printing it directly already gives useful context:

```go
_, err := os.Open("missing.txt")
fmt.Println(err) // "open missing.txt: no such file or directory"
```

**Recommended patterns:**

```go
if errors.Is(err, os.ErrNotExist) {
    // handle "doesn't exist"
}
if errors.Is(err, os.ErrPermission) {
    // handle "permission denied"
}
if errors.Is(err, os.ErrExist) {
    // handle "already exists" (relevant with O_CREATE|O_EXCL)
}
```

These `os.ErrXxx` sentinel-style values (`os.ErrNotExist`, `os.ErrPermission`, `os.ErrExist`, and others) are portable across operating systems, even though the underlying OS-level error text and codes differ between platforms — always check against these sentinels rather than trying to match error message strings, which are not a stable or portable way to detect specific failure conditions.

Always add context when propagating a file-related error further up the call stack, since the path alone is often not obvious from a bare wrapped error at a higher layer:

```go
func loadConfig(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("loadConfig: %w", err)
    }
    return data, nil
}
```

## 21. Common Mistakes and Pitfalls

### 21.1 Deferring `Close()` Before Checking the Open Error

```go
file, err := os.Open("data.txt")
defer file.Close() // PANIC if err != nil and file is nil
if err != nil {
    log.Fatal(err)
}
```

Always check the error **first**, and only `defer file.Close()` in the branch where the open actually succeeded.

### 21.2 Forgetting to Flush a `bufio.Writer`

```go
writer := bufio.NewWriter(file)
writer.WriteString("important data")
// forgot writer.Flush() — this data may never actually reach the file!
```

Buffered data sits in memory until flushed; always `defer writer.Flush()` (and ideally check its error explicitly at the point you actually need durability guarantees).

### 21.3 Ignoring That `Read` Can Return Data _and_ an Error Together

```go
n, err := file.Read(buffer)
if err != nil {
    return err // BUG: might discard the last chunk of valid data in buffer[:n]
}
// process buffer[:n]
```

Process `buffer[:n]` before checking/acting on the error, since a final read can legitimately return both some bytes and `io.EOF` in the same call.

### 21.4 Loading an Entire Huge File into Memory

```go
data, _ := os.ReadFile("100gb_file.dat") // likely exhausts available memory
```

For files whose size isn't bounded or known to be small, use streaming (`bufio`, `io.Copy`, chunked `Read` loops) instead of `os.ReadFile`.

### 21.5 Building Paths with Manual String Concatenation

```go
path := dir + "/" + filename // breaks on Windows, mishandles trailing slashes
```

Use `filepath.Join(dir, filename)` instead, which handles platform-specific separators and normalization correctly.

### 21.6 Race Between Checking Existence and Acting On It

```go
if _, err := os.Stat(path); os.IsNotExist(err) {
    // another process could create the file HERE, between the check and the next line
    os.Create(path)
}
```

When atomicity matters, use `os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm)`, which fails atomically if the file already exists, instead of a separate check-then-act sequence.

### 21.7 Using String-Matching on Error Messages Instead of Sentinel Checks

```go
if strings.Contains(err.Error(), "no such file") { // fragile, non-portable, breaks across OSes/locales
```

Use `errors.Is(err, os.ErrNotExist)` (or the equivalent `fs.ErrNotExist`) instead — it's portable and doesn't depend on the exact wording of an OS-level error message.

### 21.8 Forgetting `os.RemoveAll`'s Destructive Reach

```go
os.RemoveAll(userSuppliedPath) // extremely dangerous if userSuppliedPath is attacker- or bug-controlled
```

Validate and sanitize any path that flows into a destructive operation like `os.RemoveAll`, especially if it originates from user input, configuration, or an upstream API response.

## 22. Best Practices Summary

1. **Always check the open/create error before deferring `Close()`** — closing a nil file handle panics.
2. **Use `os.ReadFile`/`os.WriteFile` for small, complete files; stream (`bufio`, `io.Copy`, chunked reads) for large or unbounded ones.**
3. **Always flush a `bufio.Writer`** before the program or function that owns it exits, and check the flush's error when durability genuinely matters.
4. **Prefer `errors.Is` against `os.ErrNotExist`/`os.ErrPermission`/`os.ErrExist`** over string-matching error messages or the older, non-wrap-aware `os.IsNotExist`.
5. **Use `filepath.Join` and other `path/filepath` helpers** for all path construction, instead of manual string concatenation.
6. **Use `os.O_CREATE|os.O_EXCL` for atomic "create only if new"** semantics, rather than a separate existence check followed by a separate create.
7. **Use `filepath.WalkDir` (not the older `filepath.Walk`)** for directory-tree traversal, since it avoids an unnecessary `Stat` per entry.
8. **Use `os.CreateTemp`/`os.MkdirTemp`** for scratch files/directories, rather than hand-rolling unique file names.
9. **Choose conservative file permissions** (`0600` for sensitive data, `0644` for ordinary files) rather than defaulting to overly permissive modes like `0777`.
10. **Treat `os.RemoveAll` and any destructive path operation with extra caution**, especially when the path is not a hardcoded, trusted literal.

## 23. Use Case Summary Table

| Technique                                        | When to Use                                                               | Example Scenario                                                 |
| ------------------------------------------------ | ------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| `os.ReadFile` / `os.WriteFile`                   | Small, complete file reads/writes                                         | Loading a config file, writing a small report                    |
| `os.Open` + `bufio.Scanner`                      | Line-oriented text processing                                             | Reading a log file or CSV line by line                           |
| `os.Open` + chunked `Read` loop / `bufio.Reader` | Large or binary files where full in-memory load isn't feasible            | Streaming-processing a multi-gigabyte data file                  |
| `os.OpenFile` with explicit flags                | Fine-grained control over create/append/truncate behavior and permissions | Append-only log file, atomic create-if-new                       |
| `bufio.Writer`                                   | Many small writes that benefit from batching                              | Generating a large output file line by line                      |
| `io.Copy`                                        | Streaming data from any reader to any writer                              | Copying a file, proxying an HTTP response body to disk           |
| `os.Stat` / `errors.Is(err, os.ErrNotExist)`     | Checking metadata or existence                                            | Determining whether a config file is present before loading it   |
| `os.Mkdir` / `os.MkdirAll` / `os.ReadDir`        | Directory creation and listing                                            | Setting up an app's data directory structure                     |
| `filepath.WalkDir`                               | Recursively processing an entire directory tree                           | Building a file index, searching for files by pattern            |
| `filepath.Join`/`Dir`/`Base`/`Ext`               | Any path construction or decomposition                                    | Building cross-platform file paths                               |
| `os.CreateTemp` / `os.MkdirTemp`                 | Scratch files/directories for intermediate processing                     | Downloading a file to a temp location before validating it       |
| `file.Seek`                                      | Random access within a file                                               | Reading specific offsets in a binary/database-like file format   |
| `embed.FS`                                       | Bundling static assets into the compiled binary                           | Shipping HTML templates or default config inside a single binary |

## 24. References

1. Go Team — _A Tour of Go_ and _Go Documentation_ (general reference for `os`, `io`, `bufio`). https://go.dev/doc/
2. Go standard library documentation — package `os`. https://pkg.go.dev/os
3. Go standard library documentation — package `io`. https://pkg.go.dev/io
4. Go standard library documentation — package `bufio`. https://pkg.go.dev/bufio
5. Go standard library documentation — package `path/filepath`. https://pkg.go.dev/path/filepath
6. Go standard library documentation — package `io/fs`. https://pkg.go.dev/io/fs
7. Go standard library documentation — package `embed`. https://pkg.go.dev/embed
8. Go Team — _Go 1.16 Release Notes_ (`io/ioutil` deprecation, `os.ReadFile`/`WriteFile`/`ReadDir`, `embed` package). https://go.dev/doc/go1.16
9. GoLinuxCloud — *Golang os package: Open, Create, WriteFile, and _os.File_. https://www.golinuxcloud.com/golang-os/
10. Leapcell — _Writing to Files in Go: A Comprehensive Guide_. https://leapcell.io/blog/writing-to-files-in-go-a-comprehensive-guide
11. ZetCode — _Working with Files in Go_. https://zetcode.com/golang/file/
12. Honeybadger Developer Blog — _A comprehensive guide to file operations in Go_. https://www.honeybadger.io/blog/comprehensive-guide-to-file-operations-in-go/
13. Reintech — _An Overview of Go's `os` and `io` Packages_. https://reintech.io/blog/an-overview-of-gos-os-and-io-packages
14. Reintech — _Reading and Writing Files in Go_. https://reintech.io/blog/reading-writing-files-go
15. Stanza — _File Operations with os Package - Go Standard Library Mastery_. https://www.stanza.dev/courses/go-std-lib/io-files/go-std-lib-os-files
16. tutorialpedia.org — _How to Read and Write Files in Go: A Beginner's Guide to File Handling with Examples_. https://www.tutorialpedia.org/blog/how-to-read-write-from-to-a-file-using-go/
17. golang.howtos.io — _Using the bufio Package for File I/O in Go_. https://golang.howtos.io/using-the-bufio-package-for-file-i-o-in-go/
