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
	ErrNoLifecycleInCtx = errors.New("no lifecycle in context")
)

type lifecycleCtxKey struct{}

type lifecycleData struct {
	wgs  []*sync.WaitGroup
	errs []error
}

// Context returns a context with a lifecycle. The lifecycle is a slice of wait groups.
// The lifecycle is used to wait for all goroutines to finish before exiting the program.
// The wait is limited to the given timeout or a terminate or an interrupt signal after the initial
// cancelation. This means closing a running program immediately from shell, requires 2 interrupts.
// The additional CancelFunc can be used the start the shutdown process from inside the program.
// The additional Channel should be used to prevent main from exiting and receive the shutdown errors.
// This channel is closed when all errors from all closed goroutines are received. See the example below:
//
//	for err := range shutdownErrsChan {
//		if err != nil {
//			errors = append(errors, err)
//		}
//	}
//	if len(errors) > 0 {
//		logger.Error("one or more modules failed to shutdown properly", errors)
//	}
//	// End of main()
func Context(shutdownTimeout time.Duration) (context.Context, context.CancelFunc, <-chan error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	lc := &lifecycleData{
		wgs:  []*sync.WaitGroup{},
		errs: []error{},
	}
	ctx = context.WithValue(ctx, lifecycleCtxKey{}, lc)
	shutdownErrs := make(chan error)
	go func() {
		<-ctx.Done()
		go func() {
			// Force Shutdown after timeout
			time.Sleep(shutdownTimeout)
			os.Exit(1)
		}()
		go func() {
			// Allow forceful shutdown if necessary (e.g. double CTRL+C)
			force := make(chan os.Signal, 1)
			signal.Notify(force, os.Interrupt, syscall.SIGTERM)
			<-force
			os.Exit(1)
		}()
		for i := range lc.wgs {
			lc.wgs[i].Wait()
			shutdownErrs <- lc.errs[i]
		}
		close(shutdownErrs)
	}()
	return ctx, cancel, shutdownErrs
}

// RegisterCloser registers a new closer to be called when the context is cancelled.
// It uses appends a new WaitGroup to the lifecycle WaitGroups array in context.
// It is not thread-safe and should only be called from a single goroutine at a time.
// It is also up to the caller that the underlying WaitGroup counter stays positive.
// (Done must only be called once per Closer.)
func RegisterCloser(ctx context.Context) (done func(err error), err error) {
	lc, ok := ctx.Value(lifecycleCtxKey{}).(*lifecycleData)
	if !ok {
		return nil, ErrNoLifecycleInCtx
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	idx := len(lc.wgs)
	lc.wgs = append(lc.wgs, wg)
	lc.errs = append(lc.errs, nil)
	return func(err error) {
		lc.errs[idx] = err
		wg.Done()
	}, nil
}
