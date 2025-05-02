package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/pedramktb/go-base-lib/lifecycle"
	slogctx "github.com/veqryn/slog-context"
)

func HTTPS(ctx context.Context, cancel context.CancelFunc, server *http.Server) {
	var svErr = make(chan error, 1)
	go func() {
		slogctx.FromCtx(ctx).Info("starting server", slog.String("addr", server.Addr))
		svErr <- server.ListenAndServeTLS("", "")
		cancel()
	}()

	go func() {
		done, err := lifecycle.RegisterCloser(ctx)
		if err == nil {
			defer done()
		}
		<-ctx.Done()
		if err := server.Shutdown(ctx); err != nil {
			slogctx.FromCtx(ctx).Error("failed to shutdown server gracefully", slog.Any("error", err))
		}
		err = <-svErr
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slogctx.FromCtx(ctx).Error("server closed with an error", slog.Any("error", err))
		} else {
			slogctx.FromCtx(ctx).Info("server closed")
		}
	}()
}
