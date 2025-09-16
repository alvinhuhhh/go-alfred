package dinner

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	GetDinnerByDateAndChatId(ctx context.Context, chatId int64, date time.Time) (*Dinner, error)
	InsertDinner(ctx context.Context, dinner *Dinner) (int64, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{db: db}, nil
}

func (r repo) GetDinnerByDateAndChatId(ctx context.Context, chatId int64, date time.Time) (*Dinner, error) {
	query := "SELECT chatId, date, yes, no, messageIds FROM dinner WHERE chatId = ? AND date = ?"
	query = r.db.Rebind(query)
	var d Dinner
	err := r.db.GetContext(ctx, &d, query, chatId, date)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r repo) InsertDinner(ctx context.Context, dinner *Dinner) (int64, error) {
	query := "INSERT INTO dinners(chatId, date, yes, no, messageIds) VALUES (?,?,?,?,?) RETURNING id"
	query = r.db.Rebind(query)
	var id int64
	err := r.db.QueryRowContext(ctx, query, &dinner.ChatID, &dinner.Date, &dinner.Yes, &dinner.No, &dinner.MessageIds).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
