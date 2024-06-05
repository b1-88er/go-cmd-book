package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
