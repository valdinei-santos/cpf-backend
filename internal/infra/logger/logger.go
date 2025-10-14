package logger

import (
	"context"
	"log/slog"
	"os"
)

/*
	 type stackFrame struct {
		Func   string `json:"func"`
		Source string `json:"source"`
		Line   int    `json:"line"`
	}
*/
type SlogILogger struct {
	logger *slog.Logger
}

func NewSlogILogger() *SlogILogger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		//Level: slog.LevelInfo,
		Level: slog.LevelDebug,
	})
	log := slog.New(handler)
	return &SlogILogger{logger: log}
}

func (l *SlogILogger) Debug(msg string, args ...any) {
	l.logger.DebugContext(context.Background(), msg, args...)
}

func (l *SlogILogger) Info(msg string, args ...any) {
	l.logger.InfoContext(context.Background(), msg, args...)
}

func (l *SlogILogger) Warn(msg string, args ...any) {
	l.logger.WarnContext(context.Background(), msg, args...)
}

func (l *SlogILogger) Error(msg string, args ...any) {
	l.logger.ErrorContext(context.Background(), msg, args...)
}

func (l *SlogILogger) With(args ...any) ILogger {
	return &SlogILogger{logger: l.logger.With(args...)}
}

func (l *SlogILogger) WithContext(ctx context.Context) ILogger {
	return &SlogILogger{logger: l.logger}
}

func (l *SlogILogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *SlogILogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *SlogILogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *SlogILogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}
