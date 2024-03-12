package main

import (
	"bufio" // read text
	"flag"
	"fmt"
	"io" // io Reader interface
	"os"
)

// CountMode is an enum for the count mode
type CountMode int

const (
	Bytes CountMode = iota
	Lines
	Words
)

func count(r io.Reader, mode CountMode) int {
	// a scanner is used to read text from a reader such a files
	scanner := bufio.NewScanner(r)
	switch mode {
	case Bytes:
		scanner.Split(bufio.ScanBytes)
	case Lines:
		scanner.Split(bufio.ScanLines)
	case Words:
		scanner.Split(bufio.ScanWords)
	default:
		panic("Unknown mode")
	}

	wc := 0
	for scanner.Scan() {
		wc++
	}
	return wc
}
func main() {
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("c", false, "Count bytes")
	flag.Parse()
	if *lines && *bytes {
		fmt.Println("Please choose either -l or -c")
		os.Exit(1)
	}
	if *lines {
		fmt.Println(count(os.Stdin, Lines))
		return
	}
	if *bytes {
		fmt.Println(count(os.Stdin, Bytes))
		return
	}
	fmt.Println(count(os.Stdin, Words))
}
