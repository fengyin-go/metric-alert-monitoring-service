package store

import (
	"monitoring/internal/model"
	"sync"
)

type EnvelopePool struct{ pool sync.Pool }

func (p *EnvelopePool) Get() *model.RequestEnvelope {
	value, _ := p.pool.Get().(*model.RequestEnvelope)
	if value == nil {
		value = &model.RequestEnvelope{}
	}
	return value
}
func (p *EnvelopePool) Put(value *model.RequestEnvelope) { p.pool.Put(value) }
