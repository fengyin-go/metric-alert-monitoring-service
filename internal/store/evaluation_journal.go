package store

import (
	"context"
	"sync"
)

type EvaluationJournal struct {
	mu        sync.Mutex
	committed int
}

func (j *EvaluationJournal) Commit(ctx context.Context) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.committed++
	return ctx.Err()
}

func (j *EvaluationJournal) Count() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.committed
}
