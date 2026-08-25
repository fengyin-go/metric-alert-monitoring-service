package model

import "context"

// CancellationError 返回 ctx 的取消原因。
// ctx 为 nil 视为已取消；否则透传 ctx.Err()，
// 这样调用方只需判断返回值是否为 nil 即可知道是否已取消。
func CancellationError(ctx context.Context) error {
	if ctx == nil {
		return context.Canceled
	}
	return ctx.Err()
}
