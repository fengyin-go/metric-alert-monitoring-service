package store

import "sync"

type MetricBatchCache struct {
	mu     sync.Mutex
	values map[string][]float64
}

func NewMetricBatchCache() *MetricBatchCache {
	return &MetricBatchCache{values: map[string][]float64{}}
}
func (c *MetricBatchCache) Put(name string, values []float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[name] = values
}
func (c *MetricBatchCache) Get(name string) []float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.values[name]
}
