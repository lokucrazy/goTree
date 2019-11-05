package goTree

import (
	"compress/zlib"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

type GitDir struct {
	Name string
	Dir  string
}

type Commit struct {
	Tree      string
	Parent    string
	Author    string
	Committer string
}

func (git *GitDir) ReadObject(hash string) (string, error) {
	objDir := git.Dir + "objects/" + hash[:2] + "/" + hash[2:]
	file, err := os.Open(objDir)
	if err != nil {
		return "", err
	}
	defer file.Close()

	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		return "", err
	}

	var contents strings.Builder
	_, err = io.Copy(&contents, zlibReader)
	if err != nil {
		return "", err
	}

	return contents.String(), nil
}

func (*GitDir) ParseCommit(text string) (*Commit, error) {
	if text == "" {
		return nil, errors.New("cannot parse empty text")
	}
	var tree, parent, author, committer string
	treeMatcher := regexp.MustCompile(`(?:tree) (.*)`)
	parentMatcher := regexp.MustCompile(`(?:parent) (.*)`)
	authorMatcher := regexp.MustCompile(`(?:author) (.*) <`)
	committerMatcher := regexp.MustCompile(`(?:committer) (.*) <`)

	findTree := treeMatcher.FindStringSubmatch(text)
	findParent := parentMatcher.FindStringSubmatch(text)
	findAuthor := authorMatcher.FindStringSubmatch(text)
	findCommitter := committerMatcher.FindStringSubmatch(text)
	if findTree == nil && findParent == nil && findAuthor == nil && findCommitter == nil {
		return nil, errors.New("Could not parse/find commit fields")
	}
	if findTree != nil {
		tree = findTree[1]
	}
	if findParent != nil {
		parent = findParent[1]
	}
	if findAuthor != nil {
		author = findAuthor[1]
	}
	if findCommitter != nil {
		committer = findCommitter[1]
	}

	return &Commit{tree, parent, author, committer}, nil
}

func NewGitDir(project string) (*GitDir, error) {
	var dir string
	if string(project[len(project)-1:]) == "/" {
		dir = project + ".git/"
	} else {
		dir = project + "/.git/"
	}

	var name string
	tokens := strings.Split(project, "/")
	name = tokens[len(tokens)-1]

	return &GitDir{name, dir}, nil
}
