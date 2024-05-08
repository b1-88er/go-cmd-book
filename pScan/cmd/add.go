/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-cmd-book/pScan/scan"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <host1> ... <hostn>",
	Args:         cobra.MinimumNArgs(1),
	Short:        "Add host(s) to the host file",
	Aliases:      []string{"a"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}
		return addAction(os.Stdout, hostsFile, args)
	},
}

func addAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}
	for _, host := range args {
		if err := hl.Add(host); err != nil {
			return err
		}
		fmt.Fprintln(out, "Added host:", host)
	}
	return hl.Save(hostsFile)
}

func init() {
	hostsCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
