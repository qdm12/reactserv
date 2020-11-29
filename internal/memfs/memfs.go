package memfs

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/qdm12/golibs/logging"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

func New(rootPath string, oldToNew map[string]string, logger logging.Logger) (fs http.FileSystem, err error) {
	memFS := memFS{
		m:      make(map[string]memFSElement),
		mu:     &sync.RWMutex{},
		logger: logger,
	}
	err = memFS.load(rootPath, oldToNew)
	return memFS, err
}

type memFS struct {
	m      map[string]memFSElement // key is the relative path
	mu     *sync.RWMutex           // pointer to respect value receiver for Open method
	logger logging.Logger
}

func (fs memFS) Open(name string) (file http.File, err error) {
	fs.mu.RLock()
	element, ok := fs.m[name]
	fs.mu.RUnlock()
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
