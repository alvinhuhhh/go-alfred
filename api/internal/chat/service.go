package chat

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Service interface{
	ReplyHello(ctx context.Context, b *bot.Bot, update *models.Update)
}

type service struct {
	repo Repo
}

func NewService(r Repo) (Service, error) {
	return &service{repo: r}, nil
}

func (s *service) ReplyHello(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "Hello there! What can I do for you today?",
	})
}