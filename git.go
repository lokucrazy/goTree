package goTree

import (
	"compress/zlib"
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

type Tree struct {
	Commit *Commit
	Prev   *Commit
	Next   *Commit
}

func (git *GitDir) GenerateTree() (*Tree, error) {
	file, err := os.Open(git.Dir + "refs/heads/")
	if err != nil {
		return nil, err
	}
	heads, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	refs := make([]string, len(heads))
	for i, head := range heads {
		refs[i], err = git.ReadRef(head)
		if err != nil {
			refs[i] = ""
		}
	}

	return nil, nil
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
	return strings.TrimRight(contents.String(), "\n"), nil
}

func (git *GitDir) ReadRef(name string) (string, error) {
	refDir := git.Dir + "refs/heads/" + name
	file, err := os.Open(refDir)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var hash strings.Builder
	_, err = io.Copy(&hash, file)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(hash.String(), "\n"), nil
}

func (*GitDir) ParseCommit(text string) (*Commit, error) {
	if text == "" {
		return nil, parseError
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
		return nil, matchError
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
