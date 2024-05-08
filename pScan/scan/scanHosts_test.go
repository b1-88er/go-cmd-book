package scan_test

import (
	"go-cmd-book/pScan/scan"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}
	assert.Equal(t, "closed", ps.Open.String())
	ps.Open = true
	assert.Equal(t, "open", ps.Open.String())
}

func TestRun(t *testing.T) {
	host := "localhost"
	hl := &scan.HostList{}
	hl.Add(host)

	// setup listening on a random port ("0")
	ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
	assert.NoError(t, err)
	defer ln.Close()

	// get the selected port
	_, portStr, err := net.SplitHostPort(ln.Addr().String())
	assert.NoError(t, err)
	port, err := strconv.Atoi(portStr)
	assert.NoError(t, err)

	t.Run("port open", func(t *testing.T) {
		results := scan.Run(hl, []int{port})
		assert.Equal(t, []scan.Results{
			{
				Host:       host,
				NotFound:   false,
				PortStates: []scan.PortState{{Open: true, Port: port}},
			},
		}, results)
	})
	t.Run("port closed", func(t *testing.T) {
		ln.Close()
		results := scan.Run(hl, []int{port})
		assert.Equal(t, []scan.Results{
			{
				Host:       host,
				NotFound:   false,
				PortStates: []scan.PortState{{Open: false, Port: port}},
			},
		}, results)

	})
	t.Run("host not found", func(t *testing.T) {
		// invalid ip
		notFound := "257.257.257.257"
		assert.NoError(t, hl.Add(notFound))
		results := scan.Run(hl, []int{port})
		assert.Equal(t, []scan.Results{
			{
				Host:       host,
				NotFound:   false,
				PortStates: []scan.PortState{{Open: false, Port: port}},
			},
			{
				Host:     notFound,
				NotFound: true,
			},
		}, results)

	})
}
