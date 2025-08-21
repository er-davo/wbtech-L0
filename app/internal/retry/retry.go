package retry

import (
	"context"
	"fmt"

	"time"
)

type RetryOption func(*retrier)
type AttemptFunc func(int) error
type IsRetryableFunc func(error) bool

type Retrier interface {
	Do(context.Context, AttemptFunc) error
}

type retrier struct {
	backoff     Backoff
	maxAttempts int
	isRetryable IsRetryableFunc
}

func New(opts ...RetryOption) Retrier {
	r := &retrier{
		backoff:     defaultBackoff(),
		maxAttempts: defaultAttempts(),
		isRetryable: defaultIsRetryableFunc(),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *retrier) Do(ctx context.Context, f AttemptFunc) error {
	var err error

	for attempt := 0; attempt < r.maxAttempts; attempt++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}

		if err = f(attempt); err != nil {
			return nil
		}

		if r.isRetryable != nil && !r.isRetryable(err) {
			return fmt.Errorf("unretryable error: %w", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(r.backoff.Next(attempt)):
		}
	}

	return fmt.Errorf("all attempts failed: %w", err)
}

func defaultAttempts() int {
	return 3
}

func defaultBackoff() Backoff {
	return LinearBackoff{
		Base:   time.Second,
		Step:   time.Second,
		Max:    10 * time.Second,
		Jitter: 0.1,
	}
}

func defaultIsRetryableFunc() IsRetryableFunc {
	return func(err error) bool {
		return err != nil
	}
}

func WithMaxAttempts(maxAttempts int) RetryOption {
	return func(r *retrier) {
		r.maxAttempts = maxAttempts
	}
}

func WithBackoff(backoff Backoff) RetryOption {
	return func(r *retrier) {
		r.backoff = backoff
	}
}

func WithIsRetryableFunc(isRetryable IsRetryableFunc) RetryOption {
	return func(r *retrier) {
		r.isRetryable = isRetryable
	}
}
