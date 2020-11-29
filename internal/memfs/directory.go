package memfs

import (
	"os"
)

// Implements os.File.
type inMemoryDir struct {
	Name string
}

func (d *inMemoryDir) Close() error { return nil }
func (d *inMemoryDir) Stat() (os.FileInfo, error) {
	return &inMemoryFileInfo{
		name:  d.Name,
		isDir: true,
	}, nil
}

func (d *inMemoryDir) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (d *inMemoryDir) Read(b []byte) (int, error) {
	return 0, nil
}

func (d *inMemoryDir) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
