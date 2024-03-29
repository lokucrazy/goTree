package goTree

import (
	"errors"
	"os"
	"reflect"
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
	expectedContents := `commit 203` + "\x00" + `tree 380fcea7b4540e995f05504e9d2bad1eb87282bb
author Ryan Lokugamage <Ryan.Lokugamage@cerner.com> 1572387473 -0500
committer Ryan Lokugamage <Ryan.Lokugamage@cerner.com> 1572387473 -0500

initial commit`
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

func TestReadObject_FailsWhenBadGitDir(t *testing.T) {
	gitDir, _ := NewGitDir("fakeProject")
	_, err := gitDir.ReadObject("1038930c0737f81c3713fc74f3523f29614b4fdb")
	if err == nil {
		t.Error("expected err to not be nil but got nil")
	}
	var e *os.PathError
	if !errors.As(err, &e) {
		t.Errorf("expected %v but got %v", reflect.TypeOf(e), reflect.TypeOf(err))
	}
}

func TestParseCommit(t *testing.T) {
	expectedCommit := Commit{
		"380fcea7b4540e995f05504e9d2bad1eb87282bb",
		"",
		"Ryan Lokugamage",
		"Ryan Lokugamage",
	}
	text := `commit 203` + "\x00" + `tree 380fcea7b4540e995f05504e9d2bad1eb87282bb
author Ryan Lokugamage <Ryan.Lokugamage@cerner.com> 1572387473 -0500
committer Ryan Lokugamage <Ryan.Lokugamage@cerner.com> 1572387473 -0500

initial commit`
	gitDir, err := NewGitDir("./testProj")
	if err != nil {
		t.Error(err.Error())
	}
	actual, err := gitDir.ParseCommit(text)
	if err != nil {
		t.Error(err.Error())
	}
	if expectedCommit.Tree != actual.Tree {
		t.Errorf("expected %s but got %s", expectedCommit.Tree, actual.Tree)
	}
	if expectedCommit.Parent != actual.Parent {
		t.Errorf("expected %s but got %s", expectedCommit.Parent, actual.Parent)
	}
	if expectedCommit.Author != actual.Author {
		t.Errorf("expected %s but got %s", expectedCommit.Author, actual.Author)
	}
	if expectedCommit.Committer != actual.Committer {
		t.Errorf("expected %s but got %s", expectedCommit.Committer, actual.Committer)
	}
}

func TestParseCommit_FailsWhenTextIsEmpty(t *testing.T) {
	text := ""
	gitDir, err := NewGitDir("./testProj")
	if err != nil {
		t.Error(err.Error())
	}
	_, err = gitDir.ParseCommit(text)
	if err == nil {
		t.Error("expected err to not be nil but got nil")
	}
	if !errors.Is(err, parseError) {
		t.Errorf("expected %v but got %v", parseError, err)
	}
}

func TestParseCommit_FailsWhenMatchersFail(t *testing.T) {
	text := "bad commit"
	gitDir, err := NewGitDir("./testProj")
	if err != nil {
		t.Error(err.Error())
	}
	_, err = gitDir.ParseCommit(text)
	if err == nil {
		t.Error("expected err to not be nil but got nil")
	}
	if !errors.Is(err, matchError) {
		t.Errorf("expected %v but got %v", matchError, err)
	}
}

func TestReadRef(t *testing.T) {
	expectedHash := `25b7ae3966479cd84b74721fc270bc90777f6de4`
	gitDir, err := NewGitDir("./testProj")
	if err != nil {
		t.Error(err.Error())
	}
	actual, err := gitDir.ReadRef("master")
	if err != nil {
		t.Error(err.Error())
	}
	if expectedHash != actual {
		t.Errorf("expected %s but got %s", expectedHash, actual)
	}
}
