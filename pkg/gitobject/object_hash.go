package gitobject

import (
	"fmt"
	"io"

	"github.com/dibyanshu47/building-git/pkg/gitrepository"
)

func ObjectHash(fd io.Reader, format string, repo *gitrepository.GitRepository) (string, error) {
	data, err := io.ReadAll(fd)
	if err != nil {
		return "", err
	}

	var obj GitObject
	switch format {
	// case "commit":
	// 	obj = &GitCommit{}
	// case "tree":
	// 	obj = &GitTree{}
	// case "tag":
	// 	obj = &GitTag{}
	case "blob":
		obj = NewGitBlob(data)
	default:
		return "", fmt.Errorf("unknown type: %s", format)
	}

	obj.Deserialize(data)

	sha, err := ObjectWrite(obj, repo)
	if err != nil {
		return "", err
	}

	return sha, nil
}
