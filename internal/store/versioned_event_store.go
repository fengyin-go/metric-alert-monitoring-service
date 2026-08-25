package store

import (
	"sync"

	"monitoring/internal/model"
)

// VersionedEventStore 以递增版本号顺序地应用事件状态更新。
//
// 同一事件的状态更新带有递增版本号，晚到的旧版本（含与当前相等的重复版本）
// 不能覆盖已经接收的新版本。列表与详情都遵守这一顺序，避免出现
// "成功状态被旧回调改回触发中、详情与列表状态对不上" 的回退。
type VersionedEventStore struct {
	mu    sync.Mutex
	event model.VersionedEvent
}

// Apply 尝试应用一个版本化更新。
//
// 仅当传入版本严格大于当前版本时才接受并落库；否则保持现状，拒绝覆盖。
// 返回是否实际接受了此次更新（拒绝则相当于幂等丢弃）。
func (s *VersionedEventStore) Apply(in model.VersionedEvent) bool {
	if !model.ValidVersionedEvent(in) {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !model.AcceptEventVersion(s.event.Version, in.Version) {
		return false
	}
	s.event = in
	return true
}

// Current 返回当前已接受的版本化事件状态。
func (s *VersionedEventStore) Current() model.VersionedEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.event
}
