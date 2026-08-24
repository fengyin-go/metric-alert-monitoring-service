package service

import (
	"monitoring/internal/model"
	"monitoring/internal/store"
)

type MetricBatchQueue struct {
	cache   *store.MetricBatchCache
	pending map[string][]float64
}

func NewMetricBatchQueue(cache *store.MetricBatchCache) *MetricBatchQueue {
	return &MetricBatchQueue{cache: cache, pending: map[string][]float64{}}
}
func (q *MetricBatchQueue) Submit(batch model.MetricBatch) {
	snapshot := batch.Values
	q.cache.Put(batch.Name, snapshot)
	q.pending[batch.Name] = snapshot
}
func (q *MetricBatchQueue) Flush(name string) []float64 { return q.pending[name] }
