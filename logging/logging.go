package logging

import (
	"context"
	"io"
	"log/slog"

	"github.com/pedramktb/go-base-lib/env"
	slogctx "github.com/veqryn/slog-context"
)

func handler(writer io.Writer) slog.Handler {
	var handler slog.Handler
	switch env.GetEnvironment() {
	case env.EnvironmentLocal:
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	case env.EnvironmentDev, env.EnvironmentStaging:
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	case env.EnvironmentProd:
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelInfo,
		})
	}

	return handler
}

func NewLoggerCtx(ctx context.Context, writer io.Writer, prependers ...slogctx.AttrExtractor) context.Context {
	return slogctx.NewCtx(ctx, NewLogger(writer, prependers...))
}

func NewLogger(writer io.Writer, prependers ...slogctx.AttrExtractor) *slog.Logger {
	prependers = append(prependers, slogctx.ExtractPrepended)
	return slog.New(slogctx.NewHandler(handler(writer), &slogctx.HandlerOptions{
		Prependers: prependers,
		Appenders: []slogctx.AttrExtractor{
			slogctx.ExtractAppended,
		},
	}))
}
