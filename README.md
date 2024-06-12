# Code snippets from [Powerful Command-Line Applications in Go](https://pragprog.com/titles/rggo/powerful-command-line-applications-in-go/)

This repository contains the "mini projects" introduced in the book. This file contains notable findings and thoughts about the code/projects/approaches etc. It is mostly for personal use, but if you have any comments or thoughts leave an issue or a PR.

I followed the instructions mostly as presented in the book with some minor exceptions. Most notably I ditched the "table testing" approach. The tests in the book are sometimes more convoluted than the main logic. I didn't mind some repetition for the sake of the code readability. I also used testify framework to make the tests more concise.

## wc

A clone of the `wc` unix tool.

### flag returns a pointer

`flag.Bool` retuns a `*bool`. It makes sense, because there should be one instance of the argument and not bunch of copies.

### io.Reader

The `count` function doesn't take the bytes passed into the tool but `io.Reader` interface. `io.Reader` and `io.Writer` are very prominent in golang, so no wonder they appeared in the first tool.

`io.Reader` is an abstraction for the "thing I can read from". It can be either `*os.File` or `*strings.Reader` or even `net.Conn` - the most important thing is that it implements the `Read(p []byte) (n int, err error)`. It is worth noting that Read doesn't pull all the data at once. It writes number of bytes to the passed slice. So it fits both internet sockets, files and strings.
In this case `os.Stdin` is passed to the count. It is a `os.File`, specifically `NewFile(uintptr(syscall.Stdin), "/dev/stdin")`

Since count takes `io.Reader` I could use many options to "mock" the os.Stdin. Like:

```go
// b := bytes.NewBufferString("12345")
// b := strings.NewReader("12345")
b := bytes.NewReader([]byte("12345"))
assert.Equal(t, count(b, bufio.ScanBytes), 5)
```

This is another argument for using io.Reader instead of passing the data around as `[]byte` or `string`.

### bufio.Scanner

`bufio` and concept of a `Scanner` are closely realted to the `io.Reader`. By building a scanner out of the io.Reader user can process the data in chunks defined by the `Split` scanFunc.

`Scan()` advances the Scanner to the next token, which will then be available through the Scanner.Bytes or Scanner.Text method. When the Scan is done, it will return false. so

```go
wc := 0
for scanner.Scan() {
    wc++
}
```

is a nice way of counting the splits.

## walk

Walks the root directory a runs actions on files.

### io.Writer and gzip.NewWriter

gzip greatly demonstrates benefits of using io.Reader and io.Writer.
`gzip.NewWriter` accepts io.Writer, that might be `*os.File`, but might be a socket or a buffer or anything that accepts `Write(p []byte) (n int, err error)`. Accept slice of bytes and return number of bytes written from p and return an error.

To run gzip simply read bytes from `io.Reader` and write them to `io.Writer` as: `io.Copy(zw, in)`, where in is the io.Reader and ze is the io.Writer.

```golang
zw := gzip.NewWriter(out)
zw.Name = filepath.Base(path)
if _, err := io.Copy(zw, in); err != nil {
    return err
}

if err := zw.Close(); err != nil {
    return err
}
```

### use `fmt.Fprintln` to make the code testable.

Instead of using `fmt.Println` use `fmt.Fprintln` that accepts io.Writer. At the top level os.Stdout should be used as in `main` function: `run(*root, os.Stdout, c)`.
Of the fmt.Fprintln is used, it is then easily testable:

```golang
buffer := bytes.Buffer{}
if err := run(tempDir, &buffer, testCase.cfg); err != nil {
    t.Fatal(err)
}
res := buffer.String()
assert.Equal(t, testCase.expected, res)
```

This is much easier than trying to capture stdout from the testrun and somehow filter only that the code under question generated.

### filepath.Walk

Interesting way of iterating over directory tree. Nice demo of how functions are first class citizen in golang.

```golang
func filepath.Walk(root string, fn filepath.WalkFunc) error
Walks the file tree rooted at root, calling fn for each file or directory in the tree, including root.
```

## todo list

## mbp

## pScan

## colStats

## goci

## apis

## pomodoro
