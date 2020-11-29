package server

import (
	"encoding/json"
	"net/http"
)

func (h *handler) getVersion(w http.ResponseWriter) {
	b, err := json.Marshal(h.buildInfo)
	if err != nil {
		h.logger.Error(err)
		httpError(w, http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(b); err != nil {
		h.logger.Error(err)
	}
}
