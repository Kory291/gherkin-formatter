package fileHandling

import (
	"testing"
	"slices"
)

func TestWhereAmI(t *testing.T) {
	want := "/some/path"
	t.Setenv("PWD", want)
	path, err := WhereAmI()
	if path != want || err != nil {
		t.Errorf("WhereAmI returned %s, wanted %s, error %v", path, want, err)
	}
}

func TestFilterFeatureFiles(t *testing.T) {
	inputs := []string{"foo.feature", "a.feature" ,".feature", "foo"}
	want := []string{"foo.feature", "a.feature"}
	returned := filterFeatureFiles(inputs)
	if !slices.Equal(returned, want) {
		t.Errorf("filterFeatureFiles returned %v, wanted %v", returned, want)
	}
}

func TestIsFeatureFile(t *testing.T) {
	input := "test.feature"
	isFeature := isFeatureFile(input)
	if !isFeature {
		t.Errorf("isFeatureFile returned %t for %s", isFeature, input)
	} 
}

