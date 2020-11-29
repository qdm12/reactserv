package memfs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/qdm12/golibs/logging"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

func New(rootPath string, oldToNew map[string]string, logger logging.Logger) (fs http.FileSystem, err error) {
	memFS := memFS{
		m:      make(map[string]memFSElement),
		logger: logger,
	}
	err = filepath.Walk(rootPath, makeWalkFn(memFS, rootPath, oldToNew, logger))
	return memFS, err
}

func makeWalkFn(memFS memFS, rootPath string,
	oldToNew map[string]string, logger logging.Logger) filepath.WalkFunc {
	rootPath = filepath.Clean(rootPath)
	return func(path string, info os.FileInfo, err error) (newErr error) {
		if err != nil {
			return err // fails if we encounter any error previously
		}

		path = filepath.Clean(path)
		relativePath := strings.TrimPrefix(path, rootPath)
		if len(relativePath) == 0 {
			relativePath = "/"
		}

		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			memDir := memFSElement{
				name:    filepath.Base(relativePath),
				modTime: stat.ModTime(),
				isDir:   true,
			}
			memFS.m[relativePath] = memDir
			logger.Info("loading directory %s", relativePath)
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		data = processDataSpecificPath(relativePath, data)
		data = processData(data, oldToNew)

		memFile := memFSElement{
			name:    filepath.Base(relativePath),
			modTime: stat.ModTime(),
			data:    data,
		}
		memFS.m[relativePath] = memFile
		logger.Info("loaded file %s", relativePath)

		return nil
	}
}

type memFS struct {
	m      map[string]memFSElement // key is the relative path
	logger logging.Logger
}

func (fs memFS) Open(name string) (file http.File, err error) {
	element, ok := fs.m[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrFileNotFound, name)
	}
	if element.isDir {
		file = newInMemoryDirectory(element.name, element.modTime)
		return file, nil
	}
	file = newInMemoryFile(element.data, element.name, element.modTime)
	return file, nil
}
