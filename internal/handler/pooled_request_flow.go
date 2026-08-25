package handler

import (
	"monitoring/internal/model"
	"monitoring/internal/service"
	"monitoring/internal/store"
)

// PreparePooledRequest 用对象池里的信封承载一次请求的身份与标签，
// 并返回该次请求专属的延迟日志视图。
//
// 注意：CaptureRequestLog 必须在信封被复用之前完成快照，否则下一个
// 租户会覆盖同一个信封，导致延迟写入的日志串到别的租户。
func PreparePooledRequest(pool *store.EnvelopePool, tenant string, labels []string) service.DeferredRequestLog {
	value := pool.Get()
	PopulatePooledRequest(value, tenant, labels)
	log := service.CaptureRequestLog(value) // 在归还/复用之前固化本次快照
	pool.Put(value)
	return log
}

// PopulatePooledRequest 覆写式地写入身份与标签，绝不沿用信封里残留的旧值。
// 旧实现用 append，会在对象池复用场景下把上一个请求的标签累加进来。
func PopulatePooledRequest(value *model.RequestEnvelope, tenant string, labels []string) {
	value.Tenant = tenant
	value.Labels = append(value.Labels[:0], labels...)
}

func InspectPooledRequest(log service.DeferredRequestLog) model.RequestEnvelope { return log() }
