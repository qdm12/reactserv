package params

import "errors"

var (
	ErrBadFilepath = errors.New("bad filepath")
	ErrBadFile     = errors.New("bad file")
	ErrIsNotDir    = errors.New("file is not a directory")
)
