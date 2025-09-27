package secret

type Secret struct {
	ID int64
	Key string
	Value string
	ChatId int64
	KeyVersion uint64
	IV string
}