package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
