package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)
	return ts.URL, ts.Close
}

func TestComplete(t *testing.T) {
	url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todo/1", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)
		_, ok := r.URL.Query()["complete"]
		assert.Equal(t, true, ok)
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "")
	})
	defer cleanup()
	var out bytes.Buffer
	err := completeAction(&out, url, "1")
	assert.NoError(t, err)
	assert.Equal(t, "Action 1 completed\n", out.String())
}

func TestDelete(t *testing.T) {
	url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todo/1", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "")
	})
	defer cleanup()
	var out bytes.Buffer
	err := deleteAction(&out, url, "1")
	assert.NoError(t, err)
	assert.Equal(t, "Task id: 1 has been deleted\n", out.String())
}

func TestView(t *testing.T) {
	t.Run("ResultOne", func(t *testing.T) {
		url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `
{
	"date": 1717424841,
	"total_results": 1,
	"results": [
		{
			"task": "Task 1",
			"done": false,
			"created_at": "2024-06-03T16:24:49.319593+02:00",
			"completed_at": "0001-01-01T00:00:00Z"
		}
	]
}`)
		})
		defer cleanup()
		out := bytes.Buffer{}
		err := viewAction(&out, url, 1)
		assert.NoError(t, err)
		assert.Equal(t, "Task:         Task 1\nCreated:      03/06 @16:24\nCompleted:    No\n", out.String())
	})

}
func TestAdd(t *testing.T) {
	url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, "/todo")
		assert.Equal(t, http.MethodPost, r.Method)
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		r.Body.Close()
		assert.Equal(t, "{\"task\":\"Task 1\"}\n", string(body))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "created body")
	})
	defer cleanup()
	var out bytes.Buffer
	err := addAction(&out, url, []string{"Task 1"})
	assert.NoError(t, err)
	assert.Equal(t, "Added task \"Task 1\" to the list \n", out.String())
}

func TestListAction(t *testing.T) {
	t.Run("AllResults", func(t *testing.T) {
		url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `
{
	"date": 1717424841,
	"total_results": 2,
	"results": [
		{
			"task": "Task 1",
			"done": false,
			"created_at": "2024-06-03T16:24:49.319593+02:00",
			"completed_at": "0001-01-01T00:00:00Z"
		},
		{
			"task": "Task 2",
			"done": false,
			"created_at": "2024-06-03T16:24:53.015757+02:00",
			"completed_at": "0001-01-01T00:00:00Z"
		}
	]
}
`)
		})
		defer cleanup()

		out := bytes.Buffer{}
		err := listAction(&out, url)
		assert.NoError(t, err)
		assert.Equal(t, "-  1  Task 1\n-  2  Task 2\n", out.String())
	})

	t.Run("NoResults", func(t *testing.T) {
		url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `
{
	"date": 1717424841,
	"total_results": 0,
	"results": [
	]
}
`)
		})
		defer cleanup()

		out := bytes.Buffer{}
		err := listAction(&out, url)
		assert.True(t, errors.Is(err, ErrNotFound))
		assert.Equal(t, "", out.String())
	})

	t.Run("InvalidUrl", func(t *testing.T) {
		url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `
{
	"date": 1717424841,
	"total_results": 0,
	"results": [
	]
}
`)
		})

		out := &bytes.Buffer{}
		cleanup()
		err := listAction(out, url)
		assert.ErrorIs(t, err, ErrConnection)
	})
}
