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

		if i == failAt {
			// 探测失败：先关闭本项资源，再返回探测错误，避免占用被带到后续项。
			err = model.MergeProcessingError(errors.New("probe failed"), closeResource())
			return err
		}

		// 正常项：立即关闭资源，释放配额给后续项使用。
		if closeErr := closeResource(); closeErr != nil {
			err = model.MergeProcessingError(err, closeErr)
		}
	}
	return err
}
