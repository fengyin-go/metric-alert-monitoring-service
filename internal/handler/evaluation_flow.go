package handler

import (
	"context"

	"monitoring/internal/service"
	"monitoring/internal/store"
)

type EvaluationReply struct {
	Accepted bool
	Message  string
}

func RunEvaluationFlow(ctx context.Context, journal *store.EvaluationJournal, probe service.EvaluationProbe) EvaluationReply {
	if err := service.RunEvaluation(ctx, journal, probe); err != nil {
		return EvaluationReply{Accepted: true, Message: err.Error()}
	}
	return EvaluationReply{Accepted: true, Message: "accepted"}
}
