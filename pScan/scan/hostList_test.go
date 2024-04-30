package scan_test

import (
	"os"
	"testing"

	"go-cmd-book/pScan/scan"

	"github.com/stretchr/testify/assert"
)

func TestHostList(t *testing.T) {
	t.Run("addNew", func(t *testing.T) {
		hl := &scan.HostList{}
		err := hl.Add("host2")
		assert.NoError(t, err)
		assert.Len(t, hl.Hosts, 1)
	})
	t.Run("addExisting", func(t *testing.T) {
		hl := &scan.HostList{}
		hl.Add("host2")
		err := hl.Add("host2")
		assert.ErrorIs(t, err, scan.ErrExists)
		assert.Len(t, hl.Hosts, 1)
	})

	t.Run("removeExisting", func(t *testing.T) {
		hl := &scan.HostList{}
		hl.Add("host2")
		err := hl.Remove("host2")
		assert.NoError(t, err)
		assert.Len(t, hl.Hosts, 0)
	})

	t.Run("silly slice", func(t *testing.T) {
		// I am doing this for 20 years and still confuse indexes with len
		s := []string{"a", "b"}
		assert.Equal(t, []string{"a"}, s[0:1])
		assert.Equal(t, []string{s[0]}, s[0:1])
		assert.Equal(t, []string{"b"}, s[1:2])
		assert.Equal(t, []string{s[1]}, s[1:2])
		assert.Equal(t, len(s), 2)
	})

	t.Run("removeNotExisting", func(t *testing.T) {
		hl := &scan.HostList{}
		hl.Add("host1")
		err := hl.Remove("host2")
		assert.ErrorIs(t, err, scan.ErrNotExists)
		assert.Len(t, hl.Hosts, 1)
	})

	t.Run("save", func(t *testing.T) {
		f, err := os.CreateTemp("", "")
		assert.NoError(t, err)
		defer os.Remove(f.Name())

		hl1 := &scan.HostList{}
		hl2 := &scan.HostList{}

		hostName := "host1"
		hl1.Add(hostName)

		assert.NoError(t, hl1.Save(f.Name()))
		assert.NoError(t, hl2.Load(f.Name()))
		assert.Equal(t, hl1.Hosts, hl2.Hosts)
	})

	t.Run("loadNotExists", func(t *testing.T) {
		f, err := os.CreateTemp("", "")
		assert.NoError(t, err)
		assert.NoError(t, os.Remove(f.Name()))

		hl := &scan.HostList{}
		// odd, but ok
		assert.ErrorIs(t, hl.Load(f.Name()), nil)

	})
}
