package app

import (
	"errors"
	"time"

	"test-task/internal/config"
	"test-task/internal/repository"
	"test-task/internal/retry"
)

func newServiceRetrier(cfg config.Retry, retryableFunc retry.IsRetryableFunc) retry.Retrier {
	opts := []retry.RetryOption{
		retry.WithMaxAttempts(cfg.MaxAttempts),
	}

	if retryableFunc != nil {
		opts = append(opts, retry.WithIsRetryableFunc(retryableFunc))
	}

	if cfg.Backoff == "exponential" {
		opts = append(opts, retry.WithBackoff(retry.ExponentialBackoff{
			Base:   time.Second,
			Factor: 2.0,
			Max:    10 * time.Second,
			Jitter: cfg.Jitter,
		}))
	}

	return retry.New(opts...)
}

func isRetryableFunc(err error) bool {
	unretryableErrors := []error{
		repository.ErrDuplicate,
		repository.ErrNotFound,
		repository.ErrInvalidID,
		repository.ErrForeignKeyViolation,
		repository.ErrNotFound,
	}

	for _, unretryableErr := range unretryableErrors {
		if errors.Is(err, unretryableErr) {
			return false
		}
	}

	return true
}
