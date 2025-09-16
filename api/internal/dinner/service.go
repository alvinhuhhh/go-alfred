package dinner

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/alvinhuhhh/go-alfred/internal/chat"
	"github.com/alvinhuhhh/go-alfred/internal/util"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/lib/pq"
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
			ReplyMarkup: s.getKeyboard(d.ID),
		})
		return

	case "/enddinner":
		d, err := s.getOrInsertDinner(ctx, b, update)
		if err != nil {
			slog.Error("error getting dinner")
			return
		}
		err = s.repo.DeleteDinner(ctx, d.ID)
		if err != nil {
			slog.Error("unable to delete dinner")
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text: "Sorry, I can't delete tonight's dinner",
			})
			return
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: "No more dinner for tonight!",
		})
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

	callbackData := update.CallbackQuery.Data
	name := update.CallbackQuery.From.FirstName

	// Get dinner
	split := strings.Split(callbackData, "_")
	dinnerId, err := strconv.Atoi(split[1])
	if err != nil {
		slog.Error("unable to parse dinner id from callback query data")
		return
	}
	dinner, err := s.repo.GetDinnerById(ctx, int64(dinnerId))
	if err != nil {
		slog.Error("unable to get dinner")
		return
	}

	switch split[0] {
	case "joindinner":
		dinner.Yes = append(dinner.Yes, name)
		if slices.Contains(dinner.No, name) {
			dinner.No = util.Remove(dinner.No, name)
		}

	case "leavedinner":
		dinner.No = append(dinner.No, name)
		if slices.Contains(dinner.Yes, name) {
			dinner.Yes = util.Remove(dinner.Yes, name)
		}

	default:
		slog.Warn(fmt.Sprintf("unknown callback: %s", split[0]))
		return
	}

	// Send response
	message, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text: s.parseDinnerMessage(dinner),
		ParseMode: models.ParseModeHTML,
		ReplyMarkup: s.getKeyboard(dinner.ID),
	})
	if err != nil {
		slog.Error(err.Error())
	}

	// Append message ID
	dinner.MessageIds = append(dinner.MessageIds, int64(message.ID))

	// Insert db
	if err := s.repo.UpdateDinner(ctx, dinner); err != nil {
		slog.Error(err.Error())
	}
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
			MessageIds: pq.Int64Array{int64(update.Message.ID)},
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
		d.ID = id
	}
	return d, nil
}

func (s service) getKeyboard(id int64) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Join Dinner", CallbackData: fmt.Sprintf("joindinner_%d", id)},
				{Text: "Leave Dinner", CallbackData: fmt.Sprintf("leavedinner_%d", id)},
			},
		},
	}
}

func (s service) parseDinnerMessage(d *Dinner) string {
	date := d.Date
	yes := strings.Join(d.Yes, "\n")
	no := strings.Join(d.No, "\n")
	return fmt.Sprintf("\n<b>Dinner tonight:</b>\nDate: %s\n\n<u>YES:</u>\n%s\n\n<u>NO:</u>\n%s\n\n", date.Format("02/01/2006"), yes, no)
}
