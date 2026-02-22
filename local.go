package storage

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/kohirens/stdlib/fsio"
)

var ps = string(os.PathSeparator)

// LocalStorage Save data in local files.
type LocalStorage struct {
	WorkDir string
}

func NewLocalStorage(wd string) (*LocalStorage, error) {
	if !fsio.DirExist(wd) {
		return nil, &ErrDirNoExist{wd}
	}

	return &LocalStorage{
		WorkDir: wd,
	}, nil
}

// Exist Retrieve file from storage.
func (s *LocalStorage) Exist(filename string) bool {
	filePath := s.Location(filename)

	Log.Dbugf(stdout.Load, filePath)

	return fsio.Exist(filePath)
}

// List Files in a location in storage. It is not recursive, it only lists
// files in the specified directory.
func (s *LocalStorage) List(location string) ([]string, error) {
	filePath := s.Location(location)
	Log.Dbugf(stdout.Load, filePath)
	files := make([]string, 0)
	prefix := filePath + ps
	e1 := filepath.WalkDir(filePath, func(path string, d fs.DirEntry, err error) error {
		filename := strings.Replace(path, prefix, "", 1)
		if filename == "" {
			return nil
		}

		files = append(files, filename)

		return nil
	})
	if e1 != nil {
		return nil, fmt.Errorf(stderr.ListFiles, e1.Error())
	}

	if files[0] == filePath {
		return files[1:], nil
	}

	return files, nil
}

// Load Retrieve file from storage.
func (s *LocalStorage) Load(filename string) ([]byte, error) {
	filePath := s.Location(filename)

	Log.Dbugf(stdout.Load, filePath)

	if !fsio.Exist(filePath) {
		return nil, fmt.Errorf("%v %v", filePath, fs.ErrNotExist)
	}

	content, e1 := os.ReadFile(filePath)
	if e1 != nil {
		return nil, &ErrReadFile{filePath + " " + e1.Error()}
	}

	return content, nil
}

// Save Write session data to the storage medium.
func (s *LocalStorage) Save(filename string, data []byte) error {
	filePath := s.Location(filename)

	if e := os.WriteFile(filePath, data, 0774); e != nil {
		return &ErrWriteFile{e.Error()}
	}

	return nil
}

func (s *LocalStorage) Location(filename string) string {
	return s.WorkDir + ps + filename
}

// Remove Delete a file from storage.
func (s *LocalStorage) Remove(filename string) error {
	fullFilename := s.Location(filename)

	if e := os.Remove(fullFilename); e != nil {
		return fmt.Errorf(stderr.RemoveFile, e.Error())
	}

	return nil
}
