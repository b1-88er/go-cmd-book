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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an item from the list",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		return deleteAction(os.Stdout, apiRoot, args[0])
	},
}

func deleteAction(out io.Writer, apiRoot string, arg string) error {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("%w: id must be a number", err)
	}

	if err := deleteItem(apiRoot, id); err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "Task id: %d has been deleted\n", id)
	return err
}
func init() {
	rootCmd.AddCommand(deleteCmd)
}
