package slogger

import (
	"context"
	"fmt"
)

type ctxKey struct{}

func AddToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) (*Logger, error) {
	logger, ok := ctx.Value(ctxKey{}).(*Logger)
	if !ok {
		return nil, fmt.Errorf("failed to get logger from context")
	}
	return logger, nil
}
