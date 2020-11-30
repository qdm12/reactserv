package memfs

import (
	"regexp"
)

func processDataSpecificPath(path string, data []byte) (processed []byte) {
	return data
}

func processData(data []byte, oldToNew map[*regexp.Regexp]string) (processed []byte) {
	s := string(data)
	for old, new := range oldToNew {
		old.ReplaceAllString(s, new)
	}
	return []byte(s)
}
