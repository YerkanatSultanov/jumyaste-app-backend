package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

var Log *slog.Logger

func InitLogger() {
	logHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "15:04:05",
	})
	Log = slog.New(logHandler)
}
