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
func (p *EnvelopePool) Put(value *model.RequestEnvelope) {
	// 归还前清零，确保下次复用时不会残留上一个租户的身份与标签。
	value.Reset()
	p.pool.Put(value)
}
