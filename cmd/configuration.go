/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Kory291/gherkin-formatter/pkg/configuration"
)

// configurationCmd represents the configuration command
var configurationCmd = &cobra.Command{
	Use:   "configuration",
	Short: "This command is used to configure the script and read configuration options",
	Long: `Use this to read configuration options or to set configuration options`,
	Run: func(cmd *cobra.Command, args []string) {
		testRun, err := cmd.Flags().GetBool("test")
		if err != nil {
			panic(err)
		}
		if testRun {
			fmt.Println("Running configuration in test mode")
			params, err := configuration.ReadConfiguration("test_data/pyproject.toml")
			if err != nil {
				panic(err)
			}
			fmt.Println(params)
		}
	},
}

func init() {
	configurationCmd.Flags().Bool("test", false, "Use this flag if you test the script")
	rootCmd.AddCommand(configurationCmd)
}
