package chat

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Service interface {
	Start(ctx context.Context, b *bot.Bot, update *models.Update)
	ReplyHello(ctx context.Context, b *bot.Bot, update *models.Update)
}

type service struct {
	repo Repo
}

func NewService(r Repo) (Service, error) {
	return &service{repo: r}, nil
}

func (s *service) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := s.repo.GetChatByID(ctx, update.Message.Chat.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Hello there! Let me get set up and we will be ready to go",
			})

			// Insert new Chat
			id, err := s.repo.InsertChat(ctx, &Chat{
				ID:   update.Message.Chat.ID,
				Type: string(update.Message.Chat.Type),
			})
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Sorry! Having a bit of trouble, will be back soon!",
				})
				return
			}

			slog.Info(fmt.Sprintf("inserted chat id: %v", id))
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Done! How can I help?",
			})
			return
		}

		// Handle any other error
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Sorry! Having a bit of trouble, will be back soon!",
		})
		return
	}

	slog.Info("chat already exists")
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Welcome back!",
	})
}

func (s *service) ReplyHello(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hello there! What can I do for you today?",
	})
}
