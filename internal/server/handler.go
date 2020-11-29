package server

import (
	"net/http"
	"strings"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/reactserv/internal/models"
)

func newHandler(rootURL string, logger logging.Logger,
	memFS http.FileSystem, buildInfo models.BuildInformation) http.Handler {
	return &handler{
		rootURL:    rootURL,
		logger:     logger,
		fileServer: http.FileServer(memFS),
		buildInfo:  buildInfo,
	}
}

type handler struct {
	rootURL    string
	logger     logging.Logger
	fileServer http.Handler
	buildInfo  models.BuildInformation
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = strings.TrimPrefix(r.RequestURI, h.rootURL)
	h.logger.Info("%s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
	if strings.HasPrefix(r.RequestURI, "/v1/version") && r.Method == http.MethodGet {
		h.getVersion(w)
		return
	}
	h.fileServer.ServeHTTP(w, r)
}
