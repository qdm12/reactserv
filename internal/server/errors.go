package server

import (
	"errors"
	"net/http"
)

var (
	ErrBodyDecode = errors.New("cannot decode request body")
)

func httpError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
