package handlers

import (
	"embed"
	"net/http"
)

// embeds the filesystem at build time
var webFS embed.FS

type httpHandler struct{}

func NewHttpHandler() (*httpHandler, error) {
	return &httpHandler{}, nil
}

func (h httpHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) Serve(w http.ResponseWriter, r *http.Request) {
	fs := http.FS(webFS)
	path := r.URL.Path
	if f, err := fs.Open(path); err == nil {
		f.Close()
		http.FileServer(fs).ServeHTTP(w, r)
		return
	}
	// Fallback to index.html
	http.ServeFileFS(w, r, webFS, "dist/index.html")
}
