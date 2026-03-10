package example

import "log/slog"

func main() {
	var log *slog.Logger
	password := "secret"

	log.Info("Starting server!!!")
	log.Info("user password is " + password)
}
