package handler

import (
	"context"
	"monitoring/internal/service"
	"monitoring/internal/store"
)

func DispatchCancelableEvaluation(ctx context.Context, schedule *store.RetrySchedule, call service.DownstreamCall) error {
	return service.RunCancelableWorker(context.Background(), schedule, call)
}
