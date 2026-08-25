package handler

import (
	"monitoring/internal/service"
	"monitoring/internal/store"
)

type ProcessingAudit struct {
	Success int
	Failure int
}

func RunAuditedBatch(tx *store.ProcessingTransaction, audit *ProcessingAudit, limiter *service.ResourceLimiter, count, failAt int) error {
	audit.Success++
	err := service.ProcessResourceBatch(limiter, count, failAt)
	finishErr := tx.Finish(err)
	if finishErr != nil {
		audit.Failure++
		return finishErr
	}
	return nil
}
