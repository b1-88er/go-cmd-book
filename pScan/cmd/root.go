/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pScan",
	Short: "Fast TCP port scanner",
	Long: `pScan - short for Port Scanner - executes TCP port scan on a list of hosts.

pScan allows you to add, list and delete hosts from the list.
It executes a port scan on a specified TCP ports.
You can customize the target ports using a command line flag.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run:     func(cmd *cobra.Command, args []string) {},
	Version: "0.0.1",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	var cfgFile string
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pScan.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}
