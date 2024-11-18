package gitrepository

import (
	"fmt"
	"path/filepath"
)

func RepoFind(path string, required bool) (*GitRepository, error) {
	path, _ = filepath.Abs(path)

	if isDir(filepath.Join(path, ".gitt")) {
		return NewGitRepository(path, true)
	}

	parent := filepath.Dir(path)

	if parent == path {
		if required {
			return nil, fmt.Errorf("no git directory found")
		} else {
			return nil, nil
		}
	}

	return RepoFind(parent, required)
}
