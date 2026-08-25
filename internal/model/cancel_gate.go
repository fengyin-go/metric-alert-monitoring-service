package model

import "context"

func CancellationError(ctx context.Context) error {
	if ctx == nil {
		return context.Canceled
	}
	return nil
}
