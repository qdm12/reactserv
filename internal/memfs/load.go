package memfs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/qdm12/golibs/logging"
)

var (
	ErrLoading = errors.New("cannot load files in memory")
)

func (fs *memFS) load(rootPath string, oldToNew map[string]string) (err error) {
	newMapping := make(map[string]memFSElement)
	walkFn := makeWalkFn(newMapping, rootPath, oldToNew, fs.logger)
	if err := filepath.Walk(rootPath, walkFn); err != nil {
		return fmt.Errorf("%w: %s", ErrLoading, err)
	}
	fs.mu.Lock()
	fs.m = newMapping
	fs.mu.Unlock()
	return nil
}

func makeWalkFn(mapping map[string]memFSElement, rootPath string,
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
			mapping[relativePath] = memDir
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
		mapping[relativePath] = memFile
		logger.Info("loaded file %s", relativePath)

		return nil
	}
}
