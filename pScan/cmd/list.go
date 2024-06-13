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
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List hosts in the host list",
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")
		return listAction(os.Stdout, hostsFile)
	},
}

func listAction(out io.Writer, hostsFile string) error {
	hl := &scan.HostList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}
	for _, h := range *hl {
		if _, err := fmt.Fprintln(out, h); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	hostsCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
