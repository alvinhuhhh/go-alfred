package cron

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	Schedule(ctx context.Context, jobname, schedule, url string, chat_id int64) error
	Unschedule(ctx context.Context, jobname string) error
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{db: db}, nil
}

func (r repo) Schedule(ctx context.Context, jobname, schedule, url string, chat_id int64) error {
	query := `
		SELECT cron.schedule(
			?,
			?,
			$$
				SELECT
					net.http_post(
						url:=?,
						body:=?
					) AS request_id
			$$
		)
	`
	query = r.db.Rebind(query)

	bodyData := map[string]interface{}{
		"chat_id": chat_id,
	}
	body, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, query, jobname, schedule, url, body); err != nil {
		return err
	}
	return nil
}

func (r repo) Unschedule(ctx context.Context, jobname string) error {
	query := "SELECT cron.unschedule(?)"
	query = r.db.Rebind(query)
	_, err := r.db.ExecContext(ctx, query, jobname)
	if err != nil {
		return err
	}
	return nil
}
