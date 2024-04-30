package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
)

var (
	ErrExists    = errors.New("host already in the list")
	ErrNotExists = errors.New("host not in the list")
)

type HostList struct {
	Hosts []string
}

func (hl *HostList) search(host string) (bool, int) {
	idx := slices.Index(hl.Hosts, host)
	if idx >= 0 {
		return true, idx
	}
	return false, -1
}

func (hl *HostList) Add(host string) error {
	if present, _ := hl.search(host); present {
		return fmt.Errorf("%s: %w", host, ErrExists)
	}
	hl.Hosts = append(hl.Hosts, host)
	return nil
}

func (hl *HostList) Remove(host string) error {
	if present, i := hl.search(host); present {
		hl.Hosts = slices.Delete(hl.Hosts, i, i+1)
		return nil
	}
	return fmt.Errorf("%s: %w", host, ErrNotExists)
}

func (hl *HostList) Load(hostFile string) error {
	f, err := os.Open(hostFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}
	return nil
}

func (hl *HostList) Save(hostFile string) error {
	output := ""
	for _, host := range hl.Hosts {
		output += fmt.Sprintln(host)
	}
	return os.WriteFile(hostFile, []byte(output), 0644)
}
