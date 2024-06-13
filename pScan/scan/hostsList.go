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

type HostList []string

func (hl *HostList) Add(host string) error {
	if idx := slices.Index(*hl, host); idx > -1 {
		return fmt.Errorf("%s: %w", host, ErrExists)
	}
	*hl = append(*hl, host)
	return nil
}

func (hl *HostList) Remove(host string) error {
	if idx := slices.Index(*hl, host); idx > -1 {
		*hl = slices.Delete(*hl, idx, idx+1)
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
		*hl = append(*hl, scanner.Text())
	}
	return nil
}

func (hl *HostList) Save(hostFile string) error {
	output := ""
	for _, host := range *hl {
		output += fmt.Sprintln(host)
	}
	return os.WriteFile(hostFile, []byte(output), 0644)
}
