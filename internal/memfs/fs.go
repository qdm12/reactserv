package memfs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

func New(rootPath string, oldToNew map[string]string) (fs http.FileSystem, err error) {
	memFS := memFS{}
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // fails if we encounter any error previously
		}
		stat, err := os.Stat(path)
		if err != nil {
			return err
		} else if stat.IsDir() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		data = processDataSpecificPath(path, data)
		data = processData(data, oldToNew)
		memFile := &InMemoryFile{
			Name: path,
			data: data,
			fs:   memFS,
		}
		memFS[path] = memFile
		return nil
	})
	return memFS, err
}

type memFS map[string]http.File

func (fs memFS) Open(name string) (file http.File, err error) {
	file, ok := fs[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrFileNotFound, name)
	}
	return file, nil
}
