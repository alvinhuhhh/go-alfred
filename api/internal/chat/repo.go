package chat

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repo interface{
	GetChatByID(ctx context.Context, id int64) (*Chat, error)
	InsertChat(ctx context.Context, chat *Chat) (int64, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{db: db}, nil
}

func (r repo) GetChatByID(ctx context.Context, id int64) (*Chat, error) {
	query := "SELECT id, type FROM chats WHERE id = ?"
	query = r.db.Rebind(query)
	var c Chat
	if err := r.db.GetContext(ctx, &c, query, id); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r repo) InsertChat(ctx context.Context, chat *Chat) (int64, error) {
	query := "INSERT INTO chats(id, type) VALUES (?,?)"
	query = r.db.Rebind(query)
	_, err := r.db.ExecContext(ctx, query, &chat.ID, &chat.Type)
	if err != nil {
		return -1, err
	}
	return chat.ID, nil
}