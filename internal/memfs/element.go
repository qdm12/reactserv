package memfs

import "time"

type memFSElement struct {
	name    string
	modTime time.Time
	isDir   bool
	data    []byte
}
