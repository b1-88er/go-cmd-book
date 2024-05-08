package cmd

import (
	"bytes"
	"fmt"
	"go-cmd-book/pScan/scan"
	"net"
	"os"
	"strconv"
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

func TestScanAction(t *testing.T) {
	hostFile, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	hl := scan.HostList{}
	assert.NoError(t, hl.Add("localhost"))
	assert.NoError(t, hl.Save(hostFile.Name()))
	out := &bytes.Buffer{}

	listener, err := net.Listen("tcp", net.JoinHostPort("localhost", "0"))
	assert.NoError(t, err)
	defer listener.Close()

	_, portStr, err := net.SplitHostPort(listener.Addr().String())
	assert.NoError(t, err)
	port, err := strconv.ParseInt(portStr, 10, 0)
	assert.NoError(t, err)

	err = scanAction(out, hostFile.Name(), []int{int(port)})
	assert.NoError(t, err)

	expected := fmt.Sprintf("localhost: \n\t%d: open\n\n", port)
	assert.Equal(t, expected, out.String())

}
