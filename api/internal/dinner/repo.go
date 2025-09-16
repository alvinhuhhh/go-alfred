package dinner

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repo interface {
	GetDinnerById(ctx context.Context, id int64) (*Dinner, error)
	GetDinnerByDateAndChatId(ctx context.Context, chatId int64, date time.Time) (*Dinner, error)
	InsertDinner(ctx context.Context, d *Dinner) (int64, error)
	UpdateDinner(ctx context.Context, d *Dinner) error
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{db: db}, nil
}

func (r repo) GetDinnerById(ctx context.Context, id int64) (*Dinner, error) {
	query := "SELECT id, chat_id, date, yes, no, message_ids FROM dinners WHERE id = ?"
	query = r.db.Rebind(query)
	var d Dinner
	err := r.db.QueryRowContext(ctx, query, id).Scan(&d.ID, &d.ChatID, &d.Date, pq.Array(&d.Yes), pq.Array(&d.No), &d.MessageIds)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r repo) GetDinnerByDateAndChatId(ctx context.Context, chatId int64, date time.Time) (*Dinner, error) {
	query := "SELECT id, chat_id, date, yes, no, message_ids FROM dinners WHERE chat_id = ? AND date = ?"
	query = r.db.Rebind(query)
	var d Dinner
	err := r.db.QueryRowContext(ctx, query, chatId, date.Format("2006-01-02")).Scan(&d.ID, &d.ChatID, &d.Date, pq.Array(&d.Yes), pq.Array(&d.No), &d.MessageIds)
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

func (r repo) UpdateDinner(ctx context.Context, d *Dinner) error {
	query := "UPDATE dinners SET yes = ?, no = ?, message_ids = ? WHERE id = ?"
	query = r.db.Rebind(query)
	_, err := r.db.ExecContext(ctx, query, pq.Array(&d.Yes), pq.Array(&d.No), &d.MessageIds, &d.ID)
	if err != nil {
		return err
	}
	return nil
}