package config

import (
	"errors"
	"os"
)

func IsTestServer() bool {
	t, found := os.LookupEnv("TELEGRAM_TEST_SERVER")
	if !found {
		return false
	}
	return t == "1"
}

func GetBotToken() (string, error) {
	t, found := os.LookupEnv("BOT_TOKEN")
	if !found {
		return "", errors.New("BOT_TOKEN is not found")
	}
	return t, nil
}
