package secret

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	GetSecretsForChatId(ctx context.Context, id int64, limit, offet int) (*[]Secret, error)
	InsertSecret(ctx context.Context, secret *Secret) error
	DeleteSecret(ctx context.Context, id int64) error
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{db: db}, nil
}

func (r repo) GetSecretsForChatId(ctx context.Context, id int64, limit, offset int) (*[]Secret, error) {
	query := "SELECT id, key, value, chat_id, key_version, iv_b64 FROM secrets WHERE chat_id = ? ORDER BY id LIMIT ? OFFSET ?"
	query = r.db.Rebind(query)
	s := []Secret{}
	err := r.db.SelectContext(ctx, &s, query, id, limit, offset)
	return &s, err
}

func (r repo) InsertSecret(ctx context.Context, secret *Secret) error {
	query := "INSERT INTO secrets(key, value, chat_id, key_version, iv_b64) VALUES (?,?,?,?,?)"
	query = r.db.Rebind(query)
	_, err := r.db.ExecContext(ctx, query, &secret.Key, &secret.Value, &secret.ChatId, &secret.KeyVersion, &secret.IV)
	return err
}

func (r repo) DeleteSecret(ctx context.Context, id int64) error {
	query := "DELETE FROM secrets WHERE id = ?"
	query = r.db.Rebind(query)
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
