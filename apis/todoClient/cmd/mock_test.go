package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)
	return ts.URL, ts.Close
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
