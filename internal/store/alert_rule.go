package store

import (
	"monitoring/internal/model"
)

func (s *MemoryStore) CreateAlertRule(r *model.AlertRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[r.ID] = r
	return nil
}

func (s *MemoryStore) GetAlertRule(id string) (*model.AlertRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListAlertRules() []*model.AlertRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AlertRule, 0, len(s.rules))
	for _, r := range s.rules {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateAlertRule(r *model.AlertRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[r.ID]; !ok {
		return ErrNotFound
	}
	s.rules[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteAlertRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[id]; !ok {
		return ErrNotFound
	}
	delete(s.rules, id)
	return nil
}
