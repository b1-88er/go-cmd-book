package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-cmd-book/todo"
	"io"
	"os"
	"strings"
)

var todoFileName = ".todo.json"

func getTask(r io.Reader, args ...string) ([]string, error) {
	tasks := make([]string, 0)
	if len(args) > 0 {
		tasks = append(tasks, strings.Join(args, " "))
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		if len(s.Text()) == 0 {
			return nil, fmt.Errorf("task cannot be blank")

		}
		tasks = append(tasks, s.Text())

	}
	return tasks, nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool.\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}
	verbose := flag.Bool("verbose", false, "verbose output")
	add := flag.Bool("add", false, "Add task to the Todo List")
	task := flag.String("task", "", "Task to be included in the Todo list")
	list := flag.Bool("list", false, "List of all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	delete := flag.Int("delete", 0, "Delete item of a given index")
	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		if *verbose {
			for _, t := range *l {
				fmt.Println(t)
			}
		} else {
			fmt.Print(l)
		}
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *add:
		tasks, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, task := range tasks {
			l.Add(task)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *task != "":
		l.Add(*task)
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > 0:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "No flags provided")
		flag.Usage()
		os.Exit(1)
	}
}
