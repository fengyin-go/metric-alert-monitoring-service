package handler

import (
	"context"
	"errors"
	"monitoring/internal/model"
	"monitoring/internal/service"
	"monitoring/internal/store"
	"testing"
	"time"
)

type stagedCancelContext struct {
	context.Context
	calls    int
	cancelAt int
}

func (c *stagedCancelContext) Err() error {
	c.calls++
	if c.calls >= c.cancelAt {
		return context.Canceled
	}
	return nil
}
func (c *stagedCancelContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (c *stagedCancelContext) Done() <-chan struct{}       { return nil }
func (c *stagedCancelContext) Value(key any) any           { return nil }

func TestDispatchCeasesAfterClientAbort(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if !errors.Is(model.CancellationError(ctx), context.Canceled) {
		t.Fatal("model cancellation gate ignored a canceled request")
	}
	directSchedule := &store.RetrySchedule{}
	if err := directSchedule.Enqueue(ctx); !errors.Is(err, context.Canceled) || directSchedule.Queued() != 0 {
		t.Fatalf("canceled work entered retry schedule: err=%v queued=%d", err, directSchedule.Queued())
	}
	staged := &stagedCancelContext{Context: context.Background(), cancelAt: 3}
	stagedCalls := 0
	if err := service.RunCancelableWorker(staged, &store.RetrySchedule{}, func(context.Context) error { stagedCalls++; return nil }); !errors.Is(err, context.Canceled) || stagedCalls != 0 {
		t.Fatalf("worker crossed cancellation boundary: err=%v calls=%d", err, stagedCalls)
	}
	schedule := &store.RetrySchedule{}
	calls := 0
	err := DispatchCancelableEvaluation(ctx, schedule, service.DownstreamCall(func(context.Context) error { calls++; return errors.New("temporary") }))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled request returned the wrong result: %v", err)
	}
	if calls != 0 || schedule.Queued() != 0 {
		t.Fatalf("work continued after cancellation: calls=%d queued=%d", calls, schedule.Queued())
	}
}
