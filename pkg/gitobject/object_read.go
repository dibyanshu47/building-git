package gitobject

import (
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dibyanshu47/building-git/pkg/gitrepository"
)

func ObjectRead(repo *gitrepository.GitRepository, sha string) (*GitObject, error) {
	path := gitrepository.RepoPath(repo, "objects", sha[0:2], sha[2:])

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%v does not exist", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw, err := zlib.NewReader(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer raw.Close()

	rawData, err := io.ReadAll(raw)
	if err != nil {
		return nil, err
	}

	decompressedData := string(rawData)

	// Read object type
	x := strings.Index(decompressedData, " ")
	format := decompressedData[0:x]

	// Read and validate object size
	y := strings.Index(decompressedData, "\x00")
	size, err := strconv.Atoi(decompressedData[x:y])
	if err != nil {
		return nil, err
	}
	if size != len(rawData)-y-1 {
		return nil, errors.New("malformed object: bad length")
	}

	// Pick constructor
	var obj *GitObject
	switch format {
	case "commit":
		// obj = &GitCommit{content: rawData[y+1:]}
	case "tree":
		// obj = &GitTree{content: rawData[y+1:]}
	case "tag":
		// obj = &GitTag{content: rawData[y+1:]}
	case "blob":
		// obj = &GitBlob{content: rawData[y+1:]}
	default:
		return nil, fmt.Errorf("unknown type %s for object", format)
	}

	// Return object
	return obj, nil
}
