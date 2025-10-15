/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/Kory291/gherkin-formatter/pkg/fileHandling"
	"github.com/Kory291/gherkin-formatter/pkg/format"
	"github.com/Kory291/gherkin-formatter/pkg/configuration"
	"github.com/spf13/cobra"

	"fmt"
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Formats all found feature files",
	Long: `Formats all found feature files in a directory.
	
	Default behaviour is printing the formatted feature files to stdout.
	Examples:
	
	gherkin-formatter format  # this will display formatted feature files to stdout
	gherkin-formatter format --write   this will write the formatted feature files directly
	gherkin-formatter format --test  # this will run the command in test mode which uses feature files in test_data/ 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := fileHandling.WhereAmI()
		if err != nil {
			panic(err)
		}
		testFlag, err := cmd.Flags().GetBool("test")
		if err != nil {
			panic(err)
		}
		path := pwd
		if testFlag {
			path = pwd + "/test_data"
		}
		writeFlag, err := cmd.Flags().GetBool("write")
		if err != nil {
			panic(err)
		}
		featureFiles, err := fileHandling.FindFeatureFiles(path + "/features")
		if err != nil {
			panic(err)
		}
		fileContents, err := fileHandling.ReadFiles(featureFiles)
		if err != nil {
			panic(err)
		}
		configuration, err := configuration.ReadConfiguration(path)
		if err != nil {
			panic(err)
		}
		formattedFiles := make(map[string][]string, 0)
		for filePath, fileContent := range fileContents {
			formattedFile, err := format.FormatFile(fileContent, *configuration)
			if err != nil {
				fmt.Printf("Error while processing file %s - will skip this", filePath)
				continue
			}
			formattedFiles[filePath] = formattedFile
		}
		if writeFlag {
			return
		}
		for filePath, formattedFile := range formattedFiles {
			fmt.Printf("------------------------------\n%s\n", filePath)
			for _, line := range formattedFile {
				fmt.Println(line)
			} 
		}
	},
}

func init() {
	formatCmd.Flags().Bool("write", false, "Add this flag if you want to apply the changes directly")
	formatCmd.Flags().Bool("test", false, "Add this flag if you want to run the script in test mode")
	rootCmd.AddCommand(formatCmd)
}
