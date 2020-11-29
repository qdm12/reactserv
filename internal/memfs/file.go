package memfs

import (
	"io"
	"os"
)

// Implements os.File.
type inMemoryFile struct {
	Name string
	data []byte
	at   int64
}

func (f *inMemoryFile) Close() error {
	f.at = 0
	return nil
}

func (f *inMemoryFile) Stat() (os.FileInfo, error) {
	return &inMemoryFileInfo{
		name: f.Name,
		size: int64(len(f.data)),
	}, nil
}

func (f *inMemoryFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *inMemoryFile) Read(b []byte) (int, error) {
	i := 0
	for f.at < int64(len(f.data)) && i < len(b) {
		b[i] = f.data[f.at]
		i++
		f.at++
	}
	return i, nil
}

func (f *inMemoryFile) Seek(offset int64, whence int) (int64, error) {
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
