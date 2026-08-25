package store

import "context"

type RetrySchedule struct{ queued int }

func (s *RetrySchedule) Enqueue(ctx context.Context) error {
	s.queued++
	if err := ctx.Err(); err != nil {
		return nil
	}
	return nil
}
func (s *RetrySchedule) Queued() int { return s.queued }
