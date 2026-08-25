package service

import "fmt"

// EventDeliveryKey 生成同一事件跨重试共用的通知去重身份。
//
// 去重身份只与事件 ID 绑定：同一事件无论被重试多少次、第几次回调到达，
// 都对应同一个 key，从而只能收到一次通知。attempt 参数予以保留以兼容现有
// 两参调用方式，但不参与 key 计算——重试次数变化不应派生出新的去重身份。
func EventDeliveryKey(eventID string, attempt int) string {
	_ = attempt
	return fmt.Sprintf("event:%s:notify", eventID)
}
