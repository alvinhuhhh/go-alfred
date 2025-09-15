package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		uri := r.URL.Path

		slog.Info(fmt.Sprintf("%s %s", method, uri))
		next.ServeHTTP(w, r)
	})
}

func SetAccessControlHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		next.ServeHTTP(w, r)
	})
}

func LogBotRequests(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatId := update.Message.Chat.ID
		username := update.Message.From.Username
		command := update.Message.Text

		slog.Info(fmt.Sprintf("bot command %s received from %v in chat id %d", command, username, chatId))
		next(ctx, b, update)
	}
}
