package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func Initialize(level string) {
	var l slog.LevelVar
	switch level {
	case "info":
		l.Set(0)
	case "debug":
		l.Set(-4)
	case "warn":
		l.Set(4)
	case "error":
		l.Set(8)
	}
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			AddSource:  true,
			Level:      &l,
			TimeFormat: time.Kitchen,
		}),
	))
}
