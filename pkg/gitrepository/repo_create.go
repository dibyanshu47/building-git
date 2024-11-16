package gitrepository

import (
	"fmt"
	"os"
	"path/filepath"
)

func RepoCreate(path string) (*GitRepository, error) {
	// Create a new repository at path.
	repo, err := NewGitRepository(path, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create new repository: %v", err)
	}

	// First, we make sure the path either doesn't exist or is an empty dir.
	if _, err := os.Stat(repo.Worktree); err == nil {
		if !isDir(repo.Worktree) {
			return nil, fmt.Errorf("%s is not a directory", path)
		}
		if _, err := os.Stat(repo.Gitdir); err == nil {
			if entries, _ := os.ReadDir(repo.Gitdir); len(entries) > 0 {
				return nil, fmt.Errorf("%s is not empty", path)
			}
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(repo.Worktree, 0755); err != nil {
			return nil, fmt.Errorf("failed to create worktree directory: %v", err)
		}
	} else {
		return nil, fmt.Errorf("failed to check worktree: %v", err)
	}

	// Create necessary directories
	dirs := [][]string{
		{"branches"},
		{"objects"},
		{"refs", "tags"},
		{"refs", "heads"},
	}

	for _, dir := range dirs {
		if _, err := RepoDir(repo, true, dir...); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %v", filepath.Join(dir...), err)
		}
	}

	// Create .git/description file
	descPath, err := RepoFile(repo, true, "description")
	if err != nil {
		return nil, fmt.Errorf("failed to create description file: %v", err)
	}
	if err := os.WriteFile(descPath, []byte("Unnamed repository; edit this file 'description' to name the repository.\n"), 0644); err != nil {
		return nil, fmt.Errorf("failed to write to description file: %v", err)
	}

	// Create .git/HEAD file
	headPath, err := RepoFile(repo, true, "HEAD")
	if err != nil {
		return nil, fmt.Errorf("failed to create HEAD file: %v", err)
	}
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/master\n"), 0644); err != nil {
		return nil, fmt.Errorf("failed to write to HEAD file: %v", err)
	}

	// Create .git/config file
	configPath, err := RepoFile(repo, true, "config")
	if err != nil {
		return nil, fmt.Errorf("failed to create config file: %v", err)
	}
	config := RepoDefaultConfig()
	if err := config.SaveTo(configPath); err != nil {
		return nil, fmt.Errorf("failed to write to config file: %v", err)
	}

	return repo, nil
}
