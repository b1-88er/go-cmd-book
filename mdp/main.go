package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	header = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title>Markdown preview tool</title>
	</head>
	<body>
`

	footer = `
	</body>
</html>
`
)

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdownfile to preview")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string) error {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)
	outName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outName)
	return saveHTML(outName, htmlData)
}
