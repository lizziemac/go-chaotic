package logger

import (
	"context"
	"log/slog"
	"os"
)

var Logger *slog.Logger

// Add custom log levels
const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

func init() {

	var appEnv = os.Getenv("APP_ENV")

	opts := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// uses LevelNames map to print actual level instead of, for example instead of ERROR+4, it will say FATAL
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	}

	// Text output to stdout with source info
	var handler slog.Handler = slog.NewTextHandler(os.Stdout, &opts)
	if appEnv == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &opts)
	}
	Logger = slog.New(handler)

	slog.SetDefault(Logger)
}

func Info(msg string, args ...any)  { Logger.Info(msg, args...) }
func Debug(msg string, args ...any) { Logger.Debug(msg, args...) }
func Error(msg string, args ...any) { Logger.Error(msg, args...) }
func Fatal(msg string, args ...any) {
	Logger.Log(context.Background(), slog.Level(LevelFatal), msg, args...)
}
