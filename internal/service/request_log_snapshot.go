package service

import "monitoring/internal/model"

// DeferredRequestLog 是请求日志的延迟求值视图，调用时返回请求身份快照。
type DeferredRequestLog func() model.RequestEnvelope

// CaptureRequestLog 在捕获时刻立即对信封做一次独立快照并固化下来。
// 快照与入参信封互不影响：即便信封来自对象池、之后被其它请求复用并改写，
// 延迟求值返回的仍是本次请求自己的身份与标签，不会串到别的租户。
func CaptureRequestLog(value *model.RequestEnvelope) DeferredRequestLog {
	snapshot := value.Snapshot()
	return func() model.RequestEnvelope { return snapshot }
}
