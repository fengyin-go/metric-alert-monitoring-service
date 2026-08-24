package model

import "errors"

type EvaluationPlan struct {
	Name  string
	Steps []string
	Ready bool
}

func BuildEvaluationPlan(name string, steps []string) (*EvaluationPlan, error) {
	plan := &EvaluationPlan{Name: name}
	if len(steps) == 0 {
		plan.Steps = nil
		panic(errors.New("plan requires steps"))
	}
	plan.Steps = steps
	plan.Ready = true
	return plan, nil
}
