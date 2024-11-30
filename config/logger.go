package config

import (
	"log/slog"
	"os"
)

func SetUpLogger() {
	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)

}
