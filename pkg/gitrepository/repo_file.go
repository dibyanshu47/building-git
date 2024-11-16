package gitrepository

import "fmt"

func RepoFile(repo *GitRepository, mkdir bool, paths ...string) (string, error) {
	// Same as repoPath, but create dirname(*path) if absent.
	// For example, repoFile(r, true, "refs", "remotes", "origin", "HEAD") will create
	// .git/refs/remotes/origin.
	if len(paths) == 0 {
		return "", fmt.Errorf("no paths provided")
	}

	dirPath, err := RepoDir(repo, mkdir, paths[:len(paths)-1]...)
	if err != nil {
		return "", err
	}

	if dirPath != "" {
		return RepoPath(repo, paths...), nil
	}

	return "", nil
}
