package service

import "monitoring/internal/model"

type DeferredRequestLog func() model.RequestEnvelope

func CaptureRequestLog(value *model.RequestEnvelope) DeferredRequestLog {
	return func() model.RequestEnvelope { return value.Snapshot() }
}
