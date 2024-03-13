package main

import (
	"bufio" // read text
	"flag"
	"fmt"
	"io" // io Reader interface
	"os"
)

func count(r io.Reader, scanFunc bufio.SplitFunc) int {
	// a scanner is used to read text from a reader such a files
	scanner := bufio.NewScanner(r)
	scanner.Split(scanFunc)

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
		fmt.Println(count(os.Stdin, bufio.ScanLines))
		return
	}
	if *bytes {
		fmt.Println(count(os.Stdin, bufio.ScanBytes))
		return
	}
	fmt.Println(count(os.Stdin, bufio.ScanWords))
}
