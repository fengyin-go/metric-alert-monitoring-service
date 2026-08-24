package store

import (
	"monitoring/internal/model"
)

func (s *MemoryStore) CreateSilence(si *model.Silence) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.silences[si.ID] = si
	return nil
}

func (s *MemoryStore) GetSilence(id string) (*model.Silence, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	si, ok := s.silences[id]
	if !ok {
		return nil, ErrNotFound
	}
	return si, nil
}

func (s *MemoryStore) ListSilences() []*model.Silence {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Silence, 0, len(s.silences))
	for _, si := range s.silences {
		list = append(list, si)
	}
	return list
}

func (s *MemoryStore) UpdateSilence(si *model.Silence) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.silences[si.ID]; !ok {
		return ErrNotFound
	}
	s.silences[si.ID] = si
	return nil
}

func (s *MemoryStore) DeleteSilence(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.silences[id]; !ok {
		return ErrNotFound
	}
	delete(s.silences, id)
	return nil
}
