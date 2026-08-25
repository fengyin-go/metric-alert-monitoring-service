// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"monitoring/internal/model"
)

var (
	// ErrNotFound 记录不存在。
	ErrNotFound = errors.New("记录不存在")
	// ErrConflict 记录已存在或状态冲突。
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法。
type Store interface {
	// Metric
	CreateMetric(m *model.Metric) error
	GetMetric(id string) (*model.Metric, error)
	GetMetricByName(name string) (*model.Metric, error)
	ListMetrics() []*model.Metric
	UpdateMetric(m *model.Metric) error
	DeleteMetric(id string) error

	// AlertRule
	CreateAlertRule(r *model.AlertRule) error
	GetAlertRule(id string) (*model.AlertRule, error)
	ListAlertRules() []*model.AlertRule
	UpdateAlertRule(r *model.AlertRule) error
	DeleteAlertRule(id string) error

	// AlertEvent
	CreateAlertEvent(e *model.AlertEvent) error
	GetAlertEvent(id string) (*model.AlertEvent, error)
	ListAlertEvents() []*model.AlertEvent
	UpdateAlertEvent(e *model.AlertEvent) error
	DeleteAlertEvent(id string) error

	// Silence
	CreateSilence(s *model.Silence) error
	GetSilence(id string) (*model.Silence, error)
	ListSilences() []*model.Silence
	UpdateSilence(s *model.Silence) error
	DeleteSilence(id string) error
}
