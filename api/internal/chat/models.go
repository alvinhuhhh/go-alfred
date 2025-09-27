package chat

type Chat struct {
	ID int64 `db:"id"`
	Type string `db:"type"`
}