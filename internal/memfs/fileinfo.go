package memfs

import (
	"os"
	"time"
)

// Implements os.FileInfo.
type inMemoryFileInfo struct {
	name string
	size int64
}

func (i *inMemoryFileInfo) Name() string       { return i.name }
func (i *inMemoryFileInfo) Size() int64        { return i.size }
func (i *inMemoryFileInfo) Mode() os.FileMode  { return os.ModeTemporary }
func (i *inMemoryFileInfo) ModTime() time.Time { return time.Time{} }
func (i *inMemoryFileInfo) IsDir() bool        { return false }
func (i *inMemoryFileInfo) Sys() interface{}   { return nil }
