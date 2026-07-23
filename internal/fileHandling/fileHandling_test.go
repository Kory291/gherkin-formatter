package fileHandling

import (
	"testing"
)

func TestWhereAmI(t *testing.T) {
	want := "/some/path"
	t.Setenv("PWD", want)
	path, err := WhereAmI()
	if path != want || err != nil {
		t.Errorf("WhereAmI returned %s, wanted %s, error %v", path, want, err)
	}
}
