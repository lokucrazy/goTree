package goTree

import (
	"fmt"
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
	gitDir, err := NewGitDir("./testProj")
	if err != nil {
		t.Error(err.Error())
	}

	actual, err := gitDir.ReadObject("1038930c0737f81c3713fc74f3523f29614b4fdb")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(actual)
}
