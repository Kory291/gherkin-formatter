package scanFiles

import (
	"os"
	"fmt"
	"errors"
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
		return []string{}, errors.New("feature directory not found");
	}

	for _, file := range files {
		if file.IsDir() {
			fileNamesInDir, err := FindFeatureFiles(path + "/" + file.Name())
			if err != nil {
				fmt.Printf("Something went wrong while looking up %s", path)
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