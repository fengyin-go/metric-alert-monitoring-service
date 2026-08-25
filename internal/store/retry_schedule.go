package store

import (
	"context"
	"sync"
)

// RetrySchedule 跟踪等待执行的下游重试任务。
//
// 语义：任务只有未被取消时才允许入队，因此 Queued() 反映的是
// 真正待执行的下游任务数；任务入队后若立即被取消，调用方应调用
// Drain 撤销该任务，避免残留任务被后续请求接管。
type RetrySchedule struct {
	mu     sync.Mutex
	queued int
}

// Enqueue 尝试将一次下游任务登记入队。
// 已取消则不入队并返回取消原因，保证取消后不再派活。
func (s *RetrySchedule) Enqueue(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	s.queued++
	s.mu.Unlock()
	return nil
}

// Drain 撤销一次已登记的任务，用于任务刚入队即被取消的场景，
// 使队列计数回到入队前的状态，避免残留任务。
func (s *RetrySchedule) Drain() {
	s.mu.Lock()
	if s.queued > 0 {
		s.queued--
	}
	s.mu.Unlock()
}

// Queued 返回当前待执行的下游任务数。
func (s *RetrySchedule) Queued() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queued
}
