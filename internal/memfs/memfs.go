package memfs

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/qdm12/golibs/logging"
)

type MemFS interface {
	Open(name string) (file http.File, err error)
	Watch(ctx context.Context, wg *sync.WaitGroup)
}

func New(rootPath string, oldToNew map[*regexp.Regexp]string, logger logging.Logger) (fs MemFS, err error) {
	memFS := &memFS{
		mapping:  make(map[string]memFSElement),
		mu:       &sync.RWMutex{},
		logger:   logger,
		rootPath: filepath.Clean(rootPath),
		oldToNew: oldToNew,
	}
	err = memFS.loadAll()
	return memFS, err
}

type memFS struct {
	mapping     map[string]memFSElement // key is the relative path
	mu          *sync.RWMutex           // pointer to respect value receiver for Open method
	logger      logging.Logger
	rootPath    string
	oldToNew    map[*regexp.Regexp]string
	directories map[string]struct{} // key is the absolute path
}

func (fs memFS) Open(name string) (file http.File, err error) {
	fs.mu.RLock()
	element, ok := fs.mapping[name]
	fs.mu.RUnlock()
	if !ok {
		return nil, os.ErrNotExist
	}
	if element.isDir {
		file = newInMemoryDirectory(element.name, element.modTime)
		return file, nil
	}
	file = newInMemoryFile(element.data, element.name, element.modTime)
	return file, nil
}
