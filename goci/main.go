package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type executer interface {
	execute() (string, error)
}

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

	pipeline := make([]executer, 0)

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

	pipeline = append(pipeline, newExecutionStep(
		"go fmt",
		"gofmt",
		"Go fmt: SUCCESS",
		proj,
		[]string{"-l", "."},
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
