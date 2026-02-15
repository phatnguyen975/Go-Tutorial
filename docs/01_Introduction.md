<div align="center">
  <h1>Introduction</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 15, 2026</sub>
</div>

## Installation

Go binary distributions are available for all major operating systems like Linux, Windows, and MacOS. It’s super simple to install Go from the binary distributions.

If a binary distribution is not available for your operating system, you can try [installing Go from source](https://go.dev/doc/install/source).

### MacOS

**Using Homebrew**

The easiest way to install Go in MacOS is by using [Homebrew](https://brew.sh/):

```bash
brew install go
```

**Using MacOS package installer**

Download the latest Go package (`.pkg`) file from [Go's official downloads page](https://golang.org/dl/). Open the package and follow the on-screen instructions to install Go. By default, Go will be installed in `/usr/local/go`.

### Linux

Download the Linux distribution from [Go's official download page](https://golang.org/dl/) and extract it into `/usr/local` directory.

```bash
sudo tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
```

Next, add the `/usr/local/go/bin` directory to your `PATH` environment variable. You can do this by adding the following line to your `~/.bash_profile` file:

```bash
export PATH=$PATH:/usr/local/go/bin
```

You can also use any other directory like `/opt/go` instead of `/usr/local` for installing Go.

### Windows

Download the Windows MSI installer file from [Go's official download page](https://golang.org/dl/). Open the installer and follow the on-screen instructions to install Go in your windows system. By default, the installer installs Go in `C:\Program Files\Go`.

## Go Tool

The Go distribution comes bundled with the [go tool](https://golang.org/cmd/go/). It is a command line tool that lets you automate common tasks such as downloading and installing dependencies, building and testing your code, and much more.

After installing Go by following the instructions in the previous section, you should be able to run the Go tool by typing `go` in the command line

```bash
$ go
Go is a tool for managing Go source code.

Usage:

        go <command> [arguments]

The commands are:

        bug         start a bug report
        build       compile packages and dependencies
        clean       remove object files and cached files
        doc         show documentation for package or symbol
        env         print Go environment information
        fix         update packages to use new APIs
        fmt         gofmt (reformat) package sources
        generate    generate Go files by processing source
        get         add dependencies to current module and install them
        install     compile and install packages and dependencies
        list        list packages or modules
        mod         module maintenance
        work        workspace maintenance
        run         compile and run Go program
        telemetry   manage telemetry data and settings
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         report likely mistakes in packages

Use "go help <command>" for more information about a command.

Additional help topics:

        buildconstraint build constraints
        buildjson       build -json encoding
        buildmode       build modes
        c               calling between Go and C
        cache           build and test caching
        environment     environment variables
        filetype        file types
        goauth          GOAUTH environment variable
        go.mod          the go.mod file
        gopath          GOPATH environment variable
        goproxy         module proxy protocol
        importpath      import path syntax
        modules         modules, module versions, and more
        module-auth     module authentication using go.sum
        packages        package lists and patterns
        private         configuration for downloading non-public code
        testflag        testing flags
        testfunc        testing functions
        vcs             controlling version control with GOVCS

Use "go help <topic>" for more information about that topic.
```

## Code Organization

Go programs are organized into **packages**. A package is a collection of source files in the same directory that are compiled together. All the functions, types, variables, and constants defined in one source file are visible to all the other source files within the same package.

Go language has several built-in packages like:

- `fmt` package, which contains functions for formatting and printing text.
- `math` package, which provides basic constants and mathematical functions.

You need to import these packages in your program if you want to use the functions and constants defined in these packages. You can also import and use external packages built and published by other people on any source control management system like github.

Any Go source code repository contains one or more **modules**. A module is a collection of related Go packages stored in a directory with a `go.mod` file at its root. The `go.mod` file defines the module's path, which is the import path used while importing packages that are part of this module.

When you import packages contained in other modules, you manage those dependencies through your code’s own module defined by the `go.mod` file. The `go.mod` file keeps track of all the external modules that provide the packages used by your code.

```text
my-project/
├── app/
├── util/
└── go.mod
```

## Go Environment

### $GOPATH

Go is opinionated.

By convention, all Go code lives within a single workspace (folder). This workspace could be anywhere in your machine. If you don't specify, Go will assume `$HOME/go` as the default workspace. The workspace is identified (and modified) by the environment variable [GOPATH](https://pkg.go.dev/cmd/go#hdr-GOPATH_environment_variable).

You should set the environment variable so that you can use it later in scripts, shells, etc.

Go assumes that your workspace contains a specific directory structure. Go places its files in three directories: All source code lives in `src`, package objects lives in `pkg`, and the compiled programs live in `bin`. You can create these directories as follows:

```bash
mkdir -p $GOPATH/src $GOPATH/pkg $GOPATH/bin
```

At this point you can `go get` and the `src/package/bin` will be installed correctly in the appropriate `$GOPATH/xxx` directory.

### Go Modules

Go 1.11 introduced [Modules](https://github.com/golang/go/wiki/Modules), enabling an alternative workflow. This new approach will gradually [become the default](https://go.dev/blog/modules2019) mode, deprecating the use of `GOPATH`.

Modules aim to solve problems related to dependency management, version selection and reproducible builds; they also enable users to run Go code outside of `GOPATH`.

Using Modules is pretty straightforward. Select any directory outside `GOPATH` as the root of your project, and create a new module with the `go mod init` command.

A `go.mod` file will be generated, containing the module path, a Go version, and its dependency requirements, which are the other modules needed for a successful build.

If no `<module-path>` is specified, `go mod init` will try to guess the module path from the directory structure, but it can also be overrided, by supplying an argument.

```bash
mkdir my-project
cd my-project
go mod init <module-path>
```

A `go.mod` file could look like this:

```go
module module-path

go 1.25.7

require (
    github.com/google/pprof v0.0.0-20190515194954-54271f7e092f
    golang.org/x/arch v0.0.0-20190815191158-8a70ba74b3a1
    golang.org/x/tools v0.0.0-20190611154301-25a4f137592f
)
```

The built-in documentation provides an overview of all available `go mod` commands.

```bash
go help mod
go help mod init
```

## Go Debugger

A good option for debugging Go (that's integrated with VS Code) is **Delve**. This can be installed as follows:

```bash
go get -u github.com/go-delve/delve/cmd/dlv
```

For additional help configuring and running the Go debugger in VS Code, please reference the [VS Code debugging documentation](https://github.com/golang/vscode-go/blob/master/docs/debugging.md).

## Go Linting

An improvement over the default linter can be configured using [GolangCI-Lint](https://golangci-lint.run/). This can be installed as follows:

```bash
brew install golangci/tap/golangci-lint
```
