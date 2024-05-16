package globalContext

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
)

func init() {
	ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

func Context() context.Context {
	return ctx
}

func Cancel() {
	cancel()
	wg.Wait()
}

func RegisterListener() (ListenerDone func()) {
	wg.Add(1)
	return wg.Done
}
