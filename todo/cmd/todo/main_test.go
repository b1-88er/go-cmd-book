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
const testingFileName = ".testing.json"

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
	os.Remove(testingFileName)
	os.Exit(result) // have to exit on my own according to docs
}

func cmd(cmdPath string, args ...string) *exec.Cmd {
	cmd := exec.Command(cmdPath, args...)
	cmd.Env = append(os.Environ(), "TODO_FILENAME="+testingFileName)
	return cmd
}

func TestTodoCLI(t *testing.T) {
	task := "test task number 1"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		cmd := cmd(cmdPath, "-add", task)
		assert.Nil(t, cmd.Run())
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := cmd(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		expected := fmt.Sprintf(" 1: %s\n", task)

		assert.Nil(t, err)
		assert.Equal(t, expected, string(out))
	})

	t.Run("DeleteTask", func(t *testing.T) {
		cmd := cmd(cmdPath, "-delete", "1")
		assert.Nil(t, cmd.Run())

		out, err := exec.Command(cmdPath, "-list").CombinedOutput()

		assert.Nil(t, err)
		assert.Equal(t, "", string(out))

	})

	t.Run("Add task from the STDIN", func(t *testing.T) {
		const task = "task from stdin\n"
		add := cmd(cmdPath, "-add")
		cmdStdIn, err := add.StdinPipe()
		assert.Nil(t, err)
		io.WriteString(cmdStdIn, task)
		cmdStdIn.Close()
		assert.Nil(t, add.Run())

		list, err := cmd(cmdPath, "-list").Output()
		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf(" 1: %s", task), string(list))
	})
}
