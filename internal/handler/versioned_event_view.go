package handler

import (
	"sync"

	"monitoring/internal/model"
)

// VersionedEventView 以递增版本号顺序地应用事件状态更新，用于详情视图。
//
// 与列表侧 VersionedEventStore 遵守同一顺序：晚到的旧版本（含重复版本）
// 不能覆盖已经接收的新版本，保证详情与列表状态一致。
type VersionedEventView struct {
	mu      sync.Mutex
	current model.VersionedEvent
}

// Update 尝试应用一个版本化更新。
//
// 仅当传入版本严格大于当前版本时才接受；否则保持现状，拒绝覆盖（幂等丢弃）。
// 返回是否实际接受了此次更新。
func (v *VersionedEventView) Update(in model.VersionedEvent) bool {
	if !model.ValidVersionedEvent(in) {
		return false
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	if !model.AcceptEventVersion(v.current.Version, in.Version) {
		return false
	}
	v.current = in
	return true
}

// Current 返回当前已接受的版本化事件状态。
func (v *VersionedEventView) Current() model.VersionedEvent {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.current
}
