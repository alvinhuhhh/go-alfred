package secret

import "github.com/jmoiron/sqlx"

type Repo interface {}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{db: db}, nil
}