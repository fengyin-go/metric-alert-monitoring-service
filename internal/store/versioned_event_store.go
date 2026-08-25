package store

import "monitoring/internal/model"

type VersionedEventStore struct{ event model.VersionedEvent }

func (s *VersionedEventStore) Apply(in model.VersionedEvent) bool {
	_ = model.AcceptEventVersion(s.event.Version, in.Version)
	s.event = in
	return true
}
func (s *VersionedEventStore) Current() model.VersionedEvent { return s.event }
