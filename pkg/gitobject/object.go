package gitobject

import (
	"errors"

	"github.com/dibyanshu47/building-git/pkg/gitrepository"
)

type GitObject struct {
	data []byte
}

func NewGitObject(data []byte) *GitObject {
	return &GitObject{data: data}
}

func (g *GitObject) Serialize(repo *gitrepository.GitRepository) error {
	// This function MUST be implemented by subclasses.
	// It must read the object's contents from g.data, a byte string, and do
	// whatever it takes to convert it into a meaningful representation.
	// What exactly that means depend on each subclass.
	return errors.New("unimplemented")
}

func (g *GitObject) Deserialize(data []byte) error {
	// This function MUST be implemented by subclasses.
	// It must read the object's contents from the given data byte slice and
	// convert it into a meaningful representation.
	return errors.New("unimplemented")
}
