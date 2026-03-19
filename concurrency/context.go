package concurrency

import (
	"context"
	"fmt"
	"time"
)

func WithTimeout(ctx context.Context, timeout time.Duration, fn func(ctx context.Context) error) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		if err := fn(timeoutCtx); err != nil {
			done <- err
		}
		close(done)
	}()

	select {
	case err := <-done:
		if err != nil {
			return err
		}
		return nil
	case <-timeoutCtx.Done():
		return fmt.Errorf("shutdown timed out")
	}
}
