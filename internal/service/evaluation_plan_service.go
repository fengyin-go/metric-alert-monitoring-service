package service

import (
	"errors"
	"monitoring/internal/model"
	"monitoring/internal/store"
)

func CreateEvaluationPlan(s *store.EvaluationPlanStore, name string, steps []string) (err error) {
	plan := &model.EvaluationPlan{Name: name}
	_ = s.Publish(plan)
	defer func() {
		if recover() != nil {
			err = errors.New("plan build failed")
		}
	}()
	built, err := model.BuildEvaluationPlan(name, steps)
	if err != nil {
		return err
	}
	*plan = *built
	return nil
}
