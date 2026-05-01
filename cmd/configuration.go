/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Kory291/gherkin-formatter/internal/configuration"
)

// configurationCmd represents the configuration command
var configurationCmd = &cobra.Command{
	Use:   "configuration",
	Short: "This command is used to read the current configuration.",
	Long:  `Use this to read configuration options. This also displays the location where the file should be put.`,
	Run: func(cmd *cobra.Command, args []string) {
		testRun, err := cmd.Flags().GetBool("test")
		if err != nil {
			panic(err)
		}
		configDir := "."
		if testRun {
			fmt.Println("Running configuration in test mode")
			configDir = "test_data/"
		}
		params, err := configuration.ReadConfiguration(configDir)
		if err != nil {
			panic(err)
		}
		configuration.PrintConfiguration(params)
	},
}

func init() {
	configurationCmd.Flags().Bool("test", false, "Use this flag if you test the script")
	rootCmd.AddCommand(configurationCmd)
}
