package main

import (
	"fmt"
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
}
