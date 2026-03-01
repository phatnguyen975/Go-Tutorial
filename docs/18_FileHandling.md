<div align="center">
  <h1>File Handling</h1>
  <small>
    <strong>Author:</strong> Nguyễn Tấn Phát
  </small> <br />
  <sub>February 27, 2026</sub>
</div>

## Reading Files

File reading is one of the most common operations performed in any programming language. In this tutorial, we will learn about how files can be read using Go.

### Reading an Entire File into Memory

One of the most basic file operations is reading an entire file into memory. This is done with the help of the `ReadFile` function of the [os](https://pkg.go.dev/os) package.

Let's read a file and print its contents.

I have created a folder `filehandling` inside my `Documents` directory by running `mkdir ~/Documents/filehandling`.

Create a Go module named `filehandling` by running the following command from the `filehandling` directory.

```bash
go mod init filehandling
```

I have a text file `test.txt` which will be read from our Go program `filehandling.go`. The `test.txt` contains the following string:

```text
Hello World. Welcome to file handling in Go.
```

Here is my directory structure.

```text
├── Documents
│   └── filehandling
│       ├── filehandling.go
│       ├── go.mod
│       └── test.txt
```

Let's get to the code right away. Create a file `filehandling.go` with the following contents.

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	contents, err := os.ReadFile("test.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("Contents of file:", string(contents))
}
```

Please run this program from your local environment as it's not possible to read files in the playground.

Line no. 9 of the program above reads the file and returns a byte slice which is stored in `contents`. In line no. 14 we convert `contents` to a `string` and display the contents of the file.

Please run this program from the location where `test.txt` is present.

If `test.txt` is located at `~/Documents/filehandling`, then run this program using the following steps:

```bash
cd ~/Documents/filehandling/
go install
filehandling
```

This program will print,

```text
Contents of file: Hello World. Welcome to file handling in Go.
```

If this program is run from any other location, for instance, try running the program from `~/Documents/`.

```bash
cd ~/Documents/
filehandling
```

It will print the following error.

```text
File reading error open test.txt: no such file or directory
```

The reason is Go is a compiled language. What `go install` does is, it creates a binary from the source code. The binary is independent of the source code and it can be run from any location. Since `test.txt` is not found in the location from which the binary is run, the program complains that it cannot find the file specified.

There are three ways to solve this problem. Let's discuss them one by one.

**1. Using absolute file path**

The simplest way to solve this problem is to pass the absolute file path. I have modified the program and changed the path to an absolute one in line no. 9. Please change this path to the absolute location of your `test.txt`.

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	contents, err := os.ReadFile("/Users/Admin/Documents/filehandling/test.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("Contents of file:", string(contents))
}
```

Now the program can be run from any location and it will print the contents of `test.txt`.

For example, it will work even when I run it from my home directory.

```bash
cd ~/Documents/filehandling
go install
cd ~
filehandling
```

The program will print the contents of `test.txt`.

This seems to be an easy way but comes with the pitfall that the file should be located in the path specified in the program else this method will fail.

**2. Passing the file path as a command line flag**

Another way to solve this problem is to pass the file path as a command line argument. Using the [flag](https://pkg.go.dev/flag) package, we can get the file path as input argument from the command line and then read its contents.

Let's first understand how the `flag` package works. The `flag` package has a [String](https://golang.org/pkg/flag/#String) function. This function accepts 3 arguments. The first is the name of the flag, second is the default value and the third is a short description of the flag.

Let's write a small program to read the file name from the command line. Replace the contents of `filehandling.go` with the following:

```go
package main

import (
	"flag"
	"fmt"
)

func main() {
	fptr := flag.String("fpath", "test.txt", "file path to read from")
	flag.Parse()
	fmt.Println("Value of fpath is", *fptr)
}
```

Line no. 8 of the program above, creates a string flag named `fpath` with default value `test.txt` and description `file path to read from` using the `String` function. This function returns the address of the string variable that stores the value of the flag.

> `flag.Parse()` should be called before accessing any flag.

We print the value of the flag in line no. 10. When this program is run using the command:

```bash
filehandling -fpath=/path-of-file/test.txt
```

We pass `/path-of-file/test.txt` as the value of the flag `fpath`.

This program outputs:

```text
Value of fpath is /path-of-file/test.txt
```

If the program is run using just `filehandling` without passing any `fpath`, it will print:

```text
Value of fpath is test.txt
```

Since `test.txt` is the default value of `fpath`.

The `flag` also provides a nicely formatted output of the different arguments that are available. This can be displayed by running:

```bash
filehandling --help
```

This command will print the following output.

```text
Usage of filehandling:
  -fpath string
        file path to read from (default "test.txt")
```

Now that we know how to read the file path from the command line, let's go ahead and finish our file reading program.

```go
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fptr := flag.String("fpath", "test.txt", "file path to read from")
	flag.Parse()
	contents, err := os.ReadFile(*fptr)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("Contents of file:", string(contents))
}
```

The program above reads the content of the file path passed from the command line. Run this program using the command:

```bash
filehandling -fpath=/path-of-file/test.txt
```

Please replace `/path-of-file/` with the absolute path of `test.txt`. For example, in my case, I ran the command:

```bash
filehandling --fpath=/Users/naveen/Documents/filehandling/test.txt
```

and the program printed.

```text
Contents of file: Hello World. Welcome to file handling in Go.
```

**3. Bundling the text file along with the binary**

The above option of getting the file path from the command line is good but there is an even better way to solve this problem. Wouldn't it be awesome if we are able to bundle the text file along with our binary? This is what we are going to do next.

The [embed](https://pkg.go.dev/embed) package from the standard library will help us achieve this.

After importing the `embed` package, the `//go:embed` directive can be used to read the contents of the file.

A program will make us understand things better.

Replace the contents of `filehandling.go` with the following,

```go
package main

import (
	_ "embed"
	"fmt"
)

//go:embed test.txt
var contents []byte

func main() {
	fmt.Println("Contents of file:", string(contents))
}
```

In line no. 4 of the program above, we import the `embed` package with a underscore prefix. The reason is because `embed` is not explicitly used in the code but the `//go:embed` comment in line no. 8 needs some preprocessing by the compiler. Since we need to import the package without any explicit usage, we prefix it with underscore to make the compiler happy. If not, the compiler will complain stating that the package is not used anywhere.

The `//go:embed test.txt` in line no. 8 tells the compiler to read the contents of `test.txt` and assign it to the variable following that comment. In our case `contents` variable will hold the contents of the file.

Run the program using the following commands.

```bash
cd ~/Documents/filehandling
go install
filehandling
```

and the program will print:

```text
Contents of file: Hello World. Welcome to file handling in Go.
```

Now the file is bundled along with the binary and it is available to the go binary irrespective of where it's executed from. For example, try running the program from a directory where `test.txt` doesn't reside.

```bash
cd ~/Documents
filehandling
```

The above command will also print the contents of the file.

Do note that the variable to which the contents of the file should be assigned to must be at the package level. Local variables won't work. Try changing the program to the following.

```go
package main

import (
	_ "embed"
	"fmt"
)

func main() {
	//go:embed test.txt
	var contents []byte
	fmt.Println("Contents of file:", string(contents))
}
```

The above program has `contents` as a local variable. The program will now fail to compile with the following error.

```text
./filehandling.go:9:4: go:embed cannot apply to var inside func
```

### Reading a File in Small Chunks

In the last section, we learned how to load an entire file into memory. When the size of the file is extremely large it doesn't make sense to read the entire file into memory especially if you are running low on RAM. A more optimal way is to read the file in small chunks. This can be done with the help of the [bufio](https://pkg.go.dev/bufio) package.

Let's write a program that reads our `test.txt` file in chunks of 3 bytes. Replace the contents of `filehandling.go` with the following,

```go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fptr := flag.String("fpath", "test.txt", "file path to read from")
	flag.Parse()

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	r := bufio.NewReader(f)
	b := make([]byte, 3)
	for {
		n, err := r.Read(b)
		if err == io.EOF {
			fmt.Println("finished reading file")
			break
		}
		if err != nil {
			fmt.Printf("Error %s reading file", err)
			break
		}
		fmt.Println(string(b[0:n]))
	}
}
```

In line no. 16 of the program above, we open the file using the path passed from the command line flag. In line no. 20, we defer the file closing.

Line no. 26 of the program above creates a new buffered reader. In the next line, we create a byte slice of length and capacity 3 into which the bytes of the file will be read.

The `Read` method in line no. 29 reads up to `len(b)` bytes i.e up to 3 bytes and returns the number of bytes read. We store the bytes returned in a variable `n`. In line no. 38, the slice is read from index `0` to `n-1`, i.e up to the number of bytes returned by the `Read` method and printed.

Once the end of the file is reached, `Read` will return an EOF error. We check for this error in line no. 30. The rest of the program is straight forward.

If we run the program above using the commands,

```bash
cd ~/Documents/filehandling
go install
filehandling -fpath=/path-of-file/test.txt
```

The following will be output:

```text
Hel
lo
Wor
ld.
 We
lco
me
to
fil
e h
and
lin
g i
n G
o.
finished reading file
```

> You can also `Seek` to a known location in the file and `Read` from there. Let's see [Go by Example](https://gobyexample.com/reading-files).

### Reading a File Line by Line

In the section, we will discuss how to read a file line by line using Go. This can done using the [bufio](https://pkg.go.dev/bufio) package.

Please replace the contents in `test.txt` with the following:

```text
Hello World. Welcome to file handling in Go.
This is the second line of the file.
We have reached the end of the file.
```

The following are the steps involved in reading a file line by line.

1. Open the file
2. Create a new scanner from the file
3. Scan the file and read it line by line.

Replace the contents of `filehandling.go` with the following:

```go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	fptr := flag.String("fpath", "test.txt", "file path to read from")
	flag.Parse()

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
    defer func() {
	    if err = f.Close(); err != nil {
		log.Fatal(err)
	}
	}()
	s := bufio.NewScanner(f)
	for s.Scan() {
		fmt.Println(s.Text())
	}
	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}
}
```

In line no. 15 of the program above, we open the file using the path passed from the command line flag. In line no. 24, we create a new scanner using the file. The `Scan()` method in line no. 25 reads the next line of the file and the string that is read will be available through the `Text()` method.

After `Scan()` returns `false`, the `Err()` method will return any error that occurred during scanning. If the error is End of File, `Err()` will return `nil`.

If we run the program above using the commands,

```bash
cd ~/Documents/filehandling
go install
filehandling -fpath=/path-of-file/test.txt
```

the contents of the file will be printed line by line as shown below.

```text
Hello World. Welcome to file handling in Go.
This is the second line of the file.
We have reached the end of the file.
```

## Writing Files

### Writing String to a File

One of the most common file writing operations is writing a string to a file. This is quite simple to do. It consists of the following steps.

1. Create the file
2. Write the string to the file

Let's get to the code right away.

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString("Hello World")
	if err != nil {
		fmt.Println(err)
        f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

The `Create` function in line no. 9 of the program above creates a file named `test.txt`. If a file with that name already exists, then the create function truncates the file. This function returns a [File descriptor](https://pkg.go.dev/os#File).

In line no 14, we write the string `Hello World` to the file using the `WriteString` method. This method returns the number of bytes written and error if any. Finally, we close the file in line no. 21.

The above program will print:

```text
11 bytes written successfully
```

You can find a file named `test.txt` created in the directory from which this program was executed. If you open the file using any text editor, you can find that it contains the text `Hello World`.

### Writing Bytes to a File

Writing bytes to a file is quite similar to writing a string to a file. We will use the [Write](https://pkg.go.dev/os#File.Write) method to write bytes to a file. The following program writes a slice of bytes to a file.

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("/home/admin/bytes")
	if err != nil {
		fmt.Println(err)
		return
	}
	d2 := []byte{104, 101, 108, 108, 111, 32, 98, 121, 116, 101, 115}
	n2, err := f.Write(d2)
	if err != nil {
		fmt.Println(err)
        f.Close()
		return
	}
	fmt.Println(n2, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

In the program above, in line no. 15 we use the `Write` method to write a slice of bytes to a file named `bytes` in the directory `/home/admin`. You can change this directory to a different one. The remaining program is self-explanatory. This program will print `11 bytes written successfully` and it will create a file named `bytes`. Open the file and you can see that it contains the text `hello bytes`.

### Writing Strings Line by Line to a File

Another common file operation is the need to write strings to a file line by line. In this section, we will write a program to create a file with the following content.

```text
Welcome to the world of Go.
Go is a compiled language.
It is easy to learn Go.
```

Let's get to the code right away.

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("lines")
	if err != nil {
		fmt.Println(err)
                f.Close()
		return
	}
	d := []string{"Welcome to the world of Go1.", "Go is a compiled language.", "It is easy to learn Go."}

	for _, v := range d {
		if _, err := fmt.Fprintln(f, v); err != nil {
			fmt.Println(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("file written successfully")
}
```

In line no. 9 of the program above, we create a new file named `lines`. In line no. 17 we iterate through the array using a for range loop and use the [Fprintln](https://golang.org/pkg/fmt/#Fprintln) function to write the lines to a file. The `Fprintln` function takes a `io.writer` as parameter and appends a new line, just what we wanted. Running this program will print `file written successfully` and a file `lines` will be created in the current directory. The content of the file `lines` is provided below.

```text
Welcome to the world of Go1.
Go is a compiled language.
It is easy to learn Go.
```

### Appending to a File

In this section, we will append one more line to the `lines` file which we created in the previous section. We will append the line `File handling is easy` to the `lines` file.

The file has to be opened in append and write only mode. These flags are passed as parameters to the [Open](https://pkg.go.dev/os#OpenFile) function. After the file is opened in append mode, we add the new line to the file.

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.OpenFile("lines", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	newLine := "File handling is easy."
	_, err = fmt.Fprintln(f, newLine)
	if err != nil {
		fmt.Println(err)
                f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("file appended successfully")
}
```

In line no. 9 of the program above, we open the file in append and write only mode. After the file is opened successfully, we add a new line to the file in line no. 15. This program will print `file appended successfully`. After running this program, the contents of the `lines` file will be,

```text
Welcome to the world of Go1.
Go is a compiled language.
It is easy to learn Go.
File handling is easy.
```

### Writing to File Concurrently

When multiple [goroutines](https://golangbot.com/goroutines/) write to a file concurrently, we will end up with a [race condition](https://golangbot.com/mutex/#criticalsection). Hence concurrent writes to a file must be coordinated using a channel.

We will write a program that creates 100 goroutines. Each of this goroutine will generate a random number concurrently, thus generating hundred random numbers in total. These random numbers will be written to a file. We will solve the race condition problem by using the following approach.

1. Create a channel that will be used to read and write the generated random numbers.
2. Create 100 producer goroutines. Each goroutine will generate a random number and will also write the random number to a channel.
3. Create a consumer goroutine that will read from the channel and write the generated random number to the file. Thus we have only one goroutine writing to a file concurrently thereby avoiding race condition.
4. Close the file once done.

Let's write the `produce` function first which generates the random numbers.

```go
func produce(data chan int, wg *sync.WaitGroup) {
	n := rand.Intn(999)
	data <- n
	wg.Done()
}
```

The function above generates a random number and writes it to the channel `data` and then calls `Done` on the [waitgroup](https://golangbot.com/buffered-channels-worker-pools/#waitgroup) to notify that it is done with its task.

Let's move to the function which writes to the file now.

```go
func consume(data chan int, done chan bool) {
	f, err := os.Create("concurrent")
	if err != nil {
		fmt.Println(err)
		return
	}
	for d := range data {
		_, err = fmt.Fprintln(f, d)
		if err != nil {
			fmt.Println(err)
			f.Close()
			done <- false
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		done <- false
		return
	}
	done <- true
}
```

The `consume` function creates a file named `concurrent`. It then reads the random numbers from the `data` channel and writes to the file. Once it has read and written all the random numbers, it writes `true` to the `done` channel to notify that it's done with its task.

Let's write the `main` function and complete this program. I have provided the entire program below.

```go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
)

func produce(data chan int, wg *sync.WaitGroup) {
	n := rand.Intn(999)
	data <- n
	wg.Done()
}

func consume(data chan int, done chan bool) {
	f, err := os.Create("concurrent")
	if err != nil {
		fmt.Println(err)
		return
	}
	for d := range data {
		_, err = fmt.Fprintln(f, d)
		if err != nil {
			fmt.Println(err)
			f.Close()
			done <- false
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		done <- false
		return
	}
	done <- true
}

func main() {
	data := make(chan int)
	done := make(chan bool)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go produce(data, &wg)
	}
	go consume(data, done)
	go func() {
		wg.Wait()
		close(data)
	}()
	d := <-done
	if d {
		fmt.Println("File written successfully")
	} else {
		fmt.Println("File writing failed")
	}
}
```

The main function creates the `data` channel in line no. 41 from which random numbers are read from and written. The `done` channel in line no. 42 is used by the `consume` goroutine to notify `main` that it is done with its task. The `wg` waitgroup in line no. 43 is used to wait for all the 100 goroutines to finish generating random numbers.

The `for` loop in line no. 44 creates 100 goroutines. The goroutine call in line no. 49 calls `Wait()` on the waitgroup to wait for all 100 goroutines to finish creating random numbers. After that, it closes the channel. Once the channel is closed and the `consume` goroutine has finished writing all generated random numbers to the file, it writes `true` to the `done` channel in line no. 37 and the main goroutine is unblocked and prints `File written successfully`.
