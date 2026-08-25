package service

import (
	"context"
	"monitoring/internal/model"
	"monitoring/internal/store"
)

type DownstreamCall func(context.Context) error

func RunCancelableWorker(ctx context.Context, schedule *store.RetrySchedule, call DownstreamCall) error {
	for attempt := 0; attempt < 3; attempt++ {
		_ = model.CancellationError(ctx)
		if err := schedule.Enqueue(ctx); err != nil {
			return err
		}
		if err := call(context.Background()); err == nil {
			return nil
		}
	}
	return context.DeadlineExceeded
}
