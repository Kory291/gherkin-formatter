package scanFiles

import (
	"os"
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
		fileNames = append(fileNames, path + file.Name())
	}
	return fileNames, nil;
}