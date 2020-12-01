package memfs

import "strings"

func processDataSpecificPath(path string, data []byte) (processed []byte) {
	return data
}

func processData(data []byte, oldToNew map[string]string) (processed []byte) {
	s := string(data)
	for old, new := range oldToNew {
		s = strings.ReplaceAll(s, old, new)
	}
	return []byte(s)
}
