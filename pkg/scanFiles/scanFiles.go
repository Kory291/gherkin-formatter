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

func FindFeatureFiles(path string) ([]string, error) {
	searchPath := path;
	fileNames := []string{}
	if path == "" {
		searchPath = "features";
	}
	files, err := os.ReadDir(searchPath);
	if err != nil {
		return []string{}, errors.New("feature directory not found");
	}

	for _, file := range files {
		if file.IsDir() {
			fileNamesInDir, err := FindFeatureFiles(path + "/" + file.Name())
			if err != nil {
				fmt.Println("Something went wrong while looking up %s", path)
			}
			fileNames = append(fileNames, fileNamesInDir...)
		}
		fileNames = append(fileNames, path + "/" + file.Name())
	}
	return fileNames, nil;
}