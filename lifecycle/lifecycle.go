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

// Context returns a context with a lifecycle. The lifecycle is a slice of wait groups.
// The lifecycle is used to wait for all goroutines to finish before exiting the program.
// The wait is limited to the given timeout or a terminate or an interrupt signal after the initial
// cancelation. This means closing a running program immediately from shell, requires 2 interrupts.
// The additional CancelFunc can be used the start the shutdown process from inside the program.
// The additional Channel should be used to prevent main from exiting and terminating all goroutines.
func Context(shutdownTimeout time.Duration) (context.Context, context.CancelFunc, <-chan struct{}) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	lc := []*sync.WaitGroup{{}}
	lc[0].Add(1)
	ctx = context.WithValue(ctx, lifecycleCtxKey{}, &lc)
	shutdown := make(chan struct{}, 1)
	go func() {
		<-ctx.Done()
		lc[0].Done()
		go func() {
			for i := range lc {
				lc[i].Wait()
			}
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

// RegisterCloser registers a new closer to be called when the context is cancelled.
// It uses appends a new WaitGroup to the lifecycle WaitGroups array in context.
// It is not thread-safe and should only be called from a single goroutine at a time.
// It is also up to the caller that the underlying WaitGroup counter stays positive.
// (Done must only be called once per Closer.)
func RegisterCloser(ctx context.Context) (done func(), err error) {
	lc, ok := ctx.Value(lifecycleCtxKey{}).(*[]*sync.WaitGroup)
	if !ok {
		return nil, ErrNoLifeCycleInCtx
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	*lc = append(*lc, wg)
	return wg.Done, nil
}
