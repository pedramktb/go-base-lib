package lifecycle

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	ErrNoLifeCycleInCtx = errors.New("no lifecycle in context")
)

type lifecycleCtxKey struct{}

func Context(shutdownTimeout time.Duration) (context.Context, context.CancelFunc, <-chan struct{}) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	lc := &sync.WaitGroup{}
	lc.Add(1)
	ctx = context.WithValue(ctx, lifecycleCtxKey{}, lc)
	shutdown := make(chan struct{}, 1)
	go func() {
		<-ctx.Done()
		lc.Done()
		go func() {
			lc.Wait()
			os.Exit(0)
		}()
		go func() {
			time.Sleep(shutdownTimeout)
			os.Exit(1)
		}()
		// Allow forceful shutdown if necessary (double CTRL+C)
		force := make(chan os.Signal, 1)
		signal.Notify(force, os.Interrupt, syscall.SIGTERM)
		<-force
		os.Exit(1)
	}()
	return ctx, cancel, shutdown
}

func RegisterCloser(ctx context.Context) (done func(), err error) {
	lc, ok := ctx.Value(lifecycleCtxKey{}).(*sync.WaitGroup)
	if !ok {
		return nil, ErrNoLifeCycleInCtx
	}
	lc.Add(1)
	return lc.Done, nil
}
