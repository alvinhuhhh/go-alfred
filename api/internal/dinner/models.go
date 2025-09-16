package dinner

import (
	"time"

	"github.com/lib/pq"
)

type Dinner struct {
	Date       time.Time     `db:"date"`
	ChatID     int64         `db:"chat_id"`
	Yes        []string      `db:"yes"`
	No         []string      `db:"no"`
	MessageIds pq.Int64Array `db:"message_ids"`
}
