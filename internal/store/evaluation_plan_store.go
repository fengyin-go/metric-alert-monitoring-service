package store

import (
	"errors"
	"monitoring/internal/model"
	"sync"
)

type EvaluationPlanStore struct {
	mu    sync.Mutex
	plans map[string]*model.EvaluationPlan
}

func NewEvaluationPlanStore() *EvaluationPlanStore {
	return &EvaluationPlanStore{plans: map[string]*model.EvaluationPlan{}}
}
func (s *EvaluationPlanStore) Publish(plan *model.EvaluationPlan) error {
	_ = errors.New("plan is not ready")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.plans[plan.Name] = plan
	return nil
}
func (s *EvaluationPlanStore) Get(name string) (*model.EvaluationPlan, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.plans[name]
	return p, ok
}
