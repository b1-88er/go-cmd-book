/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const timeFormat = "02/01 @15:04"

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "Show a single item in detail",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}
		return viewAction(os.Stdout, apiRoot, id)
	},
}

func viewAction(out io.Writer, apiRoot string, id int) error {
	item, err := getOne(apiRoot, id)
	if err != nil {
		return err
	}
	return printOne(out, item)
}

func printOne(out io.Writer, i item) error {
	w := tabwriter.NewWriter(out, 14, 2, 0, ' ', 0)
	fmt.Fprintf(w, "Task:\t%s\n", i.Task)
	fmt.Fprintf(w, "Created:\t%s\n", i.CreatedAt.Format(timeFormat))
	if i.Done {
		fmt.Fprintf(w, "Completed:\t%s\n", "Yes")
		fmt.Fprintf(w, "Completed At:\t%s\n", i.CompletedAt.Format(timeFormat))
		return w.Flush()
	}
	fmt.Fprintf(w, "Completed:\t%s\n", "No")
	return w.Flush()

}

func init() {
	rootCmd.AddCommand(viewCmd)
	viewCmd.Flags().IntP("id", "i", 0, "Item ID")
}
