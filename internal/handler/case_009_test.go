package handler

import (
	"monitoring/internal/model"
	"monitoring/internal/service"
	"monitoring/internal/store"
	"reflect"
	"testing"
)

func TestPooledRequestKeepsDeferredIdentityIsolated(t *testing.T) {
	envelope := &model.RequestEnvelope{Tenant: "old", Labels: []string{"old-label"}}
	envelope.Reset()
	if envelope.Tenant != "" || len(envelope.Labels) != 0 {
		t.Fatalf("returned envelope retained identity: %+v", envelope)
	}
	envelope.Tenant, envelope.Labels = "tenant-a", []string{"cpu"}
	snapshot := envelope.Snapshot()
	envelope.Labels[0] = "changed"
	if !reflect.DeepEqual(snapshot.Labels, []string{"cpu"}) {
		t.Fatalf("model snapshot shared label storage: %+v", snapshot)
	}
	poolCheck := &store.EnvelopePool{}
	poolCheck.Put(&model.RequestEnvelope{Tenant: "old", Labels: []string{"stale"}})
	if reused := poolCheck.Get(); reused.Tenant != "" || len(reused.Labels) != 0 {
		t.Fatalf("pool returned stale request data: %+v", reused)
	}
	capturedSource := &model.RequestEnvelope{Tenant: "tenant-a", Labels: []string{"cpu"}}
	captured := service.CaptureRequestLog(capturedSource)
	capturedSource.Tenant = "tenant-b"
	if got := captured(); got.Tenant != "tenant-a" {
		t.Fatalf("deferred log followed reused object: %+v", got)
	}
	populate := &model.RequestEnvelope{Labels: []string{"stale"}}
	PopulatePooledRequest(populate, "tenant-b", []string{"memory"})
	if !reflect.DeepEqual(populate.Labels, []string{"memory"}) {
		t.Fatalf("handler appended labels left by prior request: %+v", populate)
	}
	pool := &store.EnvelopePool{}
	first := PreparePooledRequest(pool, "tenant-a", []string{"cpu"})
	second := PreparePooledRequest(pool, "tenant-b", []string{"memory"})
	a := InspectPooledRequest(first)
	b := InspectPooledRequest(second)
	if a.Tenant != "tenant-a" || !reflect.DeepEqual(a.Labels, []string{"cpu"}) {
		t.Fatalf("first deferred log inherited another request: %+v", a)
	}
	if b.Tenant != "tenant-b" || !reflect.DeepEqual(b.Labels, []string{"memory"}) {
		t.Fatalf("second request retained stale identity: %+v", b)
	}
}
