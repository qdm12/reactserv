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
		m:      make(map[string]http.File),
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
			memDir := &inMemoryDir{
				Name: filepath.Base(relativePath),
			}
			memFS.m[relativePath] = memDir
			logger.Info("loading directory %s", relativePath)
			return nil
		}

		memFile := &inMemoryFile{
			Name: filepath.Base(relativePath),
		}
		memFile.data, err = ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		memFile.data = processDataSpecificPath(relativePath, memFile.data)
		memFile.data = processData(memFile.data, oldToNew)

		memFS.m[relativePath] = memFile
		logger.Info("loaded file %s", relativePath)

		return nil
	}
}

type memFS struct {
	m      map[string]http.File
	logger logging.Logger
}

func (fs memFS) Open(name string) (file http.File, err error) {
	file, ok := fs.m[name]
	if !ok {
		fs.logger.Warn("%s: %q", ErrFileNotFound, name)
		return nil, fmt.Errorf("%w: %q", ErrFileNotFound, name)
	}
	return file, nil
}
