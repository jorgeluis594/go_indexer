package indexer

import (
	"os"
)

type FileReader struct {
	Path      string
	filePaths []string
}

func (r *FileReader) getFile() (*os.File, error) {
	emailFile, err := os.Open(r.getPath())
	if err != nil {
		return emailFile, err
	}
	defer emailFile.Close()
	return emailFile, nil
}

func (r *FileReader) getPath() string {
	lastIndex := len(r.filePaths) - 1
	filePath := r.filePaths[lastIndex]
	r.filePaths = r.filePaths[:lastIndex]
	return filePath
}
