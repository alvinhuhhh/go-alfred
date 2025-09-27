package secret

import (
	"net/http"
)

type Service interface {
	GetDataEncryptionKey(w http.ResponseWriter, r *http.Request)
}

type service struct {
	repo Repo
}

func NewService(r Repo) (Service, error) {
	return &service{repo: r}, nil
}

func (s *service) GetDataEncryptionKey(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
