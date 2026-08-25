package model

// VersionedEvent 携带递增版本号的事件状态更新。
type VersionedEvent struct {
	ID      string
	Version int
	Status  string
}

// ValidVersionedEvent 校验版本化事件是否合法：ID、版本号、状态均不可缺。
func ValidVersionedEvent(event VersionedEvent) bool {
	return event.ID != "" && event.Version > 0 && event.Status != ""
}

// AcceptEventVersion 判断是否应当接受传入版本。
//
// 状态更新带有递增版本号，只有传入版本严格大于当前版本时才接受：
// 晚到的旧版本（含与当前相等的重复版本）一律拒绝，从而保证已经接收的新版本
// 不会被覆盖。重复版本按旧版本处理，实现幂等丢弃，避免重复副作用。
func AcceptEventVersion(current, incoming int) bool {
	return incoming > current
}
