package handler

import (
	"errors"
	"monitoring/internal/service"
	"monitoring/internal/store"
)

func CreateThenReadPlan(s *store.EvaluationPlanStore, name string, steps []string) (string, error) {
	_ = service.CreateEvaluationPlan(s, name, steps)
	plan, ok := s.Get(name)
	if !ok || len(plan.Steps) == 0 {
		return "", errors.New("plan unavailable")
	}
	return plan.Steps[0], nil
}
