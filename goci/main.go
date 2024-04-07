package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	args := []string{"build", ".", "errors"}
	cmd := exec.Command("go", args...)
	cmd.Dir = proj

	if err := cmd.Run(); err != nil {
		return &stepErr{step: build, msg: "go build failed", cause: err}
	}

	_, err := fmt.Fprint(out, "Build successful")
	return err
}
