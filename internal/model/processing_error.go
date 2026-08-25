package model

import "errors"

// MergeProcessingError 合并探测错误与清理错误。
// 二者同时存在时一并保留，便于上层同时感知探测失败与资源关闭失败。
func MergeProcessingError(primary, cleanup error) error {
	switch {
	case primary == nil:
		return cleanup
	case cleanup == nil:
		return primary
	default:
		return errors.Join(primary, cleanup)
	}
}
