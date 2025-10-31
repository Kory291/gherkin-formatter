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
	Long: `Use this to read configuration options. This also displays the location where the file should be put.`,
	Run: func(cmd *cobra.Command, args []string) {
		testRun, err := cmd.Flags().GetBool("test")
		if err != nil {
			panic(err)
		}
		if testRun {
			fmt.Println("Running configuration in test mode")
			params, err := configuration.ReadConfiguration("test_data/")
			if err != nil {
				panic(err)
			}
			fmt.Println(*params)
			return
		}
		params, err := configuration.ReadConfiguration(".")
		if err != nil {
			panic(err)
		}
		fmt.Println("Configuration read:")
		fmt.Printf("intend-and:\t%t\n", params.IntendAnd)
		fmt.Printf("intendation:\t%d\n", params.Intendation)

	},
}

func init() {
	configurationCmd.Flags().Bool("test", false, "Use this flag if you test the script")
	rootCmd.AddCommand(configurationCmd)
}
