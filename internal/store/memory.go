package store

import (
	"sync"

	"monitoring/internal/model"
)

// MemoryStore 基于内存的 Store 实现。
type MemoryStore struct {
	mu       sync.RWMutex
	metrics  map[string]*model.Metric
	rules    map[string]*model.AlertRule
	events   map[string]*model.AlertEvent
	silences map[string]*model.Silence
}

// NewMemoryStore 创建空的内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		metrics:  make(map[string]*model.Metric),
		rules:    make(map[string]*model.AlertRule),
		events:   make(map[string]*model.AlertEvent),
		silences: make(map[string]*model.Silence),
	}
}

var _ Store = (*MemoryStore)(nil)
