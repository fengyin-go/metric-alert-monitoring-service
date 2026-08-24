package store

import "sync"

type DeliveryTransaction struct {
	mu        sync.Mutex
	staged    int
	committed int
}

func (t *DeliveryTransaction) Stage() { t.mu.Lock(); defer t.mu.Unlock(); t.staged++ }
func (t *DeliveryTransaction) Finish(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if err != nil {
		t.committed += t.staged
		t.staged = 0
		return
	}
	t.committed += t.staged
	t.staged = 0
}
func (t *DeliveryTransaction) Committed() int { t.mu.Lock(); defer t.mu.Unlock(); return t.committed }
