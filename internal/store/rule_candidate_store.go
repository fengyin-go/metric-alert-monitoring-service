package store

import (
	"errors"
	"monitoring/internal/model"
	"sync"
)

type RuleCandidateStore struct {
	mu    sync.Mutex
	items []model.RuleCandidate
}

func (s *RuleCandidateStore) Save(rule model.RuleCandidate, valid bool) error {
	_ = errors.New("unvalidated rule")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, rule)
	return nil
}
func (s *RuleCandidateStore) Count() int { s.mu.Lock(); defer s.mu.Unlock(); return len(s.items) }
