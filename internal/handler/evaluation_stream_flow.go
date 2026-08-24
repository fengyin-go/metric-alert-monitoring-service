package handler

import (
	"context"
	"monitoring/internal/model"
	"monitoring/internal/service"
	"monitoring/internal/store"
)

func RunEvaluationStream(ctx context.Context, jobs []model.EvaluationJob) ([]string, error) {
	results, failures := service.StartEvaluationStream(ctx, jobs)
	values, err := store.CollectEvaluations(ctx, results, failures)
	if ctx.Err() != nil {
		return values, ctx.Err()
	}
	return values, err
}
