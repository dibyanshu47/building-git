package gitobject

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/dibyanshu47/building-git/pkg/gitrepository"
)

func ObjectRead(repo *gitrepository.GitRepository, sha string) (GitObject, error) {
	path := gitrepository.RepoPath(repo, "objects", sha[0:2], sha[2:])

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, err
	}
	defer file.Close()

	// Read the compressed data
	rawData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Decompress the data
	reader, err := zlib.NewReader(bytes.NewReader(rawData))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Read object type (first part before space)
	spaceIndex := bytes.IndexByte(decompressedData, ' ')
	if spaceIndex == -1 {
		return nil, errors.New("malformed object: missing space delimiter")
	}
	fmtType := decompressedData[:spaceIndex]

	// Read object size (between space and null byte)
	nullIndex := bytes.IndexByte(decompressedData[spaceIndex:], '\x00')
	if nullIndex == -1 {
		return nil, errors.New("malformed object: missing null byte")
	}
	objectSize, err := strconv.Atoi(string(decompressedData[spaceIndex+1 : spaceIndex+nullIndex]))
	if err != nil {
		return nil, fmt.Errorf("invalid object size: %v", err)
	}

	// Validate size
	if objectSize != len(decompressedData)-spaceIndex-nullIndex-1 {
		return nil, fmt.Errorf("malformed object %s: bad length", sha)
	}

	// Determine the object type and create the appropriate object
	var gitObject GitObject
	switch string(fmtType) {
	// case "commit":
	// 	gitObject = &GitCommit{}
	// case "tree":
	// 	gitObject = &GitTree{}
	// case "tag":
	// 	gitObject = &GitTag{}
	case "blob":
		gitObject = NewGitBlob(decompressedData)
	default:
		return nil, fmt.Errorf("unknown type %s for object %s", fmtType, sha)
	}

	// Deserialize the object data
	gitObject.Deserialize(decompressedData[spaceIndex+nullIndex+1:])

	// Return the constructed Git object
	return gitObject, nil
}
