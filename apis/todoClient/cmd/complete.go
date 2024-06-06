/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete <id>",
	Short: "Complete a selected task",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		return completeAction(os.Stdout, apiRoot, args[0])
	},
}

func completeAction(out io.Writer, apiRoot string, arg string) error {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("%w: item id must be a number", err)
	}
	if err := completeItem(apiRoot, id); err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "Action %d completed\n", id)
	return err
}

func init() {
	rootCmd.AddCommand(completeCmd)

}
