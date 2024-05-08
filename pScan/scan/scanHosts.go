package scan

import (
	"fmt"
	"net"
	"time"
)

type state bool

func (s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

type PortState struct {
	Port int
	Open state
}

func scanPort(host string, port int) PortState {
	p := PortState{
		Port: port,
	}
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", p.Port))
	scanConn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		p.Open = false
		return p
	}
	scanConn.Close()
	p.Open = true
	return p
}

type Results struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

func Run(hl *HostList, ports []int) []Results {
	results := make([]Results, 0, len(hl.Hosts))
	// results := []Results{}
	for _, host := range hl.Hosts {
		r := Results{
			Host: host,
		}
		if _, err := net.LookupHost(host); err != nil {
			r.NotFound = true
			results = append(results, r)
			continue
		}

		for _, port := range ports {
			r.PortStates = append(r.PortStates, scanPort(host, port))
		}
		results = append(results, r)

	}
	return results
}
