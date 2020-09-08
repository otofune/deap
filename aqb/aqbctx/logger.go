package aqbctx

import (
	"context"

	"github.com/otofune/automate-eamusement-playshare/aqb/logger"
)

type loggerKey struct{}

// WithLogger embeds logger to context
func WithLogger(ctx context.Context, l logger.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// Logger extracts embed aqb Logger from context.
// Return logger.NoOpLogger if no logger in context.
func Logger(ctx context.Context) logger.Logger {
	v := ctx.Value(loggerKey{})
	if logger, ok := v.(logger.Logger); ok {
		return logger
	}
	return &logger.NoOpLogger{}
}
