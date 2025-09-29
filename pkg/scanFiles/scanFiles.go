package scanFiles

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func WhereAmI() (string, error) {
	pwd, envPresent := os.LookupEnv("PWD");
	if !envPresent {
		return "", errors.New("PWD variable not found");
	}
	return pwd, nil
}

func filterFeatureFiles(fileNames []string) []string {
	featureFiles := []string{}
	for _, fileName := range fileNames {
		if len(fileName) > 8 && fileName[len(fileName)-8:] == ".feature" {
			featureFiles = append(featureFiles, fileName)
		}
	}
	return featureFiles
}

func FindFeatureFiles(path string) ([]string, error) {
	searchPath := path
	fileNames := []string{}
	files, err := os.ReadDir(searchPath)
	if err != nil {
		return make([]string, 0), errors.New("feature directory not found");
	}

	for _, file := range files {
		if file.IsDir() {
			fileNamesInDir, err := FindFeatureFiles(path + "/" + file.Name())
			if err != nil {
				fmt.Printf("Something went wrong while looking up %s", path)
				return make([]string, 0), err
			}
			featureFiles := filterFeatureFiles(fileNamesInDir)
			fileNames = append(fileNames, featureFiles...)
		}
		if len(file.Name()) > 8 && file.Name()[len(file.Name())-8:] == ".feature" {
			fileNames = append(fileNames, path + "/" + file.Name())
		}
	}
	return fileNames, nil
}

func ReadFiles(paths []string) (map[string][]string, error) {
	resultFiles := make(map[string][]string)
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("error when opening file %s", path)
			return make(map[string][]string, 0), errors.New("file could not be opened")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		textContent := make([]string, 0)
		for scanner.Scan() {
			textContent = append(textContent, scanner.Text())
		}
		resultFiles[path] = textContent
	}
	return resultFiles, nil
}

func WriteFiles(fileContents map[string][]string) error {
	for path, fileContent := range fileContents {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		for _, line := range fileContent {
			_, err := writer.WriteString(line)
			if err != nil {
				writer.Flush()
				return err
			}
		}
		writer.Flush()
	}
	return nil
}