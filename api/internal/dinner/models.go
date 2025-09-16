package dinner

import "time"

type Dinner struct {
	Date time.Time `db:"date"`
	ChatID int64 `db:"chatId"`
	Yes []string `db:"yes"`
	No []string `db:"no"`
	MessageIds []int `db:"messageIds"`
}