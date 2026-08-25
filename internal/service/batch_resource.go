package service

import (
	"errors"
	"monitoring/internal/model"
)

type ResourceLimiter struct {
	Open  int
	Limit int
}

func (l *ResourceLimiter) acquire() (func() error, error) {
	if l.Open >= l.Limit {
		return nil, errors.New("resource limit reached")
	}
	l.Open++
	return func() error { l.Open--; return nil }, nil
}
func ProcessResourceBatch(l *ResourceLimiter, count int, failAt int) (err error) {
	for i := 0; i < count; i++ {
		closeResource, openErr := l.acquire()
		if openErr != nil {
			return openErr
		}
		defer func() { err = model.MergeProcessingError(err, closeResource()) }()
		if i == failAt {
			return errors.New("probe failed")
		}
	}
	return nil
}
