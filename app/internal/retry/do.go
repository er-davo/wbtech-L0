package retry

import (
	"context"
	"time"
)

// Do выполняет функцию f с повторными попытками.
// - ctx управляет отменой/таймаутом.
// - maxAttempts = 0 → бесконечно, пока ctx не отменится.
func Do(ctx context.Context, maxAttempts int, f AttemptFunc) error {
	var err error

	for attempt := 0; maxAttempts == 0 || attempt < maxAttempts; attempt++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}

		if err = f(attempt); err != nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt)):
		}
	}

	return err
}
