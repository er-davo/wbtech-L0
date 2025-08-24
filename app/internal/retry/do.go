package retry

import (
	"context"
)

// Do выполняет функцию f с повторными попытками.
// - ctx управляет отменой/таймаутом.
// - maxAttempts = 0 → бесконечно, пока ctx не отменится.
func Do(ctx context.Context, maxAttempts int, f AttemptFunc) error {
	return New(
		WithMaxAttempts(maxAttempts),
	).Do(ctx, f)
}
