package gitrepository

import "path/filepath"

func RepoPath(repo *GitRepository, paths ...string) string {
	return filepath.Join(append([]string{repo.Gitdir}, paths...)...)
}
