package secret

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alvinhuhhh/go-alfred/internal/util"
	"github.com/gorilla/mux"
)

type Service interface {
	GetDataEncryptionKey(w http.ResponseWriter, r *http.Request)
	GetSecretsForChatId(w http.ResponseWriter, r *http.Request)
	InsertSecret(w http.ResponseWriter, r *http.Request)
	DeleteSecret(w http.ResponseWriter, r *http.Request)
}

type service struct {
	repo Repo
}

func NewService(r Repo) (Service, error) {
	return &service{repo: r}, nil
}

func (s *service) GetDataEncryptionKey(w http.ResponseWriter, r *http.Request) {
	kv := r.URL.Query().Get("keyVersion")
	keyVersion, err := strconv.ParseUint(kv, 10, 64)
	if err != nil {
		slog.Error("unable to parse keyVersion from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c := r.URL.Query().Get("chatId")
	chatId, err := strconv.ParseInt(c, 10, 64)
	if err != nil {
		slog.Error("unable to parse chatId from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	master, err := util.GetMasterKey(keyVersion)
	if err != nil {
		slog.Error("unable to get master key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dek, err := util.DeriveDEK(master, keyVersion, chatId)
	if err != nil {
		slog.Error("unable to derive encryption key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base64.StdEncoding.EncodeToString(dek)))
}

func (s *service) GetSecretsForChatId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["chatId"]
	chatId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("unable to parse chatId from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	secrets, err := s.repo.GetSecretsForChatId(r.Context(), chatId, 100, 0)
	if err != nil {
		slog.Error("error fetching secrets")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(secrets)
}

func (s service) InsertSecret(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var secret Secret
	if err := decoder.Decode(&secret); err != nil {
		slog.Error("error parsing request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := s.repo.InsertSecret(r.Context(), &secret); err != nil {
		slog.Error("unable to insert secret")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s service) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	id, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		slog.Error("unable to parse id from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := s.repo.DeleteSecret(r.Context(), id); err != nil {
		slog.Error("error deleting secret")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
