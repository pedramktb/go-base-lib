package ctx

import (
	"context"
	"os/signal"
	"syscall"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func init() {
	ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

func Ctx() context.Context {
	return ctx
}

func Cancel() {
	cancel()
}
