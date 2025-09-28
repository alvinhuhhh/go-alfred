package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"

	"github.com/alvinhuhhh/go-alfred/internal/chat"
	"github.com/alvinhuhhh/go-alfred/internal/config"
	"github.com/alvinhuhhh/go-alfred/internal/dinner"
	"github.com/alvinhuhhh/go-alfred/internal/handlers"
	"github.com/alvinhuhhh/go-alfred/internal/middleware"
	"github.com/alvinhuhhh/go-alfred/internal/secret"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	port, err := getPort()
	if err != nil {
		log.Fatal(err)
	}

	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = runMigrations(db.DB)
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	// Bot handler
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update.Message != nil {
				slog.Warn(fmt.Sprintf("unhandled message with id: %s", strconv.Itoa(update.Message.ID)))
			}
			if update.CallbackQuery != nil {
				slog.Warn(fmt.Sprintf("unhandled callback with id: %v", update.CallbackQuery.ID))

				// Answer callback query first so that Telegram stops spamming updates
				b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
					CallbackQueryID: update.CallbackQuery.ID,
					ShowAlert:       false,
				})
			}
		}),
		bot.WithMiddlewares(middleware.LogBotRequests),
	}

	if config.IsTestServer() {
		opts = append(opts, bot.UseTestEnvironment())
	}
	token, err := config.GetBotToken()
	if err != nil {
		log.Fatal(err)
	}
	b, err := bot.New(token, opts...)
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

	dinnerRepo, err := dinner.NewRepo(db)
	if err != nil {
		log.Fatal(err)
	}
	dinnerService, err := dinner.NewService(b, dinnerRepo, chatRepo)
	if err != nil {
		log.Fatal(err)
	}

	secretRepo, err := secret.NewRepo(db)
	if err != nil {
		log.Fatal(err)
	}
	secretService, err := secret.NewService(secretRepo)
	if err != nil {
		log.Fatal(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, chatService.Start)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/app", bot.MatchTypePrefix, chatService.ReplyHello)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/hello", bot.MatchTypePrefix, chatService.ReplyHello)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/getdinner", bot.MatchTypePrefix, dinnerService.HandleDinner)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/enddinner", bot.MatchTypePrefix, dinnerService.HandleDinner)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "joindinner", bot.MatchTypePrefix, dinnerService.HandleCallbackQuery)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "leavedinner", bot.MatchTypePrefix, dinnerService.HandleCallbackQuery)

	go b.StartWebhook(ctx)
	slog.Info("Bot webhook listener started")

	// HTTP handler
	httpHandler, err := handlers.NewHttpHandler()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	// API router
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.SetAccessControlHeaders)
	api.Use(middleware.LogRequests)
	api.Use(middleware.Auth)

	api.HandleFunc("/webhook", b.WebhookHandler()).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions) // routes to Bot handlers
	api.HandleFunc("/ping", httpHandler.Ping).Methods(http.MethodGet)
	api.HandleFunc("/cron", dinnerService.CronTrigger).Methods(http.MethodPost)

	api.HandleFunc("/encryption/key", secretService.GetDataEncryptionKey).Methods(http.MethodGet)
	api.HandleFunc("/secrets/{chatId}", secretService.GetSecretsForChatId).Methods(http.MethodGet)
	api.HandleFunc("/secrets", secretService.InsertSecret).Methods(http.MethodPost)
	api.HandleFunc("/secrets/{id}", secretService.DeleteSecret).Methods(http.MethodDelete)

	// Web router
	web := router.NewRoute().Subrouter()
	web.Use(middleware.LogRequests)
	web.PathPrefix("/").HandlerFunc(httpHandler.Serve)

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

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	mig, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}
	return mig.Up()
}
