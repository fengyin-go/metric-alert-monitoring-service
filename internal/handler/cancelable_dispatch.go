package handler

import (
	"context"
	"monitoring/internal/service"
	"monitoring/internal/store"
)

// DispatchCancelableEvaluation 将可取消的规则评估分派到后台 worker。
// 透传 ctx，使请求取消/服务关闭能即时传导到 worker 与下游调用。
func DispatchCancelableEvaluation(ctx context.Context, schedule *store.RetrySchedule, call service.DownstreamCall) error {
	return service.RunCancelableWorker(ctx, schedule, call)
}
