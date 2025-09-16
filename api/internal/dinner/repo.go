package dinner

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repo interface {
	GetDinnerByDateAndChatId(ctx context.Context, chatId int64, date time.Time) (*Dinner, error)
	InsertDinner(ctx context.Context, d *Dinner) (int64, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{db: db}, nil
}

func (r repo) GetDinnerByDateAndChatId(ctx context.Context, chatId int64, date time.Time) (*Dinner, error) {
	query := "SELECT chat_id, date, yes, no, message_ids FROM dinners WHERE chat_id = ? AND date = ?"
	query = r.db.Rebind(query)
	var d Dinner
	err := r.db.QueryRowContext(ctx, query, chatId, date.Format("2006-01-02")).Scan(&d.ChatID, &d.Date, pq.Array(&d.Yes), pq.Array(&d.No), &d.MessageIds)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r repo) InsertDinner(ctx context.Context, d *Dinner) (int64, error) {
	query := "INSERT INTO dinners(chat_id, date, yes, no, message_ids) VALUES (?,?,?,?,?) RETURNING id"
	query = r.db.Rebind(query)
	var id int64
	err := r.db.QueryRowContext(ctx, query, &d.ChatID, d.Date.Format("2006-01-02"), pq.Array(&d.Yes), pq.Array(&d.No), &d.MessageIds).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
