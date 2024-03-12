package main

import (
	"bufio" // read text
	"flag"
	"fmt"
	"io" // io Reader interface
	"os"
)

func count(r io.Reader, countBytes bool, countLines bool) int {
	// a scanner is used to read text from a reader such a files
	scanner := bufio.NewScanner(r)
	if !countLines {
		// split by words, default is by lines
		scanner.Split(bufio.ScanWords)
	}
	if countBytes {
		scanner.Split(bufio.ScanBytes)
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
	fmt.Println(count(os.Stdin, *bytes, *lines))
}
