package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/pedramktb/go-base-lib/lifecycle"
	slogctx "github.com/veqryn/slog-context"
)

// HTTPS starts the provided HTTP server that must be configured with TLS.
// It uses an optional cancel function that is called when the server is no longer running.
func HTTPS(ctx context.Context, cancel context.CancelFunc, server *http.Server) {
	var svErr = make(chan error, 1)
	go func() {
		slogctx.FromCtx(ctx).Info("starting server", slog.String("addr", server.Addr))
		svErr <- server.ListenAndServeTLS("", "")
		if cancel != nil {
			cancel()
		}
	}()

	go func() {
		if done, err := lifecycle.RegisterCloser(ctx); err == nil {
			defer func() {
				err := server.Shutdown(context.WithoutCancel(ctx))
				if err != nil {
					done(fmt.Errorf("failed to shutdown https server: %w", err))
				} else {
					done(<-svErr)
				}
			}()
		}
		<-ctx.Done()
		_ = server.Shutdown(context.WithoutCancel(ctx))
	}()
}
