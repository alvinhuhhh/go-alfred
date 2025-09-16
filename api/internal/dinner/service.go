package dinner

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/alvinhuhhh/go-alfred/internal/chat"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Service interface {
	HandleDinner(ctx context.Context, b *bot.Bot, update *models.Update)
	HandleCallbackQuery(ctx context.Context, b *bot.Bot, update *models.Update)
}

type service struct {
	repo     Repo
	chatRepo chat.Repo
}

func NewService(r Repo, cr chat.Repo) (Service, error) {
	return &service{
		repo:     r,
		chatRepo: cr,
	}, nil
}

func (s service) HandleDinner(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := s.verifyChat(ctx, b, update)
	if err != nil {
		slog.Error("error verifying chat")
		return
	}
	command := update.Message.Text

	switch command {
	case "/getdinner":
		d, err := s.getOrInsertDinner(ctx, b, update)
		if err != nil {
			slog.Error("error getting dinner")
			return
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        s.parseDinnerMessage(d),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: s.getKeyboard(),
		})
		return

	case "/joindinner":
		slog.Info("Handle joindinner")
		return

	case "/leavedinner":
		slog.Info("Handle leavedinner")
		return

	case "/enddinner":
		return

	default:
		return
	}
}

func (s service) HandleCallbackQuery(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Answer callback query first so that Telegram stops spamming updates
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	s.HandleDinner(ctx, b, update)
}

func (s service) verifyChat(ctx context.Context, b *bot.Bot, update *models.Update) (*chat.Chat, error) {
	c, err := s.chatRepo.GetChatByID(ctx, update.Message.Chat.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			// Handle any other error
			slog.Error(err.Error())
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Sorry! Having a bit of trouble, will be back soon!",
			})
			return nil, err
		}

		// Insert new Chat
		c = &chat.Chat{
			ID:   update.Message.Chat.ID,
			Type: string(update.Message.Chat.Type),
		}
		id, err := s.chatRepo.InsertChat(ctx, c)
		if err != nil {
			slog.Error(err.Error())
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Sorry! Having a bit of trouble, will be back soon!",
			})
			return nil, err
		}
		slog.Info(fmt.Sprintf("inserted chat id: %v", id))
	}
	return c, nil
}

func (s service) getOrInsertDinner(ctx context.Context, b *bot.Bot, update *models.Update) (*Dinner, error) {
	d, err := s.repo.GetDinnerByDateAndChatId(ctx, update.Message.Chat.ID, time.Now())
	if err != nil {
		if err != sql.ErrNoRows {
			// Any other error
			slog.Error(err.Error())
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Sorry! Having a bit of trouble, will be back soon!",
			})
			return nil, err
		}

		// Insert new Dinner
		d = &Dinner{
			ChatID:     update.Message.Chat.ID,
			Date:       time.Now(),
			Yes:        []string{update.Message.From.FirstName},
			No:         []string{},
			MessageIds: []int{update.Message.ID},
		}
		id, err := s.repo.InsertDinner(ctx, d)
		if err != nil {
			slog.Error(err.Error())
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Sorry! Having a bit of trouble, will be back soon!",
			})
			return nil, err
		}
		slog.Info(fmt.Sprintf("inserted dinner id: %v", id))
	}
	return d, nil
}

func (s service) getKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Join Dinner", CallbackData: "Yes"},
				{Text: "Leave Dinner", CallbackData: "No"},
			},
		},
	}
}

func (s service) parseDinnerMessage(d *Dinner) string {
	date := d.Date
	var yes, no string
	for _, v := range d.Yes {
		yes += fmt.Sprintf("%s\n", v)
	}
	for _, v := range d.No {
		no += fmt.Sprintf("%s\n", v)
	}
	return fmt.Sprintf("\n<b>Dinner tonight:</b>\nDate: %s\n\n<u>YES:</u>\n%s\n\n<u>NO:</u>\n%s", date, yes, no)
}
