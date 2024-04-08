package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	proj := flag.String("project", "", "project name")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("project name is required: %w", ErrValidation)
	}

	pipeline := make([]step, 0)

	pipeline = append(pipeline, newStep(
		"go build",
		"go",
		"Go build: SUCCESS\n",
		proj,
		[]string{"build", ".", "errors"},
	))

	pipeline = append(pipeline, newStep(
		"go test",
		"go",
		"Go test: SUCCESS\n",
		proj,
		[]string{"test", "./...", "-v"},
	))

	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}
		if _, err := fmt.Fprint(out, msg); err != nil {
			return err
		}
	}
	return nil
}
