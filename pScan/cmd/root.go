/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pScan.yaml)")
	rootCmd.PersistentFlags().StringP("hosts-file", "f", "pScan.hosts", "pScan hosts file")

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("PSCAN")

	viper.BindPFlag("hosts-file", rootCmd.PersistentFlags().Lookup("hosts-file"))

	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		homedir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("homedir missing: %w", err))
			os.Exit(1)
		}

		viper.AddConfigPath(homedir)
		viper.SetConfigName(".pScan")
	}

	viper.AutomaticEnv()

	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Fprintln(os.Stderr, "using config file: ", viper.ConfigFileUsed())
	// }
	fmt.Println(viper.ReadInConfig())
}
