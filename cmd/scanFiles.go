/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Kory291/gherkin-formatter/pkg/scanFiles"
	"github.com/spf13/cobra"
)

// scanFilesCmd represents the scanFiles command
var scanFilesCmd = &cobra.Command{
	Use:   "scanFiles",
	Short: "Returns all found Feature files",
	Long: `This command can be used to retrieve a list of all feature files that are discovered by gherkin-formatter`,
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := scanFiles.WhereAmI()
		if err != nil {
			os.Exit(1);
		}
		testRun, err := cmd.Flags().GetBool("test");
		if err != nil {
			os.Exit(1);
		}
		path := pwd
		if testRun {
			path = pwd + "/tests" 
		}
		fileNames, err := scanFiles.FindFeatureFiles(path + "/features");
		if err != nil {
			os.Exit(1);
		}
		fmt.Println(fileNames);
	},
}

func init() {
	scanFilesCmd.Flags().Bool("test", false, "Enable this when testing the application")
	rootCmd.AddCommand(scanFilesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanFilesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanFilesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
