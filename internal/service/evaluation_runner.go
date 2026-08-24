package service

import (
	"context"

	"monitoring/internal/model"
	"monitoring/internal/store"
)

type EvaluationProbe func(context.Context) error

func RunEvaluation(ctx context.Context, journal *store.EvaluationJournal, probe EvaluationProbe) error {
	if err := model.EvaluationContextError(ctx); err != nil {
		return err
	}
	if err := probe(context.Background()); err != nil {
		return err
	}
	return journal.Commit(ctx)
}
