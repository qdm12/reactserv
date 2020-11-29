package memfs

import (
	"os"
	"time"
)

// Implements os.FileInfo.
type inMemoryFileInfo struct {
	name    string
	size    int64
	modTime time.Time
	isDir   bool
}

func (i *inMemoryFileInfo) Name() string       { return i.name }
func (i *inMemoryFileInfo) Size() int64        { return i.size }
func (i *inMemoryFileInfo) Mode() os.FileMode  { return os.ModeTemporary }
func (i *inMemoryFileInfo) ModTime() time.Time { return i.modTime }
func (i *inMemoryFileInfo) IsDir() bool        { return i.isDir }
func (i *inMemoryFileInfo) Sys() interface{}   { return nil }
