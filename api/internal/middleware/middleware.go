package middleware

import (
	"slices"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/alvinhuhhh/go-alfred/internal/config"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	initdata "github.com/telegram-mini-apps/init-data-golang"
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
		if update.Message != nil {
			slog.Info(fmt.Sprintf("bot command %s received from %v in chat id %d", update.Message.Text, update.Message.From.Username, update.Message.Chat.ID))
		}
		if update.CallbackQuery != nil {
			slog.Info(fmt.Sprintf("bot callback %s received from %v", update.CallbackQuery.Data, update.CallbackQuery.From.Username))
		}
		next(ctx, b, update)
	}
}

func Auth(next http.Handler) http.Handler {
	whitelist := []string{
		"/api/webhook",
		"/api/ping",
		"/api/cron",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow routes in whitelist
		if slices.Contains(whitelist, r.URL.Path) {
			next.ServeHTTP(w, r)
			return 
		}

		// Get raw init data
		auth := r.Header.Get("Authorization")
		authSplit := strings.Split(auth, " ")
		if authSplit[0] != "tma" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Get Bot token
		token, err := config.GetBotToken()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Set expiry time
		exp := 24 * time.Hour

		if err := initdata.Validate(authSplit[1], token, exp); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
