package memfs

import (
	"fmt"
	"os"
)

// Implements os.File.
type inMemoryDir struct {
	Name string
}

func (d *inMemoryDir) Close() error { return nil }
func (d *inMemoryDir) Stat() (os.FileInfo, error) {
	return &inMemoryDirInfo{
		name: d.Name,
	}, nil
}

func (d *inMemoryDir) Readdir(count int) ([]os.FileInfo, error) {
	fmt.Println("dir readdir")
	return nil, nil
}

func (d *inMemoryDir) Read(b []byte) (int, error) {
	fmt.Println("dir read")
	return 0, nil
}

func (d *inMemoryDir) Seek(offset int64, whence int) (int64, error) {
	fmt.Println("dir seek")
	return 0, nil
}
