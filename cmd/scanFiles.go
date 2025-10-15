/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Kory291/gherkin-formatter/pkg/fileHandling"
	"github.com/spf13/cobra"
)

// scanFilesCmd represents the scanFiles command
var scanFilesCmd = &cobra.Command{
	Use:   "scanFiles",
	Short: "Returns all found Feature files",
	Long: `This command can be used to retrieve a list of all feature files that are discovered by gherkin-formatter`,
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := fileHandling.WhereAmI()
		if err != nil {
			panic("Could not determine current directory")
		}
		testRun, err := cmd.Flags().GetBool("test");
		if err != nil {
			panic("Flag content for 'test' could not be collected")
		}
		path := pwd
		if testRun {
			fmt.Println("Running scanFiles in test mode")
			path = pwd + "/test_data" 
		}
		fileNames, err := fileHandling.FindFeatureFiles(path + "/features")
		if err != nil {
			panic("Feature files could not be found")
		}
		fileContents, err := fileHandling.ReadFiles(fileNames)
		if err != nil {
			panic("Could not read files")
		}
		for filePath, fileContent := range fileContents {
			fmt.Println("\n" + filePath)
			for _, line := range fileContent {
				fmt.Println(line)
			}
		}
	},
}

func init() {
	scanFilesCmd.Flags().Bool("test", false, "Enable this when testing the application")
	rootCmd.AddCommand(scanFilesCmd)
}
