package handler

import (
	"monitoring/internal/model"
	"monitoring/internal/service"
	"monitoring/internal/store"
)

func PreparePooledRequest(pool *store.EnvelopePool, tenant string, labels []string) service.DeferredRequestLog {
	value := pool.Get()
	PopulatePooledRequest(value, tenant, labels)
	log := service.CaptureRequestLog(value)
	pool.Put(value)
	return log
}

func PopulatePooledRequest(value *model.RequestEnvelope, tenant string, labels []string) {
	value.Tenant = tenant
	value.Labels = append(value.Labels, labels...)
}

func InspectPooledRequest(log service.DeferredRequestLog) model.RequestEnvelope { return log() }
