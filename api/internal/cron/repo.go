package cron

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	Schedule(ctx context.Context, job *CronJob) error
	Unschedule(ctx context.Context, jobname string) error
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{db: db}, nil
}

func (r repo) Schedule(ctx context.Context, job *CronJob) error {
	// Build the JSON body
	bodyData := map[string]interface{}{
		"chat_id": job.ChatId,
	}
	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}
	body := string(bodyBytes)

	// Construct the command string that pg_cron will execute
	command := fmt.Sprintf(
		"SELECT * FROM http_post('%s', '%s', 'application/json')",
		job.URL, body,
	)

	// Outer query uses placeholders for the cron function arguments
	query := `SELECT cron.schedule(?, ?, ?)`
	query = r.db.Rebind(query)

	if _, err := r.db.ExecContext(ctx, query, job.JobName, job.Schedule, command); err != nil {
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
