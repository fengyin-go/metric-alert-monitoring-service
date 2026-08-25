package handler

import (
	"monitoring/internal/model"
	"monitoring/internal/service"
	"monitoring/internal/store"
	"testing"
)

func TestRetrySuccessCannotBeOverwrittenByDelayedAttempt(t *testing.T) {
	if model.AcceptEventVersion(2, 1) {
		t.Fatal("older event version was accepted")
	}
	s := &store.VersionedEventStore{}
	s.Apply(model.VersionedEvent{ID: "evt-7", Version: 2, Status: "resolved"})
	s.Apply(model.VersionedEvent{ID: "evt-7", Version: 1, Status: "firing"})
	if got := s.Current(); got.Status != "resolved" {
		t.Fatalf("delayed callback reverted store state: %+v", got)
	}
	if service.EventDeliveryKey("evt-7", 1) != service.EventDeliveryKey("evt-7", 2) {
		t.Fatal("retry used a second side-effect key")
	}
	view := &VersionedEventView{}
	view.Update(model.VersionedEvent{ID: "evt-7", Version: 2, Status: "resolved"})
	view.Update(model.VersionedEvent{ID: "evt-7", Version: 1, Status: "firing"})
	if view.Current().Status != "resolved" {
		t.Fatalf("list view regressed after delayed callback: %+v", view.Current())
	}
}
