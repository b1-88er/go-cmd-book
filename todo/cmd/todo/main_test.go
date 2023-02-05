package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const binName = "todo"

// executes once per test suite
func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	// executes in the path of this file
	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Connot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests ...")
	result := m.Run()
	fmt.Println("Cleaning up ...")
	os.Remove(binName)
	os.Remove(todoFileName)
	os.Exit(result) // have to exit on my own according to docs
}

func TestTodoCLI(t *testing.T) {
	task := "test task number 1"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-task", task)
		assert.Nil(t, cmd.Run())
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		expected := fmt.Sprintf(" 1: %s\n", task)

		assert.Nil(t, err)
		assert.Equal(t, expected, string(out))
	})

	t.Run("DeleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-delete", "1")
		assert.Nil(t, cmd.Run())

		out, err := exec.Command(cmdPath, "-list").CombinedOutput()

		assert.Nil(t, err)
		assert.Equal(t, "", string(out))

	})

	t.Run("Add task from the STDIN", func(t *testing.T) {
		const task = "task from stdin"
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		cmdStdIn.Close()
		assert.Nil(t, err)
		io.WriteString(cmdStdIn, task)
		assert.NotNil(t, cmd.Run())
	})
}
