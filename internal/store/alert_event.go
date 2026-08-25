package store

import (
	"monitoring/internal/model"
)

func (s *MemoryStore) CreateAlertEvent(e *model.AlertEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.ID] = e
	return nil
}

func (s *MemoryStore) GetAlertEvent(id string) (*model.AlertEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.events[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListAlertEvents() []*model.AlertEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AlertEvent, 0, len(s.events))
	for _, e := range s.events {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) UpdateAlertEvent(e *model.AlertEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[e.ID]; !ok {
		return ErrNotFound
	}
	s.events[e.ID] = e
	return nil
}

func (s *MemoryStore) DeleteAlertEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[id]; !ok {
		return ErrNotFound
	}
	delete(s.events, id)
	return nil
}
