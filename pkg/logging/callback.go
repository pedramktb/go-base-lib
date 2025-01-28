package logging

import (
	"context"
)

type CBLogger struct {
	callback func(ctx context.Context, level, msg string, fields map[string]any)
}

// CallbackLogger creates a new logger that logs to the given callback function.
// The callback function is called with the context, log level (one of "debug", "info", "warn", "error", "panic", or "fatal"), message, and fields
func CallbackLogger(callback func(ctx context.Context, level, msg string, fields map[string]any)) *CBLogger {
	return &CBLogger{callback: callback}
}

func (l *CBLogger) Debug(ctx context.Context, msg string, fields map[string]any) {
	if l.callback != nil {
		l.log(ctx, "debug", msg, fields)
	} else {
		FromContext(ctx).Debug(msg)
	}
}

func (l *CBLogger) Info(ctx context.Context, msg string, fields map[string]any) {
	if l.callback != nil {
		l.log(ctx, "info", msg, fields)
	} else {
		FromContext(ctx).Info(msg)
	}
}

func (l *CBLogger) Warn(ctx context.Context, msg string, fields map[string]any) {
	if l.callback != nil {
		l.log(ctx, "warn", msg, fields)
	} else {
		FromContext(ctx).Warn(msg)
	}
}

func (l *CBLogger) Error(ctx context.Context, msg string, fields map[string]any) {
	if l.callback != nil {
		l.log(ctx, "error", msg, fields)
	} else {
		FromContext(ctx).Error(msg)
	}
}

func (l *CBLogger) Panic(ctx context.Context, msg string, fields map[string]any) {
	if l.callback != nil {
		l.log(ctx, "panic", msg, fields)
	} else {
		FromContext(ctx).Panic(msg)
	}
}

func (l *CBLogger) Fatal(ctx context.Context, msg string, fields map[string]any) {
	if l.callback != nil {
		l.log(ctx, "fatal", msg, fields)
	} else {
		FromContext(ctx).Fatal(msg)
	}
}

func (l *CBLogger) log(ctx context.Context, level, msg string, fields map[string]any) {
	m := ctxToMap(ctx)
	for k, v := range fields {
		m[k] = v
	}
	l.callback(ctx, level, msg, m)
}
