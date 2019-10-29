package goTree

import (
	"testing"
)

func TestNewGitDir(t *testing.T) {
	path := "/fakeProject"
	expectedDir := "/fakeProject/.git/"
	expectedName := "fakeProject"
	actual, err := NewGitDir(path)
	if err != nil {
		t.Error("NewGitDir returned with err ", err.Error())
	}
	if actual.Dir != expectedDir {
		t.Errorf("expected %s but got %s", expectedDir, actual.Dir)
	}
	if actual.Name != expectedName {
		t.Errorf("expected %s but got %s", expectedName, actual.Name)
	}
}

func TestReadObject(t *testing.T) {
	expectedContents := `commit 203tree 380fcea7b4540e995f05504e9d2bad1eb87282bb
	author Ryan Lokugamage <Ryan.Lokugamage@cerner.com> 1572387473 -0500
	committer Ryan Lokugamage <Ryan.Lokugamage@cerner.com> 1572387473 -0500
	
	initial commit
	`
	gitDir, err := NewGitDir("./testProj")
	if err != nil {
		t.Error(err.Error())
	}

	actual, err := gitDir.ReadObject("1038930c0737f81c3713fc74f3523f29614b4fdb")
	if err != nil {
		t.Error(err.Error())
	}
	if expectedContents != actual {
		t.Errorf("expected %s but got %s", expectedContents, actual)
	}
}
