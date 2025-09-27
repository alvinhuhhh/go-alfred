package secret

import (
	"encoding/base64"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alvinhuhhh/go-alfred/internal/util"
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
