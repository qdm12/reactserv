package memfs

import (
	"io"
	"os"
)

// Implements os.File.
type InMemoryFile struct {
	at   int64
	Name string
	data []byte
	fs   memFS
}

func (f *InMemoryFile) Close() error               { return nil }
func (f *InMemoryFile) Stat() (os.FileInfo, error) { return &InMemoryFileInfo{f}, nil }

func (f *InMemoryFile) Readdir(count int) ([]os.FileInfo, error) {
	res := make([]os.FileInfo, len(f.fs))
	i := 0
	for _, file := range f.fs {
		res[i], _ = file.Stat()
		i++
	}
	return res, nil
}

func (f *InMemoryFile) Read(b []byte) (int, error) {
	i := 0
	for f.at < int64(len(f.data)) && i < len(b) {
		b[i] = f.data[f.at]
		i++
		f.at++
	}
	return i, nil
}

func (f *InMemoryFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.at = offset
	case io.SeekCurrent:
		f.at += offset
	case io.SeekEnd:
		f.at = int64(len(f.data)) + offset
	}
	return f.at, nil
}
