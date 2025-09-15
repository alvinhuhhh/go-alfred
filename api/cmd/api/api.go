package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/alvinhuhhh/go-alfred/internal/chat"
	"github.com/alvinhuhhh/go-alfred/internal/handlers"
	"github.com/alvinhuhhh/go-alfred/internal/middleware"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	port, err := getPort()
	if err != nil {
		log.Fatal(err)
	}

	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}

	httpHandler, err := handlers.NewHttpHandler()
	if err != nil {
		log.Fatal(err)
	}

	chatRepo, err := chat.NewRepo(db)
	if err != nil {
		log.Fatal(err)
	}
	chatService, err := chat.NewService(chatRepo)
	if err != nil {
		log.Fatal(err)
	}

	// Bot handler
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Hello there! What can I do for you today?",
			})
		}),
	}
	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		log.Fatal(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "hello", bot.MatchTypeCommand, chatService.ReplyHello)

	go b.StartWebhook(ctx)
	slog.Info("Bot webhook listener started")

	// HTTP handler
	router := mux.NewRouter()
	router.Use(middleware.SetAccessControlHeaders)
	router.Use(middleware.LogRequests)

	router.HandleFunc("/ping", httpHandler.Ping).Methods(http.MethodGet)
	router.HandleFunc("/cron", httpHandler.CronTrigger).Methods(http.MethodGet)

	slog.Info(fmt.Sprintf("Alfred has started listening on port: %s", port))
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

func getPort() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("unable to get PORT from env")
	}
	return port, nil
}

func getDB() (*sqlx.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	parseConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parsing DSN failed")
	}
	connConfig := stdlib.RegisterConnConfig(parseConfig)
	db, err := sqlx.Open("pgx", connConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
