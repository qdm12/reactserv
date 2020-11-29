package memfs

import (
	"io"
	"os"
	"time"
)

func newInMemoryFile(data []byte, name string, modTime time.Time) *inMemoryFile {
	return &inMemoryFile{
		fileInfo: inMemoryFileInfo{
			name:    name,
			size:    int64(len(data)),
			modTime: modTime,
		},
		data: data,
	}
}

func newInMemoryDirectory(name string, modTime time.Time) *inMemoryFile {
	return &inMemoryFile{
		fileInfo: inMemoryFileInfo{
			name:    name,
			modTime: modTime,
			isDir:   true,
		},
	}
}

// Implements os.File.
type inMemoryFile struct {
	fileInfo inMemoryFileInfo
	data     []byte
	at       int64
}

func (f *inMemoryFile) Close() error {
	f.at = 0
	return nil
}

func (f *inMemoryFile) Stat() (os.FileInfo, error) {
	fileInfo := f.fileInfo
	return &fileInfo, nil
}

// Readdir is never run.
func (f *inMemoryFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

// Read only applies to files and not directories.
func (f *inMemoryFile) Read(b []byte) (int, error) {
	i := 0
	for f.at < int64(len(f.data)) && i < len(b) {
		b[i] = f.data[f.at]
		i++
		f.at++
	}
	return i, nil
}

// Seek only applies to files and not directories.
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
