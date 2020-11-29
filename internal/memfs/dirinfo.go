package memfs

import (
	"os"
	"time"
)

// Implements os.FileInfo.
type inMemoryDirInfo struct {
	name string
}

func (i *inMemoryDirInfo) Name() string       { return i.name }
func (i *inMemoryDirInfo) Size() int64        { return 0 }
func (i *inMemoryDirInfo) Mode() os.FileMode  { return os.ModeTemporary }
func (i *inMemoryDirInfo) ModTime() time.Time { return time.Time{} }
func (i *inMemoryDirInfo) IsDir() bool        { return true }
func (i *inMemoryDirInfo) Sys() interface{}   { return nil }
