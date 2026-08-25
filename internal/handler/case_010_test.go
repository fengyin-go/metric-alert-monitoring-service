package handler

import (
	"errors"
	"monitoring/internal/model"
	"monitoring/internal/service"
	"monitoring/internal/store"
	"strings"
	"testing"
)

func TestFailedBatchReleasesResourcesAndPreservesFailure(t *testing.T) {
	primary := errors.New("probe failed")
	cleanup := errors.New("close failed")
	if merged := model.MergeProcessingError(primary, cleanup); !errors.Is(merged, primary) || !errors.Is(merged, cleanup) {
		t.Fatalf("error merge discarded processing or cleanup failure: %v", merged)
	}
	stepLimiter := &service.ResourceLimiter{Limit: 1}
	if stepErr := service.ProcessResourceBatch(stepLimiter, 2, 1); stepErr == nil || !strings.Contains(stepErr.Error(), "probe failed") || stepLimiter.Open != 0 {
		t.Fatalf("loop retained a resource before the next item: err=%v open=%d", stepErr, stepLimiter.Open)
	}
	tx := &store.ProcessingTransaction{}
	audit := &ProcessingAudit{}
	limiter := &service.ResourceLimiter{Limit: 1}
	err := RunAuditedBatch(tx, audit, limiter, 2, 0)
	if err == nil || !strings.Contains(err.Error(), "probe failed") {
		t.Fatalf("original processing failure was lost: %v", err)
	}
	if tx.Committed() || !tx.RolledBack() {
		t.Fatalf("failed batch transaction state is wrong: committed=%v rolledBack=%v", tx.Committed(), tx.RolledBack())
	}
	if audit.Success != 0 || audit.Failure != 1 {
		t.Fatalf("audit was published before outcome: %+v", audit)
	}
	if limiter.Open != 0 {
		t.Fatalf("batch resources were not released: open=%d", limiter.Open)
	}
}
