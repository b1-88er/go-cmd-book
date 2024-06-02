package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-cmd-book/todo"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupApi(t *testing.T) (string, func()) {
	t.Helper()
	tempTodoFile, err := os.CreateTemp("", "todotest")
	assert.NoError(t, err)

	ts := httptest.NewServer(newMux(tempTodoFile.Name()))

	for i := 1; i < 3; i++ {
		var body bytes.Buffer
		taskName := fmt.Sprintf("Task number %d", i)
		item := struct {
			Task string `json:"task"`
		}{
			Task: taskName,
		}
		err := json.NewEncoder(&body).Encode(item)
		assert.NoError(t, err)
		r, err := http.Post(ts.URL+"/todo", "application/json", &body)
		assert.NoError(t, err)
		assert.Equal(t, r.StatusCode, http.StatusCreated)

	}
	return ts.URL, func() {
		ts.Close()
	}
}

var (
	resp struct {
		Results      todo.List `json:"results"`
		Date         int64     `json:"date"`
		TotalResults int       `json:"total_results"`
	}
	body []byte
	err  error
)

func TestGet(t *testing.T) {
	url, cleanUp := setupApi(t)
	defer cleanUp()

	t.Run("GetRoot", func(t *testing.T) {
		r, err := http.Get(url + "/")
		assert.NoError(t, err)
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, "Hello, from the the api", string(body))

		defer r.Body.Close()
		assert.Equal(t, 200, r.StatusCode)
	})
	t.Run("GetAll", func(t *testing.T) {
		r, err := http.Get(url + "/todo")
		assert.NoError(t, err)
		defer r.Body.Close()

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, r.StatusCode)

		err = json.NewDecoder(r.Body).Decode(&resp)
		assert.NoError(t, err)

		assert.Equal(t, 2, resp.TotalResults)
		assert.Equal(t, "Task number 1", resp.Results[0].Task)
		assert.Equal(t, "Task number 2", resp.Results[1].Task)
	})

	t.Run("GetOne", func(t *testing.T) {
		r, err := http.Get(url + "/todo/1")
		assert.NoError(t, err)
		defer r.Body.Close()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, r.StatusCode)

		err = json.NewDecoder(r.Body).Decode(&resp)
		assert.NoError(t, err)

		assert.Equal(t, 1, resp.TotalResults)
		assert.Equal(t, "Task number 1", resp.Results[0].Task)
	})

	t.Run("NotFound", func(t *testing.T) {
		r, err := http.Get(url + "/todo/404")
		assert.NoError(t, err)
		defer r.Body.Close()
		assert.Equal(t, 404, r.StatusCode)
	})

}
