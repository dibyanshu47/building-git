package gitrepository

import (
	"fmt"
	"os"
)

func RepoDir(repo *GitRepository, mkdir bool, paths ...string) (string, error) {
	// Same as repoPath, but mkdir *path if absent if mkdir is true.
	path := RepoPath(repo, paths...)

	fileInfo, err := os.Stat(path)
	if err == nil {
		if fileInfo.IsDir() {
			return path, nil
		}
		return "", fmt.Errorf("not a directory: %s", path)
	}

	if os.IsNotExist(err) && mkdir {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %s", err)
		}
		return path, nil
	}

	return "", nil
}
