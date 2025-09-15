package handlers

import "net/http"

type httpHandler struct {
}

func NewHttpHandler() (*httpHandler, error) {
	return &httpHandler{}, nil
}

func (h httpHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h httpHandler) CronTrigger(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
