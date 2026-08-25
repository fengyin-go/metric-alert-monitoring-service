package service

import (
	"context"
	"errors"

	"monitoring/internal/model"
	"monitoring/internal/store"
)

// DownstreamCall 是一次下游调用。
type DownstreamCall func(context.Context) error

// RunCancelableWorker 按照重试计划执行下游调用，并尊重取消语义：
//   - 分派前已取消：直接返回取消原因，既不入队也不发起下游调用；
//   - 刚入队即被取消：撤销刚登记的任务并返回取消原因，不启动下游调用；
//   - 下游执行过程中取消：透传 ctx，下游可及时返回，且不再发起后续重试；
//   - 每次尝试结束时移出队列，保证 Queued() 只反映在途任务，不残留已完成/已取消的任务。
func RunCancelableWorker(ctx context.Context, schedule *store.RetrySchedule, call DownstreamCall) error {
	if err := model.CancellationError(ctx); err != nil {
		return err
	}
	for attempt := 0; attempt < 3; attempt++ {
		// 入队前再次确认未取消，避免取消后继续派活。
		if err := model.CancellationError(ctx); err != nil {
			return err
		}
		if err := schedule.Enqueue(ctx); err != nil {
			// Enqueue 在取消时不递增计数，此处无需撤销。
			return err
		}
		// 入队成功后、启动下游前再确认：刚入队即被取消也不启动下游。
		if err := model.CancellationError(ctx); err != nil {
			schedule.Drain()
			return err
		}

		// 下游调用复用同一 ctx，取消/关闭信号能即时传导到下游。
		err := call(ctx)
		// 本次尝试结束，移出队列，避免残留被后续请求接管。
		schedule.Drain()

		switch {
		case err == nil:
			return nil
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			// 取消/超时属于终态，停止后续重试。
			return err
		}
		// 普通失败：进入下一轮重试（循环顶部会再次检查取消）。
	}
	return context.DeadlineExceeded
}
