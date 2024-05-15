/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-cmd-book/pScan/scan"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run the scan for the hosts list",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")

		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}
		portRange, err := cmd.Flags().GetString("port-range")
		if err != nil {
			return err
		}
		startPort, endPort, err := parsePortRange(portRange)
		if err != nil {
			return err
		}
		for i := startPort; i <= endPort; i++ {
			ports = append(ports, i)
		}

		return scanAction(os.Stdout, hostsFile, ports)
	},
}

func parsePortRange(portRange string) (int, int, error) {
	if portRange == "" {
		return 0, 0, nil
	}
	parts := strings.Split(portRange, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid port range: %s", portRange)
	}
	lowerBound, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid port range: %s", portRange)
	}
	upperBound, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid port range: %s", portRange)
	}
	return lowerBound, upperBound, nil
}

func scanAction(out io.Writer, hostsFile string, ports []int) error {
	hl := &scan.HostList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	return printResults(out, scan.Run(hl, ports))
}

func printResults(out io.Writer, results []scan.Results) error {
	for _, r := range results {
		message := ""
		message += fmt.Sprintf("%s: ", r.Host)
		if r.NotFound {
			message += "host not found \n\n"
			continue
		}
		message += "\n"
		for _, p := range r.PortStates {
			message += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
		}

		_, err := fmt.Fprintln(out, message)
		if err != nil {
			return nil
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntSliceP("ports", "p", []int{22, 80, 443}, "ports to scan")
	scanCmd.Flags().StringP("port-range", "r", "", "port range, ex (1-1024)")
}
