package cron

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Service interface {
	CreateCronJob(w http.ResponseWriter, r *http.Request)
	RemoveCronJob(w http.ResponseWriter, r *http.Request)
}

type service struct {
	repo Repo
}

func NewService(r Repo) (Service, error) {
	return &service{
		repo: r,
	}, nil
}

func (s service) CreateCronJob(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var job CronJob
	if err := decoder.Decode(&job); err != nil {
		slog.Error(err.Error())
		slog.Error("error parsing request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := s.repo.Schedule(r.Context(), &job); err != nil {
		slog.Error(err.Error())
		slog.Error("unable to schedule CronJob")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s service) RemoveCronJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["jobName"]
	if err := s.repo.Unschedule(r.Context(), jobName); err != nil {
		slog.Error(err.Error())
		slog.Error("error deleting secret")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
