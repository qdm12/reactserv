package params

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qdm12/golibs/logging"
	libparams "github.com/qdm12/golibs/params"
)

type Reader interface {
	GetListeningPort() (listeningPort uint16, warning string, err error)
	GetLoggerConfig() (encoding logging.Encoding, level logging.Level, err error)
	GetRootURL(setters ...libparams.GetEnvSetter) (rootURL string, err error)
	GetRootDir() (rootDir string, err error)
}

type reader struct {
	envParams libparams.EnvParams
}

func NewReader() Reader {
	return &reader{
		envParams: libparams.NewEnvParams(),
	}
}

func (r *reader) GetListeningPort() (listeningPort uint16, warning string, err error) {
	return r.envParams.GetListeningPort("LISTENING_PORT", libparams.Default("8000"))
}

func (r *reader) GetLoggerConfig() (encoding logging.Encoding, level logging.Level, err error) {
	return r.envParams.GetLoggerConfig()
}

func (r *reader) GetRootURL(setters ...libparams.GetEnvSetter) (rootURL string, err error) {
	return r.envParams.GetRootURL()
}

func (r *reader) GetRootDir() (rootDir string, err error) {
	rootDir, err = r.envParams.GetEnv("ROOT_DIR", libparams.Default("srv"))
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(rootDir)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrBadFilepath, err)
	}
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrBadFilepath, err)
	}
	stats, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrBadFile, err)
	} else if !stats.IsDir() {
		return "", fmt.Errorf("%w: %s", ErrIsNotDir, err)
	}
	return rootDir, nil
}
