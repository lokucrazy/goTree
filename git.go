package goTree

import (
	"compress/zlib"
	"io"
	"os"
	"strings"
)

type GitDir struct {
	Name string
	Dir  string
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
