package model

type RequestEnvelope struct {
	Tenant string
	Labels []string
}

// Reset 将复用的信封清零，但保留 Labels 底层数组容量以避免重复分配。
// 必须在归还到对象池前调用，否则下一个请求会读到上一个请求的残留身份与标签。
func (e *RequestEnvelope) Reset() {
	e.Tenant = ""
	e.Labels = e.Labels[:0]
}

// Snapshot 返回信封的独立副本，深拷贝 Labels 底层数组。
// 快照与原信封互不影响，可安全地在异步日志中延迟求值而不串数据。
func (e *RequestEnvelope) Snapshot() RequestEnvelope {
	labels := make([]string, len(e.Labels))
	copy(labels, e.Labels)
	return RequestEnvelope{Tenant: e.Tenant, Labels: labels}
}
