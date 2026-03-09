package linter_test

import (
	"log/slog"
)

type zapLogger struct{}

func (z *zapLogger) Info(msg string, args ...interface{}) {}

func test() {
	var log *slog.Logger
	var zap *zapLogger
	password := "12345"

	// fine messages
	log.Info("starting server")
	zap.Info("process finished")

	// rule 1 - first letter
	log.Info("Starting server")   // want "log message should start with a lowercase letter"
	zap.Info("Failed to connect") // want "log message should start with a lowercase letter"

	// rule 2 - cyrillic symbols
	log.Info("запуск сервера") // want "log message should be in English only"

	// rule 3 - specials
	log.Info("server started!")         // want "log message contains forbidden symbols or emojis"
	log.Info("connection failed!!!")    // want "log message contains forbidden symbols or emojis"
	log.Info("we are flying 🚀")         // want "log message contains forbidden symbols or emojis"
	log.Info("something went wrong...") // want "log message contains forbidden symbols or emojis"

	// rule 4 - sensitive data
	log.Info("user password is " + password) // want "potential sensitive data exposure: password" "avoid logging sensitive variables like 'password'"
	log.Info("api key is missing")           // want "potential sensitive data exposure: api key"

	token := "a"
	log.Info("validating", "token", token) // want "potential sensitive data exposure: token" "avoid logging sensitive variables like 'token'"
}
