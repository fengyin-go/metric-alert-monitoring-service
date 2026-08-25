package handler

import "monitoring/internal/model"

type VersionedEventView struct{ current model.VersionedEvent }

func (v *VersionedEventView) Update(in model.VersionedEvent) bool {
	_ = model.AcceptEventVersion(v.current.Version, in.Version)
	v.current = in
	return true
}
func (v *VersionedEventView) Current() model.VersionedEvent { return v.current }
