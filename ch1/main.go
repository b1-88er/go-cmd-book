package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func count(r io.Reader, countBytes bool, countLines bool) int {
	scanner := bufio.NewScanner(r)
	if !countLines {
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
