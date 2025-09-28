package secret

type Secret struct {
	ID int64 `db:"id" json:"id,omitempty"`
	Key string `db:"key" json:"key"`
	Value string `db:"value" json:"value"`
	ChatId int64 `db:"chat_id" json:"chatId"`
	KeyVersion uint64 `db:"key_version" json:"keyVersion"`
	IV string `db:"iv_b64" json:"ivB64"`
}