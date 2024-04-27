package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type executer interface {
	execute() (string, error)
}

func main() {
	proj := flag.String("project", "", "project name")
	flag.Parse()

	if err := run(*proj, os.Stdout, exec.CommandContext); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj string, out io.Writer, command Command) error {
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
		"Go fmt: SUCCESS\n",
		proj,
		[]string{"-l", "."},
	))

	pipeline = append(pipeline, newTimeoutStep(
		"git push",
		"git",
		"Git push: SUCCESS\n",
		proj,
		[]string{"push", "origin", "master"},
		5*time.Second,
		command,
	))

	sig := make(chan os.Signal, 1)

	errCh := make(chan error)
	// struct is a 0 byte type?
	done := make(chan struct{})

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for _, s := range pipeline {
			msg, err := s.execute()
			if err != nil {
				errCh <- err
				return
			}
			if _, err := fmt.Fprint(out, msg); err != nil {
				errCh <- err
				return
			}
		}
		close(done)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case rec := <-sig:
			signal.Stop(sig)
			return fmt.Errorf("%s Exiting: %w", rec, ErrSignal)
		case <-done:
			return nil
		}
	}
}
