package cmd

import (
	"bytes"
	"go-cmd-book/pScan/scan"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostActions(t *testing.T) {
	hostFile, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	t.Run("add hosts", func(t *testing.T) {
		out := bytes.Buffer{}
		err := addAction(&out, hostFile.Name(), []string{"host1", "host2"})
		assert.NoError(t, err)
		assert.Equal(t, "Added host: host1\nAdded host: host2\n", out.String())
	})

	t.Run("add empty hosts", func(t *testing.T) {
		out := &bytes.Buffer{}
		addAction(out, hostFile.Name(), []string{})
		assert.NoError(t, err)
		assert.Equal(t, "", out.String())
	})

	t.Run("list hosts", func(t *testing.T) {
		err := addAction(&bytes.Buffer{}, hostFile.Name(), []string{"host3"})
		assert.NoError(t, err)

		out := &bytes.Buffer{}
		err = listAction(out, hostFile.Name())
		assert.NoError(t, err)
		assert.Equal(t, "host1\nhost2\nhost3\n", out.String())
	})

	t.Run("delete host", func(t *testing.T) {
		out := &bytes.Buffer{}
		err := deleteAction(out, hostFile.Name(), []string{"host1"})
		assert.NoError(t, err)
		assert.Equal(t, "Host deleted: host1\n", out.String())

		err = deleteAction(out, hostFile.Name(), []string{"host1"})
		assert.ErrorIs(t, err, scan.ErrNotExists)
	})

	db, err := os.ReadFile(hostFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "host2\nhost3\n", string(db))
}
