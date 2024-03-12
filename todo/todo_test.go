package todo_test

import (
	"go-cmd-book/todo"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	l := todo.List{}
	taskName := "New Task"
	l.Add(taskName)
	assert.Equal(t, l[0].Task, taskName)
}

func TestComplete(t *testing.T) {
	l := todo.List{}
	taskName := "New Task"
	l.Add(taskName)
	assert.Equal(t, l[0].Task, taskName)
	assert.Equal(t, l[0].Done, false)
	l.Complete(1)
	assert.Equal(t, l[0].Done, true)
}

func TestDelete(t *testing.T) {
	l := todo.List{}
	tasks := []string{
		"New Task 1",
		"New Task 2",
		"New Task 3",
	}
	for _, v := range tasks {
		l.Add(v)
	}

	assert.Equal(t, l[0].Task, tasks[0])
	l.Delete(2)
	assert.Equal(t, len(l), 2)
}

func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	taskName := "New Task"
	l1.Add(taskName)

	assert.Equal(t, l1[0].Task, taskName)

	tf, err := os.CreateTemp("", "")

	if err != nil {
		t.Fatalf("Error creating tmp file %s", err)
	}

	defer os.Remove(tf.Name())
	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}

	l2 := todo.List{}
	if err := l2.Get(tf.Name()); err != nil {
		t.Fatalf("Error getting list from a file :%s", err)
	}

	assert.Equal(t, l2[0].Task, l1[0].Task)
}
