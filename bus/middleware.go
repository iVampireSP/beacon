package bus

import (
	"context"
	"fmt"
	"time"

	"github.com/iVampireSP/beacon/logger"
)

type envelopeContextKey struct{}

func ContextWithEnvelope(ctx context.Context, env *Envelope) context.Context {
	return context.WithValue(ctx, envelopeContextKey{}, env)
}

func EnvelopeFromContext(ctx context.Context) *Envelope {
	if env, ok := ctx.Value(envelopeContextKey{}).(*Envelope); ok {
		return env
	}
	return nil
}

func recovery() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, payload []byte) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic recovered: %v", r)
					logger.Error("event handler panic", "panic", r)
				}
			}()
			return next(ctx, payload)
		}
	}
}

func logging() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, payload []byte) error {
			env := EnvelopeFromContext(ctx)
			start := time.Now()
			err := next(ctx, payload)
			duration := time.Since(start)
			if env != nil {
				if err != nil {
					logger.Warn("event processing failed",
						"name", env.Name,
						"id", env.ID,
						"duration", duration,
						"error", err,
					)
				} else {
					logger.Debug("event processed",
						"name", env.Name,
						"id", env.ID,
						"duration", duration,
					)
				}
			}
			return err
		}
	}
}
