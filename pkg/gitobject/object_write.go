package gitobject

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dibyanshu47/building-git/pkg/gitrepository"
)

func ObjectWrite(obj GitObject, repo *gitrepository.GitRepository) (string, error) {
	data := obj.Serialize()

	lengthStr := strconv.Itoa(len(data))
	lengthBytes := []byte(lengthStr)

	// Concatenate the parts into a single byte slice
	fmt.Println(obj.GetFormat())
	header := append(obj.GetFormat(), []byte(" ")...) // Append the format and a space
	header = append(header, lengthBytes...)           // Append the length in bytes
	header = append(header, []byte("\x00")...)        // Append the null byte
	header = append(header, data...)                  // Append the actual data

	hash := sha1.New()
	hash.Write(header)
	sha := fmt.Sprintf("%x", hash.Sum(nil))

	// If repo is provided, write the object to the file system
	if repo != nil {
		path := gitrepository.RepoPath(repo, "objects", sha[:2], sha[2:])

		// Create the directory if it doesn't exist
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", err
		}

		// Check if the object already exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// If not, write the compressed object to the file
			file, err := os.Create(path)
			if err != nil {
				return "", err
			}
			defer file.Close()

			// Compress the header data and write it to the file
			compressed := new(bytes.Buffer)
			writer := zlib.NewWriter(compressed)
			_, err = writer.Write(header)
			if err != nil {
				return "", err
			}
			writer.Close()

			// Write the compressed data to the file
			_, err = file.Write(compressed.Bytes())
			if err != nil {
				return "", err
			}
		}
	}

	return sha, nil
}
