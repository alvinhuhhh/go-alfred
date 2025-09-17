package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
)

type httpHandler struct{}

func NewHttpHandler() (*httpHandler, error) {
	return &httpHandler{}, nil
}

func (h httpHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) Serve(w http.ResponseWriter, r *http.Request) {
	fs := http.Dir("./dist")

	// Never let /api be handled here
	if strings.HasPrefix(r.URL.Path, "/api/") {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Try file
	reqPath := r.URL.Path
	if strings.HasSuffix(reqPath, "/") {
		reqPath += "index.html"
	}
	full := filepath.Clean("/" + reqPath)
	f, err := fs.Open(full)
	if err == nil {
		defer f.Close()
		// File exists - serve it
		http.FileServer(fs).ServeHTTP(w, r)
		return
	}

	// Otherwise serve index.html
	http.ServeFile(w, r, "./dist/index.html")
}
