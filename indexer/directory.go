package indexer

import (
	"fmt"
	"os"
	"path/filepath"
)

type Directory struct {
	rootPath       string
	filePaths      []string
	subDirectories []Directory
}

func (d *Directory) GetPaths() []string {
	var filePaths []string
	filePaths = append(filePaths, d.filePaths...)
	dirsToProcess := d.subDirectories
	for len(dirsToProcess) > 0 {
		dir := dirsToProcess[0]
		dirsToProcess = dirsToProcess[1:]
		filePaths = append(filePaths, dir.filePaths...)
		dirsToProcess = append(dirsToProcess, dir.subDirectories...)
	}

	return filePaths
}

func InitDirectory(path string) (*Directory, error) {
	directory := Directory{rootPath: path}
	if !fileExists(path) {
		return &directory, fmt.Errorf("error file %s doesn't exists", path)
	}
	directory.loadPaths()
	return &directory, nil
}

func (d *Directory) hasSubDirectories() bool {
	return len(d.subDirectories) > 0
}

func (d *Directory) loadPaths() {
	files, _ := os.ReadDir(d.rootPath)

	for _, fileInfo := range files {
		childPath := filepath.Join(d.rootPath, fileInfo.Name())
		if fileInfo.IsDir() {
			subDirectory, _ := InitDirectory(childPath)
			d.subDirectories = append(d.subDirectories, *subDirectory)
		} else {
			d.filePaths = append(d.filePaths, childPath)
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
