package gitrepository

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type GitRepository struct {
	Worktree string
	Gitdir   string
	Conf     *ini.File
}

func NewGitRepository(path string, force bool) (*GitRepository, error) {
	gitdir := filepath.Join(path, ".git")
	if !(force || isDir(gitdir)) {
		return nil, fmt.Errorf("not a Git repository: %s", path)
	}

	// Read configuration file in .git/config
	conf, err := readConfigFile(path, force)
	if err != nil {
		return nil, err
	}

	repo := &GitRepository{
		Worktree: path,
		Gitdir:   gitdir,
		Conf:     conf,
	}

	if !force {
		// Check repository format version
		vers, err := repo.getRepositoryFormatVersion()
		if err != nil {
			return nil, err
		}
		if vers != 0 {
			return nil, fmt.Errorf("unsupported repositoryformatversion %d", vers)
		}
	}

	return repo, nil
}

// Check if the directory exists
func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Read configuration file
func readConfigFile(path string, force bool) (*ini.File, error) {
	confPath := filepath.Join(path, ".git", "config")
	if !isDir(confPath) {
		if !force {
			return nil, fmt.Errorf("configuration file missing: %s", confPath)
		}
		return nil, nil
	}
	conf, err := ini.Load(confPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}
	return conf, nil
}

// Get repository format version from the config
func (repo *GitRepository) getRepositoryFormatVersion() (int, error) {
	section := repo.Conf.Section("core")
	versStr := section.Key("repositoryformatversion").String()
	if versStr == "" {
		return 0, fmt.Errorf("repositoryformatversion not found in config")
	}
	var version int
	_, err := fmt.Sscanf(versStr, "%d", &version)
	if err != nil {
		return 0, fmt.Errorf("invalid repositoryformatversion: %s", versStr)
	}
	return version, nil
}
